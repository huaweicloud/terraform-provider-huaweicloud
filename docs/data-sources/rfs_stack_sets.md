---
subcategory: "Resource Formation (RFS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rfs_stack_sets"
description: |-
  Use this data source to list all RFS stack sets.
---

# huaweicloud_rfs_stack_sets

Use this data source to list all RFS stack sets.

## Example Usage

```hcl
data "huaweicloud_rfs_stack_sets" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `filter` - (Optional, String) Specifies the filter condition.
  + The AND operator is defined by commas (,), and the OR operator is defined by vertical bars (|).
  + The OR operator has higher priority than the AND operator.
  + Parentheses are not supported.
  + Only double equal signs (==) are supported for filter operators.
  + Filter parameter names and their values only support uppercase and lowercase English letters, numbers, and underscores.
  + Semicolons are prohibited in filter conditions. If a semicolon is present, the filter will be ignored.
  + A filter parameter can only be associated with one AND condition, and multiple OR conditions in an AND condition can
  only be associated with one filter parameter.

* `sort_key` - (Optional, String) Specifies the sorting field, only supports **create_time**.

* `sort_dir` - (Optional, String) Specifies whether to sort in ascending or descending order.
  Valid values are:
  + **asc**: Ascending order
  + **desc**: Descending order

* `call_identity` - (Optional, String) Specifies the identity used to call the stack set.
  This parameter is only supported when the permission model of the stack set is **SERVICE_MANAGED**.
  Valid values are:
  + **SELF**: Call as a management account.
  + **DELEGATED_ADMIN**: Call as a delegated administrator.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `stack_sets` - The list of stack sets.

  The [stack_sets](#stack_sets_struct) structure is documented below.

<a name="stack_sets_struct"></a>
The `stack_sets` block supports:

* `stack_set_id` - The unique ID of the stack set.

* `stack_set_name` - The name of the stack set.

* `stack_set_description` - The description of the stack set.

* `permission_model` - The permission model of the stack set.
  Valid values are:
  + **SELF_MANAGED**: Self-managed mode. Users need to manually create delegations in advance.
  + **SERVICE_MANAGED**: Service-managed mode. RFS automatically creates all required IAM delegations.

* `status` - The status of the stack set.
  Valid values are:
  + **IDLE**: The stack set is idle.
  + **OPERATION_IN_PROGRESS**: The stack set is in operation.
  + **DEACTIVATED**: The stack set is deactivated.

* `create_time` - The creation time of the stack set, in RFC3339 format (yyyy-mm-ddTHH:MM:SSZ),
  such as **1970-01-01T00:00:00.000Z**.

* `update_time` - The update time of the stack set, in RFC3339 format (yyyy-mm-ddTHH:MM:SSZ),
  such as **1970-01-01T00:00:00.000Z**.
