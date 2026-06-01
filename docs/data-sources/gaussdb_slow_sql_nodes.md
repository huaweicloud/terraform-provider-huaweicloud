---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_slow_sql_nodes"
description: |-
  Use this data source to query the slow SQL nodes of GaussDB instances within HuaweiCloud.
---

# huaweicloud_gaussdb_slow_sql_nodes

Use this data source to query the slow SQL nodes of GaussDB instances within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_gaussdb_slow_sql_nodes" "test" {
  instance_id = var.instance_id
  action      = "slow"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the slow SQL nodes.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the instance ID.

* `action` - (Required, String) Specifies the type. The valid value is **slow**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `nodes` - The list of nodes with slow SQL.

  The [nodes](#gaussdb_slow_sql_nodes_struct) structure is documented below.

<a name="gaussdb_slow_sql_nodes_struct"></a>
The `nodes` block supports:

* `node_id` - The node ID.

* `node_name` - The node name.

* `role` - The node role. The valid values are as follows:
  + **master**: Primary node.
  + **slave**: Standby node.
  + **secondary**: Log node.
  + **readreplica**: Read replica.

* `instance_id` - The instance ID.

* `component_type` - The component type. The valid values are as follows:
  + **CN**: CN component.
  + **DN**: DN component.
