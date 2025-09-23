---
subcategory: "Virtual Private Network (VPN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpn_users"
description: |-
  Use this data source to get the list of VPN users.
---

# huaweicloud_vpn_users

Use this data source to get the list of VPN users.

## Example Usage

```hcl
variable "vpn_server_id" {}

data "huaweicloud_vpn_users" "test" {
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

* `users` - The user list.

  The [users](#users_struct) structure is documented below.

<a name="users_struct"></a>
The `users` block supports:

* `id` - The user ID.

* `name` - The username.

* `description` - The user description.

* `user_group_id` - The ID of the user group to which a user belongs.

* `user_group_name` - The name of the user group to which a user belongs.

* `created_at` - The creation time.

* `updated_at` - The update time.
