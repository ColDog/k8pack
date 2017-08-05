module "vpcdata" {
  source = "./vpcdata"

  vpc_id       = "${var.vpc_id}"
  cluster_name = "${var.cluster_name}"
}

module "etcd" {
  source = "./etcd"

  base_domain  = "${var.base_domain}"
  cluster_name = "${var.cluster_name}"

  subnets = ["${module.vpcdata.subnet_ids}"]

  security_groups = "${compact([
    "${module.vpcdata.worker_sg}",
    "${var.ssh_enabled ? module.vpcdata.ssh_sg : ""}",
  ])}"

  instances     = "${var.etcd_instances}"
  instance_type = "${var.etcd_instance_size}"
  ssh_key       = "${var.ssh_key}"

  container_image = "${var.etcd_container_image}"
}

module "config" {
  source = "./config"

  cluster_name   = "${var.cluster_name}"
  api_host       = "${var.api_dns_name}.${var.base_domain}"
  node_cidr      = "${var.node_cidr}"
  pod_cidr       = "${var.pod_cidr}"
  service_cidr   = "${var.service_cidr}"
  dns_service_ip = "${var.dns_service_ip}"
  api_service_ip = "${var.api_service_ip}"
  asset_bucket   = "${var.asset_bucket}"
  etcd_urls      = "${module.etcd.etcd_urls}"
}

module "master" {
  source = "./master"

  ami           = "${var.ami}"
  base_domain   = "${var.base_domain}"
  cluster_name  = "${var.cluster_name}"
  api_dns_name  = "${var.api_dns_name}"
  instance_size = "${var.master_instance_size}"
  ssh_key       = "${var.ssh_key}"

  max     = "${var.master_instances["max"]}"
  min     = "${var.master_instances["min"]}"
  desired = "${var.master_instances["desired"]}"

  elb_sgs = ["${module.vpcdata.master_lb_sg}"]

  autoscaling_sgs = "${compact([
    "${module.vpcdata.master_sg}",
    "${module.vpcdata.worker_sg}",
    "${var.ssh_enabled ? module.vpcdata.ssh_sg : ""}",
  ])}"

  subnets = ["${module.vpcdata.subnet_ids}"]

  user_data =<<EOF
#!/bin/sh
/opt/bin/kubesetup -config-uri=${module.config.master_config_uri}
EOF
}

module "worker" {
  source = "./worker"

  ami           = "${var.ami}"
  cluster_name  = "${var.cluster_name}"
  instance_size = "${var.worker_instance_size}"
  ssh_key       = "${var.ssh_key}"

  max     = "${var.worker_instances["max"]}"
  min     = "${var.worker_instances["min"]}"
  desired = "${var.worker_instances["desired"]}"

  autoscaling_sgs = "${compact([
    "${module.vpcdata.worker_sg}",
    "${var.ssh_enabled ? module.vpcdata.ssh_sg : ""}",
  ])}"

  subnets = ["${module.vpcdata.subnet_ids}"]

  user_data =<<EOF
#!/bin/sh
/opt/bin/kubesetup -config-uri=${module.config.worker_config_uri}
EOF
}
