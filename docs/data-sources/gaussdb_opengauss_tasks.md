---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_opengauss_tasks"
description: |-
  Use this data source to get the list of GaussDB OpenGauss tasks.
---

# huaweicloud_gaussdb_opengauss_tasks

Use this data source to get the list of GaussDB OpenGauss tasks.

## Example Usage

```hcl
data "huaweicloud_gaussdb_opengauss_tasks" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `status` - (Optional, String) Specifies the task execution status. Value options:
  + **Running**: The task is being executed.
  + **Completed**: The task is successfully executed.
  + **Failed**: The task failed to be executed.

* `name` - (Optional, String) Specifies the task name. Value options:
  + **CreateGaussDBV5Instance**: Creating a DB instance.
  + **BackupSnapshotGaussDBV5InInstance**: Creating a manual backup.
  + **CloneGaussDBV5NewInstance**: Restoring data to a new DB instance.
  + **RestoreGaussDBV5InInstance**: Restoring data to the original DB instance.
  + **RestoreGaussDBV5InInstanceToExistedInst**: Restoring data to an existing DB instance.
  + **DeleteGaussDBV5Instance**: Deleting a DB instance.
  + **EnlargeGaussDBV5Volume**: Scaling up storage.
  + **ResizeGaussDBV5Flavor**: Changing specifications.
  + **GaussDBV5ExpandClusterCN**: Adding coordinator nodes.
  + **GaussDBV5ExpandClusterDN**: Adding shards.

* `start_time` - (Optional, String) Specifies the start time in the **yyyy-mm-ddThh:mm:ssZ** format.

* `end_time` - (Optional, String) Specifies the end time in the **yyyy-mm-ddThh:mm:ssZ** format.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `tasks` - Indicates the task list.

  The [tasks](#tasks_struct) structure is documented below.

<a name="tasks_struct"></a>
The `tasks` block supports:

* `instance_id` - Indicates the instance ID.

* `instance_name` - Indicates the instance name.

* `instance_status` - Indicates the iInstance status.

* `job_id` - Indicates the task ID.

* `name` - Indicates the task name.

* `status` - Indicates the task execution status.

* `process` - Indicates the task progress.

* `fail_reason` - Indicates the task failure cause.

* `created_at` - Indicates the task creation time in the **yyyy-mm-ddThh:mm:ssZ** format.

* `ended_at` - Indicates the task end time in the **yyyy-mm-ddThh:mm:ssZ** format.
