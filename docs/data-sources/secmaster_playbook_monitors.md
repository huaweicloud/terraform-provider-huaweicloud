---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_playbook_monitors"
description: |-
  Use this data source to get the list of SecMaster playbook monitor details.
---

# huaweicloud_secmaster_playbook_monitors

Use this data source to get the list of SecMaster playbook monitor details.

## Example Usage

```hcl
variable "workspace_id" {}
variable "playbook_id" {}
variable "start_time" {}
variable "end_time" {}
variable "version_query_type" {}

data "huaweicloud_secmaster_playbook_monitors" "test" {
  workspace_id       = var.workspace_id
  playbook_id        = var.playbook_id
  start_time         = var.start_time
  end_time           = var.end_time
  version_query_type = var.version_query_type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID to which the playbook belongs.

* `playbook_id` - (Required, String) Specifies the playbook ID.

* `start_time` - (Required, String) Specifies the start time.
  For example, **2021-01-30T23:00:00Z+0800**.

* `end_time` - (Required, String) Specifies the end time.
  For example, **2021-01-30T23:00:00Z+0800**.

* `version_query_type` - (Required, String) Specifies the playbook version type.
  The value can be **ALL**, **VALID**, or **DELETED**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data` - The playbook running monitor details.

  The [data](#data_struct) structure is documented below.

<a name="data_struct"></a>
The `data` block supports:

* `total_instance_run_num` - The total running times.

* `average_run_time` - The average duration.

* `min_run_time_instance` - The workflow with the shortest running duration.

  The [min_run_time_instance](#run_time_instance_struct) structure is documented below.

* `success_instance_num` - The number of successful instances.

* `terminate_instance_num` - The number of terminated instances.

* `running_instance_num` - The number of running instances.

* `schedule_instance_run_num` - The number of scheduled trigger executions.

* `event_instance_run_num` - The time-triggered executions.

* `max_run_time_instance` - The workflow with the longest running duration.

  The [max_run_time_instance](#run_time_instance_struct) structure is documented below.

* `total_instance_num` - The total number of playbook instances.

* `fail_instance_num` - The number of failed instances.

<a name="run_time_instance_struct"></a>
The `min_run_time_instance` and `max_run_time_instance` block supports:

* `playbook_instance_id` - The playbook instance ID.

* `playbook_instance_name` - The playbook instance name.

* `playbook_instance_run_time` - The playbook instance running time.
