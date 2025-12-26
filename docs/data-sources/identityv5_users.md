---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identityv5_users"
description: |-
  Use this data source to get the list of user in the Identity and Access Management V5 service.
---

# huaweicloud_identityv5_users

Use this data source to get the list of user in the Identity and Access Management V5 service.

## Example Usage

```hcl
data "huaweicloud_identityV5_users" "users" {}
```

### ShowUser

```hcl
resource "huaweicloud_identityv5_user" "user_1" {
name        = "TestUser"
description = "tested by terraform"
}

data "huaweicloud_identityv5_users" "test" {
user_id = huaweicloud_identityv5_user.user_1.id
}
```

## Argument Reference

The following arguments are supported:

* `group_id` - (Optional, String) Specifies the group ID of the users.

* `user_id` - (Optional, String) Specifies the id of the IAM user.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `users` - Indicates the user information.
The [users](#IdentityV5User_Users) structure is documented below.

<a name="IdentityV5User_Users"></a>
The `users` block contains:

* `id` - Indicates the group ID of this membership.

* `user_name` - Indicates the IAM username.

* `is_root_user` - Indicates whether the user is root user or user.

* `created_at` - Indicates the time when the IAM user was created.

* `urn` - Indicates the uniform resource name.

* `user_id` - Indicates the user ID in this group.

* `description` - Indicates the description of the user.

* `enabled` - Indicates whether the user is enabled or disabled. Valid values are **true** and **false**.

* `last_login_at` - Indicates the last login time of the IAM user. If null, it indicates that they have never logged in.

* `tags` - Indicates the Custom Tag List.
  The [tags](#IdentityV5User_Tags) structure is documented below.

* `password_reset_required` - Indicates whether the password reset is required.

* `password_expires_at` - Indicates the time when the password expired.

<a name="IdentityV5User_Tags"></a>
The `tags` block contains:

* `tag_key` - Indicates Tag keys which contain any combination of letters, numbers, spaces, and the symbols "_", ".",
  ":","=", "-", "@", but cannot start or end with a space, and cannot begin with "sys". The length must be between 1
  and 64 characters.

* `tag_value` - Indicates Tag keys can contain any combination of letters, numbers, spaces, and the symbols "_", ".",
  ":", "=", "-", "@", but cannot start or end with a space, and cannot begin with "sys". The length must be between 1
  and 64 characters.
