---
subcategory: "GeminiDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_geminidb_scheduled_tasks"
description: |-
  Use this data source to query the list of scheduled tasks for GeminiDB within HuaweiCloud.
---

# huaweicloud_geminidb_scheduled_tasks

Use this data source to query the list of scheduled tasks for GeminiDB within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_geminidb_scheduled_tasks" "test" {
  job_status  = "Pending"
  job_name    = "REBOOT"
  instance_id = var.instance_id
  start_time  = "2019-05-27T03:38:51+0000"
  end_time    = "2019-05-28T03:38:51+0000"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) The region in which to query the scheduled jobs.
  If omitted, the provider-level region will be used.

* `job_name` - (Optional, String) The name of the scheduled job. Valid values are:
  - **REBOOT**: Reboot instance
  - **RESIZE_FLAVOR**: Change instance CPU and memory specifications
  - **UPGRADE_DATABASE**: Patch upgrade

* `job_status` - (Optional, String) The execution status of the job. Valid values are:
  - **Pending**: Task not executed
  - **Running**: Task is executing
  - **Completed**: Task executed successfully
  - **Failed**: Task execution failed
  - **Canceled**: Task canceled

* `instance_id` - (Optional, String) The ID of the GeminiDB instance.
  If not specified, all instances that meet the criteria will be queried.

* `start_time` - (Optional, String) The start time for task creation. Format is **yyyy-mm-ddThh:mm:ssZ**.
  If not specified, defaults to 7 days before the current time.

* `end_time` - (Optional, String) The end time for task creation. Format is **yyyy-mm-ddThh:mm:ssZ**.
  If not specified, defaults to the current time.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `schedules` - The list of parameter template application histories.
  The [schedules](#geminidb_schedules) structure is documented below.

<a name="geminidb_schedules"></a>
The `schedules` block supports:

* `job_id` - The job ID.

* `job_name` - The job name.

* `job_status` - The job execution status.

* `instance_id` - The instance ID.

* `instance_name` - The instance name.

* `instance_status` - The instance status.

* `datastore_type` - The database type.

* `create_time` - The task creation time. Format is **yyyy-mm-ddThh:mm:ssZ**.

* `start_time` - The task start time. Format is **yyyy-mm-ddThh:mm:ssZ**.

* `end_time` - The task end time. Format is **yyyy-mm-ddThh:mm:ssZ**.
