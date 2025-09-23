---
subcategory: "Virtual Private Cloud (VPC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpc_subnet_private_ip"
description: |-
  Manages a VPC subnet private IP resource within HuaweiCloud.
---

# huaweicloud_vpc_subnet_private_ip

Manages a VPC subnet private IP resource within HuaweiCloud.

## Example Usage

```hcl
variable "subnet_id" {}
variable "ip_address" {}

resource "huaweicloud_vpc_subnet_private_ip" "test" {
  subnet_id  = var.subnet_id
  ip_address = var.ip_address
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `subnet_id` - (Required, String, NonUpdatable) Specifies the ID of the subnet to which the private IP belongs.

* `ip_address` - (Optional, String, NonUpdatable) Specifies the IP address. The value must be an unused address
  within the subnet cidr. If it is not specified, the system automatically assigns an IP address.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - The status of the private IP address. Possible values are **ACTIVE** and **DOWN**.

* `device_owner` - The resource using the private IP address. The parameter is left blank if it is not used.

## Import

The private IP can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_vpc_subnet_private_ip.test <id>
```
