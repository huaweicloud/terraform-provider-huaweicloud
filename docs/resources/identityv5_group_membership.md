---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identityv5_group_membership"
description: |-
  Manages an IAM group membership resource within HuaweiCloud.
---

# huaweicloud_identityv5_group_membership

Manages an IAM group membership resource within HuaweiCloud.

## Example Usage

```hcl
variable "user_group_id" {}
variable "user_ids" {
  type = list(string)
}

resource "huaweicloud_identityv5_group_membership" "test" {
  group_id     = var.user_group_id
  user_id_list = var.user_ids
}
```

## Argument Reference

The following arguments are supported:

* `group_id` - (Required, String, NonUpdatable) Specifies the ID of the user group.

* `user_id_list` - (Required, List) Specifies the list of user IDs to associate with the group.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `users` - The list of users associated with the group.  
  The [users](#v5_group_membership_users) structure is documented below.

<a name="v5_group_membership_users"></a>
The `users` block supports:

* `user_id` - The ID of the user.

* `user_name` - The name of the user.

* `enabled` - Whether the user is enabled.

* `description` - The description of the user.

* `is_root_user` - Whether the user is a root user.

* `created_at` - The creation time of the user.

* `urn` - The uniform resource name of the user.

## Import

The resource can be imported using group ID (`id`), e.g.

```bash
$ terraform import huaweicloud_identityv5_group_membership.test <id>
```
