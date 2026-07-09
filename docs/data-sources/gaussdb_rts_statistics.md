---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_rts_statistics"
description: |-
  Use this data source to query the real time session statistics of a GaussDB instance within HuaweiCloud.
---

# huaweicloud_gaussdb_rts_statistics

Use this data source to query the real time session statistics of a GaussDB instance within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
variable "instance_id" {}

data "huaweicloud_gaussdb_rts_statistics" "test" {
  instance_id = var.instance_id
  dimension   = "usename"
}
```

### Query with sort

```hcl
variable "instance_id" {}

data "huaweicloud_gaussdb_rts_statistics" "test" {
  instance_id = var.instance_id
  dimension   = "usename"
  order_field = "active_num"
  order       = "DESC"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the instance ID.

* `dimension` - (Required, String) Specifies the dimension of the session statistics.
  The valid values are as follows:
  + **usename**
  + **client_addr**
  + **datname**

* `order_field` - (Optional, String) Specifies the field to sort the results.
  The valid values are as follows:
  + **active_num**
  + **total_num**

* `order` - (Optional, String) Specifies the sort order.
  The valid values are as follows:
  + **ASC**
  + **DESC**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `statistics_list` - The list of session statistics.
  The [statistics_list](#rts_statistics_list) structure is documented below.

<a name="rts_statistics_list"></a>
The `statistics_list` block supports:

* `name` - The dimension name.

* `active_num` - The number of active sessions.

* `total_num` - The total number of sessions.
