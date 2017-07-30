provider "aws" {
  region = "us-west-2"
}

module "vpc" {
  source = "../../vpc"

  asset_bucket = "coldog-k8s-cluster"
  vpc_name     = "k8s"
}
