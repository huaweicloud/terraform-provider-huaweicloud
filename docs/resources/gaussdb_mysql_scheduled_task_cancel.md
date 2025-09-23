---
subcategory: "GaussDB(for MySQL)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_mysql_scheduled_task_cancel"
description: |-
  Manages a GaussDB MySQL scheduled task cancel resource within HuaweiCloud.
---

# huaweicloud_gaussdb_mysql_scheduled_task_cancel

Manages a GaussDB MySQL scheduled task cancel resource within HuaweiCloud.

-> This resource is only a one-time action resource for operating the API.
Deleting this resource will not clear the corresponding request record,
but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "job_id" {}

resource "huaweicloud_gaussdb_mysql_scheduled_task_cancel" "test" {
  job_id = var.job_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `job_id` - (Required, String, ForceNew) Specifies the task ID. Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID. The value is `job_id`.

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
