---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_delete_nodes"
description: |-
  Manages a delete nodes resource within HuaweiCloud.
---

# huaweicloud_secmaster_delete_nodes

Manages a delete nodes resource within HuaweiCloud.

-> 1. This resource is a one-time action resource used to delete SecMaster nodes. Deleting this resource will not
  restore the deleted nodes, but will only remove the resource information from the tfstate file.
  <br/>2. A successful API request does not guarantee that all nodes have been deleted successfully. Please check
  the `delete_fail_list` and `delete_success_list` for details.

## Example Usage

```hcl
variable "workspace_id" {}
variable "delete_ids" {
  type = list(string)
}

resource "huaweicloud_secmaster_delete_nodes" "test" {
  workspace_id = var.workspace_id
  delete_ids   = var.delete_ids
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `workspace_id` - (Required, String, NonUpdatable) Specifies the ID of the workspace to which the nodes belong.

* `delete_ids` - (Required, List, NonUpdatable) Specifies the list of node IDs to be deleted.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `delete_fail_list` - The list of nodes that failed to be deleted.
  The [delete_fail_list](#delete_nodes_node) structure is documented below.

* `delete_success_list` - The list of nodes that were successfully deleted.
  The [delete_success_list](#delete_nodes_node) structure is documented below.

<a name="delete_nodes_node"></a>
The `delete_fail_list` and `delete_success_list` block supports:

* `create_by` - The IAM user ID who created the node.

* `create_time` - The creation time, in milliseconds timestamp.

* `description` - The error object returned when the deletion fails.
  The [description](#delete_nodes_isap_error_rsp) structure is documented below.

* `device_type` - The device type of the node.

* `ip_address` - The IP address of the node.

* `monitor` - The monitor information of the node.
  The [monitor](#delete_nodes_monitor) structure is documented below.

* `node_expansion` - The node expansion information.
  The [node_expansion](#delete_nodes_isap_node_expansion) structure is documented below.

* `node_id` - The ID of the node, in UUID format.

* `node_name` - The name of the node.

* `os_type` - The operating system type of the node.

* `private_ip_address` - The private IP address of the node.

* `region` - The region of the node.

* `specification` - The specification of the node.

* `subnet_id` - The subnet ID of the node.

* `update_time` - The update time, in milliseconds timestamp.

* `vpc_endpoint_address` - The VPC endpoint address of the node.

* `vpc_endpoint_id` - The VPC endpoint ID of the node.

* `vpc_id` - The VPC ID of the node, in UUID format.

* `vpcep_service_ip` - The VPC endpoint service IP address of the node.

<a name="delete_nodes_isap_error_rsp"></a>
The `description` block supports:

* `error_code` - The error code.

* `error_msg` - The error message.

<a name="delete_nodes_monitor"></a>
The `monitor` block supports:

* `cpu_idle` - The percentage of CPU idle time.

* `cpu_usage` - The current CPU usage.

* `disk_count` - The number of disk devices in the system.

* `disk_usage` - The current disk space usage.

* `down_pps` - The number of download packets per second.

* `health_status` - The health status of the node.
  The valid values are as follows:
  + **NORMAL**: normal.
  + **ANOMALIES**: abnormal.
  + **FAULTS**: faulty.
  + **LOST_CONTACT**: lost contact.

* `heart_beat` - Whether the heartbeat signal is successfully received.
  The valid values are as follows:
  + **ONLINE**: online.
  + **OFFLINE**: offline.

* `heart_beat_time` - The time when the last heartbeat signal was received.

* `memory_cache` - The memory size used for caching data.

* `memory_count` - The number of physical memory modules.

* `memory_free` - The current free physical memory.

* `memory_shared` - The total amount of memory shared by multiple processes.

* `memory_usage` - The used physical memory.

* `mini_on_online` - Whether the node is online.

* `read_rate` - The disk read rate.

* `up_pps` - The number of upload packets per second.

* `write_rate` - The disk write rate.

<a name="delete_nodes_isap_node_expansion"></a>
The `node_expansion` block supports:

* `custom_label` - The custom label of the node.

* `data_center` - The data center of the node.

* `description` - The description of the node.

* `maintainer` - The maintainer of the node.

* `network_plane` - The network plane of the node.

* `node_id` - The ID of the node, in UUID format.
