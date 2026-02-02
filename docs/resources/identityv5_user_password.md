---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identityv5_user_password"
description: |-
  Use this resource to modify IAM user's own password within HuaweiCloud.
---

# huaweicloud_identityv5_user_password

Use this resource to modify IAM user's own password within HuaweiCloud.

-> This resource is only a one-time action resource for changing user password. Deleting this resource will
   not clear the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "new_password" {}
variable "old_password" {}

resource "huaweicloud_identityv5_user_password" "test" {
  new_password = var.new_password
  old_password = var.old_password
}
```

## Argument Reference

* `new_password` - (Required, String, NonUpdatable) Specifies the new password of the user.

* `old_password` - (Required, String, NonUpdatable) Specifies the old password of the user.

## Attribute Reference

* `id` - The resource ID, also the user ID.
