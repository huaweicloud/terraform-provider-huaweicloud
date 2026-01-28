---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_group"
description: |-
  Use this data source to get details of the specified IAM user group.
---

# huaweicloud_identity_group

Use this data source to get details of the specified IAM user group.

-> **NOTE:** You *must* have IAM read privileges to use this data source.

## Example Usage

```hcl
variable "group_name" {}

data "huaweicloud_identity_group" "group" {
  name = var.group_name
}
```

## Argument Reference

* `name` - (Optional, String) Specifies the name of the identity group.

* `id` - (Optional, String) Specifies the ID of the identity group.

* `description` - (Optional, String) Specifies the description of the identity group.

## Attribute Reference

* `domain_id` - Indicates the domain the group belongs to.

* `users` - Indicates the users the group contains.  
  The [users](#identity_group_users) structure is documented below.

<a name="identity_group_users"></a>
The `users` block contains:

* `id` - Indicates the ID of the IAM user.

* `name` - Indicates the name of the IAM user.

* `description` - Indicates the description of the IAM user.

* `enabled` - Whether the IAM user is enabled.

* `password_expires_at` - Indicates the time when the password will expire.
  If this value is not set, the password will not expire.

* `password_status` - Indicates the password status. True means that the password needs to be changed,
  and false means that the password is normal.

* `password_strength` - Indicates the password strength. The value can be high, mid, or low.
