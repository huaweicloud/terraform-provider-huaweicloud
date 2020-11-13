---
subcategory: "Virtual Private Cloud (VPC)"
---

# huaweicloud\_vpc\_subnet

Provides an VPC subnet resource.
This is an alternative to `huaweicloud_vpc_subnet_v1`

# Example Usage

```hcl
resource "huaweicloud_vpc" "vpc" {
  name = var.vpc_name
  cidr = var.vpc_cidr
}

resource "huaweicloud_vpc_subnet" "subnet" {
  name       = var.subnet_name
  cidr       = var.subnet_cidr
  gateway_ip = var.subnet_gateway_ip
  vpc_id     = huaweicloud_vpc.vpc.id
}

resource "huaweicloud_vpc_subnet" "subnet_with_tags" {
  name       = var.subnet_name
  cidr       = var.subnet_cidr
  gateway_ip = var.subnet_gateway_ip
  vpc_id     = huaweicloud_vpc.vpc.id

  tags = {
    foo = "bar"
    key = "value"
  }
}

 ```

# Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the vpc subnet. If omitted, the provider-level region will work as default. hanging this creates a new Subnet resource.

* `name` (Required) - The subnet name. The value is a string of 1 to 64 characters that can contain letters, digits, underscores (_), and hyphens (-).

* `cidr` (Required) - Specifies the network segment on which the subnet resides. The value must be in CIDR format. The value must be within the CIDR block of the VPC. The subnet mask cannot be greater than 28. Changing this creates a new Subnet.

* `gateway_ip` (Required) - Specifies the gateway of the subnet. The value must be a valid IP address. The value must be an IP address in the subnet segment. Changing this creates a new Subnet.

* `vpc_id` (Required) - Specifies the ID of the VPC to which the subnet belongs. Changing this creates a new Subnet.

* `dhcp_enable` (Optional) - Specifies whether the DHCP function is enabled for the subnet. The value can be true or false. If this parameter is left blank, it is set to true by default.

* `primary_dns` (Optional) - Specifies the IP address of DNS server 1 on the subnet. The value must be a valid IP address.

* `secondary_dns` (Optional) - Specifies the IP address of DNS server 2 on the subnet. The value must be a valid IP address.

* `dns_list` (Optional) - Specifies the DNS server address list of a subnet. This field is required if you need to use more than two DNS servers. This parameter value is the superset of both DNS server address 1 and DNS server address 2.

* `availability_zone` (Optional) - Identifies the availability zone (AZ) to which the subnet belongs. The value must be an existing AZ in the system. Changing this creates a new Subnet.

* `tags` - (Optional) The key/value pairs to associate with the subnet.

# Attributes Reference

All of the argument attributes are also exported as
result attributes:

* `id` - Specifies a resource ID in UUID format.
 
* `status` - Specifies the status of the subnet. The value can be ACTIVE, DOWN, UNKNOWN, or ERROR.

* `subnet_id` - Specifies the subnet (Native OpenStack API) ID.

# Import

Subnets can be imported using the `subnet id`, e.g.

```
$ terraform import huaweicloud_vpc_subnet 4779ab1c-7c1a-44b1-a02e-93dfc361b32d
```

## Timeouts
This resource provides the following timeouts configuration options:
- `create` - Default is 10 minute.
- `delete` - Default is 10 minute.

