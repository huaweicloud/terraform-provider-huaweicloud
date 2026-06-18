---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_sql_execution_plan_status"
description: |-
  Use this data source to query the binding status of a SQL execution plan of a GaussDB instance within HuaweiCloud.
---

# huaweicloud_gaussdb_sql_execution_plan_status

Use this data source to query the binding status of a SQL execution plan of a GaussDB instance within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "node_id" {}
variable "sql_id" {}

data "huaweicloud_gaussdb_sql_execution_plan_status" "test" {
  instance_id = var.instance_id
  node_id     = var.node_id
  sql_id      = var.sql_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the SQL execution plan status.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the GaussDB instance.

* `node_id` - (Required, String) Specifies the node ID of the GaussDB instance.

* `sql_id` - (Required, String) Specifies the SQL ID to query the execution plan status.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `sql_plan_bind_state_list` - The list of SQL execution plan bind states.
  The [sql_plan_bind_state_list](#sql_plan_bind_state_list_struct) structure is documented below.

<a name="sql_plan_bind_state_list_struct"></a>
The `sql_plan_bind_state_list` block supports:

* `outline` - The outline of the current execution plan.

* `cost` - The cost of the SQL execution plan.

* `status` - The binding status of the execution plan.
  The valid values are as follows:
  + **bind**: The execution plan is bound.
  + **drop**: The execution plan is unbound.

* `sql_hash` - The hash value of the SQL text.

* `plan_id` - The ID of the SQL execution plan.
