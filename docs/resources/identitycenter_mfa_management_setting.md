---
subcategory: "IAM Identity Center"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identitycenter_mfa_management_setting"
description: |-
  Manages an Identity Center mfa management setting resource within HuaweiCloud.
---

# huaweicloud_identitycenter_mfa_management_setting

Manages an Identity Center mfa management setting resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "user_permission " {}
variable "identity_store_id " {}

resource "huaweicloud_identitycenter_mfa_management_setting" "test" {
  instance_id       = var.instance_id
  user_permission   = var.user_permission
  identity_store_id = var.identity_store_id
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the Identity Center instance.

* `identity_store_id` - (Required, String) Specifies the ID of the identity store.

* `user_permission` - (Required, String) Whether to allow users to manage MFA themselves.
  Value Options: **READ_ACTIONS**, **ALL_ACTIONS**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
