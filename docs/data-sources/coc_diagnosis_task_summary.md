---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_diagnosis_task_summary"
description: |-
  Use this data source to get the list of COC diagnosis task summary.
---

# huaweicloud_coc_diagnosis_task_summary

Use this data source to get the list of COC diagnosis task summary.

## Example Usage

```hcl
variable "task_id" {}

data "huaweicloud_coc_diagnosis_task_summary" "test" {
  task_id = var.task_id
}
```

## Argument Reference

The following arguments are supported:

* `task_id` - (Required, String) Specifies the diagnostic task ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `status` - Indicates the work order execution status.

* `region` - Indicates the region where the diagnosed instance is located.

* `start_time` - Indicates the start time.

* `instance_summary` - Indicates the list of summary information about the diagnosed instance.

  The [instance_summary](#instance_summary_struct) structure is documented below.

<a name="instance_summary_struct"></a>
The `instance_summary` block supports:

* `instance_id` - Indicates the ID of the diagnosed instance.

* `instance_name` - Indicates the name of the instance being diagnosed.

* `progress` - Indicates the execution progress of the diagnostic task.

* `status` - Indicates the execution status of the diagnostic task.

* `normal_item_num` - Indicates the number of normal diagnostic items of the instance.

* `abnormal_item_num` - Indicates the number of abnormal diagnosis items of the instance.
