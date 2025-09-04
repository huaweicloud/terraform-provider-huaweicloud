---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_component_running_nodes"
description: |-
  Use this data source to get the list of SecMaster component running nodes within HuaweiCloud.
---

# huaweicloud_secmaster_component_running_nodes

Use this data source to get the list of SecMaster component running nodes within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}
variable "component_id" {}

data "huaweicloud_secmaster_component_running_nodes" "test" {
  workspace_id = var.workspace_id
  component_id = var.component_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID.

* `component_id` - (Required, String) Specifies the component ID.

* `node_id` - (Optional, String) Specifies the node ID.

* `node_name` - (Optional, String) Specifies the node name.

* `sort_key` - (Optional, String) Specifies the attribute fields for sorting.

* `sort_dir` - (Optional, String) Specifies the sorting order. Supported values are **ASC** and **DESC**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `records` - The component running nodes list.

  The [records](#records_struct) structure is documented below.

<a name="records_struct"></a>
The `records` block supports:

* `component_id` - The component ID.

* `component_name` - The component name.

* `node_id` - The node ID.

* `create_time` - The creation time (timestamp in milliseconds).

* `node_name` - The node name.

* `specification` - The specification.

* `config_status` - The node configuration status. Valid values are **UN_SAVED**, **SAVE_AND_UN_APPLY**,
  **MOVE_AND_UN_APPLY**, **APPLYING**, **FAIL_APPLY**, **APPLIED**.

* `fail_deploy_message` - The deployment failure message.

* `ip_address` - The IP address.

* `private_ip_address` - The private IP address.

* `region` - The region.

* `vpc_endpoint_id` - The VPC endpoint ID.

* `vpc_endpoint_address` - The VPC endpoint address.

* `monitor` - The monitor information.

  The [monitor](#records_monitor_struct) structure is documented below.

* `node_expansion` - The node expansion information.

  The [node_expansion](#records_node_expansion_struct) structure is documented below.

* `node_apply_fail_enum` - The node application success or failure status and reason.
  The valid values are as follows:
  + **COLLECTOR_USE**: The collector is in use and cannot be removed.
  + **NODE_OFFLINE**: Node is in a disconnected state and cannot be applied.

* `list` - The component configuration parameter list.

  The [list](#records_list_struct) structure is documented below.

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

<a name="records_node_expansion_struct"></a>
The `node_expansion` block supports:

* `node_id` - The node ID.

* `data_center` - The data center.

* `custom_label` - The custom label.

* `network_plane` - The network plane.

* `description` - The description information.

* `maintainer` - The maintainer.

<a name="records_list_struct"></a>
The `list` block supports:

* `configuration_id` - The configuration ID.

* `component_id` - The component ID.

* `node_id` - The node ID.

* `file_name` - The file name.

* `file_path` - The file path.

* `file_type` - The file type. Valid values are **JVM**, **LOG4J2**, **YML**.

* `param` - The parameter.

* `version` - The version.

* `type` - The configuration type. Valid values are **HISTORY**, **CURRENT_SAVE**, **CURRENT_APPLY**.
