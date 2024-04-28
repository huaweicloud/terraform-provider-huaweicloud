---
subcategory: "IAM Identity Center"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identitycenter_users"
description: ""
---

# huaweicloud_identitycenter_users

Use this data source to get the Identity Center users.

## Example Usage

```hcl
data "huaweicloud_identitycenter_instance" "system" {}

data "huaweicloud_identitycenter_users" "test"{
  identity_store_id = data.huaweicloud_identitycenter_instance.system.identity_store_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `identity_store_id` - (Required, String) Specifies the ID of the identity store that associated with IAM Identity
  Center.

* `user_name` - (Optional, String) Specifies the name of the user.

* `family_name` - (Optional, String) Specifies the family name of the user.

* `given_name` - (Optional, String) Specifies the given name of the user.

* `display_name` - (Optional, String) Specifies the display name of the user.

* `email` - (Optional, String) Specifies the email of the user.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `users` - Indicates the list of IdentityCenter user.
  The [users](#IdentityCenterUsers_User) structure is documented below.

<a name="IdentityCenterUsers_User"></a>
The `users` block supports:

* `id` - Indicates the ID of the user.

* `user_name` - Indicates the name of the user.

* `family_name` - Indicates the family name of the user.

* `given_name` - Indicates the given name of the user.

* `display_name` - Indicates the display name of the user.

* `email` - Indicates the email of the user.
