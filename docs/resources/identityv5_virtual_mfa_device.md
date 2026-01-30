---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identityv5_virtual_mfa_device"
description: |-
  Manages an IAM virtual MFA device resource within HuaweiCloud.
---

# huaweicloud_identityv5_virtual_mfa_device

Manages an IAM virtual MFA device resource within HuaweiCloud.

## Example Usage

```hcl
variable "device_name" {}
variable "user_id" {}

resource "huaweicloud_identityv5_virtual_mfa_device" "test" {
  name    = var.device_name
  user_id = var.user_id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String, NonUpdatable) Specifies the name of the MFA device.
   The name maximum length is `64` characters.  
   Only letters, digits, underscores   (_) and hyphens (-) are allowed.

* `user_id` - (Required, String, NonUpdatable) Specifies the ID of the user.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which is the serial number of the MFA device.

* `base32_string_seed` - The key information used for third-party generation of image verification codes.

* `enabled` - Whether the MFA device is enabled.

## Import

The IAM virtual MFA device can be imported using the `user_id`, e.g:

```bash
$ terraform import huaweicloud_identityv5_virtual_mfa_device.test <user_id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response. The missing attributes include: `name`. It is generally recommended
running `terraform plan` after importing the resource. You can then decide if changes should be applied to the resource,
or the resource definition should be updated to align with the resource. Also, you can ignore changes as below.

```hcl
resource "huaweicloud_identityv5_virtual_mfa_device" "test" {
  ...

  lifecycle {
    ignore_changes = [
      name,
    ]
  }
}
```
