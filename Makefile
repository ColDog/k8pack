export CGOENABLED=0
export GOOS=linux

build-image:
	@go build -o kubesetup/kubesetup ./kubesetup/cmd/kubesetup
	@packer build image/node.json

module-docs:
	@terraform-docs md ./cluster > docs/modules/cluster.md
	@terraform-docs md ./cluster/config > docs/modules/cluster-config.md
	@terraform-docs md ./cluster/etcd > docs/modules/cluster-etcd.md
	@terraform-docs md ./cluster/master > docs/modules/cluster-master.md
	@terraform-docs md ./cluster/vpcdata > docs/modules/cluster-vpcdata.md
	@terraform-docs md ./cluster/worker > docs/modules/cluster-worker.md

kubesetup-docs:
	@godocdown kubesetup/pkg/signer > docs/kubesetup/signer.md
	@godocdown kubesetup/pkg/config > docs/kubesetup/config.md
	@godocdown kubesetup/pkg/getter > docs/kubesetup/getter.md

docs: kubesetup-docs module-docs

.PHONY: build-image kubesetup-docs module-docs
