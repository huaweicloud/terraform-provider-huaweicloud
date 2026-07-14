---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_sql_stack_information"
description: |-
  Use this data source to query the SQL stack information of a specified GaussDB instance within HuaweiCloud.
---

# huaweicloud_gaussdb_sql_stack_information

Use this data source to query the SQL stack information of a specified GaussDB instance within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
variable "instance_id" {}
variable "pid" {}
variable "node_id" {}

data "huaweicloud_gaussdb_sql_stack_information" "test" {
  instance_id = var.instance_id
  pid         = var.pid
  node_id     = var.node_id
}
```

### Query with Component ID

```hcl
variable "instance_id" {}
variable "pid" {}
variable "node_id" {}
variable "comp_id" {}

data "huaweicloud_gaussdb_sql_stack_information" "test" {
  instance_id = var.instance_id
  pid         = var.pid
  node_id     = var.node_id
  comp_id     = var.comp_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the SQL stack information. If omitted, the
  provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the GaussDB instance.

* `pid` - (Required, String) Specifies the thread ID.

* `node_id` - (Required, String) Specifies the node ID.

* `comp_id` - (Optional, String) Specifies the component ID on the node.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `gs_stack` - The SQL stack information.
