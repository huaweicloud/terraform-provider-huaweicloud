---
subcategory: "Document Database Service (DDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dds_bind_gateway"
description: |-
  Manages a resource to bind the public gateway to the DDS instance node within HuaweiCloud.
---

# huaweicloud_dds_readonly_node

Manages a resource to bind the public gateway to the DDS instance node within HuaweiCloud.

-> Before use this resource, you need to pay attention to the following:
  <br/>1. Only the primary or secondary node of a replica set instance is supported.
  <br/>2. Only the mongos node of a cluster instance is supported.
  <br/>3. This operation cannot be performed on frozen or abnormal instances.

## Example Usage

```hcl
variable "instance_id" {}
variable "node_id" {}
variable "nat_gateway_id" {}
variable "public_ip_id" {}
variable "external_service_port" {}

resource "huaweicloud_dds_bind_gateway" "test"{
  instance_id           = var.instance_id
  node_id               = var.node_id
  nat_gateway_id        = var.nat_gateway_id
  public_ip_id          = var.public_ip_id
  external_service_port = var.external_service_port
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the DDS instance.

* `node_id` - (Required, String, NonUpdatable) Specifies the ID of the node.

* `nat_gateway_id` - (Required, String, NonUpdatable) Specifies the ID of the public NAT gateway.

* `public_ip_id` - (Required, String, NonUpdatable) Specifies the EIP ID.

* `external_service_port` - (Required, Int, NonUpdatable) Specifies the port of the EIP for providing services
  for external systems.
  The valid value ranges from `1` to `65,535`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, also is the node ID.

## Import

The resource can be imported using the `instance_id` and the `id` separated by a slash (/), e.g.

```bash
$ terraform import huaweicloud_dds_bind_gateway.test <instance_id>/<id>
```
