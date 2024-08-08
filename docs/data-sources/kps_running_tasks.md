---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_kps_running_tasks"
description: |-
  Use this data source to get a list of running tasks.
---

# huaweicloud_kps_running_tasks

Use this data source to get a list of running tasks.

## Example Usage

```hcl
data "huaweicloud_kps_running_tasks" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tasks` - The list of the running tasks.

  The [tasks](#tasks_struct) structure is documented below.

<a name="tasks_struct"></a>
The `tasks` block supports:

* `id` - The ID of the task.

* `server_id` - The ID of the instance associated with the task.

* `server_name` - The name of the instance associated with the task.

* `operate_type` - The operation type of the task.
  The value can be **RUNNING**.

* `keypair_name` - The name of the keypair associated with the task.

* `task_time` - The start time of the task, in RFC3339 format.
