---
subcategory: "Virtual Private Cloud (VPC)"
---

# huaweicloud\_vpc\_subnet

Provides details about a specific VPC subnet.
This is an alternative to `huaweicloud_vpc_subnet_v1`

This resource can prove useful when a module accepts a subnet id as
an input variable and needs to, for example, determine the id of the
VPC that the subnet belongs to.

## Example Usage

```hcl
data "huaweicloud_vpc_subnet" "subnet" {
  id = var.subnet_id
}

output "subnet_vpc_id" {
  value = data.huaweicloud_vpc_subnet.subnet.vpc_id
}
```

## Argument Reference

The arguments of this data source act as filters for querying the available
subnets in the current tenant. The given filters must match exactly one
subnet whose data will be exported as attributes.

* `region` - (Optional, String) The region in which to obtain the subnet. If omitted, the provider-level region will be used.

* `id` - (Optional, String) - Specifies a resource ID in UUID format.

* `name` - (Optional, String) The name of the specific subnet to retrieve.

* `cidr` - (Optional, String) The network segment of specific subnet to retrieve. The value must be in CIDR format.

* `status` - (Optional, String) The value can be ACTIVE, DOWN, UNKNOWN, or ERROR.

* `vpc_id` - (Optional, String) The id of the VPC that the desired subnet belongs to.

* `gateway_ip` - (Optional, String) The subnet gateway address of specific subnet.

* `primary_dns` - (Optional, String) The IP address of DNS server 1 on the specific subnet.

* `secondary_dns` - (Optional, String) The IP address of DNS server 2 on the specific subnet.

* `availability_zone` - (Optional, String) The availability zone (AZ) to which the subnet should belong.

## **Attributes Reference**

In addition to all arguments above, the following attributes are exported:

* `dns_list` - The IP address list of DNS servers on the subnet.

* `dhcp_enable` - DHCP function for the subnet.

* `subnet_id` - Specifies the subnet (Native OpenStack API) ID.
