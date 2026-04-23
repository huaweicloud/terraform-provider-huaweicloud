---
subcategory: "Resource Formation (RFS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rfs_stack_set_operations"
description: |-
  Use this datasource to get the list of stack set operations.
---

# huaweicloud_rfs_stack_set_operations

Use this datasource to get the list of stack set operations.

## Example Usage

```hcl
variable "stack_set_name" {}

data "huaweicloud_rfs_stack_set_operations" "test" {
  stack_set_name = var.stack_set_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `stack_set_name` - (Required, String) Specifies the name of the stack set.

* `stack_set_id` - (Optional, String) Specifies the ID of the stack set.

* `filter` - (Optional, String) Specifies the filter condition.  
  The use of this parameter must comply with the following conditions:
  + The AND operator is defined with a comma (`,`).
  + The OR operator is defined with a vertical line (`|`), and the OR operator has higher priority than the
    AND operator.
  + Does not support parentheses.
  + The filter only supports the double equals (`==`) operator.
  + The filter parameter name and its value can only contain English letters, digits, and underscores.
  + Semicolons are not allowed; if a semicolon exists, the filter entry is ignored.  
  + A filtering parameter can only be related to one condition, and one condition or multiple conditions in the
    condition can only be related to one filtering parameter.

* `sort_key` - (Optional, String) Specifies the sort field, only supports **create_time**.

* `sort_dir` - (Optional, String) Specifies the ascending or descending order.
  Valid values are **asc** and **desc**.

* `call_identity` - (Optional, String) This parameter is only used when the stack set permission model is
  **SERVICE_MANAGED**. It specifies whether the call is made as the organization management account or as the
  delegated service administrator in a member account. The default is **SELF**.

  The valid values are as follows:
  + **SELF**: Call as an organizational management account.
  + **DELEGATED_ADMIN**: Call as a service delegation administrator. The user's Huawei Cloud account must have been
    registered as a delegated administrator for the "Resource orchestration Resource Stack Service" within the
    organization.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `stack_set_operations` - The list of stack set operations.

  The [stack_set_operations](#stack_set_operations_struct) structure is documented below.

<a name="stack_set_operations_struct"></a>
The `stack_set_operations` block supports:

* `operation_id` - The ID of the stack set operation.

* `stack_set_id` - The unique ID of the stack set.

* `stack_set_name` - The name of the stack set.

* `action` - The current user action.

* `status` - The status of the stack set operation.

* `status_message` - A brief message when the stack set operation fails.

* `create_time` - The creation time of a stack set operation. It is represented in UTC format
  (YYYY-MM-DDTHH:mm:ss.SSSZ), such as **1970-01-01T00:00:00.000Z**.

* `update_time` - The update time of a stack set operation. It is represented in UTC format
  (YYYY-MM-DDTHH:mm:ss.SSSZ), such as **1970-01-01T00:00:00.000Z**.
