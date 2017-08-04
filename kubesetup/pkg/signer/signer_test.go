package signer

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSigner_SignWithoutInit(t *testing.T) {
	s := &Signer{
		CaURI:  "file://test-ca.pem",
		KeyURI: "file://test-ca-key.pem",
		OutDir: "./",
	}

	err := s.Sign(CertConfig{
		CN: "test",
	})
	assert.NotNil(t, err)
}

func TestSigner_Sign(t *testing.T) {
	s := &Signer{
		CaURI:   "file://test-ca.pem",
		KeyURI:  "file://test-ca-key.pem",
		OutDir:  "./",
		Cluster: "default",
		APIHost: "test.k8s",
	}

	err := s.Init()
	assert.Nil(t, err)

	err = s.Sign(CertConfig{
		CN: "test",
	})
	assert.Nil(t, err)
}

func TestSigner_SignWithData(t *testing.T) {
	s := &Signer{
		CaURI:   "file://test-ca.pem",
		KeyURI:  "file://test-ca-key.pem",
		OutDir:  "./",
		Cluster: "default",
		APIHost: "test.k8s",
	}

	err := s.Init()
	assert.Nil(t, err)

	conf, err := s.SignWithData(CertConfig{
		CN: "test",
	})
	assert.Nil(t, err)
	fmt.Printf("%s", conf)
}
