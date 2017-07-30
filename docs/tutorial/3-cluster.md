# Creating The Cluster

## Terraform Cluster

This will create the cluster by using the terraform module `cluster`.

Copy the following configuration into a new folder and fill out the required variables. You may need to create a new SSH key using the aws console.

To find the most up to date AMI, follow the docs [here](../amis.md). Or to build your own AMI go [here](../building.md).

```terraform
provider "aws" {
  region = "us-west-2"
}

module "cluster" {
  source = "github.com/coldog/k8pack/cluster"

  ami          = "<ami>"
  base_domain  = "<domain>"
  api_dns_name = "<api_dns_name>"
  cluster_name = "<cluster_name>"
  vpc_id       = "<vpc_id>"
  asset_bucket = "<asset_bucket>"
  ssh_key      = "<ssh_key>"

  master_instances {
    min     = 0
    max     = 1
    desired = 1
  }
  worker_instances {
    min     = 0
    max     = 10
    desired = 10
  }
}
```

Run `terraform apply`. Your cluster should be starting up! The API will be available at the following endpoint once it's ready.

Api Endpoint:

```
<api_dns_name>.<domain>
```

## Accessing the Cluster

Run the following commands to create the necessary local kubeconfig:

```bash
kubectl config set-cluster <cluster_name> \
  --certificate-authority=ca.pem \
  --embed-certs=true \
  --server=https://${api_endpoint}
```

```bash
kubectl config set-credentials admin \
  --client-certificate=admin.pem \
  --client-key=admin-key.pem
```

```bash
kubectl config set-context <cluster_name> \
  --cluster=<cluster_name> \
  --user=admin
```

```bash
kubectl config use-context <cluster_name>
```

Try connecting to the API:

```bash
kubectl get no
```
