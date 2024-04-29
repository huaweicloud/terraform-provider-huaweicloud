---
subcategory: "GeminiDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_redis_eip_associate"
description: ""
---

# huaweicloud_gaussdb_redis_eip_associate

GeminiDB Redis node EIP associate resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "node_id" {}
variable "public_ip" {}

resource "huaweicloud_gaussdb_redis_eip_associate" "test"{
  instance_id = var.instance_id
  node_id     = var.node_id
  public_ip   = var.public_ip
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of a GaussDB Redis instance.

  Changing this parameter will create a new resource.

* `node_id` - (Required, String, ForceNew) Specifies the ID of a GaussDB Redis node.

  Changing this parameter will create a new resource.

* `public_ip` - (Required, String, ForceNew) Specifies the EIP address to associate. The value must be in the
  standard IP address format.

  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID which is formatted `<instance_id>/<node_id>`.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `delete` - Default is 10 minutes.

## Import

The GaussDB Redis node EIP associate can be imported using `instance_id` and `node_id` separated by a slash, e.g.

```bash
$ terraform import huaweicloud_gaussdb_redis_eip_associate.test <instance_id>/<node_id>
```
