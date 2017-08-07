provider "aws" {
  region = "us-west-2"
}

variable "ami" {}
variable "domain" {}
variable "api_dns_name" {}
variable "cluster_name" {}
variable "vpc_id" {}
variable "asset_bucket" {}
variable "ssh_key" {}

module "cluster" {
  source = "../../cluster"

  ami          = "${var.ami}"
  base_domain  = "${var.domain}"
  api_dns_name = "${var.api_dns_name}"
  cluster_name = "${var.cluster_name}"
  vpc_id       = "${var.vpc_id}"
  asset_bucket = "${var.asset_bucket}"
  ssh_key      = "${var.ssh_key}"

  master_instances {
    min     = 0
    max     = 1
    desired = 1
  }
  worker_instances {
    min     = 0
    max     = 2
    desired = 2
  }
}
