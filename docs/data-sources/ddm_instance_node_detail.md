---
subcategory: "Distributed Database Middleware (DDM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ddm_instance_node_detail"
description: |-
  Use this data source to obtain the detail info of the DDM instance node.
---

# huaweicloud_ddm_instance_node_detail

Use this data source to obtain the detail info of the DDM instance node.

## Example Usage

```hcl
variable "instance_id" {}
variable "node_id" {}

data "huaweicloud_ddm_instance_node_detail" "test" {
  instance_id = var.instance_id
  node_id     = var.node_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the DDM instance.

* `node_id` - (Required, String) Specifies the ID of the DDM instance node.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `status` - Indicates the node status.

* `name` - Indicates the node name.

* `private_ip` - Indicates the private IP address of the node.

* `floating_ip` - Indicates the floating IP address of the node.

* `server_id` - Indicates the VM ID.

* `subnet_name` - Indicates the subnet name.

* `datavolume_id` - Indicates the data disk ID.

* `res_subnet_ip` - Indicates the IP address provided by the resource subnet.

* `systemvolume_id` - Indicates the system disk ID.
