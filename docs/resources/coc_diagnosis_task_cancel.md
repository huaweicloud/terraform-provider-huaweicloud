---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_diagnosis_task_cancel"
description: |-
  Manages a COC diagnosis task cancel resource within HuaweiCloud.
---

# huaweicloud_coc_diagnosis_task_cancel

Manages a COC diagnosis task cancel resource within HuaweiCloud.

~> Deleting diagnosis task cancel resource is not supported, it will only be removed from the state.

## Example Usage

```hcl
variable "task_id" {}

resource "huaweicloud_coc_diagnosis_task_cancel" "test" {
  task_id = var.task_id
}
```

## Argument Reference

The following arguments are supported:

* `task_id` - (Required, String, NonUpdatable) Specifies the diagnostic task ticket ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which equals to `task_id`.
