---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identityv5_user"
description: |-
  Manages an IAM v5 user resource within HuaweiCloud.
---

# huaweicloud_identityv5_user

Manages an IAM v5 user resource within HuaweiCloud.

-> NOTE: You must have admin privileges to use this resource.

## Example Usage

```hcl
variable "name" {}

resource "huaweicloud_identityv5_user" "user" {
  name        = var.name
  enabled     = true
  description = "tested by terraform"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Specifies the name of the user. The username consists of `1` to `64` characters. It can
  contain only uppercase letters, lowercase letters, digits, spaces, and special characters (-_) and cannot start with a
  digit.

* `description` - (Optional, String) Specifies the description of the user.

* `enabled` - (Optional, Bool) Specifies whether the user is enabled or disabled. Valid values are **true** and **false**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `name` - Indicates the IAM username.

* `is_root_user` - Indicates whether the user is root user or user.

* `created_at` - Indicates the time when the IAM user was created.

* `urn` - Indicates the uniform resource name.

* `last_login_at` - Indicates the last login time of the IAM user. If null, it indicates that they have never logged in.

* `tags` - Indicates the Custom Tag List.
  The [tags](#IdentityV5User_Tags) structure is documented below.

<a name="IdentityV5User_Tags"></a>
The `tags` block contains:

* `tag_key` - Indicates Tag keys which contain any combination of letters, numbers, spaces, and the symbols "_", ".",
  ":","=", "-", "@", but cannot start or end with a space, and cannot begin with "sys". The length must be between 1
  and 64 characters.

* `tag_value` - Indicates Tag keys can contain any combination of letters, numbers, spaces, and the symbols "_", ".",
  ":", "=", "-", "@", but cannot start or end with a space, and cannot begin with "sys". The length must be between 1
  and 64 characters.

## Import

The IAM v5 user can be imported using the id, e.g:

```bash
$ terraform import huaweicloud_identityv5_user.test <id>
```
