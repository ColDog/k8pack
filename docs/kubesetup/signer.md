# signer
--
    import "."


## Usage

```go
const (
	// RSAKeySize for all generated keys.
	RSAKeySize = 2048

	// Duration365d represents the default expiry for all certs.
	Duration365d = time.Hour * 24 * 365
)
```

#### type CertConfig

```go
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
```

CertConfig is used to sign certificates.

#### type Signer

```go
type Signer struct {
	CaURI   string
	KeyURI  string
	APIHost string
	Cluster string
	OutDir  string
}
```

Signer is the top level signing object, initialize this object to sign certs.

#### func (*Signer) Init

```go
func (s *Signer) Init() error
```
Init does setup on the signer and validates the fields.

#### func (*Signer) Sign

```go
func (s *Signer) Sign(certConfig CertConfig) error
```
Sign will generate a certificate and kubeconfig using the provided CertConfig.

#### func (*Signer) SignWithData

```go
func (s *Signer) SignWithData(certConfig CertConfig) ([]byte, error)
```
SignWithData will generate a certificate and kubeconfig including the data as
base64 encoded blobs.
