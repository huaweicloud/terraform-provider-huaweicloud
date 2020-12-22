---
subcategory: "Intelligent EdgeCloud (IEC)"
---

# huaweicloud\_iec\_vip

Manages a VIP resource within HuaweiCloud IEC.

## Example Usage

```hcl
data "huaweicloud_iec_sites" "iec_sites" {}

resource "huaweicloud_iec_vpc" "vpc_test" {
  name = "iec-vpc-test"
  cidr = "192.168.0.0/16"
  mode = "SYSTEM"
}

resource "huaweicloud_iec_vpc_subnet" "subnet_test" {
  name        = "iec-subnet-test"
  cidr        = "192.168.199.0/24"
  vpc_id      = huaweicloud_iec_vpc.vpc_test.id
  site_id     = data.huaweicloud_iec_sites.iec_sites.sites[0].id
}

resource "huaweicloud_iec_vip" "vip_test" {
  subnet_id = huaweicloud_iec_vpc_subnet.subnet_test.id
}
```

## Argument Reference

The following arguments are supported:

* `subnet_id` - (Required, ForceNew) Specifies the subnet in which to allocate 
    IP address for this vip. Changing this parameter creates a new vip resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the vip.

* `mac_address` - The MAC address of the vip.

* `fixed_ips` - An Array of one or more network. The fixed_ips object structure is documented below.

The `fixed_ips` block supports:

* `subnet_id` - The id of the subnet network.
* `ip_address` - The ip address of the subnet network.

## Timeouts

This resource provides the following timeouts configuration options:
- `create` - Default is 10 minute.
- `delete` - Default is 10 minute.