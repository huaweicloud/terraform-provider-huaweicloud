---
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpc_subnet_v1"
sidebar_current: "docs-huaweicloud-resource-vpc-subnet-v1"
description: |-
  Provides an VPC subnet resource.
---

# huaweicloud_vpc_subnet_v1

Provides an VPC subnet resource.

# Example Usage

 ```hcl
resource "huaweicloud_vpc_v1" "vpc_v1" {
  name = "${var.vpc_name}"
  cidr = "${var.vpc_cidr}"
}


resource "huaweicloud_vpc_subnet_v1" "subnet_v1" {
  name = "${var.subnet_name}"
  cidr = "${var.subnet_cidr}"
  gateway_ip = "${var.subnet_gateway_ip}"
  vpc_id = "${huaweicloud_vpc_v1.vpc_v1.id}"
}
 ```

# Argument Reference

The following arguments are supported:

* `name` (Required) - The subnet name. The value is a string of 1 to 64 characters that can contain letters, digits, underscores (_), and hyphens (-).

* `cidr` (Required) - Specifies the network segment on which the subnet resides. The value must be in CIDR format. The value must be within the CIDR block of the VPC. The subnet mask cannot be greater than 28. Changing this creates a new Subnet.

* `gateway_ip` (Required) - Specifies the gateway of the subnet. The value must be a valid IP address. The value must be an IP address in the subnet segment. Changing this creates a new Subnet.

* `vpc_id` (Required) - Specifies the ID of the VPC to which the subnet belongs. Changing this creates a new Subnet.

* `dhcp_enable` (Optional) - Specifies whether the DHCP function is enabled for the subnet. The value can be true or false. If this parameter is left blank, it is set to true by default.

* `primary_dns` (Optional) - Specifies the IP address of DNS server 1 on the subnet. The value must be a valid IP address.

* `secondary_dns` (Optional) - Specifies the IP address of DNS server 2 on the subnet. The value must be a valid IP address.

* `dns_list` (Optional) - Specifies the DNS server address list of a subnet. This field is required if you need to use more than two DNS servers. This parameter value is the superset of both DNS server address 1 and DNS server address 2.

* `availability_zone` (Optional) - Identifies the availability zone (AZ) to which the subnet belongs. The value must be an existing AZ in the system. Changing this creates a new Subnet.


# Attributes Reference

All of the argument attributes are also exported as
result attributes:

* `id` - The ID of the subnet.
 
* `status` - Specifies the status of the subnet. The value can be ACTIVE, DOWN, UNKNOWN, or ERROR.

# Import

Subnets can be imported using the `subnet id`, e.g.

> $ terraform import huaweicloud_vpc_subnet_v1 4779ab1c-7c1a-44b1-a02e-93dfc361b32d