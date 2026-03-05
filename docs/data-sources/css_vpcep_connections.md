---
subcategory: "Cloud Search Service (CSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_css_vpcep_connections"
description: |-
  Use this data source to get the list of endpoint connections.
---

# huaweicloud_css_vpcep_connections

Use this data source to get the list of endpoint connections.

## Example Usage

```hcl
variable "cluster_id" {}

data "huaweicloud_css_vpcep_connections" "test" {
  cluster_id = var.cluster_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `cluster_id` - (Required, String) Specifies the cluster ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `connections` - The cluster VPC endpoint connection information.
  The [connections](#connections_struct) structure is documented below.

* `permissions` - The permissions of the cluster VPC endpoint whitelist.
  The [permissions](#permissions_struct) structure is documented below.

* `vpc_service_name` - The name of the VPC endpoint service.

* `vpcep_update_switch` - Whether to update VPC endpoint.

<a name="connections_struct"></a>
The `connections` block supports:

* `id` - The cluster VPC endpoint ID.

* `status` - The cluster VPC endpoint status.

* `max_session` - The maximum number of connections to a VPC endpoint.

* `specification_name` - The cluster VPC endpoint name.

* `created_at` - The creation time.

* `update_at` - The update time.

* `domain_id` - The account ID of the owner.

* `vpcep_ip` - The IPv4 address of a cluster VPC endpoint.

* `vpcep_ipv6_address` - The IPv6 address of a cluster VPC endpoint.

* `vpcep_dns_name` - The private domain name of a cluster VPC endpoint.

<a name="permissions_struct"></a>
The `permissions` block supports:

* `id` - The ID of the permission.

* `permission` - The permission details for the VPCEP connection whitelist.

* `permission_type` - The VPC endpoint permission type.

* `created_at` - The creation time.
