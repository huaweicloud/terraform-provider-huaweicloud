---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_diagnosis_task_retry"
description: |-
  Manages a COC diagnosis task retry resource within HuaweiCloud.
---

# huaweicloud_coc_diagnosis_task_retry

Manages a COC diagnosis task retry resource within HuaweiCloud.

~> Deleting diagnosis task retry resource is not supported, it will only be removed from the state.

## Example Usage

```hcl
variable "task_id" {}
variable "instance_id" {}

resource "huaweicloud_coc_diagnosis_task_retry" "test" {
  task_id     = var.task_id
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `task_id` - (Required, String, NonUpdatable) Specifies the diagnostic task ticket ID.

* `instance_id` - (Optional, String, NonUpdatable) Specifies the instance ID to be retried.

  -> The instance ID must be in the work order submitted for this task. The diagnostic task corresponding to this
  instance ID must be in a failed state.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which equals to `task_id`.
