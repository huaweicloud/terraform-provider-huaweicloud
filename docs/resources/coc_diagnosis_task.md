---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_diagnosis_task"
description: |-
  Manages a COC diagnosis task resource within HuaweiCloud.
---

# huaweicloud_coc_diagnosis_task

Manages a COC diagnosis task resource within HuaweiCloud.

~> Deleting diagnosis task resource is not supported, it will only be removed from the state.

## Example Usage

```hcl
variable "resource_id" {}

resource "huaweicloud_coc_diagnosis_task" "test" {
  resource_id = var.resource_id
  type        = "ECS"
}
```

## Argument Reference

The following arguments are supported:

* `resource_id` - (Required, String, NonUpdatable) Specifies the ID of the instance to be diagnosed.

* `type` - (Required, String, NonUpdatable) Specifies the type of instance being diagnosed.
  Values can be **ECS**, **RDS**, **DCS**, **DMS** or **ELB**.

* `extra_properties` - (Optional, String, NonUpdatable) Specifies the additional parameters applicable to RDS, DMS, DCS,
  ELB, etc.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `code` - Indicates the step code.

* `project_id` - Indicates the project ID to which the diagnosed instance belongs.

* `user_id` - Indicates the user ID to which the diagnostic record belongs.

* `user_name` - Indicates the user name to which the diagnostic record belongs.

* `progress` - Indicates the progress of the diagnostic task execution.

* `work_order_id` - Indicates the diagnostic task ticket ID.

* `instance_name` - Indicates the name of the instance being diagnosed.

* `status` - Indicates the execution status of the diagnostic task.

* `start_time` - Indicates the start time.

* `end_time` - Indicates the end time.

* `instance_num` - Indicates the number of instances included in the diagnostic task.

* `os_type` - Indicates the operating system type of the diagnosed instance.

* `region` - Indicates the region to which the diagnostic resource belongs.

* `node_list` - Indicates the diagnostic step structure object.

  The [node_list](#node_list_struct) structure is documented below.

* `message` - Indicates the diagnostic report.

<a name="node_list_struct"></a>
The `node_list` block supports:

* `id` - Indicates the diagnostic task node ID.

* `code` - Indicates the step code.

* `name` - Indicates the diagnosis step name.

* `name_zh` - Indicates the Chinese name of the diagnostic step.

* `diagnosis_task_id` - Indicates the diagnostic task ID.

* `status` - Indicates the execution status of the diagnostic task.

## Import

The COC group can be imported using `resource_id` and `id`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_coc_group.test <resource_id>/<id>
```
