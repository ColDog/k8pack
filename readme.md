# K8Pack

This is a toolkit for bringing up Kubernetes clusters on AWS. It's not a one step command for bringing up clusters, instead it focuses on providing the core tools and terraform modules needed.

## Quickstart

Read through the [tutorial](docs/tutorial/0-prerequisites.md)

## Goals

1. Toolkit for Kubernetes clusters.
2. Security as a priority
3. Hardened setup that avoids beta features.
4. Fast boot times.

## Design

The core parts of this are terraform modules that configure the underlying resources. ETCD is completely configured usign terraform, while the Kubernetes cluster is composed of a worker autoscaling group and a master autoscaling group.

The individual nodes are first provisioned using Packer. This provides a faster boot experience as the nodes don't need to download big containers or binaries. On startup, the `kubesetup` binary is run. This downloads various `systemd` units, prepares certificates and fetches metadata. The `kubesetup` has a flexible configuration that uses go templates to inject variables at startup. Once the `kubelet` systemd unit is started, the node effectively has booted and `kubesetup` shuts down.

Further docs: [tutorial](docs/design/)
