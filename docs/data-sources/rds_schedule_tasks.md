---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_schedule_tasks"
description: |-
  Use this data source to get the list of RDS schedule tasks.
---

# huaweicloud_rds_schedule_tasks

Use this data source to get the list of RDS schedule tasks.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_rds_schedule_tasks" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Optional, String) Specifies the ID of the instance.

* `instance_name` - (Optional, String) Specifies the name of the instance.

* `status` - (Optional, String) Specifies the status of the task. Value options: **Initing**, **Pending**, **Running**,
  **Completed**, **Failed**, **Unauthorized**, **Canceled**.

* `start_time` - (Optional, String) Specifies the start time in UTC timestamp format (milliseconds since epoch).
  `end_time` is mandatory if it is not empty.

* `end_time` - (Optional, String) Specifies the end time in UTC timestamp format (milliseconds since epoch).
  `start_time` is mandatory if it is not empty.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `schedule_tasks` - Indicates the list of schedule tasks.

  The [schedule_tasks](#schedule_tasks_struct) structure is documented below.

<a name="schedule_tasks_struct"></a>
The `schedule_tasks` block supports:

* `id` - Indicates the task ID.

* `name` - Indicates the task name.

* `instance_id` - Indicates the instance ID.

* `instance_name` - Indicates the instance name.

* `instance_status` - Indicates the instance status. The value can be: **BUILD**, **CREATE FAIL**, **ACTIVE**, **FAILED**,
  **FROZEN**, **MODIFYING**, **REBOOTING**, **RESTORING**, **MODIFYING INSTANCE TYPE**, **SWITCHOVER**, **MIGRATING**,
  **BACKING UP**, **MODIFYING DATABASE PORT**, **STORAGE FULL**.

* `datastore_type` - Indicates the DB type.

* `status` - Indicates the task execution status. The value can be: **Initing**, **Pending**, **Running**, **Completed**,
  **Failed**, **Unauthorized**, **Canceled**.

* `order` - Indicates the task order. Value ranges: **1-100**.

* `volume_type` - Indicates the volume type.

* `create_time` -  Indicates the creation time. The value is in the **yyyy-mm-ddThh:mm:ssZ** format.

* `start_time` - Indicates the start time. The value is in the **yyyy-mm-ddThh:mm:ssZ** format.

* `end_time` - Indicates the end time. The value is in the **yyyy-mm-ddThh:mm:ssZ** format.

* `target_config` - Indicates the target config.
  The [target_config](#schedule_tasks_target_config_struct) structure is documented below.

<a name="schedule_tasks_target_config_struct"></a>
The `target_config` block supports:

* `cpu` - Indicates the target cpu for flavor update task when name is **RESIZE_FLAVOR**.

* `flavor` - Indicates the target flavor for flavor update task when name is **RESIZE_FLAVOR**.

* `mem` - Indicates the target mem for flavor update task when name is **RESIZE_FLAVOR**.
