---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_collector_channel_instances"
description: |-
  Use this data source to get the list of SecMaster collector channel instances within HuaweiCloud.
---

# huaweicloud_secmaster_collector_channel_instances

Use this data source to get the list of SecMaster collector channel instances within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}

data "huaweicloud_secmaster_collector_channel_instances" "test" {
  workspace_id = var.workspace_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID.

* `channel_id` - (Optional, String) Specifies the channel ID.

* `node_id` - (Optional, String) Specifies the node ID.

* `node_name` - (Optional, String) Specifies the node name.

* `sort_key` - (Optional, String) Specifies the sort key.

* `sort_dir` - (Optional, String) Specifies the sort direction. Supported values are **asc** and **desc**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `records` - The collector channel instances list.

  The [records](#records_struct) structure is documented below.

<a name="records_struct"></a>
The `records` block supports:

* `channel_name` - The channel name.

* `config_status` - The collector channel configuration status. Valid values are **OK**, **CHANGE**.

* `create_by` - The IAM user ID.

* `node_name` - The node name.

* `region` - The region.

* `mini_on_online` - Whether online.

* `public_ip_address` - The public IP address.

* `private_ip_address` - The private IP address.

* `monitor` - The monitor information.

  The [monitor](#records_monitor_struct) structure is documented below.

* `read_write` - The read write record information.

  The [read_write](#records_read_write_struct) structure is documented below.

<a name="records_monitor_struct"></a>
The `monitor` block supports:

* `mini_on_online` - Whether online.

* `memory_count` - The number of physical memory modules.

* `memory_usage` - The amount of physical memory used.

* `memory_free` - The amount of currently free physical memory.

* `memory_shared` - The total amount of memory shared by multiple processes.

* `memory_cache` - The memory size of cached data.

* `cpu_usage` - The current CPU usage rate.

* `cpu_idle` - The percentage of CPU idle time.

* `up_pps` - The number of upload data packets per second.

* `down_pps` - The number of download data packets per second.

* `write_rate` - The disk write rate.

* `read_rate` - The disk read rate.

* `disk_count` - The number of disk devices in the system.

* `disk_usage` - The current disk space usage.

* `heart_beat_time` - The time when the last heartbeat signal was received, in ISO 8601 format.

* `health_status` - The health status of the node. Valid values are **NORMAL**, **ANOMALIES**, **FAULTS**,
  **LOST_CONTACT**.

* `heart_beat` - Whether the node successfully received heartbeat signals. Valid values are **ONLINE**, **OFFLINE**.

<a name="records_read_write_struct"></a>
The `read_write` block supports:

* `channel_id` - The channel ID (UUID).

* `minion_id` - The minion ID (UUID).

* `accept_count` - The accept count.

* `send_count` - The send count.

* `accept_rate` - The accept rate.

* `send_rate` - The send rate.

* `heart_beat_time` - The time when the last heartbeat signal was received.

* `latest_transmission_time` - The time of the last transmission.

* `channel_instance_count` - The number of collector channel instances.

* `heart_beat` - Whether the node successfully received heartbeat signals. Valid values are **ONLINE**, **OFFLINE**.
