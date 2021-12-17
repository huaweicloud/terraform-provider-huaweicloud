---
subcategory: "Virtual Private Cloud (VPC)"
---

# huaweicloud_vpc_address_group

Provide a resource to manage the IP address group of VPC

## Example Usage

### Basic IP Address Group

```hcl
resource "huaweicloud_vpc_address_group" "test" {
	dry_run = false
	name = "test1"
	ip_version = 4
	description  =  "vpc test"
	ip_set	=	[
		"192.168.5.0/24",
		"192.168.9.0/24"
	]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) The region in which to create the IP Address Group. If omitted, the
  provider-level region will be used. Changing this creates a new IP Address Group.
  
* `name` - (Required, String) The name of the IP Address Group.

* `dry_run` - (Optional, Bool) The function is to pre check only this request. If the value is true: Send a check request, the address group will not be created. The check items include whether the required parameters, request format and business restrictions are filled in. If the check fails, the corresponding error is returned. If the check passes, the response code 202 is returned. If the value is false (default): Send a normal request and directly create an address group.
  
* `ip_version` - (Required，Int) This parameter is the IP version of the address group. Value range: 4, indicating the IPv4 address group; 6. Indicates IPv6 address group
  
* `description` - (Optional, String) This parameter is the address group description information. Value range: 0-255 characters, cannot contain "<" and ">".
  
* `ip_set` - (Optional, List) This parameter is an address group and can contain an address set; Value range: can be a single IP address, IP address range, IP address CIDR. Constraint: current IP address group_ The default value of set quantity limit is 20, that is, the default limit of the total number of configured IP addresses, IP address ranges or IP address cidrs is 20

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

*  `id` - ID of the IP Address Group.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minute.
* `update` - Default is 10 minute.
* `delete` - Default is 10 minute.