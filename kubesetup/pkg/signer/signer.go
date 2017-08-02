package signer

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"errors"
	"io/ioutil"
	"math"
	"math/big"
	"net"
	"time"
)

const (
	// RSAKeySize for all generated keys.
	RSAKeySize = 2048

	// Duration365d represents the default expiry for all certs.
	Duration365d = time.Hour * 24 * 365
)

// Signer is the top level signing object, initialize this object to sign certs.
type Signer struct {
	CaURI   string
	KeyURI  string
	APIHost string
	Cluster string
	OutDir  string

	caCert *x509.Certificate
	caKey  *rsa.PrivateKey
}

// CertConfig is used to sign certificates.
type CertConfig struct {
	// Common name, this is interpreted by Kubernetes as the username.
	CN string

	// Organization, these are groups in Kubernetes.
	Org []string

	IPs []string
	DNS []string

	// Controls the expiry of the cert. Defaults to 1 year.
	ExpiresIn time.Duration
}

// Init does setup on the signer and validates the fields.
func (s *Signer) Init() error {
	if s.Cluster == "" {
		return errors.New("cluster is required")
	}
	if s.APIHost == "" {
		return errors.New("api-host is required")
	}

	var err error
	s.caCert, err = getCa(s.CaURI)
	if err != nil {
		return err
	}
	s.caKey, err = getKey(s.KeyURI)
	if err != nil {
		return err
	}
	return nil
}

func (s *Signer) signCert(certConfig CertConfig) (*rsa.PrivateKey, *x509.Certificate, error) {
	if certConfig.ExpiresIn == 0 {
		certConfig.ExpiresIn = Duration365d
	}

	key, err := rsa.GenerateKey(rand.Reader, RSAKeySize)
	if err != nil {
		return nil, nil, err
	}

	serial, err := rand.Int(rand.Reader, new(big.Int).SetInt64(math.MaxInt64))
	if err != nil {
		return nil, nil, err
	}

	if certConfig.IPs == nil {
		certConfig.IPs = []string{"127.0.0.1"}
	}

	ips := []net.IP{}
	for _, ip := range certConfig.IPs {
		ips = append(ips, net.ParseIP(ip))
	}

	certTmpl := x509.Certificate{
		Subject: pkix.Name{
			CommonName:   certConfig.CN,
			Organization: certConfig.Org,
		},
		DNSNames:     certConfig.DNS,
		IPAddresses:  ips,
		SerialNumber: serial,
		NotBefore:    s.caCert.NotBefore,
		NotAfter:     time.Now().Add(certConfig.ExpiresIn).UTC(),
		KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
	}
	certDERBytes, err := x509.CreateCertificate(rand.Reader, &certTmpl, s.caCert, key.Public(), s.caKey)
	if err != nil {
		return nil, nil, err
	}
	cert, err := x509.ParseCertificate(certDERBytes)
	if err != nil {
		return nil, nil, err
	}
	return key, cert, nil
}

// Sign will generate a certificate and kubeconfig using the provided CertConfig.
func (s *Signer) Sign(certConfig CertConfig) error {
	const kubeconfig = `
apiVersion: v1
kind: Config
current-context: default
clusters:
- name: {{ .clusterName }}
  cluster:
    server: https://{{ .apiHost }}
    certificate-authority: {{ .outDir }}{{ .commonName }}-ca.pem
users:
- name: {{ .commonName }}
  user:
    client-certificate: {{ .outDir }}{{ .commonName }}.pem
    client-key: {{ .outDir }}{{ .commonName }}-key.pem
contexts:
- context:
    cluster: {{ .clusterName }}
    user: {{ .commonName }}
  name: default
`

	if certConfig.CN == "" {
		return errors.New("cn is required")
	}

	outDir := "/etc/kubernetes/secrets/"
	if s.OutDir != "" {
		outDir = s.OutDir
	}

	key, cert, err := s.signCert(certConfig)
	if err != nil {
		return err
	}

	encodedKey := encodePrivateKey(key)
	encodedCert := encodeCertificate(cert)
	encodedCaCert := encodeCertificate(s.caCert)

	kubeData, err := formatMap(kubeconfig, map[string]string{
		"commonName":  certConfig.CN,
		"clusterName": s.Cluster,
		"apiHost":     s.APIHost,
		"outDir":      outDir,
	})
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(outDir+certConfig.CN+"-key.pem", encodedKey, 0600)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(outDir+certConfig.CN+".pem", encodedCert, 0600)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(outDir+certConfig.CN+"-ca.pem", encodedCaCert, 0600)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(outDir+certConfig.CN+".kubeconfig", kubeData, 0600)
	if err != nil {
		return err
	}
	return nil
}

// SignWithData will generate a certificate and kubeconfig including the data as base64 encoded blobs.
func (s *Signer) SignWithData(certConfig CertConfig) ([]byte, error) {
	const kubeconfig = `
apiVersion: v1
kind: Config
current-context: default
clusters:
- name: {{ .clusterName }}
  cluster:
    server: https://{{ .apiHost }}
    certificate-authority-data: "{{ .caData }}"
users:
- name: {{ .commonName }}
  user:
    client-certificate-data: "{{ .certData }}"
    client-key-data: "{{ .keyData }}"
contexts:
- context:
    cluster: {{ .clusterName }}
    user: {{ .commonName }}
  name: default
`

	if certConfig.CN == "" {
		return nil, errors.New("cn is required")
	}

	key, cert, err := s.signCert(certConfig)
	if err != nil {
		return nil, err
	}

	encodedKey := encodePrivateKey(key)
	encodedCert := encodeCertificate(cert)
	encodedCaCert := encodeCertificate(s.caCert)

	kubeData, err := formatMap(kubeconfig, map[string]string{
		"commonName":  certConfig.CN,
		"clusterName": s.Cluster,
		"apiHost":     s.APIHost,
		"keyData":     base64.StdEncoding.EncodeToString(encodedKey),
		"caData":      base64.StdEncoding.EncodeToString(encodedCaCert),
		"certData":    base64.StdEncoding.EncodeToString(encodedCert),
	})
	if err != nil {
		return nil, err
	}
	return kubeData, nil
}
