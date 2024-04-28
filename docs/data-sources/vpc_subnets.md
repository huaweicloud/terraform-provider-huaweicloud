---
subcategory: "Virtual Private Cloud (VPC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpc_subnets"
description: ""
---

# huaweicloud_vpc_subnets

Use this data source to get a list of VPC subnet.

## Example Usage

An example filter by name and tag

```hcl
data "huaweicloud_vpc_subnets" "subnet" {
  name = var.subnet_name

  tags = {
    foo = "bar"
  }
}

output "subnet_vpc_ids" {
  value = data.huaweicloud_vpc_subnets.subnet.subnets[*].vpc_id
}
```

## Argument Reference

The arguments of this data source act as filters for querying the available subnets in the current tenant.
All subnets that meet the filter criteria will be exported as attributes.

* `region` - (Optional, String) Specifies the region in which to obtain the subnet. If omitted, the provider-level
  region will be used.

* `id` - (Optional, String) - Specifies the id of the desired subnet.

* `name` - (Optional, String) Specifies the name of the desired subnet.

* `cidr` - (Optional, String) Specifies the network segment of desired subnet. The value must be in CIDR format.

* `status` - (Optional, String) Specifies the current status of the desired subnet.
  the value can be ACTIVE, DOWN, UNKNOWN, or ERROR.

* `vpc_id` - (Optional, String) Specifies the id of the VPC that the desired subnet belongs to.

* `gateway_ip` - (Optional, String) Specifies the subnet gateway address of desired subnet.

* `primary_dns` - (Optional, String) Specifies the IP address of DNS server 1 on the desired subnet.

* `secondary_dns` - (Optional, String) Specifies the IP address of DNS server 2 on the desired subnet.

* `availability_zone` - (Optional, String) Specifies the availability zone (AZ) to which the desired subnet belongs to.

* `tags` - (Optional, Map) Specifies the included key/value pairs which associated with the desired subnet.

 -> A maximum of 10 tag keys are allowed for each query operation. Each tag key can have up to 10 tag values.
  The tag key cannot be left blank or set to an empty string. Each tag key must be unique, and each tag value in a
  tag must be unique, use commas(,) to separate the multiple values. An empty for values indicates any value.
  The values are in the OR relationship.

## **Attributes Reference**

The following attributes are exported:

* `id` - Indicates a data source ID.
* `subnets` - Indicates a list of all subnets found. Structure is documented below.

The `subnets` block supports:

* `id` - Indicates the ID of the subnet.
* `name` - Indicates the name of the subnet.
* `description` - Indicates the description of the subnet.
* `cidr` - Indicates the cidr block of the subnet.
* `status` - Indicates the current status of the subnet.
* `vpc_id` - Indicates the Id of the VPC that the subnet belongs to.
* `gateway_ip` - Indicates the subnet gateway address of the subnet.
* `primary_dns` - Indicates the IP address of DNS server 1 on the subnet.
* `secondary_dns` - Indicates the IP address of DNS server 2 on the subnet.
* `availability_zone` - Indicates the availability zone (AZ) to which the subnet belongs to.
* `dhcp_enable` - Indicates whether the DHCP is enabled.
* `dns_list` - Indicates The IP address list of DNS servers on the subnet.
* `ipv4_subnet_id` - Indicates the ID of the IPv4 subnet (Native OpenStack API).
* `ipv6_enable` - Indicates whether the IPv6 is enabled.
* `ipv6_subnet_id` - Indicates the ID of the IPv6 subnet (Native OpenStack API).
* `ipv6_cidr` - Indicates the IPv6 subnet CIDR block.
* `ipv6_gateway` - Indicates the IPv6 subnet gateway.
* `tags` - Indicates the key/value pairs which associated with the subnet.
