---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_user"
description: ""
---

# huaweicloud_workspace_user

Manages a Workspace user resource within HuaweiCloud.

## Example Usage

### Create a user that never expires

```hcl
variable "user_name" {}
variable "email_address" {}

resource "huaweicloud_workspace_user" "test" {
  name  = var.user_name
  email = var.email_address

  account_expires            = "0"
  password_never_expires     = false
  enable_change_password     = true
  next_login_change_password = true
  disabled                   = false
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the Workspace user resource.
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the user name.
  + Pure numeric: the valid length is between `1` and `20`.
  + Non-pure numeric: the name can contain `1` to `20` characters, only letters, digits, hyphens (-), underscore (_) and
  dots (.) are allowed. The name must start with a letter.

  Changing this will create a new resource.

* `active_type` - (Optional, String) Specifies the activation mode of the user. Defaults to **USER_ACTIVATE**.  
  The valid values are as follows:
  + **USER_ACTIVATE**: Activated by the user.
  + **ADMIN_ACTIVATE**: Activated by the administator.

* `email` - (Optional, String) Specifies the email address of user. The value can contain `1` to `64` characters.

* `phone` - (Optional, String) Specifies the phone number of user.

-> At least one of `email` and `phone` parameters must be provided.

* `password` - (Optional, String) Specifies the initial passowrd of user.  
  This parameter is available and required when `active_type` is set to **ADMIN_ACTIVATE**.

* `description` - (Optional, String) Specifies the description of user. The maximum length is `255` characters.

* `account_expires` - (Optional, String) Specifies the user's valid period configuration.
  Defaults to "0".
  + Never expires: **0**.
  + Expires at a certain time: account expires must in RFC3339 format like `yyyy-MM-ddTHH:mm:ssZ`.
    The times is in local time, depending on the timezone.

  -> Only support the hours timezones, e.g. **+04:00 Baku, Tbilisi, Yerevan** or **+05:00 Ekaterinburg** is supported,
     but **+04:30 Kabul** is not supported.

* `password_never_expires` - (Optional, Bool) Specifies whether the password will never expires.
  Defaults to **false**.

* `enable_change_password` - (Optional, Bool) Specifies whether to allow password modification.
  Defaults to **true**.

* `next_login_change_password` - (Optional, Bool) Specifies whether the next login requires a password reset.
  Defaults to **true**.

* `disabled` - (Optional, Bool) Specifies whether the user is disabled.
  Defaults to **false**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The user ID in UUID format.

* `locked` - Whether the user is locked.

* `total_desktops` - The number of desktops the user has.

## Import

Users can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_workspace_user.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason.
The missing attributes include: `password`.
It is generally recommended running `terraform plan` after importing the resource.
You can then decide if changes should be applied to the user, or the resource definition should be updated to
align with the user. Also you can ignore changes as below.

```hcl
resource "huaweicloud_workspace_user" "test" {
  ...

  lifecycle {
    ignore_changes = [
      password,
    ]
  }
}
```
