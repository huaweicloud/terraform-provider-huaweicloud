---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identityv5_password_policy"
description: |-
  Manages the account password policy resource within HuaweiCloud.
---

# huaweicloud_identityv5_password_policy

Manages the account password policy resource within HuaweiCloud.

-> 1. You **must** have admin privileges to use this resource.
  <br>2. This resource overwrites an existing configuration, make sure one resource per account.  
  <br>3. During action `terraform destroy` it sets values the same as defaults for this resource.

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
  allowed to consecutively present in a password.  
  The valid value ranges from `0` to `32`.  
  Defaults to `0`, it means consecutive identical characters are allowed.

* `minimum_password_age` - (Optional, Int) Specifies the minimum password usage time, in minutes.  
  The valid value ranges from `0` to `1,440`, defaults to `0`.

* `minimum_password_length` - (Optional, Int) Specifies the minimum number of characters that a password
  must contain.  
  The valid value ranges from `8` to `32`, defaults to `8`.

* `password_reuse_prevention` - (Optional, Int) Specifies the password cannot be repeated with historical passwords
  for a certain number of times.  
  The valid value ranges from `0` to `24`.  
  Defaults to `1`, it means the user cannot set the last one password that the user has previously used when setting
  a new password.

* `password_not_username_or_invert` - (Optional, Bool) Specifies whether the password can be the username or the
  username spelled backwards.  
  Defaults to `true`, it means the username or the inversion of username is not allowed to be used as a password.

* `password_validity_period` - (Optional, Int) Specifies the password validity period, in days.
  The valid value ranges from `0` to `180`.  
  Defaults to `0`, it means this requirement does not apply.

* `password_char_combination` - (Optional, Int) Specifies the minimum number of character types that a password must
  contain.  
  The valid value ranges from `2` to `4`.  
  Defaults to `2`, it means a password must contain at least two of the following: uppercase letters, lowercase
  letters, digits, and special characters.

* `allow_user_to_change_password` - (Optional, Bool) Specifies whether IAM users are allowed to change their own
  passwords.  
  Defaults to **true**.  
  This parameter does not apply to the root user.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, also the domain (account) ID.

* `maximum_password_length` - The maximum number of characters that a password can contain.

* `password_requirements` - The requirements of character that passwords must include.

## Import

The IAM v5 password policy can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_identityv5_password_policy.test <id>
```
