---
subcategory: "Elastic IP (EIP)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpc_eipv3_associate"
description: |-
  Manages a VPC EIP associates with an instance resource within HuaweiCloud.
---

# huaweicloud_vpc_eipv3_associate

Manages a VPC EIP associates with an instance resource within HuaweiCloud.

## Example Usage

```hcl
variable "publicip_id" {}
variable "elb_id" {}

resource "huaweicloud_vpc_eipv3_associate" "test" {
  publicip_id             = var.publicip_id
  associate_instance_type = "ELB"
  associate_instance_id   = var.elb_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the VPC EIP associate resource. If
  omitted, the provider-level region will be used. Changing this creates a new resource.

* `publicip_id` - (Required, String, ForceNew) Specifies the ID of a EIP. Changing this creates a new resource.

* `associate_instance_type` - (Required, String, ForceNew) Specifies the type of the instance that the port belongs to.
  Value options: **PORT**, **NATGW**, **VPN** and **ELB**. Changing this creates a new resource.

* `associate_instance_id` - (Required, String, ForceNew) Specifies the ID of the instance that the port belongs to.
  Changing this creates a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID. The value is the `publicip_id`.

## Import

The VPC EIP associations can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_vpc_eipv3_associate.eip <id>
```
