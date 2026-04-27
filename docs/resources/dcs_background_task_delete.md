---
subcategory: "Distributed Cache Service (DCS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dcs_background_task_delete"
description: |-
  Manages a DCS background task delete resource within HuaweiCloud.
---

# huaweicloud_dcs_background_task_delete

Manages a DCS background task delete resource within HuaweiCloud.

## Example Usage

### Delete a background task

```hcl
variable "instance_id" {}
variable "task_id" {}

resource "huaweicloud_dcs_background_task_delete" "test" {
  instance_id = var.instance_id
  task_id     = var.task_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to manage the background task.
  If omitted, the provider-level region will be used. This parameter is non-updatable.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the DCS instance.

* `task_id` - (Required, String, ForceNew) Specifies the ID of the background task.

## Attributes

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which is the `task_id`.
