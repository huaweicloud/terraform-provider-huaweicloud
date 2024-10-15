---
subcategory: "Virtual Private Network (VPN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpn_user_group"
description: |-
  Manages a VPN user group resource within HuaweiCloud.
---

# huaweicloud_vpn_user_group

Manages a VPN user group resource within HuaweiCloud.

## Example Usage

### Create a basic VPN user group

```hcl
variable "vpn_server_id" {}
variable "name" {}
variable "description" {}

resource "huaweicloud_vpn_user_group" "test" {
  vpn_server_id = var.vpn_server_id
  name          = var.name
  description   = var.description
}
```

### Create a VPN user group with users

```hcl
variable "vpn_server_id" {}
variable "name" {}
variable "description" {}
variable "user1_id" {}
variable "user2_id" {}

resource "huaweicloud_vpn_user_group" "test" {
  vpn_server_id = var.vpn_server_id
  name          = var.name
  description   = var.description

  users {
    id = var.user1_id
  }

  users {
    id = var.user2_id
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `vpn_server_id` - (Required, String, NonUpdatable) Specifies the VPN server ID.

* `name` - (Required, String) Specifies the name of the user group.

* `description` - (Optional, String) Specifies the description of the user group.

* `users` - (Optional, List) Specifies the user list.
  The [users](#Users) structure is documented below.

<a name="Users"></a>
The `users` block supports:

* `id` - (Required, String) Specifies the user ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `type` - The user group type.

* `created_at` - The creation time.

* `updated_at` - The update time.

* `user_number` - The number of users.

* `users` - The user list.
  The [users](#UsersAttr) structure is documented below.

<a name="UsersAttr"></a>
The `users` block supports:

* `name` - The username.

* `description` - The user description.

## Import

The user group can be imported using `vpn_server_id` and `id`, separated by a slash (/), e.g.

```bash
$ terraform import huaweicloud_vpn_user_group.test <vpn_server_id>/<id>
```
