---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_scheduled_task_record_details"
description: |-
  Use this data source to get the details list of the scheduled task execution record within HuaweiCloud.
---

# huaweicloud_workspace_scheduled_task_record_details

Use this data source to get the details list of the scheduled task execution record within HuaweiCloud.

## Example Usage

```hcl
variable "task_id" {}
variable "record_id" {}

data "huaweicloud_workspace_scheduled_task_record_details" "test" {
  task_id   = var.task_id
  record_id = var.record_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the scheduled task record details are located.  
  If omitted, the provider-level region will be used.

* `task_id` - (Required, String) Specifies the ID of the scheduled task to be queried.

* `record_id` - (Required, String) Specifies the ID of the scheduled task execution record to be queried.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `details` - The list of scheduled task execution record details.  
  The [details](#scheduled_task_record_details_attr) structure is documented below.

<a name="scheduled_task_record_details_attr"></a>
The `details` block supports:

* `id` - The ID of the scheduled task execution record detail.

* `record_id` - The ID of the scheduled task execution record.

* `desktop_id` - The ID of the desktop.

* `desktop_name` - The name of the desktop.

* `exec_status` - The execution status.

* `exec_script_id` - The ID of the execution script.

* `result_code` - The error code of the failure or skip reason.

* `fail_reason` - The failure or skip reason.

* `start_time` - The execution start time, in RFC3339 format.

* `end_time` - The execution end time, in RFC3339 format.

* `time_zone` - The time zone information.
