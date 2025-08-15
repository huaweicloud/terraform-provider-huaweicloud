---
subcategory: "Virtual Private Cloud (VPC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpc_subnet_cidr_reservation"
description: |-
  Manages a VPC subnet CIDR reservation resource within HuaweiCloud.
---

# huaweicloud_vpc_subnet_cidr_reservation

Manages a VPC subnet CIDR reservation resource within HuaweiCloud.

## Example Usage

```hcl
variable "subnet_id" {}

resource "huaweicloud_vpc_subnet_cidr_reservation" "test" {
  subnet_id   = var.subnet_id
  ip_version  = 4
  mask        = 24
  name        = "test-reservation"
  description = "test subnet CIDR reservation"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `subnet_id` - (Required, String, NonUpdatable) Specifies the ID of the virtual subnet to which the
  CIDR reservation belongs.

* `ip_version` - (Required, Int, NonUpdatable) Specifies the IP version of the subnet CIDR reservation (4 or 6).

* `cidr` - (Optional, String, NonUpdatable) Specifies the reserved CIDR block in CIDR notation.
  Conflicts with `mask`.

* `mask` - (Optional, Int, NonUpdatable) Specifies the subnet mask length.
  Conflicts with `cidr`.

* `name` - (Optional, String) Specifies the name of the subnet CIDR reservation (1-64 characters).

* `description` - (Optional, String) Specifies the description of the subnet CIDR reservation (max 255 characters).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `vpc_id` - The ID of the VPC to which the subnet belongs.

* `project_id` - The project ID of the subnet CIDR reservation.

* `created_at` - The creation time of the subnet CIDR reservation.

* `updated_at` - The last update time of the subnet CIDR reservation.

## Import

Subnet CIDR reservations can be imported using their `id`:

```bash
terraform import huaweicloud_vpc_subnet_cidr_reservation.test <id>
```
