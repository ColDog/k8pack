
## Inputs

| Name | Description | Default | Required |
|------|-------------|:-----:|:-----:|
| ami | AMI to use, must be built with the image script. | - | yes |
| api_dns_name | API dns name to use, will have base_domain appended. | - | yes |
| api_service_ip | The VIP (Virtual IP) address of the Kubernetes API Service. If the SERVICE_IP_RANGE is changed above, this must be set to the first IP in that range. | `10.3.0.1` | no |
| asset_bucket | Asset bucket name. Must be created already. | - | yes |
| base_domain | Route53 public zone to setup entries. | - | yes |
| cluster_name | Kubernetes cluster name, this will prefix the domain name and all resources. | - | yes |
| dns_service_ip | The VIP (Virtual IP) address of the cluster DNS service. This IP must be in the range of the SERVICE_IP_RANGE and cannot be the first IP in the range. This same IP must be configured on all worker nodes to enable DNS service discovery. This should be left blank if you are not going to install cluster dns. | `` | no |
| etcd_container_image | Etcd container image, used to launch the etcd containers. Kept at latest by default. | `quay.io/coreos/etcd:latest` | no |
| etcd_instance_size | AWS EC2 instance size. | `t2.small` | no |
| etcd_instances | Count of etcd instances to launch, each will be identified by index etc0 - n. | `2` | no |
| master_instance_size | AWS EC2 instance size. | `t2.small` | no |
| master_instances | Master autoscaling group defaults. | `<map>` | no |
| node_cidr | The CIDR network to use for the entire cluster. This must contain the node IPs and the pod IPs. | `10.0.0.0/16` | no |
| pod_cidr | The CIDR network to use for pod IPs. Each pod launched in the cluster will be assigned an IP out of this range. This network must be routable between all hosts in the cluster. In a default installation, the flannel overlay network will provide routing to this network. | `10.2.0.0/16` | no |
| service_cidr | The CIDR network to use for service cluster VIPs (Virtual IPs). Each service will be assigned a cluster IP out of this range. This must not overlap with any IP ranges assigned to the POD_NETWORK, or other existing network infrastructure. Routing to these VIPs is handled by a local kube-proxy service to each host, and are not required to be routable between hosts. | `10.3.0.0/24` | no |
| ssh_key | SSH key name, this should be created in the aws dashboard. | - | yes |
| vpc_id | VPC id, the VPC must be correctly set up but is created outside of the cluster. | - | yes |
| worker_instance_size | AWS EC2 instance size. | `t2.small` | no |
| worker_instances | Worker autoscaling group defaults. | `<map>` | no |

## Outputs

| Name | Description |
|------|-------------|
| etcd_urls |  |
| master_config_uri |  |
| master_lb_sg |  |
| master_sg |  |
| public_sg |  |
| ssh_sg |  |
| subnet_ids |  |
| worker_config_uri |  |
| worker_sg |  |

