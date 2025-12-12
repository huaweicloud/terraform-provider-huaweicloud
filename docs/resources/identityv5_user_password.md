---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identityv5_user_password"
description: |-
  Modify IAM user's own password within HuaweiCloud.
---

# huaweicloud_identityv5_user_password

Modify IAM user's own password within HuaweiCloud.

->**Note** The password can not be destroyed.

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

* `new_password` - (Required, String, NonUpdatable) Specifies the IAM user new password.

* `old_password` - (Required, String, NonUpdatable) Specifies the IAM user old password.

## Attribute Reference

* `id` - Resource ID in format `<user_id>`.
