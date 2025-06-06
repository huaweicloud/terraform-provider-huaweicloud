---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_kps_failed_task_delete"
description: |-
  Manages a KPS failed task delete resource within HuaweiCloud.
---

# huaweicloud_kps_failed_task_delete

Manages a KPS failed task delete resource within HuaweiCloud.

-> The current resource is a one-time resource, and destroying this resource will not recover the deleted task,
but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "task_id" {}

resource "huaweicloud_kps_failed_task_delete" "test" {
  task_id = var.task_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `task_id` - (Required, String, NonUpdatable) Specifies the task ID, which identifies the failed task to be deleted.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
