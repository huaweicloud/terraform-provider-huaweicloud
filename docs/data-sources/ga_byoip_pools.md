---
subcategory: "Global Accelerator (GA)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ga_byoip_pools"
description: |-
  Use this data source to get the list of GA BYOIP pools within HuaweiCloud.
---

# huaweicloud_ga_byoip_pools

Use this data source to get the list of GA BYOIP pools within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_ga_byoip_pools" "test" {}
```

## Attribute Reference

The following attributes are exported:

* `id` - The data source ID.

* `byoip_pools` - The list of BYOIP pools.

  The [byoip_pools](#byoip_pools_struct) structure is documented below.

<a name="byoip_pools_struct"></a>
The `byoip_pools` block supports:

* `id` - The ID of the BYOIP pool.

* `cidr` - The CIDR block of the BYOIP pool.

* `ip_type` - The IP address version. The value can be **IPV4** or **IPV6**.

* `created_at` - The creation time of the BYOIP pool.

* `updated_at` - The latest update time of the BYOIP pool.

* `area` - The acceleration area. The value can be one of the following:
  + **OUTOFCM**: Outside the Chinese mainland.
  + **CM**: Chinese mainland.

* `domain_id` - The domain ID.
