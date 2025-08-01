---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_task_delete"
description: |-
  Manages RDS task delete resource within HuaweiCloud.
---

# huaweicloud_rds_task_delete

Manages RDS task delete resource within HuaweiCloud.

## Example Usage

```hcl
variable "job_id" {}

resource "huaweicloud_rds_task_delete" "test" {
  job_id = var.job_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `job_id` - (Required, String, NonUpdatable) Specifies the ID of the RDS task.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID which equals the `job_id`.
