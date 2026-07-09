---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_rts_overview"
description: |-
  Use this data source to query the real-time session overview of a GaussDB instance within HuaweiCloud.
---

# huaweicloud_gaussdb_rts_overview

Use this data source to query the real-time session overview of a GaussDB instance within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_gaussdb_rts_overview" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the session overview.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the GaussDB instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `active_num` - The number of active real-time sessions.

* `total_num` - The total number of real-time sessions.

* `slow_sql_num` - The number of slow SQL statements in real-time sessions.

* `lock_num` - The number of lock-waiting sessions.
