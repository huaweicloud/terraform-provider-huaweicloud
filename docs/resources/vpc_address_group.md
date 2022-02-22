---
subcategory: "Virtual Private Cloud (VPC)"
---

# huaweicloud_vpc_address_group

Manages a VPC IP address group resource within HuaweiCloud.

## Example Usage

```hcl
resource "huaweicloud_vpc_address_group" "test" {
  name      = "group-test"
  addresses = [
    "192.168.10.10",
    "192.168.1.1-192.168.1.50"
  ]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies tThe region in which to create the IP address group. If omitted, the
  provider-level region will be used. Changing this creates a new address group.
  
* `name` - (Required, String) Specifies the IP address group name. The value is a string of 1 to 64 characters that can contain
  letters, digits, underscores (_), hyphens (-) and periods (.).

* `addresses` - (Required, List) Specifies an array of one or more IPv4 addresses. The address can be a single IP
  address (such as 192.168.10.10), IP address range (such as 192.168.1.1-192.168.1.50) or IP address CIDR (such as 192.168.0.0/16).
  The maximum length is 20.
  
* `description` - (Optional, String) Specifies the supplementary information about the IP address group.
  The value is a string of no more than 255 characters and cannot contain angle brackets (< or >).
  
## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.
* `ip_version` - The IP version of the address group. The value is 4.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minute.
* `update` - Default is 5 minute.
* `delete` - Default is 5 minute.

## Import

IP address groups can be imported using the `id`, e.g.

```
$ terraform import huaweicloud_vpc_address_group.test bc96f6b0-ca2c-42ee-b719-0f26bc9c8661
```
