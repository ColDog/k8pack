# Kube Image

This is a packer builder that builds up a CoreOS instance with the necessary binaries and docker containers for Kubernetes.

## Installed

- `flanneld`: Flannel binary is installed.
- `kubelet`: Kubelet binary.
- `kubeproxy`: Kube proxy binary.
- `cni`: Container networking interface.
- `hyperkube`: Hyperkube image (coreos version) is downloaded to the local docker cache.
- `kubesetup`: Kubesetup binary from this repo.

## Building

Run `make build-image` from the root of this repo. Modify the version tag in the `node.json` file.
