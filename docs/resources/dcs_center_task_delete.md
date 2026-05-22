---
subcategory: "Distributed Cache Service (DCS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dcs_center_task_delete"
description: |-
  Manages a DCS center task removal resource within HuaweiCloud.
---

# huaweicloud_dcs_center_task_delete

Manages a DCS center task removal resource within HuaweiCloud.

## Example Usage

```hcl
variable "task_id" {}

resource "huaweicloud_dcs_center_task_delete" "test" {
  task_id = var.task_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to delete the task. If omitted, the
  provider-level region will be used. Changing this parameter will create a new resource.

* `task_id` - (Required, String, NonUpdatable) Specifies the ID of the task to be deleted.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.
