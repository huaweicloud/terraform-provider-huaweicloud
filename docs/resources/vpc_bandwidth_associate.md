---
subcategory: "Elastic IP (EIP)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpc_bandwidth_associate"
description: ""
---

# huaweicloud_vpc_bandwidth_associate

Associates an EIP to a specified **shared** bandwidth.

-> Yearly/monthly EIPs cannot be added to a shared bandwidth. After an EIP is removed from a shared bandwidth,
  a dedicated bandwidth will be allocated to the EIP. By default, the dedicated bandwidth will be billed by bandwidth
  and the size is 5 Mbit/s. You can configure the bandwidth as needed.

## Example Usage

### Associate an EIP

```hcl
variable "public_id" {}

resource "huaweicloud_vpc_bandwidth" "test" {
  name = "bandwidth_1"
  size = 100
}

resource "huaweicloud_vpc_bandwidth_associate" "test" {
  bandwidth_id = huaweicloud_vpc_bandwidth.test.id
  eip_id       = var.public_id
}
```

### Associate multiple EIPs

```hcl
variable "multiple_eips" {
  type = list(string)
}

resource "huaweicloud_vpc_bandwidth" "test" {
  name = "bandwidth_1"
  size = 100
}

resource "huaweicloud_vpc_bandwidth_associate" "test" {
  count = length(var.multiple_eips)

  bandwidth_id = huaweicloud_vpc_bandwidth.test.id
  eip_id       = var.multiple_eips[count.index]
}
```

### Associate an EIP managed by Terraform

```hcl
resource "huaweicloud_vpc_eip" "dedicated" {
  publicip {
    type = "5_bgp"
  }

  bandwidth {
    share_type  = "PER"
    name        = "dedicated"
    size        = 5
    charge_mode = "traffic"
  }

  lifecycle {
    ignore_changes = [ bandwidth ]
  }
}

resource "huaweicloud_vpc_bandwidth" "test" {
  name = "bandwidth_1"
  size = 100
}

resource "huaweicloud_vpc_bandwidth_associate" "test" {
  bandwidth_id = huaweicloud_vpc_bandwidth.test.id
  eip_id       = huaweicloud_vpc_eip.dedicated.id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to associate the bandwidth. If omitted,
  the provider-level region will be used. Changing this creates a new resource.

* `bandwidth_id` - (Required, String, ForceNew) Specifies the shared bandwidth ID to associate.
  Changing this creates a new resource.

* `eip_id` - (Required, String) Specifies the ID of the EIP that uses the bandwidth.

* `bandwidth_charge_mode` - (Optional, String) Specifies the billing mode of the dedicated bandwidth used by the EIP that
  has been removed from a shared bandwidth. The value can be **bandwidth** or **traffic**. If not specified, the dedicated
  bandwidth will be billed by bandwidth.

* `bandwidth_size` - (Optional, Int) Specifies the size (Mbit/s) of the dedicated bandwidth used by the EIP that
  has been removed from a shared bandwidth. The default bandwidth size is 5 Mbit/s.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in format of `<bandwidth_id>/<eip_id>`.
* `public_ip` - The EIP address.
* `bandwidth_name` - The shared bandwidth name.

## Import

Bandwidth associations can be imported using the `bandwidth_id` and `eip_id` separated by a slash, e.g.:

```bash
$ terraform import huaweicloud_vpc_bandwidth_associate.eip <bandwidth_id>/<eip_id>
```
