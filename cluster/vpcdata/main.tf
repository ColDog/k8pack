variable "vpc_id" {}

variable "cluster_name" {}

data "aws_subnet_ids" "main" {
  vpc_id = "${var.vpc_id}"
}

output "subnet_ids" {
  value = "${data.aws_subnet_ids.main.ids}"
}

output "master_sg" {
  value = "${aws_security_group.master.id}"
}

output "master_lb_sg" {
  value = "${aws_security_group.master_lb.id}"
}

output "worker_sg" {
  value = "${aws_security_group.worker.id}"
}

output "ssh_sg" {
  value = "${aws_security_group.ssh.id}"
}

output "public_sg" {
  value = "${aws_security_group.public.id}"
}
