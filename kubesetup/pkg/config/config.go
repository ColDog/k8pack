package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/coldog/k8pack/kubesetup/pkg/getter"
	"github.com/coldog/k8pack/kubesetup/pkg/signer"
)

// Asset is a downloadable file. Assets are run through the GO template engine
// with Config.Config as the context.
type Asset struct {
	// Name is the name of the asset on disk.
	Name string

	// URI can either specify an http:// resource or a file:// resource.
	URI string
}

// Systemd unit file.
type Systemd struct {
	Asset

	// Should this service be started.
	Start bool
}

// CNIConfig represents CNI configuration.
type CNIConfig struct {
	Name   string
	Config map[string]interface{}
}

// Config represents all the configuration needed to create a k8s cluster.
type Config struct {
	// BaseURI is prepended to all URI's before processing.
	BaseURI string

	// Signer represents the necessary config to sign certificates.
	Signer *signer.Signer

	// Certs are signed by the signer config.
	Certs []signer.CertConfig

	// Manifests are Kubernetes manifests and are installed in /etc/kubernetes/manifests/
	Manifests []Asset

	// Systemd are systemd services and are installed in /etc/systemd/system/
	Systemd []Systemd

	// Secrets are files placed into the /etc/kubernetes/secrets/ directory.
	Secrets []Asset

	// Logrotate are files placed into the /etc/logrotate.d/ directory.
	Logrotate []Asset

	// MetadataProviders represents metadata providers that will load data
	// into the Config map.
	MetadataProviders []string

	// Configuration used to render templates.
	Config map[string]string

	// Configuration override, merged into Config on Run.
	Override map[string]string

	// CNIConfig holds CNI configuration placed in /etc/cni/net.d/<name>.
	CNIConfig *CNIConfig
}

// Run applies the configuration to the local machine.
func (c *Config) Run() error {
	if c.Signer == nil {
		return errors.New("signer must be present")
	}

	for k, v := range c.Override {
		c.Config[k] = v
	}

	c.Signer.CaURI = c.BaseURI + c.Signer.CaURI
	c.Signer.KeyURI = c.BaseURI + c.Signer.KeyURI
	c.Signer.Cluster = c.Config["ClusterName"]
	c.Signer.APIHost = c.Config["APIHost"]

	log.Printf("initializing signing")
	err := c.Signer.Init()
	if err != nil {
		return err
	}

	for _, providerName := range c.MetadataProviders {
		log.Printf("loading metadata from %s", providerName)

		provider, ok := Providers[providerName]
		if !ok {
			return fmt.Errorf("metadata provider %s does not exist", providerName)
		}
		meta, err := provider()
		if err != nil {
			return err
		}
		for k, v := range meta {
			c.Config[k] = v
		}
	}

	data, _ := json.MarshalIndent(c, "", "  ")
	println(string(data))

	for _, cert := range c.Certs {
		log.Printf("writing cert %s", cert.CN)

		cert.CN = format(cert.CN, c.Config)
		cert.Org = formatList(cert.Org, c.Config)
		cert.IPs = formatList(cert.IPs, c.Config)
		cert.DNS = formatList(cert.DNS, c.Config)

		err = c.Signer.Sign(cert)
		if err != nil {
			return err
		}
	}

	for _, manifest := range c.Manifests {
		log.Printf("writing manifest unit %s", manifest.Name)

		err = c.getAsset("/etc/kubernetes/manifests/", ".yml", manifest, 0755)
		if err != nil {
			return err
		}
	}

	for _, systemd := range c.Systemd {
		log.Printf("writing systemd unit %s", systemd.Name)

		err = c.getAsset("/etc/systemd/system/", "", systemd.Asset, 0755)
		if err != nil {
			return err
		}
	}

	for _, secret := range c.Secrets {
		log.Printf("writing secret %s", secret.Name)

		err = c.getAsset("/etc/kubernetes/secrets/", "", secret, 0600)
		if err != nil {
			return err
		}
	}

	for _, logrotate := range c.Logrotate {
		log.Printf("writing logrotate %s", logrotate.Name)

		err = c.getAsset("/etc/logrotate.d/", "", logrotate, 0700)
		if err != nil {
			return err
		}
	}

	if c.CNIConfig != nil {
		log.Printf("writing cni config %s", c.CNIConfig.Name)

		data, _ := json.MarshalIndent(c.CNIConfig.Config, "", "  ")
		err = ioutil.WriteFile("/etc/cni/net.d/"+c.CNIConfig.Name, data, 0600)
		if err != nil {
			return err
		}
	}

	log.Println("systemctl reload")
	err = exec.Command("systemctl", "daemon-reload").Run()
	if err != nil {
		return err
	}

	for _, systemd := range c.Systemd {
		if !systemd.Start {
			continue
		}
		log.Printf("systemctl start %s", systemd.Name)
		err := exec.Command("systemctl", "start", "--no-block", systemd.Name).Run()
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Config) getAsset(base, suffix string, asset Asset, perm os.FileMode) error {
	data, err := getter.Get(c.BaseURI + asset.URI)
	if err != nil {
		return err
	}
	data, err = render(data, c.Config)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(base+asset.Name+suffix, data, perm)
	if err != nil {
		return err
	}
	return nil
}

func render(data []byte, m map[string]string) ([]byte, error) {
	tpl, err := template.New("format").Parse(string(data))
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

func format(format string, m map[string]string) string {
	tpl, err := template.New("format").Parse(format)
	if err != nil {
		return format
	}
	var buf bytes.Buffer
	err = tpl.Execute(&buf, m)
	if err != nil {
		return format
	}
	return buf.String()
}

func formatList(list []string, m map[string]string) (out []string) {
	for _, item := range list {
		out = append(out, format(item, m))
	}
	return out
}
