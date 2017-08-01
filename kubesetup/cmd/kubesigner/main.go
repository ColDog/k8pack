package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/coldog/k8pack/kubesetup/pkg/signer"
)

var (
	caURI   string
	keyURI  string
	apiHost string
	cluster string
	cn      string
	org     string
	dir     string
)

func main() {
	flag.StringVar(&caURI, "ca-uri", "file://ca.pem", "Certificate authority file")
	flag.StringVar(&keyURI, "key-uri", "file://ca-key.pem", "Private key file")
	flag.StringVar(&apiHost, "api", "", "Kubernetes api host")
	flag.StringVar(&cluster, "cluster", "", "Kubernetes cluster name")
	flag.StringVar(&cn, "name", "", "Certificate common name (username)")
	flag.StringVar(&org, "org", "", "Certificate organizations (groups)")
	flag.StringVar(&dir, "dir", "", "Out directory")
	flag.Parse()

	sign := &signer.Signer{
		CaURI:   caURI,
		KeyURI:  keyURI,
		APIHost: apiHost,
		Cluster: cluster,
		OutDir:  dir,
	}

	err := sign.Init()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	err = sign.Sign(signer.CertConfig{
		CN:  cn,
		Org: strings.Split(org, ","),
	})
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
