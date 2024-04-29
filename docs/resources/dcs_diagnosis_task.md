---
subcategory: "Distributed Cache Service (DCS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dcs_diagnosis_task"
description: ""
---

# huaweicloud_dcs_diagnosis_task

Manages a DCS diagnosis task resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_dcs_diagnosis_task" "test" {
  instance_id  = var.instance_id
  begin_time   = "2024-03-11T01:17:48.998Z"
  end_time     = "2024-03-11T01:27:48.998Z"
  node_ip_list = ["10.168.179.171"]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the DCS instance.
  Changing this creates a new resource.

* `begin_time` - (Required, String, ForceNew) Specifies the start time of the diagnosis task, in RFC3339 format.
  Changing this creates a new resource.

* `end_time` - (Required, String, ForceNew) Specifies the end time of the diagnosis task, in RFC3339 format.
  Changing this creates a new resource.

* `node_ip_list` - (Optional, List, ForceNew) Specifies the IP addresses of diagnosed nodes.
  By default, all nodes are diagnosed. Changing this creates a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID. The value is the dianosis report ID.

* `abnormal_item_sum` - Indicates the total number of abnormal diagnosis items.

* `failed_item_sum` - Indicates the total number of failed diagnosis items.

* `diagnosis_node_report_list` - Indicates the list of node diagnosis report
  The [diagnosis_node_report_list](#diagnosis_node_report_list) structure is documented below.

<a name="diagnosis_node_report_list"></a>
The `diagnosis_node_report_list` block supports:

* `node_ip` - Indicates the IP address of the node diagnosed.

* `role` - Indicates the node role. The value can be **master** or **slave**.

* `group_name` - Indicates the name of the shard where the node is.

* `az_code` - Indicates the code of the AZ where the node is.

* `abnormal_sum` - Indicates the total number of abnormal diagnosis items.

* `failed_sum` - Indicates the total number of failed diagnosis items.

* `is_faulted` - Indicates whether the node is faulted.

* `diagnosis_dimension_list` - Indicates the diagnosis dimension list.
  The [diagnosis_dimension_list](#diagnosis_dimension_list) structure is documented below.

* `command_time_taken_list` - Indicates the command execution duration list.
  The [command_time_taken_list](#command_time_taken_list) structure is documented below.

<a name="diagnosis_dimension_list"></a>
The `diagnosis_dimension_list` block supports:

* `name` - Indicates the diagnosis dimension name. The value can be **network**, **storage** or **load**.

* `abnormal_num` - Indicates the total number of abnormal diagnosis items.

* `failed_num` - Indicates the total number of failed diagnosis items.

* `diagnosis_item_list` - Indicates the diagnosis items.
  The [diagnosis_item_list](#diagnosis_item_list) structure is documented below.

<a name="command_time_taken_list"></a>
The `command_time_taken_list` block supports:

* `total_num` - Indicates the total number of times that commands are executed.

* `total_usec_sum` - Indicates the total duration of command execution.

* `result` - Indicates the command execution latency result. The value can be **succeed** or **failed**.

* `error_code` - Indicates the error code for the command time taken.

* `command_list` - Indicates the command execution latency statistics.
  The [command_list](#command_list) structure is documented below.

<a name="diagnosis_item_list"></a>
The `diagnosis_item_list` block supports:

* `name` - Indicates the diagnosis item name.
  The value can be **connection_num**, **rx_controlled**, **persistence**, **centralized_expiration**,
  **inner_memory_fragmentation**, **time_consuming_commands**, **hit_ratio**, **memory_usage** or **cpu_usage**.

* `result` - Indicates the diagnosis result. The value can be **failed**, **abnormal** or **normal**.

* `error_code` - Indicates the error code for the diagnosis item.

* `advice_ids` - Indicates the list of suggestion IDs.
  The [advice_ids](#conclusion_item) structure is documented below.

* `cause_ids` - Indicates the list of cause IDs.
  The [cause_ids](#conclusion_item) structure is documented below.

* `impact_ids` - Indicates the list of impact IDs.
  The [impact_ids](#conclusion_item) structure is documented below.

<a name="command_list"></a>
The `command_list` block supports:

* `command_name` - Indicates the command name.

* `average_usec` - Indicates the average duration of calls.

* `calls_sum` - Indicates the number of calls.

* `per_usec` - Indicates the duration percentage.

* `usec_sum` - Indicates the total time consumed.

<a name="conclusion_item"></a>
The `conclusion_item` block supports:

* `id` - Indicates the conclusion ID.

* `params` - Indicates the conclusion parameters.

## Timeouts

This resource provides the following timeout configuration option:

* `create` - Default is 30 minutes.

## Import

The DCS diagnosis task can be imported using `instance_id` and `id` separated by a slash, e.g.

```bash
$ terraform import huaweicloud_dcs_diagnosis_task.test <instance_id>/<id>
```
