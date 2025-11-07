---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_kms_dedicated_keystore_action"
description: |-
  Manages actions (enable/disable) on a dedicated keystore within HuaweiCloud KMS.
---

# huaweicloud_kms_dedicated_keystore_action

Manages actions (enable/disable) on a dedicated keystore within HuaweiCloud KMS.

-> Destroying this resource will not affect the actual status of the dedicated keystore, but will only remove the
  resource information from the tfstate file.

## Example Usage

```hcl
variable "keystore_id" {}
variable "action" {}

resource "huaweicloud_kms_dedicated_keystore_action" "test" {
  keystore_id = var.keystore_id
  action      = var.action
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to manage the dedicated keystore.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `keystore_id` - (Required, String, NonUpdatable) Specifies the ID of the dedicated keystore to manage.

* `action` - (Required, String) Specifies the action to perform on the dedicated keystore.
  Valid values are **enable** to enable the keystore or **disable** to disable it.

  -> The dedicated keystore can only be disabled when it is in the enabled state. Similarly, it can only be enabled when
    it is in the disabled state.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID which is the same as the `keystore_id`.
