---
subcategory: "Virtual Private Network (VPN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpn_user_groups"
description: |-
  Use this data source to get the list of VPN user groups.
---

# huaweicloud_vpn_user_groups

Use this data source to get the list of VPN user groups.

## Example Usage

```hcl
variable "vpn_server_id" {}

data "huaweicloud_vpn_user_groups" "test" {
  vpn_server_id = var.vpn_server_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `vpn_server_id` - (Required, String) Specifies the ID of a VPN server.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `user_groups` - The user group list information.

  The [user_groups](#user_groups_struct) structure is documented below.

<a name="user_groups_struct"></a>
The `user_groups` block supports:

* `id` - The user group ID.

* `name` - The user group name.

* `description` - The user group description.

* `type` - The user group type.

* `user_number` - The number of users.

* `created_at` - The creation time.

* `updated_at` - The update time.
