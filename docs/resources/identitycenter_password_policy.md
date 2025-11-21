---
subcategory: "IAM Identity Center"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identitycenter_password_policy"
description: |-
  Manages an Identity Center password policy resource within HuaweiCloud.
---

# huaweicloud_identitycenter_password_policy

Manages an Identity Center password policy resource within HuaweiCloud.

## Example Usage

```hcl
variable "identity_store_id " {}
variable "max_password_age" {}
variable "minimum_password_length" {}
variable "password_reuse_prevention" {}
variable "require_lowercase_characters" {}
variable "require_numbers" {}
variable "require_symbols" {}
variable "require_uppercase_characters" {}

resource "huaweicloud_identitycenter_password_policy" "test" {
  identity_store_id            = var.identity_store_id
  max_password_age             = var.max_password_age
  minimum_password_length      = var.minimum_password_length
  password_reuse_prevention    = var.password_reuse_prevention
  require_lowercase_characters = var.require_lowercase_characters
  require_numbers              = var.require_numbers
  require_symbols              = var.require_symbols
  require_uppercase_characters = var.require_uppercase_characters
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `identity_store_id` - (Required, String, NonUpdatable) Specifies the ID of the identity store.

* `max_password_age` - (Optional, Int) Specifies the max password age.

* `minimum_password_length` - (Optional, Int) Specifies the minimum password length.

* `password_reuse_prevention` - (Optional, Int) Specifies the password reuse times. Value Options: **1**.

* `require_lowercase_characters` - (Optional, Bool) Whether the password require lowercase characters.

* `require_numbers` - (Optional, Bool) Whether the password require numbers.

* `require_symbols` - (Optional, Bool) Whether the password require symbols.

* `require_uppercase_characters` - (Optional, Bool) Whether the password require uppercase characters.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
