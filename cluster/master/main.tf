variable "ami" {}

variable "base_domain" {}

variable "cluster_name" {}

variable "api_dns_name" {}

variable "instance_size" {}

variable "ssh_key" {}

variable "user_data" {
  default = ""
}

variable "autoscaling_sgs" {
  type = "list"
}

variable "elb_sgs" {
  type = "list"
}

variable "subnets" {
  type = "list"
}

variable "max" {
  default = 1
}

variable "min" {
  default = 1
}

variable "desired" {
  default = 1
}
