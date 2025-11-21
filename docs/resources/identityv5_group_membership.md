---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identityv5_group_membership"
description: |-
  Manages an IAM v5 group membership resource within HuaweiCloud.
---

# huaweicloud_identityv5_group_membership

Manages an IAM v5 group membership resource within HuaweiCloud.

## Example Usage

```hcl
resource "huaweicloud_identityv5_group" "group_1" {
  group_name = "test_group"
}

resource "huaweicloud_identityv5_user" "user_1" {
  name        = "test_name1"
  description = "tested by terraform"
  enabled     = true
}

resource "huaweicloud_identityv5_user" "user_2" {
  name        = "test_name2"
  description = "tested by terraform"
  enabled     = true
}

resource "huaweicloud_identityv5_group_membership" "membership_1" {
  group_id     = huaweicloud_identityv5_group.group_1.id
  user_id_list = [huaweicloud_identityv5_user.user_1.id, huaweicloud_identityv5_user.user_2.id]
}
```

## Argument Reference

The following arguments are supported:

* `group_id` - (Required, String) Specifies the group ID of this membership.

* `user_id_list` - (Required, List) Specifies a list of IAM user IDs to associate to the group.

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
