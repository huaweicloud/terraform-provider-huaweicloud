---
subcategory: "Virtual Private Network (VPN)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vpn_user"
description: |-
  Manages a VPN user resource within HuaweiCloud.
---

# huaweicloud_vpn_user

Manages a VPN user resource within HuaweiCloud.

## Example Usage

```hcl
variable "vpn_server_id" {}
variable "name" {}
variable "password" {}

resource "huaweicloud_vpn_user" "test" {
  vpn_server_id = var.vpn_server_id
  name          = var.name
  password      = var.password
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `vpn_server_id` - (Required, String, NonUpdatable) Specifies the VPN server ID.

* `name` - (Required, String, NonUpdatable) Specifies the user name.

* `password` - (Required, String) Specifies the user password.

* `description` - (Optional, String) Specifies the description of the user.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `user_group_id` - The user group ID.

* `user_group_name` - The user group name.

* `created_at` - The creation time.

* `updated_at` - The update time.

## Import

The user can be imported using `vpn_server_id` and `id`, separated by a slash (/), e.g.

```bash
$ terraform import huaweicloud_vpn_user.test <vpn_server_id>/<id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attribute is `password`. It is generally recommended running
`terraform plan` after importing the resource. You can then decide if changes should be applied to the user, or the
resource definition should be updated to align with the user. Also you can ignore changes as below.

```hcl
resource "huaweicloud_vpn_user" "test" {
    ...

  lifecycle {
    ignore_changes = [
      password
    ]
  }
}
```
