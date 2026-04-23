---
subcategory: "Resource Formation (RFS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rfs_execution_plans"
description: |-
  Use this data source to list all execution plans of a specified stack.
---

# huaweicloud_rfs_execution_plans

Use this data source to list all execution plans of a specified stack.

## Example Usage

```hcl
variable "stack_name" {}

data "huaweicloud_rfs_execution_plans" "test" {
  stack_name = var.stack_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `stack_name` - (Required, String) Specifies the name of the resource stack.

* `stack_id` - (Optional, String) Specifies the unique ID of the resource stack for strong matching.
  If it does not match the current resource stack, the API will return an error.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `execution_plans` - The list of execution plans.

  The [execution_plans](#execution_plans_struct) structure is documented below.

<a name="execution_plans_struct"></a>
The `execution_plans` block supports:

* `stack_name` - The name of the stack.

* `stack_id` - The unique ID of the stack.

* `execution_plan_id` - The ID of the execution plan.

* `execution_plan_name` - The name of the execution plan.

* `description` - The description of the execution plan.

* `status` - The status of the execution plan. Valid values are:
  + **CREATION_IN_PROGRESS**: The execution plan is being created.
  + **CREATION_FAILED**: Creating the execution plan failed. Please retrieve the error message summary from status_message.
  + **AVAILABLE**: The execution plan has been created and is ready to be applied.
  + **APPLY_IN_PROGRESS**: The execution plan is being applied.
  + **APPLIED**: The execution plan has been applied.

* `status_message` - The status message of the execution plan, which provides a brief error summary for debugging
  when the status is **CREATION_FAILED** or **DELETE_FAILED**.

* `create_time` - The creation time of the execution plan, in RFC3339 format (yyyy-mm-ddTHH:MM:SSZ),
  such as **1970-01-01T00:00:00Z**.

* `apply_time` - The application time of the execution plan, in RFC3339 format (yyyy-mm-ddTHH:MM:SSZ),
  such as **1970-01-01T00:00:00Z**.
