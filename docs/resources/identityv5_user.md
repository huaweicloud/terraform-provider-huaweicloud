---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identityv5_user"
description: |-
  Manages an IAM user resource within HuaweiCloud.
---

# huaweicloud_identityv5_user

Manages an IAM user resource within HuaweiCloud.

-> You must have admin privileges to use this resource.

## Example Usage

```hcl
variable "user_name" {}

resource "huaweicloud_identityv5_user" "test" {
  name = var.user_name
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) Specifies the name of the user.  
  The username consists of `1` to `64` characters.  
  It can contain only letters, digits, spaces, and special characters (-_) and cannot start with a digit.

* `enabled` - (Optional, Bool) Specifies whether the user is enabled.  
  Default value is **true**.
  
* `description` - (Optional, String) Specifies the description of the user.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, also the user ID.

* `is_root_user` - Whether the user is root user.

* `created_at` - The creation time of the user.

* `urn` - The uniform resource name of the user.

* `last_login_at` - The last login time of the user.  
  If omitted, it means that the user has never logged in.

* `tags` - The list of tags associated with the user.  
  The [tags](#IdentityV5User_Tags) structure is documented below.

<a name="IdentityV5User_Tags"></a>
The `tags` block supports:

* `tag_key` - The key of the tag.

* `tag_value` - The value of the tag.

## Import

The user can be imported using the `id`, e.g:

```bash
$ terraform import huaweicloud_identityv5_user.test <id>
```
