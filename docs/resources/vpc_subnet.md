---
subcategory: "Virtual Private Cloud (VPC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpc_subnet"
description: ""
---

# huaweicloud_vpc_subnet

Provides a VPC subnet resource within HuaweiCloud.

## Example Usage

```hcl

variable "vpc_name" {}
variable "vpc_cidr" {}
variable "subnet_name" {}
variable "subnet_cidr" {}
variable "subnet_gateway_ip" {}
variable "availability_zone" {}

resource "huaweicloud_vpc" "vpc" {
  name = var.vpc_name
  cidr = var.vpc_cidr
}

resource "huaweicloud_vpc_subnet" "subnet" {
  name              = var.subnet_name
  cidr              = var.subnet_cidr
  gateway_ip        = var.subnet_gateway_ip
  vpc_id            = huaweicloud_vpc.vpc.id
  availability_zone = var.availability_zone
}

resource "huaweicloud_vpc_subnet" "subnet_with_tags" {
  name              = var.subnet_name
  cidr              = var.subnet_cidr
  gateway_ip        = var.subnet_gateway_ip
  vpc_id            = huaweicloud_vpc.vpc.id
  availability_zone = var.availability_zone

  tags = {
    foo = "bar"
    key = "value"
  }
}

resource "huaweicloud_vpc_subnet" "subnet_with_dhcp" {
  name              = var.subnet_name
  cidr              = var.subnet_cidr
  gateway_ip        = var.subnet_gateway_ip
  vpc_id            = huaweicloud_vpc.vpc.id
  availability_zone = var.availability_zone

  dhcp_lease_time      = "24h"
  dhcp_ipv6_lease_time = "4h"
  ntp_server_address   = "10.100.0.33,10.100.0.34"
  dhcp_domain_name     = "test.domainnanme"
}

 ```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the vpc subnet. If omitted, the
  provider-level region will be used. Changing this creates a new subnet.

* `name` - (Required, String) Specifies the subnet name. The value is a string of 1 to 64 characters that can contain
  letters, digits, underscores (_), and hyphens (-).

* `cidr` - (Required, String, ForceNew) Specifies the network segment on which the subnet resides. The value must be in
  CIDR format and within the CIDR block of the VPC. The subnet mask cannot be greater than 28. Changing this creates a
  new subnet.

* `gateway_ip` - (Required, String, ForceNew) Specifies the gateway of the subnet. The value must be a valid IP address
  in the subnet segment. Changing this creates a new subnet.

* `vpc_id` - (Required, String, ForceNew) Specifies the ID of the VPC to which the subnet belongs. Changing this creates
  a new subnet.

* `availability_zone` - (Optional, String, ForceNew) Specifies the availability zone (AZ) to which the subnet belongs.
  The value must be an existing AZ in the system. Changing this creates a new subnet.

* `description` - (Optional, String) Specifies supplementary information about the subnet. The value is a string of
  no more than 255 characters and cannot contain angle brackets (< or >).

* `ipv6_enable` - (Optional, Bool) Specifies whether the IPv6 function is enabled for the subnet. Defaults to false.

* `dhcp_enable` - (Optional, Bool) Specifies whether the DHCP function is enabled for the subnet. Defaults to true.

* `primary_dns` - (Optional, String) Specifies the IP address of DNS server 1 on the subnet. The value must be a valid
  IP address.

* `secondary_dns` - (Optional, String) Specifies the IP address of DNS server 2 on the subnet. The value must be a valid
  IP address.

* `dns_list` - (Optional, List) Specifies the DNS server address list of a subnet. This field is required if you need to
  use more than two DNS servers. This parameter value is the superset of both DNS server address 1 and DNS server
  address 2.

* `ntp_server_address` - (Optional, String) Specifies the NTP server address. Currently only IPv4 addresses are supported.
  A maximum of four IP addresses can be configured, and each address must be unique. Multiple IP addresses must be
  separated using commas(,). Removing this parameter indicates that no NTP server is configured.

* `dhcp_lease_time` - (Optional, String) Specifies the DHCP lease expiration time. The value can be -1, which indicates
  unlimited lease time, or Number+h. the number ranges from 1 to 30,000. For example, the value can be 5h. The default
  value is 24h.

* `dhcp_ipv6_lease_time` - (Optional, String) Specifies the DHCP lease expiration time of the IPv6 subnet. The value can
  be -1, which indicates unlimited lease time, or Number+h. the number ranges from 1 to 175200. For example, the value
  can be 5h. The default value is 2h.

* `dhcp_domain_name` - (Optional, String) Specifies the domain name configured for DNS and is used to obtain the IP address
  from the DNS server. A domain name can contain only letters, digits, and hyphens (-) and cannot start or end with a
  hyphen (-). Each domain name contains at least two labels separated by periods (.). Max total: 254 characters. Max
  label: 63 characters.

* `tags` - (Optional, Map) The key/value pairs to associate with the subnet.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `status` - The status of the subnet. The value can be ACTIVE, DOWN, UNKNOWN, or ERROR.

* `ipv4_subnet_id` - The ID of the IPv4 subnet (Native OpenStack API).

* `ipv6_subnet_id` - The ID of the IPv6 subnet (Native OpenStack API).

* `ipv6_cidr` - The IPv6 subnet CIDR block.

* `ipv6_gateway` - The IPv6 subnet gateway.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minutes.
* `delete` - Default is 10 minutes.

## Import

Subnets can be imported using the subnet `Network ID`, e.g.

```bash
$ terraform import huaweicloud_vpc_subnet.subnet 4779ab1c-7c1a-44b1-a02e-93dfc361b32d
```
