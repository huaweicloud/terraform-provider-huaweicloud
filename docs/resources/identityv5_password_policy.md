---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identityv5_password_policy"
description: |-
  Manages the account password policy v5 resource within HuaweiCloud.
---

# huaweicloud_identityv5_password_policy

Manages the account password policy v5 resource within HuaweiCloud.

-> **NOTE:**
  You *must* have admin privileges to use this resource.  
  This resource overwrites an existing configuration, make sure one resource per account.  
  During action `terraform destroy` it sets values the same as defaults for this resource.

## Example Usage

```hcl
resource "huaweicloud_identityv5_password_policy" "test" {
  maximum_consecutive_identical_chars = 1
  minimum_password_age                = 30
  minimum_password_length             = 10
  password_reuse_prevention           = 2
  password_not_username_or_invert     = false
  password_validity_period            = 90
  password_char_combination           = 3
  allow_user_to_change_password       = false
}
```

## Argument Reference

The following arguments are supported:

* `maximum_consecutive_identical_chars` - (Optional, Int) Specifies the maximum number of times that a character is
  allowed to consecutively present in a password. The value ranges from `0` to `32` and defaults to `0` which
  indicates that consecutive identical characters are allowed in a password. For example, value `2` indicates that
  two consecutive identical characters are not allowed in a password.

* `minimum_password_age` - (Optional, Int) Specifies the minimum period (minutes) after which users are allowed to make
  a password change. The value ranges from `0` to `1,440` and defaults to `0`.

* `minimum_password_length` - (Optional, Int) Specifies the minimum number of characters that a password must contain.
  The value ranges from `8` to `32` and defaults to `8`.

* `password_reuse_prevention` - (Optional, Int) Specifies the password cannot be repeated with historical passwords
  for a certain number of times. The value ranges from `0` to `24` and defaults to `1`. For example, value `3`
  indicates that the user cannot set the last three passwords that the user has previously used when setting a new
  password.

* `password_not_username_or_invert` - (Optional, Bool) Specifies whether the password can be the username or the
  username spelled backwards. Defaults to `true`, which indicates that the username or the inversion of username
  is not allowed to be used as a password.

* `password_validity_period` - (Optional, Int) Specifies the password validity period (days).
  The value ranges from `0` to `180` and defaults to `0` which indicates that this requirement does not apply.

* `password_char_combination` - (Optional, Int) Specifies the minimum number of character types that a password must
  contain. The value ranges from `2` to `4` and defaults to `2` which indicates that a password must contain at least
  two of the following: uppercase letters, lowercase letters, digits, and special characters.

* `allow_user_to_change_password` - (Optional, Bool) Specifies whether IAM users are allowed to change their own
  passwords, which does not apply to the root user. Defaults to `true`, which indicates IAM users are allowed to
  change their own passwords.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates the ID of account password policy, which is the same as the account ID.

* `maximum_password_length` - Indicates the maximum number of characters that a password can contain.

* `password_requirements` - Indicates the requirements of character that passwords must include.

## Import

The IAM v5 password policy can be imported using the account ID or domain ID, e.g.

```bash
$ terraform import huaweicloud_identityv5_password_policy.test <id>
```
