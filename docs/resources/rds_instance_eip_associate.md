---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_instance_eip_associate"
description: |-
  Manages an RDS instance EIP associate resource within HuaweiCloud.
---

# huaweicloud_rds_instance_eip_associate

Manages an RDS instance EIP associate resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "public_ip" {}
variable "public_ip_id" {}

resource "huaweicloud_rds_instance_eip_associate" "test" {
  instance_id  = var.instance_id
  public_ip    = var.public_ip
  public_ip_id = var.node_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of an RDS instance.

* `public_ip` - (Required, String, NonUpdatable) Specifies the EIP address to be bound.

* `public_ip_id` - (Required, String, NonUpdatable) Specifies the EIP ID.

## Attribute Reference

* `id` - The resource ID. The value is `instance_id`.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.

* `delete` - Default is 30 minutes.

## Import

The RDS instance eip associate can be imported using the `id`, e.g.

```bash
terraform import huaweicloud_rds_instance_eip_associate.test <id>
```
