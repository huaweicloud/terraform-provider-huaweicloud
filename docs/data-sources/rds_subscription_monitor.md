---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_subscription_monitor"
description: |-
  Use this data source to query the monitor information of an RDS subscription within HuaweiCloud.
---

# huaweicloud_rds_subscription_monitor

Use this data source to query the monitor information of an RDS subscription within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
variable "instance_id" {}
variable "subscription_id" {}

data "huaweicloud_rds_subscription_monitor" "test" {
  instance_id     = var.instance_id
  subscription_id = var.subscription_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the subscription monitor.  
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the RDS instance.

* `subscription_id` - (Required, String) Specifies the ID of the subscription.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `status` - The running status of the snapshot agent associated with the subscription.  
  The valid values are as follows:
  + **started**: Started.
  + **succeeded**: Succeeded.
  + **in_progress**: In progress.
  + **idle**: Idle.
  + **retrying**: Retrying.
  + **failed**: Failed.

* `latency` - The longest latency of data changes, in seconds.

* `agent_not_running` - The duration that the agent has not been running, in hours.

* `pending_cmd_count` - The number of unexecuted commands for the subscription.

* `last_dist_sync` - The last time the distribution agent ran. The format is yyyy-mm-ddThh:mm:ssZ.

* `estimated_process_time` - The estimated time to complete the unexecuted commands, in seconds.
