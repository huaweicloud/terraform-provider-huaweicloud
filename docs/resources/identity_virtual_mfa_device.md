---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_virtual_mfa_device"
description: ""
---

# huaweicloud_identity_virtual_mfa_device

Manages an IAM virtual MFA device resource within HuaweiCloud.

-> **NOTE:** You *must* have admin privileges to use this resource.

-> **NOTE:** A virtual MFA device cannot be directly associated with an IAM User from Terraform. To associate the
virtual MFA device with a user and enable it, use the code returned in either `base32_string_seed` or `qr_code_png` to
generate TOTP authentication codes.

## Example Usage

```hcl
variable "name" {}
variable "user_id" {}

resource "huaweicloud_identity_virtual_mfa_device" "test" {
  name    = var.name
  user_id = var.user_id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String, ForceNew) Specifies the virtual MFA device name. Changing this will create a new virtual
  MFA device.

* `user_id` - (Required, String, ForceNew) Specifies the user ID which the virtual MFA device belongs to.
  Changing this will create a new virtual MFA device.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates the serial number.

* `base32_string_seed` - The base32 seed, which a third-patry system can use to generate a `CAPTCHA` code.

* `qr_code_png` - A QR code PNG image that encodes `otpauth://totp/huawei:$domainName@$userName?secret=$Base32String`
  where `$domainName` is IAM domain name, `$userName` is IAM user name, and `Base32String` is the seed in base32 format.

## Import

The virtual MFA device can be imported using the `user_id`, e.g.

```bash
$ terraform import huaweicloud_identity_virtual_mfa_device.test <user_id>
```
