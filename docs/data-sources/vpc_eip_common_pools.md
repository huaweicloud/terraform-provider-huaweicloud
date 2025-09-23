---
subcategory: "Elastic IP (EIP)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpc_eip_common_pools"
description: |-
  Use this data source to get a list of EIP common pools.
---

# huaweicloud_vpc_eip_common_pools

Use this data source to get a list of EIP common pools.

## Example Usage

```hcl
data "huaweicloud_vpc_eip_common_pools" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the common pool name.

* `public_border_group` - (Optional, String) Specifies whether the common pool is at the center or at the edge.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `common_pools` - Indicates the common pools.

  The [common_pools](#common_pools_struct) structure is documented below.

<a name="common_pools_struct"></a>
The `common_pools` block supports:

* `id` - Indicates the common pool ID.

* `name` - Indicates the common pool name。

* `type` - Indicates the common pool type, such as **bgp** and **sbgp**.

* `status` - Indicates the common pool status.

* `allow_share_bandwidth_types` - Indicates the list of shared bandwidth types that the public IP address can be added to.

* `public_border_group` - Indicates whether the common pool is at the center or at the edge.

* `used` - Indicates the number of used IP addresses.

* `available` - Indicates the number of available IP addresses.
