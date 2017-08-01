package config

import (
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// MetadataProvider represents the interface for a metadata fetching function.
type MetadataProvider func() (map[string]string, error)

// Providers is a map of name to metadata provider.
var Providers = map[string]MetadataProvider{
	"EC2": EC2Metadata,
}

// EC2Metadata represents a MetadataProvider for EC2 metadata.
func EC2Metadata() (meta map[string]string, err error) {
	meta = map[string]string{}

	meta["LocalHostname"], err = getEc2Metadata("local-hostname")
	if err != nil {
		return nil, err
	}
	meta["PublicHostname"], err = getEc2Metadata("public-hostname")
	if err != nil {
		return nil, err
	}
	meta["LocalIPv4"], err = getEc2Metadata("local-ipv4")
	if err != nil {
		return nil, err
	}
	meta["PublicIPv4"], err = getEc2Metadata("public-ipv4")
	if err != nil {
		return nil, err
	}
	meta["InstanceID"], err = getEc2Metadata("instance-id")
	if err != nil {
		return nil, err
	}
	return
}

func getEc2Metadata(path string) (string, error) {
	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	url := "http://169.254.169.254/latest/meta-data/" + path
	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}
