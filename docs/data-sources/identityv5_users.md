---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identityv5_users"
description: |-
  Use this data source to get the list of IAM users within HuaweiCloud.
---

# huaweicloud_identityv5_users

Use this data source to get the list of IAM users within HuaweiCloud.

## Example Usage

### Query All Users

```hcl
data "huaweicloud_identityv5_users" "test" {}
```

### Query the specified user by user ID

```hcl
variable "user_id" {}

data "huaweicloud_identityv5_users" "test" {
  user_id = var.user_id
}
```

## Argument Reference

The following arguments are supported:

* `group_id` - (Optional, String) Specifies the ID of the user group to which the users belong.

* `user_id` - (Optional, String) Specifies the ID of the user.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `users` - The list of users that matched filter parameters.  
The [users](#v5_users) structure is documented below.

<a name="v5_users"></a>
The `users` block supports:

* `user_id` - The ID of the user.

* `user_name` - The name of the user.

* `enabled` - Whether the user is enabled.

* `description` - The description of the user.

* `is_root_user` - Whether the user is root user.

* `created_at` - The creation time of the user.

* `urn` - The uniform resource name of the user.

* `tags` - The list of tags associated with the user.  
  The [tags](#v5_users_tags) structure is documented below.

* `last_login_at` - The last login time of the user.  
  If omitted, it means that the user has never logged in.

* `password_reset_required` - Whether the password needs to be reset when the user logs in next time.

* `password_expires_at` - The expiration time of the password.

<a name="v5_users_tags"></a>
The `tags` block supports:

* `tag_key` - The key of the tag.

* `tag_value` - The value of the tag.
