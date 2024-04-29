---
subcategory: "Cloud Container Engine (CCE)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cce_clusters"
description: ""
---

# huaweicloud_cce_clusters

Use this data source to get a list of CCE clusters.

## Example Usage

```hcl
variable "cluster_name" {}

data "huaweicloud_cce_clusters" "clusters" {
  name   = var.cluster_name
  status = "Available"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to obtain the CCE clusters. If omitted, the
  provider-level region will be used.

* `name` - (Optional, String) Specifies the name of the cluster.

* `cluster_id` - (Optional, String) Specifies the ID of the cluster.

* `cluster_type` - (Optional, String) Specifies the type of the cluster. Possible values: **VirtualMachine**, **BareMetal**.

* `vpc_id` - (Optional, String) Specifies the VPC ID to which the cluster belongs.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID of the cluster.

* `status` - (Optional, String) Specifies the status of the cluster.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates a data source ID.

* `ids` - Indicates a list of IDs of all CCE clusters found.

* `clusters` - Indicates a list of CCE clusters found. Structure is documented below.

The `clusters` block supports:

* `name` - The name of the cluster.

* `id` - The ID of the cluster.

* `cluster_type` - The type of the cluster. Possible values: **VirtualMachine**, **ARM64**.

* `status` - The status of the cluster.

* `flavor_id` - The specification of the cluster.

* `cluster_version` - The version of the cluster.

* `description` - The description of the cluster.

* `billing_mode` - The charging mode of the cluster.

* `container_network_cidr` - The container network segment.

* `container_network_type` - The container network type: **overlay_l2** , **underlay_ipvlan**, **vpc-router** or **eni**.

* `eni_subnet_id` - The **IPv4 subnet ID** of the subnet where the ENI resides.

* `eni_subnet_cidr` - The ENI network segment.

* `service_network_cidr` - The service network segment.

* `authentication_mode` - The authentication mode of the cluster, possible values are x509 and rbac. Defaults to **rbac**.

* `masters` - The advanced configuration of master nodes. Structure is documented below.

* `security_group_id` - The security group ID of the cluster.

* `vpc_id` - The vpc ID of the cluster.

* `subnet_id` - The ID of the subnet used to create the node.

* `highway_subnet_id` - The ID of the high speed network used to create bare metal nodes.

* `enterprise_project_id` - The enterprise project ID of the CCE cluster.

* `endpoints` - The access addresses of kube-apiserver in the cluster. Structure is documented below.

* `certificate_clusters` - The certificate clusters. Structure is documented below.

* `certificate_users` - The certificate users. Structure is documented below.

* `kube_config_raw` - The raw Kubernetes config to be used by kubectl and other compatible tools.

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
