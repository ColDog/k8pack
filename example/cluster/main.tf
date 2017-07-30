provider "aws" {
  region = "us-west-2"
}

module "cluster" {
  source = "../../cluster"

  ami          = "ami-2218fd5a"
  base_domain  = "coldog.xyz"
  api_dns_name = "k8s.default"
  cluster_name = "default"
  vpc_id       = "vpc-0c835b6a"
  asset_bucket = "coldog-k8s-cluster"
  ssh_key      = "default_key"

  master_instances {
    min     = 0
    max     = 1
    desired = 1
  }
  worker_instances {
    min     = 0
    max     = 1
    desired = 1
  }
}
