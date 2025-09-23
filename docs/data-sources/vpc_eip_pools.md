---
subcategory: "Elastic IP (EIP)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpc_eip_pools"
description: |-
  Use this data source to get a list of EIP pools.
---

# huaweicloud_vpc_eip_pools

Use this data source to get a list of EIP pools.

## Example Usage

```hcl
data "huaweicloud_vpc_eip_pools" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the pool name.

* `size` - (Optional, Int) Specifies the pool size.

* `status` - (Optional, String) Specifies the pool status.

* `type` - (Optional, String) Specifies the pool type.

* `description` - (Optional, String) Specifies the pool description.

* `public_border_group` - (Optional, String) Specifies whether the pool is at the center or at the edge.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `pools` - Indicates the public network pools.

  The [pools](#pools_struct) structure is documented below.

<a name="pools_struct"></a>
The `pools` block supports:

* `id` - Indicates the pool ID.

* `name` - Indicates the pool name.

* `type` - Indicates the pool type.

* `size` - Indicates the pool size.

* `status` - Indicates the pool status.

* `used` - indicates the number of used IP addresses.

* `shared` - Indicates whether to share the pool.

* `description` - Indicates the description.

* `public_border_group` - Indicates the pool at the central site or edge site.

* `allow_share_bandwidth_types` - Indicates the list of shared bandwidth types to which the public IP address can be added.

* `billing_info` - Indicates the order information. If an order is available, it indicates a yearly/monthly pool.

  The [billing_info](#pools_billing_info_struct) structure is documented below.

* `tags` - The key/value pairs which associated with the EIP pool.

* `enterprise_project_id` - Indicates the ID of an enterprise project.

* `created_at` - Indicates the create time of the pool.

* `updated_at` - Indicates the update time of the pool.

<a name="pools_billing_info_struct"></a>
The `billing_info` block supports:

* `order_id` - Indicates the order ID.

* `product_id` - Indicates the product ID
