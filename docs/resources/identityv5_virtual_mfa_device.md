---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identityv5_virtual_mfa_device"
description: |-
  Manages an IAM v5 virtual mfa device resource within HuaweiCloud.
---

# huaweicloud_identityv5_virtual_mfa_device

Manages an IAM v5 virtual mfa device resource within HuaweiCloud.

## Example Usage

```hcl
variable "name" {}
variable "user_id" {}

resource "huaweicloud_identityv5_virtual_mfa_device" "test" {
  name    = var.name
  user_id = var.user_id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String, NonUpdatable) Specifies the name of the user. The username consists of 1 to 64 characters.
  It can contain only uppercase letters, lowercase letters, digits, spaces, and special characters (-_) and cannot
  start with a digit.

* `user_id` - (Required, String, NonUpdatable) Specifies the ID of the user.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID of mfa device.

* `base32_string_seed` - Indicates key information used for third-party generation of image verification codes.

* `enabled` - Indicates whether the user is enabled or disabled. Valid values are **true** and **false**.

## Import

The IAM v5 virtual mfa device can be imported using the id, e.g:

```bash
$ terraform import huaweicloud_identityv5_virtual_mfa_device.test <id>
```
