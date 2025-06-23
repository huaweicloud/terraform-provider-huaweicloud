---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_wal_log_replay_delay_status"
description: |-
  Use this data source to query the WAL log replay delay status of a specified RDS instance.
---

# huaweicloud_rds_wal_log_replay_delay_status

Use this data source to query the WAL log replay delay status of a specified RDS instance.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_rds_wal_log_replay_delay_status" "test" {
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

* `cur_delay_time_mills` - Indicates the current WAL log replay delay time in milliseconds.

* `delay_time_value_range` - Indicates the valid range for the WAL log delay time.

* `real_delay_time_mills` - Indicates the actual replay delay time of the WAL log in milliseconds.

* `cur_log_replay_paused` - Indicates whether WAL log replay is currently paused. Value can be as follow:
  + **true**: Indicates that the WAL log replay is paused.
  + **false**: Indicates that the WAL log replay is actively running.

* `latest_receive_log` - Indicates the latest WAL log that has been received by the RDS instance.

* `latest_replay_log` - Indicates the latest WAL log that has been replayed by the RDS instance.
