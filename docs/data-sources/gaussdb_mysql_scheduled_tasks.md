---
subcategory: "GaussDB(for MySQL)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_mysql_scheduled_tasks"
description: |-
  Use this data source to get the list of GaussDB MySQL scheduled tasks.
---

# huaweicloud_gaussdb_mysql_scheduled_tasks

Use this data source to get the list of GaussDB MySQL scheduled tasks.

## Example Usage

```hcl
data "huaweicloud_gaussdb_mysql_scheduled_tasks" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `status` - (Optional, String) Specifies the task execution status. Value options:
  + **Running**: The task is being executed.
  + **Completed**: The task is successfully executed.
  + **Failed**: The task failed to be executed.
  + **Pending**: The task is not executed.
  + **Canceled**: The task is canceled.

* `job_id` - (Optional, String) Specifies the task ID.

* `job_name` - (Optional, String) Specifies the task scheduling type.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tasks` - Indicates the list of scheduled task details.

  The [tasks](#tasks_struct) structure is documented below.

<a name="tasks_struct"></a>
The `tasks` block supports:

* `job_id` - Indicates the task ID.

* `instance_id` - Indicates the instance ID.

* `instance_name` - Indicates the instance name.

* `instance_status` - Indicates the instance status.

* `project_id` - Indicates the project ID of a tenant in a region.

* `job_name` - Indicates the task name.

* `create_time` - Indicates the task creation time in the **yyyy-mm-ddThh:mm:ssZ** format.

* `start_time` - Indicates the task start time in the **yyyy-mm-ddThh:mm:ssZ** format.

* `end_time` - Indicates the task end time in the **yyyy-mm-ddThh:mm:ssZ** format.

* `job_status` - Indicates the task execution status.

* `datastore_type` - Indicates the database type.
