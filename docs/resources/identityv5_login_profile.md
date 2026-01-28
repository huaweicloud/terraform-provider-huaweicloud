---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identityv5_login_profile"
description: |-
  Manages an IAM login profile resource within HuaweiCloud.
---

# huaweicloud_identityv5_login_profile

Manages an IAM login profile resource within HuaweiCloud.

-> You must have admin privileges to use this resource.

## Example Usage

```hcl
variable "user_id" {}

resource "huaweicloud_identityv5_login_profile" "test" {
  user_id                 = var.user_id
  password_reset_required = true
}
```

## Argument Reference

The following arguments are supported:

* `user_id` - (Required, String, NonUpdatable) Specifies the ID of the user.

* `password_reset_required` - (Optional, Bool) Specifies whether the user needs to reset the password at the
  next login.  
  Default: **false**.
  
* `password` - (Optional, String) Specifies the password of the user login.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, also the user ID.

* `password_expires_at` - The password expiration time of the user.

* `created_at` - The creation time of the login profile.

## Import

The login profile can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_identityv5_login_profile.test <id>
```
