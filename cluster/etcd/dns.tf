data "aws_route53_zone" "main" {
  name = "${var.base_domain}"
}

resource "aws_route53_record" "etcd_srv_discover" {
  count   = 1
  name    = "_etcd-server._tcp"
  type    = "SRV"
  zone_id = "${data.aws_route53_zone.main.zone_id}"
  records = ["${formatlist("0 0 2380 %s", aws_route53_record.etc_a_nodes.*.fqdn)}"]
  ttl     = "300"
}

resource "aws_route53_record" "etcd_srv_client" {
  count   = 1
  name    = "_etcd-client._tcp"
  type    = "SRV"
  zone_id = "${data.aws_route53_zone.main.zone_id}"
  records = ["${formatlist("0 0 2379 %s", aws_route53_record.etc_a_nodes.*.fqdn)}"]
  ttl     = "60"
}

resource "aws_route53_record" "etc_a_nodes" {
  count   = "${var.instances}"
  type    = "A"
  ttl     = "60"
  zone_id = "${data.aws_route53_zone.main.zone_id}"
  name    = "etcd${count.index}.${var.cluster_name}"
  records = ["${aws_instance.etcd_node.*.private_ip[count.index]}"]
}
