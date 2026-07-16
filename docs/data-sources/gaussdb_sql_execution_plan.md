---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_sql_execution_plan"
description: |-
  Use this data source to query the SQL execution plan information of a GaussDB instance within HuaweiCloud.
---

# huaweicloud_gaussdb_sql_execution_plan

Use this data source to query the SQL execution plan information of a GaussDB instance within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "pid" {}
variable "node_id" {}
variable "comp_id" {}

data "huaweicloud_gaussdb_sql_execution_plan" "test" {
  instance_id = var.instance_id
  pid         = var.pid
  node_id     = var.node_id
  comp_id     = var.comp_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the SQL execution plan.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the GaussDB instance.

* `pid` - (Required, String) Specifies the thread ID.

* `node_id` - (Required, String) Specifies the node ID.

* `comp_id` - (Required, String) Specifies the component ID on the node.
  For distributed instances, only CN components are supported.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `gs_explain` - The SQL execution plan information.
