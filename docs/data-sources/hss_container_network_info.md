---
subcategory: "Host Security Service (HSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_hss_container_network_info"
description: |-
  Use this data source to get the network information of a container cluster within HuaweiCloud.
---

# huaweicloud_hss_container_network_info

Use this data source to get the network information of a container cluster within HuaweiCloud.

## Example Usage

```hcl
variable "cluster_id" {}

data "huaweicloud_hss_container_network_info" "test" {
  cluster_id = var.cluster_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `cluster_id` - (Required, String) Specifies the cluster ID.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.  
  This parameter is valid only when the enterprise project is enabled.
  The default value is **0**, indicating the default enterprise project.
  If you need to query data for all enterprise projects, the value is **all_granted_eps**.
  If you only have permissions for a specific enterprise project, you need to set the enterprise project ID. Otherwise,
  the operation may fail due to insufficient permissions.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `mode` - The network model.

* `vpc` - The VPC ID.

* `subnet` - The subnet ID.

* `security_group` - The security group.

* `ipv4_cidr` - The IPv4 service network segment.

* `cidrs` - The container network segment.

* `kube_proxy_mode` - The service forwarding mode.  
  The valid values are as follows:
  + **iptables**
  + **ipvs**

* `is_support_egress` - Whether egress rule configuration is supported.  
  The valid values are as follows:
  + **true**: Supported.
  + **false**: Not supported.
