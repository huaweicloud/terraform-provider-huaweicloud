---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_login_protects"
description: |-
  Use this data source to query the list of login protection statuses for IAM users within HuaweiCloud.
---

# huaweicloud_identity_login_protects

Use this data source to query the list of login protection statuses for IAM users within HuaweiCloud.

## Example Usage

```hcl
variable "user_id" {} 

data "huaweicloud_identity_login_protects" "test" {
  user_id = var.user_id
}
```

## Argument Reference

* `user_id` - (Optional, String) Specifies the user id.

## Attribute Reference

* `login_protects` - Indicates the login status protection information list.  
  The [login_protects](#identity_login_protects_login_protects) structure is documented below.

<a name="identity_login_protects_login_protects"></a>
The `login_protects` block supports:

* `user_id` - Indicates the user id.

* `enabled` - Indicates whether to enable login protection.

* `verification_method` - Indicates the login verification method.
