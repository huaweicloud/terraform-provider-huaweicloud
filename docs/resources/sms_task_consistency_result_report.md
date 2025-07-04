---
subcategory: "Server Migration Service (SMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_sms_task_consistency_result_report"
description: |-
  Manages an SMS task consistency result report resource within HuaweiCloud.
---

# huaweicloud_sms_task_consistency_result_report

Manages an SMS task consistency result report resource within HuaweiCloud.

~> Deleting task consistency result report resource is not supported, it will only be removed from the state.

## Example Usage

```hcl
variable "task_id" {}

resource "huaweicloud_sms_task_consistency_result_report" "test" {
  task_id = var.task_id

  consistency_result {
    dir_check             = "/root/dev"
    num_total_files       = 1
    num_different_files   = 1
    num_target_miss_files = 1
    num_target_more_files = 1
  }
}
```

## Argument Reference

The following arguments are supported:

* `task_id` - (Required, String, NonUpdatable) Specifies the task ID.

* `consistency_result` - (Optional, List, NonUpdatable) Specifies the consistency verification results.
  The [consistency_result](#consistency_result_struct) structure is documented below.

<a name="consistency_result_struct"></a>
The `consistency_result` block supports:

* `dir_check` - (Required, String, NonUpdatable) Specifies the directory verified.

* `num_total_files` - (Required, Int, NonUpdatable) Specifies the total number of files verified.

* `num_different_files` - (Required, Int, NonUpdatable) Specifies the number of files inconsistent.

* `num_target_miss_files` - (Required, Int, NonUpdatable) Specifies the number of files missing at the target.

* `num_target_more_files` - (Required, Int, NonUpdatable) Specifies the number of files redundant at the target.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which equals to `task_id`.
