# Kubesetup

This is a GO program that initializes systemd units, Kubernetes manifests and Kubeconfig files. It takes a single configuration file that can be stored remotely and accessed over HTTP. When working with terraform the `cluster/config` module handles creating the configuration file and the necessary Systemd unit files. However, the `kubesetup` package can be used on it's own and this configuration file can be used to create a custom kubernetes distribution.

All assets are run through the GO templating engine using the `Config` map as the context. In addition, EC2 or other metadata providers can be used to fill the `Config` map before rendering.

## Config File

The configuration file is specified by the following GO structs.

```go
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

	// CNIConfig holds CNI configuration placed in /etc/cni/net.d/<name>.
	CNIConfig *CNIConfig
}
```

## Running

The program is run with a single argument `-config-uri` which should point to a JSON config file.
