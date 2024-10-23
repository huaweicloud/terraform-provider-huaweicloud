---
subcategory: "Distributed Database Middleware (DDM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ddm_instance_nodes"
description: ""
---

# huaweicloud_ddm_instance_nodes

Use this data source to get the list of DDM instance nodes.

## Example Usage

```hcl
variable "ddm_instance_id" {}

data "huaweicloud_ddm_instance_nodes" "test" {
  instance_id = var.ddm_instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of DDM instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `nodes` - Indicates the list of DDM instance node.
  The [Node](#DdmInstanceNodes_Node) structure is documented below.

<a name="DdmInstanceNodes_Node"></a>
The `Node` block supports:

* `id` - Indicates the ID of the DDM instance node.

* `ip` - Indicates the IP address of the DDM instance node.

* `port` - Indicates the port of the DDM instance node.

* `status` - Indicates the status of the DDM instance node.
