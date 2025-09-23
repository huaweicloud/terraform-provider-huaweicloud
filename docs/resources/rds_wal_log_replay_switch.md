---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_wal_log_replay_switch"
description: |-
  Manages pausing or resuming WAL replay on a PostgreSQL read replica resource within HuaweiCloud.
---

# huaweicloud_rds_wal_log_replay_switch

Manages pausing or resuming WAL replay on a PostgreSQL read replica resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_rds_wal_log_replay_switch" "test" {
  instance_id       = var.instance_id
  pause_log_replay  = "true"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource. If omitted, the
  provider-level region will be used. Changing this creates a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the RDS PostgreSQL instance.

* `pause_log_replay` - (Required, String, NonUpdatable) Specifies whether to pause or resume WAL replay.
  Valid values are:
  + **true**: Specifies that WAL replay will be paused.
  + **false**: Specifies that WAL replay will be resumed.

## Attribute Reference

In addition to all arguments above, the following attribute is exported:

* `id` - The resource ID. The value is the instance ID.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 60 minutes.
