---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_wal_log_recovery_time_window"
description: |-
  Use this data source to query the WAL log recovery time window of a specified RDS read-replica instance.
---

# huaweicloud_rds_wal_log_recovery_time_window

Use this data source to query the WAL log recovery time window of a specified RDS read-replica instance.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_rds_wal_log_recovery_time_window" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the RDS instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `recovery_min_time` - Indicates the start of the WAL log recovery window (exclusive), formatted as a timestamp string.

* `recovery_max_time` - Indicates the end of the WAL log recovery window (inclusive), formatted as a timestamp string.
