export CGOENABLED=0
export GOOS=linux

build-image:
	@go build -o kubesetup/kubesetup ./kubesetup/cmd/kubesetup
	@packer build image/node.json

.PHONY: build-image
