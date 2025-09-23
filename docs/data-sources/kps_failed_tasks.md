---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_kps_failed_tasks"
description: |-
  Use this data source to get a list of the tasks that failed to bind, unbind, reset or replace key pairs.
---

# huaweicloud_kps_failed_tasks

Use this data source to get a list of the tasks that failed to bind, unbind, reset or replace key pairs.

## Example Usage

```hcl
data "huaweicloud_kps_failed_tasks" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tasks` - The list of the failed tasks.

  The [tasks](#tasks_struct) structure is documented below.

<a name="tasks_struct"></a>
The `tasks` block supports:

* `id` - The ID of the task.

* `server_id` - The ID of the instance associated with the task.

* `server_name` - The name of the instance associated with the task.

* `operate_type` - The operation type of the task.
  The value can be **FAILED_RESET**, **FAILED_REPLACE** or **FAILED_UNBIND**.

* `keypair_name` - The name of the keypair associated with the task.

* `task_error_msg` - The error information of the task execution failure.

* `task_error_code` - The error code of the task execution failure.

* `task_time` - The start time of the task, in RFC3339 format.
