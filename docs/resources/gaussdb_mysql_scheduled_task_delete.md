---
subcategory: "GaussDB(for MySQL)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_mysql_scheduled_task_delete"
description: |-
  Manages a GaussDB MySQL scheduled task delete resource within HuaweiCloud.
---

# huaweicloud_gaussdb_mysql_scheduled_task_delete

Manages a GaussDB MySQL scheduled task delete resource within HuaweiCloud.

-> This resource is only a one-time action resource for operating the API.
Deleting this resource will not clear the corresponding request record,
but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "instance_id" {}
variable "job_id" {}

resource "huaweicloud_gaussdb_mysql_scheduled_task_delete" "test" {
  instance_id = var.instance_id
  job_id      = var.job_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the instance ID. Changing this parameter will create a new resource.

* `job_id` - (Required, String, ForceNew) Specifies the task ID. Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID. The value is `job_id`.
