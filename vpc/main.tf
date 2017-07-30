variable "vpc_name" {}

variable "asset_bucket" {}

variable "cidr" {
  default = "10.0.0.0/16"
}

output "subnet_ids" {
  value = ["${aws_subnet.main_subnets.*.id}"]
}

output "id" {
  value = "${aws_vpc.main_vpc.id}"
}
