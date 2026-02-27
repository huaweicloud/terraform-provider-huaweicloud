---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_user_password"
description: |-
  Manages an IAM user's password resource within HuaweiCloud.
---

# huaweicloud_identity_user_password

Manages an IAM user's password resource within HuaweiCloud.

-> This resource is a one-time action resource for modifying user password. Deleting this resource will
   not clear the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "original_password" {}
variable "password" {}

resource "huaweicloud_identity_user_password" "test" {
  original_password = var.original_password
  password          = var.password
}
```

## Argument Reference

* `original_password` - (Required, String, NonUpdatable) Specifies the original password of the IAM user.

* `password` - (Required, String, NonUpdatable) Specifies the new password of the IAM user.

## Attribute Reference

* `id` - Resource ID in format `<user_id>`.
