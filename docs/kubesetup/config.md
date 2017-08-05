# config
--
    import "."


## Usage

```go
var Providers = map[string]MetadataProvider{
	"EC2": EC2Metadata,
}
```
Providers is a map of name to metadata provider.

#### func  EC2Metadata

```go
func EC2Metadata() (meta map[string]string, err error)
```
EC2Metadata represents a MetadataProvider for EC2 metadata.

#### type Asset

```go
type Asset struct {
	// Name is the name of the asset on disk.
	Name string

	// URI can either specify an http:// resource or a file:// resource.
	URI string
}
```

Asset is a downloadable file. Assets are run through the GO template engine with
Config.Config as the context.

#### type CNIConfig

```go
type CNIConfig struct {
	Name   string
	Config map[string]interface{}
}
```

CNIConfig represents CNI configuration.

#### type Config

```go
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
```

Config represents all the configuration needed to create a k8s cluster.

#### func (*Config) Run

```go
func (c *Config) Run() error
```
Run applies the configuration to the local machine.

#### type MetadataProvider

```go
type MetadataProvider func() (map[string]string, error)
```

MetadataProvider represents the interface for a metadata fetching function.

#### type Systemd

```go
type Systemd struct {
	Asset

	// Should this service be started.
	Start bool
}
```

Systemd unit file.
