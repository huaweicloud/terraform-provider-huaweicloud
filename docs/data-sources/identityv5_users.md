---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identityv5_users"
description: |-
  Use this data source to get the list of user in the Identity and Access Management V5 service.
---

# huaweicloud_identityv5_virtual_mfa_devices

Use this data source to get the list of user in the Identity and Access Management V5 service.

## Example Usage

```hcl
data "huaweicloud_identityV5_users" "devices" {}
```

## Argument Reference

The following arguments are supported:

* `group_id` - (Optional, String) Specifies the group ID of the users.

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
