resource "aws_launch_configuration" "worker" {
  name                 = "${var.cluster_name}_${var.name}_lc_${uuid()}"
  image_id             = "${var.ami}"
  instance_type        = "${var.instance_size}"
  key_name             = "${var.ssh_key}"
  iam_instance_profile = "${aws_iam_instance_profile.worker_profile.id}"
  security_groups      = ["${var.autoscaling_sgs}"]

  associate_public_ip_address = true
  user_data                   = "${var.user_data}"

  lifecycle {
    ignore_changes        = ["name"]
    create_before_destroy = true
  }
}

resource "aws_autoscaling_group" "worker" {
  name                 = "${var.cluster_name}_${var.name}_asg"
  max_size             = "${var.max}"
  min_size             = "${var.min}"
  desired_capacity     = "${var.desired}"
  launch_configuration = "${aws_launch_configuration.worker.id}"
  force_delete         = true
  vpc_zone_identifier  = ["${var.subnets}"]
  termination_policies = ["OldestLaunchConfiguration"]

  tag {
    key                 = "cluster"
    value               = "${var.cluster_name}"
    propagate_at_launch = true
  }

  tag {
    key                 = "Name"
    value               = "${var.cluster_name}_${var.name}"
    propagate_at_launch = true
  }
}
