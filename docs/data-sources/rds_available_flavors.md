---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_available_flavors"
description: |-
  Use this data source to get the specifications that a RDS instance can be changed to.
---

# huaweicloud_rds_available_flavors

Use this data source to get the specifications that a RDS instance can be changed to.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_rds_available_flavors" "test" {
  instance_id           = var.instance_id
  availability_zone_ids = "cn-north-4a"
  ha_mode               = "ha"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the RDS instance.

* `availability_zone_ids` - (Required, String) Specifies the availability zone.

* `ha_mode` - (Required, String) Specifies the HA mode. Value options: **single**, **ha**, **replica**.

* `spec_code_like` - (Optional, String) Specifies the resource specification code, fuzzy matching is supported.

* `flavor_category_type` - (Optional, String) Specifies the flavor category type.

* `is_rha_flavor` - (Optional, Bool) Specifies whether display highly available read-only types.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `optional_flavors` - Indicates the list of available flavors.

  The [optional_flavors](#optional_flavors_struct) structure is documented below.

<a name="optional_flavors_struct"></a>
The `optional_flavors` block supports:

* `vcpus` - Indicates the CPU size.

* `ram` - Indicates the memory size, in GB.

* `spec_code` - Indicates the resource specification code.

* `is_ipv6_supported` - Indicates whether supported ipv6.

* `type_code` - Indicates the resource type.

* `az_status` - Indicates the az status.

* `group_type` - Indicates the performance specifications. Its value can be any of the following:
  + **normal**: general-enhanced
  + **normal2**: general-enhanced II
  + **armFlavors**: Kunpeng general-enhanced
  + **dedicicatenormal**: exclusive x86
  + **armlocalssd**: standard Kunpeng
  + **normallocalssd**: standard x86
  + **general**: general-purpose
  + **dedicated**: dedicated, which is only supported for cloud SSDs
  + **rapid**: dedicated, which is only supported for extreme SSDs
  + **bigmen**: Large-memory

* `max_connection` - Indicates the max connection.

* `tps` - Indicates the number of transactions executed by the database per second, each containing 18 SQL statements.

* `qps` - Indicates the number of SQL statements executed by the database per second, including **insert**, **select**,
  **update**, **delete** and so on.

* `min_volume_size` - Indicates the minimum disk capacity in GB.

* `max_volume_size` - Indicates the maximum disk capacity in GB.
