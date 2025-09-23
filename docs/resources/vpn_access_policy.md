---
subcategory: "Virtual Private Network (VPN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpn_access_policy"
description: |-
  Manages a VPN access policy within HuaweiCloud.
---

# huaweicloud_vpn_access_policy

Manages a VPN access policy within HuaweiCloud.

## Example Usage

```hcl
variable "vpn_server_id" {}
variable "name" {}
variable "user_group_id" {}
variable "dest_ip_cidr" {}

resource "huaweicloud_vpn_access_policy" "test" {
  vpn_server_id = var.vpn_server_id
  name          = var.name
  user_group_id = var.user_group_id
  dest_ip_cidrs = [var.dest_ip_cidr]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `vpn_server_id` - (Required, String, NonUpdatable) Specifies the VPN server ID.

* `name` - (Required, String) Specifies the access policy name.

* `user_group_id` - (Required, String) Specifies the user group ID.

* `dest_ip_cidrs` - (Required, List) Specifies the list of destination IP CIDRs.

* `description` - (Optional, String) Specifies the description of the access policy.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `user_group_name` - The user group name.

* `created_at` - The creation time.

* `updated_at` - The update time.

## Import

The access policy can be imported using `vpn_server_id` and `id`, separated by a slash (/), e.g.

```bash
$ terraform import huaweicloud_vpn_access_policy.test <vpn_server_id>/<id>
```
