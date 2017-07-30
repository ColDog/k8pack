variable "base_domain" {}

variable "cluster_name" {}

variable "instance_type" {}

variable "ssh_key" {}

variable "subnets" {
  type = "list"
}

variable "security_groups" {
  type = "list"
}

variable "container_image" {}

variable "instances" {}

output "etcd_urls" {
  value = ["${formatlist("http://%s:2379", aws_route53_record.etc_a_nodes.*.fqdn)}"]
}
