# Security

## Certificates

* Currently the ca-key file must be downloaded by every node on startup. This means that full access to the cluster is possible as long as you are inside the cluster you can download the ca-key and sign a certificate.
    - Vault can potentially be used as an external source for signing certificates, this requires having vault running.
    - Storing secret key and ca files in a separate bucket and using temporary access token from AWS to access them.
    - Using signed AWS url's with a short expiry.

* Certificates are also not rotated. The best way to do this would probably be just bringing up a new instance.

TODO
