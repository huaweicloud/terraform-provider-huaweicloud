---
subcategory: "Intelligent EdgeCloud (IEC)"
---

# huaweicloud\_iec\_vpc\_subnet

Manages a VPC subnet resource within HuaweiCloud IEC.

## Example Usage

```hcl
data "huaweicloud_iec_sites" "sites_test" {}

resource "huaweicloud_iec_vpc" "vpc_test" {
  name = "vpc_demo"
  cidr = "192.168.0.0/16"
  mode = "CUSTOMER"
}

resource "huaweicloud_iec_vpc_subnet" "subnet_test" {
  name       = "subnet_demo"
  cidr       = "192.168.128.0/18"
  vpc_id     = huaweicloud_iec_vpc.vpc_test.id
  site_id    = data.huaweicloud_iec_sites.sites_test.sites[0].id
  gateway_ip = "192.168.128.1"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) The name of the iec subnet. Changing this updates the 
    name of the existing subnet.

* `cidr` - (Required, String, ForceNew) CIDR representing IP range for this 
    subnet, based on IP version. Changing this parameter creates a new subnet 
    resource.

* `vpc_id` - (Required, String, ForceNew) Specifies the ID of the iec vpc. Changing 
    this parameter creates a new subnet resource.

* `site_id` - (Required, String, ForceNew) Specifies the ID of the iec site. 
    Changing this parameter creates a new subnet resource.

* `gateway_ip` - (Required, String, ForceNew)  Default gateway used by devices in 
    this subnet. Changing this parameter creates a new subnet resource.

* `dhcp_enable` - (Optional, Bool) The administrative state of the network.
    The value must be "true".

* `dns_list` - (Optional, List) An array of DNS name server names used by hosts
    in this subnet.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies a resource ID in UUID format.

* `site_info` - Specifies the information of the iec site.

* `status` - Specifies the status of the subnet.

## Timeouts

This resource provides the following timeouts configuration options:
- `create` - Default is 10 minute.
- `delete` - Default is 3 minute.

## Import

IEC subnet can be imported using the `id`, e.g.

```
$ terraform import huaweicloud_iec_vpc_subnet.subnet_demo 51be9f2b-5a3b-406a-9271-36f0c929fbcc
```
