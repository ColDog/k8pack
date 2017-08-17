resource "aws_s3_bucket_object" "systemd_apiserver" {
  bucket  = "${var.asset_bucket}"
  key     = "${var.cluster_name}/systemd/apiserver.service"
  content = "${file("${path.module}/systemd/apiserver.service")}"
}

resource "aws_s3_bucket_object" "systemd_controllermanager" {
  bucket  = "${var.asset_bucket}"
  key     = "${var.cluster_name}/systemd/controllermanager.service"
  content = "${file("${path.module}/systemd/controllermanager.service")}"
}

resource "aws_s3_bucket_object" "systemd_flannel" {
  bucket  = "${var.asset_bucket}"
  key     = "${var.cluster_name}/systemd/flannel.service"
  content = "${file("${path.module}/systemd/flannel.service")}"
}

resource "aws_s3_bucket_object" "systemd_calico" {
  bucket  = "${var.asset_bucket}"
  key     = "${var.cluster_name}/systemd/calico.service"
  content = "${file("${path.module}/systemd/calico.service")}"
}

resource "aws_s3_bucket_object" "systemd_healthzmaster" {
  bucket  = "${var.asset_bucket}"
  key     = "${var.cluster_name}/systemd/healthzmaster.service"
  content = "${file("${path.module}/systemd/healthzmaster.service")}"
}

resource "aws_s3_bucket_object" "systemd_kubelet" {
  bucket  = "${var.asset_bucket}"
  key     = "${var.cluster_name}/systemd/kubelet.service"
  content = "${file("${path.module}/systemd/kubelet.service")}"
}

resource "aws_s3_bucket_object" "systemd_kubeproxy" {
  bucket  = "${var.asset_bucket}"
  key     = "${var.cluster_name}/systemd/kubeproxy.service"
  content = "${file("${path.module}/systemd/kubeproxy.service")}"
}

resource "aws_s3_bucket_object" "systemd_scheduler" {
  bucket  = "${var.asset_bucket}"
  key     = "${var.cluster_name}/systemd/scheduler.service"
  content = "${file("${path.module}/systemd/scheduler.service")}"
}

resource "aws_s3_bucket_object" "systemd_logrotate" {
  bucket  = "${var.asset_bucket}"
  key     = "${var.cluster_name}/systemd/logrotate.service"
  content = "${file("${path.module}/systemd/logrotate.service")}"
}

resource "aws_s3_bucket_object" "systemd_logrotate_timer" {
  bucket  = "${var.asset_bucket}"
  key     = "${var.cluster_name}/systemd/logrotate.timer"
  content = "${file("${path.module}/systemd/logrotate.timer")}"
}

resource "aws_s3_bucket_object" "systemd_logger" {
  bucket  = "${var.asset_bucket}"
  key     = "${var.cluster_name}/systemd/logger.service"
  content = "${file("${path.module}/systemd/logger.service")}"
}

data "template_file" "cni_config" {
  template = "${file("${path.module}/cni/${var.network_plugin}.json")}"
  vars {
    etcd_urls = "${join(",", var.etcd_urls)}"
  }
}

data "template_file" "worker_config" {
  template = "${file("${path.module}/worker_config.json")}"
  vars {
    base_uri       = "https://s3-${data.aws_region.current.name}.amazonaws.com/${var.asset_bucket}/${var.cluster_name}/"
    cluster_name   = "${var.cluster_name}"
    api_host       = "${var.api_host}"
    node_cidr      = "${var.node_cidr}"
    pod_cidr       = "${var.pod_cidr}"
    service_cidr   = "${var.service_cidr}"
    dns_service_ip = "${var.dns_service_ip}"
    api_service_ip = "${var.api_service_ip}"
    asset_bucket   = "${var.asset_bucket}"
    etcd_urls      = "${join(",", var.etcd_urls)}"
    aws_region     = "${data.aws_region.current.name}"
    cni_config     = "${data.template_file.cni_config.rendered}",
    network_plugin = "${var.network_plugin}",
  }
}

resource "aws_s3_bucket_object" "worker_config" {
  bucket  = "${var.asset_bucket}"
  key     = "${var.cluster_name}/worker_config_${md5(data.template_file.worker_config.rendered)}.json"
  content = "${data.template_file.worker_config.rendered}"
}

data "template_file" "master_config" {
  template = "${file("${path.module}/master_config.json")}"
  vars {
    base_uri       = "https://s3-${data.aws_region.current.name}.amazonaws.com/${var.asset_bucket}/${var.cluster_name}/"
    cluster_name   = "${var.cluster_name}"
    api_host       = "${var.api_host}"
    node_cidr      = "${var.node_cidr}"
    pod_cidr       = "${var.pod_cidr}"
    service_cidr   = "${var.service_cidr}"
    dns_service_ip = "${var.dns_service_ip}"
    api_service_ip = "${var.api_service_ip}"
    asset_bucket   = "${var.asset_bucket}"
    etcd_urls      = "${join(",", var.etcd_urls)}"
    aws_region     = "${data.aws_region.current.name}"
    cni_config     = "${data.template_file.cni_config.rendered}",
    network_plugin = "${var.network_plugin}",
  }
}

resource "aws_s3_bucket_object" "master_config" {
  bucket  = "${var.asset_bucket}"
  key     = "${var.cluster_name}/master_config_${md5(data.template_file.master_config.rendered)}.json"
  content = "${data.template_file.master_config.rendered}"
}
