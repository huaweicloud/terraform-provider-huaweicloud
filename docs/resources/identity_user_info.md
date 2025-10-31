---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_user_info"
description: |-
  Modify IAM user's own information within HuaweiCloud.
---

# huaweicloud_identity_user_info

Modify IAM user's own information within HuaweiCloud.

->**Note** The information can not be destroyed.

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

* `email` - (Optional, String) Specifies the IAM user email. The email must conform to the email format and be no longer
  than 255 characters.

* `mobile` - (Optional, String) Specifies the IAM user mobile. The mobile format is `<country code>-<phone number>`,
  such as 0086-123456789.

## Attribute Reference

* `id` - Resource ID in format `<user_id>`.
