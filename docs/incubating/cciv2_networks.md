---
subcategory: "Cloud Container Instance (CCI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cciv2_networks"
description: |-
  Use this data source to get the list of CCI networks within HuaweiCloud.
---

# huaweicloud_cciv2_networks

Use this data source to get the list of CCI networks within HuaweiCloud.

## Example Usage

```hcl
variable "namespace" {}

data "huaweicloud_cciv2_networks" "test" {
  namespace = var.namespace
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `namespace` - (Required, String) Specifies the namespace of the CCI network.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `networks` - The networks.
  The [networks](#networks) structure is documented below.

<a name="networks"></a>
The `networks` block supports:

* `annotations` - The annotations of the CCI network.

* `creation_timestamp` - The creation timestamp of the CCI network.

* `ip_families` - The IP families of the CCI network.

* `labels` - The labels of the CCI network.

* `name` - The name of the CCI network.

* `namespace` - The namespace of the CCI network.

* `resource_version` - The resource version of the CCI network.

* `security_group_ids` - The security group IDs of the CCI network.

* `finalizers` - The finalizers of the CCI network.

* `status` - The status of the CCI network.

* `subnets` - The subnets of the CCI network.
  The [subnets](#networks_subnets) structure is documented below.

* `status` - The status of the CCI network.
  The [status](#status) structure is documented below.

* `uid` - The uid of the CCI network.

<a name="networks_subnets"></a>
The `subnets` block supports:

* `subnet_id` - The subnet ID of the CCI network.

<a name="status"></a>
The `status` block supports:

* `conditions` - The conditions of the CCI network.
  The [conditions](#status_conditions) structure is documented below.

* `status` - The status of the CCI network.

* `subnet_attrs` - The subnet attributes of the CCI network.
  The [subnet_attrs](#status_subnet_attrs) structure is documented below.

<a name="status_conditions"></a>
The `conditions` block supports:

* `last_transition_time` - The last transition time of the CCI network conditions.

* `message` - The message of the CCI network conditions.

* `reason` - The reason of the CCI network conditions.

* `status` - The status of the CCI network conditions.

* `type` - The type of the CCI network conditions.

<a name="status_subnet_attrs"></a>
The `subnet_attrs` block supports:

* `network_id` - The ID of the CCI network.

* `subnet_v4_id` - The subnet IPv4 ID of the CCI network.

* `subnet_v6_id` - The subnet IPv6 ID of the CCI network.
