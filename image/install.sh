#!/usr/bin/env bash

set -e

if [ -z ${CNI_VERSION} ]; then
    echo "CNI_VERSION is required"
    exit 1
fi

if [ -z ${K8S_VERSION} ]; then
    echo "K8S_VERSION is required"
    exit 1
fi

if [ -z ${FLANNELD_VERSION} ]; then
    echo "FLANNELD_VERSION is required"
    exit 1
fi

mkdir -p /etc/kubernetes/secrets
mkdir -p /etc/kubernetes/manifests
mkdir -p /opt/downloads/ /opt/cni/bin/ /opt/bin
mkdir -p /etc/cni/net.d

# Hyperkube config.
hyperkube="quay.io/coreos/hyperkube:${K8S_VERSION}_coreos.0"
docker pull $hyperkube

# Create versions file.
cat > /etc/kubernetes/versions.env <<EOF
HYPERKUBE_IMAGE=$hyperkube
FLANNELD_VERSION=$FLANNELD_VERSION
K8S_VERSION=$K8S_VERSION
CNI_VERSION=$K8S_VERSION
EOF

# Downloads.
for name in "kubelet" "kube-proxy"; do
    curl -L -o ${name} https://storage.googleapis.com/kubernetes-release/release/${K8S_VERSION}/bin/linux/amd64/${name}
    chmod +x ${name}
    mv ${name} /opt/bin/
done

# Flannel binary.
curl -L -o flanneld https://github.com/coreos/flannel/releases/download/${FLANNELD_VERSION}/flanneld-amd64
chmod +x flanneld
mv flanneld /opt/bin/

# CNI binaries.
curl -L -o cni.tgz https://github.com/containernetworking/cni/releases/download/${CNI_VERSION}/cni-amd64-${CNI_VERSION}.tgz
tar -xvf cni.tgz -C /opt/cni/bin/
rm cni.tgz

# Kubesetup binary.
mv /home/core/kubesetup /opt/bin/kubesetup
chmod +x /opt/bin/kubesetup

# This is done because of https://github.com/kubernetes/kubernetes/issues/32728.
sudo mount -o remount,rw '/sys/fs/cgroup'
sudo ln -s /sys/fs/cgroup/cpu,cpuacct /sys/fs/cgroup/cpuacct,cpu