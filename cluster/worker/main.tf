variable "name" {
  default = "worker"
}

variable "ami" {}

variable "user_data" {
  default = ""
}

variable "cluster_name" {}

variable "instance_size" {}

variable "ssh_key" {}

variable "autoscaling_sgs" {
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
