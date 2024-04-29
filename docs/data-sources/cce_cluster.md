---
subcategory: "Cloud Container Engine (CCE)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cce_cluster"
description: ""
---

# huaweicloud_cce_cluster

Provides details about the cluster and obtains certificate for accessing cluster information.

## Example Usage

```hcl
variable "cluster_name" {}

data "huaweicloud_cce_cluster" "cluster" {
  name   = var.cluster_name
  status = "Available"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to obtain the CCE cluster. If omitted, the provider-level
  region will be used.

* `name` - (Optional, String) Specifies the name of the cluster.

* `id` - (Optional, String) Specifies the ID of the cluster.

* `status` - (Optional, String) Specifies the status of the cluster.

* `cluster_type` - (Optional, String) Specifies the type of the cluster. Possible values: **VirtualMachine**, **ARM64**.

* `vpc_id` - (Optional, String) Specifies the VPC ID to which the cluster belongs.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `billing_mode` - Charging mode of the cluster.

* `description` - Cluster description.

* `name` - The name of the cluster in string format.

* `flavor_id` - The cluster specification in string format.

* `cluster_version` - The version of cluster in string format.

* `container_network_cidr` - The container network segment.

* `container_network_type` - The container network type: **overlay_l2** , **underlay_ipvlan**, **vpc-router** or **eni**.

* `eni_subnet_id` - The **IPv4 subnet ID** of the subnet where the ENI resides.
  Specified when creating a CCE Turbo cluster.

* `eni_subnet_cidr` - ENI network segment. Specified when creating a CCE Turbo cluster.

* `service_network_cidr` - The service network segment.

* `authentication_mode` - Authentication mode of the cluster, possible values are x509 and rbac. Defaults to **rbac**.

* `masters` - Advanced configuration of master nodes. Structure is documented below.

* `security_group_id` - Security group ID of the cluster.

* `subnet_id` - The ID of the subnet used to create the node.

* `highway_subnet_id` - The ID of the high speed network used to create bare metal nodes.

* `enterprise_project_id` - The enterprise project ID of the CCE cluster.

* `endpoints` - The access addresses of kube-apiserver in the cluster. Structure is documented below.

* `certificate_clusters` - The certificate clusters. Structure is documented below.

* `certificate_users` - The certificate users. Structure is documented below.

* `kube_config_raw` - Raw Kubernetes config to be used by kubectl and other compatible tools.

The `masters` block supports:

* `availability_zone` - The availability zone (AZ) of the master node.

The `endpoints` block supports:

* `url` - The URL of the cluster access address.

* `type` - The type of the cluster access address.
  + **Internal**: The user's subnet access address.
  + **External**: The public network access address.

The `certificate_clusters` block supports:

* `name` - The cluster name.

* `server` - The server IP address.

* `certificate_authority_data` - The certificate data.

The `certificate_users` block supports:

* `name` - The user name.

* `client_certificate_data` - The client certificate data.

* `client_key_data` - The client key data.
