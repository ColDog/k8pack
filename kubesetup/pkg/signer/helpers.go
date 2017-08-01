package signer

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"text/template"

	"github.com/coldog/k8pack/kubesetup/pkg/getter"
)

func getKey(uri string) (*rsa.PrivateKey, error) {
	data, err := getter.Get(uri)
	if err != nil {
		return nil, err
	}
	decoded, _ := pem.Decode(data)
	if decoded == nil {
		return nil, errors.New("no PEM data found")
	}
	return x509.ParsePKCS1PrivateKey(decoded.Bytes)
}

func getCa(uri string) (*x509.Certificate, error) {
	data, err := getter.Get(uri)
	if err != nil {
		return nil, err
	}
	decoded, _ := pem.Decode(data)
	if decoded == nil {
		return nil, errors.New("no PEM data found")
	}
	return x509.ParseCertificate(decoded.Bytes)
}

func formatMap(format string, m map[string]string) ([]byte, error) {
	tpl, err := template.New("format").Parse(format)
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	err = tpl.Execute(&buf, m)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil

}

func encodePrivateKey(key *rsa.PrivateKey) []byte {
	block := pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}
	return pem.EncodeToMemory(&block)
}

func encodeCertificate(cert *x509.Certificate) []byte {
	block := pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert.Raw,
	}
	return pem.EncodeToMemory(&block)
}
