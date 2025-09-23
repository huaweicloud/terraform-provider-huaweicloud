---
subcategory: "Document Database Service (DDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dds_scheduled_tasks"
description: |-
  Use this data source to get the list of DDS scheduled tasks.
---

# huaweicloud_dds_scheduled_tasks

Use this data source to get the list of DDS scheduled tasks.

## Example Usage

```hcl
variable "instance_id" {}
variable "job_status" {}
variable "job_name" {}
variable "start_time" {}
variable "end_time" {}

data "huaweicloud_dds_scheduled_tasks" "test" {
  instance_id = var.instance_id
  job_status  = var.job_status
  job_name    = var.job_name
  start_time  = var.start_time
  end_time    = var.end_time
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Optional, String) Specifies the instance ID.

* `job_name` - (Optional, String) Specifies the task name.
  Value can be as follows:
  + **RESIZE_FLAVOR**: Change vCPUs and memory of an instance.

* `job_status` - (Optional, String) Specifies the task execution status.
  Value can be as follows:
  + **Pending**: The task is not executed.
  + **Running**: The task is being executed.
  + **Completed**: The task is successfully executed.
  + **Failed**: The task fails to be executed.
  + **Canceled**: The task is canceled.

* `start_time` - (Optional, String) Specifies the start time.
  The format of the start time is **yyyy-mm-ddThh:mm:ssZ**.
  Defaults to seven days before the current time.

* `end_time` - (Optional, String) Specifies the end time.
  The format of the end time is **yyyy-mm-ddThh:mm:ssZ**.
  Defaults to current time.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `schedules` - Indicates the tasks list.

  The [schedules](#schedules_struct) structure is documented below.

<a name="schedules_struct"></a>
The `schedules` block supports:

* `job_id` - Indicates the task ID.

* `job_name` - Indicates the task name.

* `create_time` - Indicates the create time.

* `start_time` - Indicates the start time.

* `end_time` - Indicates the end time.

* `job_status` - Indicates the task execution status.

* `instance_id` - Indicates the instance ID.

* `instance_name` - Indicates the instance name.

* `instance_status` - Indicates the instance status.
  Values can be as follows:
  + **createfail**: The instance failed to be created.
  + **creating**: The instance is being created.
  + **normal**: The instance is running properly.
  + **abnormal**: The instance is abnormal.
  + **deleted**: The instance has been deleted.
