---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_csms_tasks"
description: |-
  Use this data source to get the list of CSMS tasks within HuaweiCloud.
---

# huaweicloud_csms_tasks

Use this data source to get the list of CSMS tasks within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_csms_tasks" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `secret_name` - (Optional, String) Specifies the name of the secret.

* `status` - (Optional, String) Specifies the task status. Valid values are:
  + **SUCCESS**: Task rotation successful.
  + **FAILED**: Task rotation failed.

* `task_id` - (Optional, String) Specifies the task ID. This parameter cannot exist at the same time as other parameters.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tasks` - The secret task list.

  The [tasks](#tasks_struct) structure is documented below.

<a name="tasks_struct"></a>
The `tasks` block supports:

* `operate_type` - The rotation type.

* `task_error_msg` - The task error information.

* `task_id` - The task ID.

* `rotation_func_urn` - The URN of a FunctionGraph function.

* `task_status` - The task status.

* `task_error_code` - The task error code.

* `secret_name` - The secret name.

* `attempt_nums` - The number of task attempts.

* `task_time` - The time when a task is created.
