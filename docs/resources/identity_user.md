---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_user"
description: ""
---

# huaweicloud_identity_user

Manages an IAM user resource within HuaweiCloud.

-> **NOTE:** You *must* have admin privileges to use this resource.

## Example Usage

```hcl
variable "user_1_password" {}

resource "huaweicloud_identity_user" "user_1" {
  name        = "user_1"
  description = "A user"
  password    = var.user_1_password
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Specifies the name of the user. The user name consists of `1` to `32` characters. It can
  contain only uppercase letters, lowercase letters, digits, spaces, and special characters (-_) and cannot start with a
  digit.

* `description` - (Optional, String) Specifies the description of the user.

* `email` - (Optional, String) Specifies the email address with a maximum of `255` characters.

* `phone` - (Optional, String) Specifies the mobile number with a maximum of `32` digits. This parameter must be used
  together with `country_code`.

* `country_code` - (Optional, String) Specifies the country code. The country code of the Chinese mainland is 0086. This
  parameter must be used together with `phone`.

* `password` - (Optional, String) Specifies the password for the user with `6` to `32` characters. It must contain at least
  two of the following character types: uppercase letters, lowercase letters, digits, and special characters.

* `pwd_reset` - (Optional, Bool) Specifies whether or not the password should be reset. By default, the password is asked
   to reset at the first login.

* `enabled` - (Optional, Bool) Specifies whether the user is enabled or disabled. Valid values are **true** and **false**.

* `access_type` - (Optional, String) Specifies the access type of the user. Available values are:
  + **default**: support both programmatic and management console access.
  + **programmatic**: only support programmatic access.
  + **console**: only support management console access.

* `external_identity_id` - (Optional, String) Specifies the ID of the IAM user in the external system.
  This parameter is used for IAM user SSO type, make sure that the **IAM_SAML_Attributes_xUserId** of the federated user
  is the same as the `external_identity_id` of the corresponding IAM user.

* `external_identity_type` - (Optional, String) Specifies the type of the IAM user in the external system.
  Only **TenantIdp** is supported now. This parameter must be used together with `external_identity_id`.

* `login_protect_verification_method` - (Optional, String) Specifies the verification method of login protect. If it is
  empty, the login protection will be disabled.
  
  Valid values are as follows:
  + **sms**: Use phone number to verify.
  + **email**: Use email to verify.
  + **vmfa**: Use MFA to verify.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `password_strength` - Indicates the password strength.

* `create_time` - The time when the IAM user was created.

* `last_login` - The time when the IAM user last login.

## Import

Users can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_identity_user.user_1 89c60255-9bd6-460c-822a-e2b959ede9d2
```

But due to the security reason, `password` can not be imported, you can ignore it as below.

```hcl
resource "huaweicloud_identity_user" "user_1" {
  ...

  lifecycle {
    ignore_changes = [
      password,
    ]
  }
}
```
