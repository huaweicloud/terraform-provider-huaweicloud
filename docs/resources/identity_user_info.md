---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_user_info"
description: |-
  Manages an IAM user info resource within HuaweiCloud.
---

# huaweicloud_identity_user_info

Manages an IAM user info resource within HuaweiCloud.

-> This resource is a one-time action resource for creating user information. Deleting this resource will
   not clear the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "email" {}
variable "mobile" {}

resource "huaweicloud_identity_user_info" "test" {
  email  = var.email
  mobile = var.mobile
}
```

## Argument Reference

* `email` - (Optional, String, NonUpdatable) Specifies the email of the user. The email must conform to the email
  format and be no longer than `255` characters.

* `mobile` - (Optional, String, NonUpdatable) Specifies the mobile phone number of the user. The mobile format is
  `<country code>-<phone number>`. For example: `0086-123456789`.

## Attribute Reference

* `id` - The resource ID in format `<user_id>`.
