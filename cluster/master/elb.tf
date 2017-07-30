data "aws_route53_zone" "main" {
  name = "${var.base_domain}"
}

resource "aws_route53_record" "api_server" {
  zone_id = "${data.aws_route53_zone.main.zone_id}"
  name    = "k8s.${var.cluster_name}"
  type    = "A"

  alias {
    name                   = "${aws_elb.master_api.dns_name}"
    zone_id                = "${aws_elb.master_api.zone_id}"
    evaluate_target_health = false
  }
}

resource "aws_elb" "master_api" {
  name            = "${var.cluster_name}-api-public" // only allows hyphens
  subnets         = ["${var.subnets}"]
  internal        = false
  security_groups = ["${var.elb_sgs}"]

  listener {
    instance_port     = 443
    instance_protocol = "tcp"
    lb_port           = 443
    lb_protocol       = "tcp"
  }

  health_check {
    healthy_threshold   = 2
    unhealthy_threshold = 2
    timeout             = 5
    target              = "HTTP:6199/healthz" // uses kube-master-healthz
    interval            = 15
  }

  tags {
    cluster = "${var.cluster_name}"
  }
}
