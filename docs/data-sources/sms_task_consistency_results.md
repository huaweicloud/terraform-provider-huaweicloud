---
subcategory: "Server Migration Service (SMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_sms_task_consistency_results"
description: |-
  Use this data source to get the list of SMS task consistency results.
---

# huaweicloud_sms_task_consistency_results

Use this data source to get the list of SMS task consistency results.

## Example Usage

```hcl
variable "task_id" {}

data "huaweicloud_sms_task_consistency_results" "test" {
  task_id = var.task_id
}
```

## Argument Reference

The following arguments are supported:

* `task_id` - (Required, String) Specifies the task ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `result_list` - Indicates the consistency verification results.

  The [result_list](#result_list_struct) structure is documented below.

<a name="result_list_struct"></a>
The `result_list` block supports:

* `finished_time` - Indicates the verification completion time.

* `check_result` - Indicates the verification execution result.

* `consistency_result` - Indicates the verification results.

  The [consistency_result](#result_list_consistency_result_struct) structure is documented below.

<a name="result_list_consistency_result_struct"></a>
The `consistency_result` block supports:

* `dir_check` - Indicates the directory verified.

* `num_total_files` - Indicates the total number of files verified.

* `num_different_files` - Indicates the number of files inconsistent.

* `num_target_miss_files` - Indicates the number of files missing at the target.

* `num_target_more_files` - Indicates the number of files redundant at the target.
