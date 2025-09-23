---
subcategory: "Dedicated Load Balance (Dedicated ELB)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_elb_asynchronous_tasks"
description: |-
  Use this data source to get the list of asynchronous tasks.
---

# huaweicloud_elb_asynchronous_tasks

Use this data source to get the list of asynchronous tasks.

## Example Usage

```hcl
data "huaweicloud_elb_asynchronous_tasks" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `job_id` - (Optional, String) Specifies the task ID.

* `job_type` - (Optional, String) Specifies the task type.

* `status` - (Optional, String) Specifies the task status.
  Value options: **INIT**, **RUNNING**, **FAIL**, **SUCCESS**, **ROLLBACKING**, **COMPLETE**, **ROLLBACK_FAIL**, and **CANCEL**.

* `error_code` - (Optional, String) Specifies the error code of the task.

* `resource_id` - (Optional, String) Specifies the resource ID.

* `begin_time` - (Optional, String) Specifies the time when the task started, in the format of **yyyy-MM-dd'T'HH:mm:ss**.
  The tasks that started on or after the specified time will be returned.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `jobs` - Indicates the task list.

  The [jobs](#jobs_struct) structure is documented below.

<a name="jobs_struct"></a>
The `jobs` block supports:

* `job_id` - Indicates the task ID.

* `job_type` - Indicates the task type.

* `status` - Indicates the task status.

* `resource_id` - Indicates  the resource ID.

* `project_id` - Indicates the project ID.

* `begin_time` - Indicates the time when the task was started.

* `end_time` - Indicates the time when the task was ended.

* `error_code` - Indicates the task error code.

* `error_msg` - Indicates the task error message.

* `sub_jobs` - Indicates the subtask list.

  The [sub_jobs](#jobs_sub_jobs_struct) structure is documented below.

<a name="jobs_sub_jobs_struct"></a>
The `sub_jobs` block supports:

* `job_id` - Indicates  the task ID.

* `job_type` - Indicates the task type.

* `status` - Indicates the task status.

* `resource_id` - Indicates the resource ID.

* `project_id` - Indicates the project ID.

* `begin_time` - Indicates the time when the task was started.

* `end_time` - Indicates  the time when the task was ended.

* `error_code` - Indicates the task error code.

* `error_msg` - Indicates  the task error message.

* `entities` - Indicates the resource to be processed in a subtask.

  The [entities](#sub_jobs_entities_struct) structure is documented below.

<a name="sub_jobs_entities_struct"></a>
The `entities` block supports:

* `resource_id` - Indicates the ID of the resource associated with a subtask.

* `resource_type` - Indicates the type of the resource associated with a subtask.
