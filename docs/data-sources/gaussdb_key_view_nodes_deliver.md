---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_key_view_nodes_deliver"
description: |-
  Use this data source to query the nodes that can deliver key views within HuaweiCloud.
---

# huaweicloud_gaussdb_key_view_nodes_deliver

Use this data source to query the nodes that can deliver key views within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_gaussdb_key_view_nodes_deliver" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the nodes.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the GaussDB instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `nodes` - The list of nodes that can deliver key views.
  The [nodes](#key_view_nodes_deliver_nodes_attr) structure is documented below.

<a name="key_view_nodes_deliver_nodes_attr"></a>
The `nodes` block supports:

* `node_id` - The ID of the node.

* `node_name` - The name of the node.

* `role` - The role of the node.
  The valid values are as follows:
  + **master**: The master node.
  + **slave**: The slave node.
  + **secondary**: The log node.
  + **readreplica**: The read-only node.

* `type` - The component type of the node.
  The valid values are as follows:
  + **CN**: The CN component.
  + **DN**: The DN component.

* `distributed_id` - The distributed ID of the node.

* `component_id` - The component ID of the node.

* `detail` - The description of the node.
