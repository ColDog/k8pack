variable "cluster_name" {}

variable "api_host" {}

variable "node_cidr" {}

variable "pod_cidr" {}

variable "service_cidr" {}

variable "dns_service_ip" {}

variable "api_service_ip" {}

variable "asset_bucket" {}

variable "etcd_urls" {
  type = "list"
}

variable "network_plugin" {
  default = "flannel"
}

data "aws_region" "current" {
  current = true
}

output "master_config_uri" {
  value = "https://s3-${data.aws_region.current.name}.amazonaws.com/${aws_s3_bucket_object.master_config.bucket}/${aws_s3_bucket_object.master_config.key}"
}

output "worker_config_uri" {
  value = "https://s3-${data.aws_region.current.name}.amazonaws.com/${aws_s3_bucket_object.master_config.bucket}/${aws_s3_bucket_object.worker_config.key}"
}
