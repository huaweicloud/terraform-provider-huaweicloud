---
subcategory: "GaussDB(for MySQL)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_mysql_eip_associate"
description: |-
  Manages GaussDB mysql EIP associate resource within HuaweiCloud.
---

# huaweicloud_gaussdb_mysql_eip_associate

Manages GaussDB mysql EIP associate resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "public_ip" {}
variable "public_ip_id" {}

resource "huaweicloud_gaussdb_mysql_eip_associate" "test"{
  instance_id  = var.instance_id
  public_ip    = var.public_ip
  public_ip_id = var.public_ip_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of a GaussDB MySQL instance. Changing this parameter
  will create a new resource.

* `public_ip` - (Required, String, ForceNew) Specifies the EIP address to be bound. The value must be in the standard IP
  address format. Changing this parameter will create a new resource.

* `public_ip_id` - (Required, String, ForceNew) Specifies the ID of the EIP. Changing this parameter will create a new
  resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID which equals the `instance_id`.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
* `delete` - Default is 30 minutes.

## Import

The GaussDB MySQL EIP associate can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_gaussdb_mysql_eip_associate.test <id>
```
