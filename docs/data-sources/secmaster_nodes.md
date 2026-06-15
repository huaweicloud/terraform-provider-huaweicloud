---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_nodes"
description: |-
  Use this data source to query SecMaster nodes within HuaweiCloud.
---

# huaweicloud_secmaster_nodes

Use this data source to query SecMaster nodes within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}

data "huaweicloud_secmaster_nodes" "test" {
  workspace_id = var.workspace_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID.

* `node_id` - (Optional, String) Specifies the node ID.

* `node_name` - (Optional, String) Specifies the node name.

* `sort_key` - (Optional, String) Specifies the sorting field.

* `sort_dir` - (Optional, String) Specifies the sorting direction.
  The valid values are as follows:
  + **asc**: Ascending.
  + **desc**: Descending.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `count` - The total number of nodes.

* `records` - The list of nodes.

  The [records](#records_struct) structure is documented below.

<a name="records_struct"></a>
The `records` block supports:

* `create_by` - The IAM user ID.

* `create_time` - The creation time in milliseconds.

* `description` - The error description when node is abnormal.

  The [description](#description_struct) structure is documented below.

* `device_type` - The device type.

* `ip_address` - The IP address.

* `monitor` - The monitoring metrics.

  The [monitor](#monitor_struct) structure is documented below.

* `node_expansion` - The node expansion information.

  The [node_expansion](#node_expansion_struct) structure is documented below.

* `node_id` - The node UUID.

* `node_name` - The tenant name.

* `os_type` - The operating system type.

* `private_ip_address` - The private IP address.

* `region` - The region.

* `specification` - The node specification.

* `subnet_id` - The subnet ID.

* `update_time` - The update time in milliseconds.

* `vpc_endpoint_address` - The VPC endpoint address.

* `vpc_endpoint_id` - The VPC endpoint ID.

* `vpc_id` - The VPC ID.

* `vpcep_service_ip` - The VPC endpoint service IP address.

<a name="description_struct"></a>
The `description` block supports:

* `error_code` - The error code.

* `error_msg` - The error description.

<a name="monitor_struct"></a>
The `monitor` block supports:

* `cpu_idle` - The CPU idle ratio.

* `cpu_usage` - The CPU usage.

* `disk_count` - The number of disk devices.

* `disk_usage` - The disk usage.

* `down_pps` - The downstream packets per second.

* `health_status` - The node health status.
  The valid values are **NORMAL**, **ANOMALIES**, **FAULTS** and **LOST_CONTACT**.

* `heart_beat` - The node heartbeat status.
  The valid values are **ONLINE** and **OFFLINE**.

* `heart_beat_time` - The last heartbeat time.

* `memory_cache` - The memory cache size.

* `memory_count` - The number of physical memory modules.

* `memory_free` - The free physical memory size.

* `memory_shared` - The shared memory size.

* `memory_usage` - The used physical memory size.

* `mini_on_online` - The online status identifier.

* `read_rate` - The disk read rate.

* `up_pps` - The upstream packets per second.

* `write_rate` - The disk write rate.

<a name="node_expansion_struct"></a>
The `node_expansion` block supports:

* `custom_label` - The custom label.

* `data_center` - The data center.

* `description` - The description.

* `maintainer` - The maintainer.

* `network_plane` - The network plane.

* `node_id` - The node UUID.
