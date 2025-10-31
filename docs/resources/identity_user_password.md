---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_user_password"
description: |-
  Modify IAM user's own passwords within HuaweiCloud.
---

# huaweicloud_identity_user_password

Modify IAM user's own passwords within HuaweiCloud.

->**Note** The password can not be destroyed.

## Example Usage

```hcl
variable "password" {}
variable "original_password" {}

resource "huaweicloud_identity_user_password" "test" {
  password          = var.password
  original_password = var.original_password
}
```

## Argument Reference

* `password` - (Required, String, ForceNew) Specifies the IAM user password.

* `original_password` - (Required, String, ForceNew) Specifies the IAM user original password.

## Attribute Reference

* `id` - Resource ID in format `<user_id>`.
