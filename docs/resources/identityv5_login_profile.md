---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identityv5_login_profile"
description: |-
  Manages an IAM v5 login profile resource within HuaweiCloud.
---

# huaweicloud_identityv5_login_profile

Manages an IAM v5 login profile resource within HuaweiCloud.

-> NOTE: You must have admin privileges to use this resource.

## Example Usage

```hcl
variable "name" {}

resource "huaweicloud_identityv5_user" "user" {
  name        = var.name
  enabled     = true
  description = "tested by terraform"
}

resource "huaweicloud_identityv5_login_profile" "login_profile" {
  user_id                 = huaweicloud_identityv5_user.user.id
  password                = "default888"
  password_reset_required = true
}
```

## Argument Reference

The following arguments are supported:

* `user_id` - (Required, String， NonUpdatable) Specifies the ID of the user.

* `password` - (Optional, String) Specifies the password of the user login.

* `password_reset_required` - (Optional, Bool) Specifies whether the IAM user need to change their password when next
  time they log in. The value can be: **true**, **false**. Default: **true**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates the ID of the user.

* `password_expires_at` - Indicates IAM user password expiration time.

* `created_at` - Indicates the creation time of the IAM user login information.

## Import

The IAM v5 login profile can be imported using the id, e.g:

```bash
$ terraform import huaweicloud_identityv5_login_profile.login_profile <user_id>
```
