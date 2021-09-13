---
subcategory: "Identity and Access Management (IAM)"
---

# huaweicloud_identity_group

Use this data source to get details of the specified IAM user group.

## Example Usage

```hcl
data "huaweicloud_identity_group" "group" {
  name = "my_group"
}
```

## Argument Reference

* `name` - (Optional, String) Specifies the name of the identity group.

* `id` - (Optional, String) Specifies the ID of the identity group.

* `description` - (Optional, String) Specifies the description of the identity group.

## Attributes Reference

* `domain_id` - Indicates the domain the group belongs to.

* `users` - Indicates the users the group contains. Structure is documented below.

The `users` block contains:

* `name` - Indicates the IAM user name.

* `id` - Indicates the ID of the User.

* `description` - Indicates the description of the IAM user.

* `enabled` - Indicates the whether the IAM user is enabled.

* `password_expires_at` - Indicates the time when the password will expire.
  Null indicates that the password has unlimited validity.

* `password_status` - Indicates the password status. True means that the password needs to be changed,
  and false means that the password is normal.

* `password_strength` - Indicates the password strength. The value can be high, mid, or low.
