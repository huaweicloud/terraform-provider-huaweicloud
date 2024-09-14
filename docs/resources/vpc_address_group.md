---
subcategory: "Virtual Private Cloud (VPC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpc_address_group"
description: ""
---

# huaweicloud_vpc_address_group

Manages a VPC IP address group resource within HuaweiCloud.

## Example Usage

### IPv4 Address Group

```hcl
resource "huaweicloud_vpc_address_group" "ipv4" {
  name = "group-ipv4"

  addresses = [
    "192.168.10.10",
    "192.168.1.1-192.168.1.50"
  ]
}
```

### IPv6 Address Group

```hcl
resource "huaweicloud_vpc_address_group" "ipv6" {
  name       = "group-ipv6"
  ip_version = 6

  addresses = [
    "2001:db8:a583:6e::/64"
  ]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the IP address group. If omitted, the
  provider-level region will be used. Changing this creates a new address group.
  
* `name` - (Required, String) Specifies the IP address group name.  
  The value is a string of `1` to `64` characters that can contain letters, digits, underscores (_), hyphens (-) and
  periods (.).

* `addresses` - (Required, List) Specifies an array of one or more IP addresses. The address can be a single IP
  address, IP address range or IP address CIDR. The maximum length is 20.

* `ip_version` - (Optional, Int, ForceNew) Specifies the IP version, either `4` (default) or `6`.
  Changing this creates a new address group.

* `description` - (Optional, String) Specifies the supplementary information about the IP address group.
  The value is a string of no more than `255` characters and cannot contain angle brackets (< or >).

* `max_capacity` - (Optional, Int) Specifies the maximum number of addresses that an address group can contain.
  The valid value is range from `1` to `20`, the default value is `20`.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project ID.
  Changing this creates a new address group.

* `force_destroy` - (Optional, Bool) Specifies whether to forcibly destroy the address group if it is associated with
  a security group rule, the address group and the associated security group rule will be deleted together.
  The default value is **false**.
  
## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minutes.
* `update` - Default is 5 minutes.
* `delete` - Default is 5 minutes.

## Import

IP address groups can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_vpc_address_group.test bc96f6b0-ca2c-42ee-b719-0f26bc9c8661
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response. The missing attributes include: `force_destroy`. It is generally recommended running `terraform plan` after
importing the image. You can then decide if changes should be applied to the image, or the resource
definition should be updated to align with the image. Also you can ignore changes as below.

```hcl
resource "huaweicloud_vpc_address_group" "test" {
  ...

  lifecycle {
    ignore_changes = [
      force_destroy,
    ]
  }
}
```
