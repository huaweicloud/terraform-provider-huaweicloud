---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_protection_policy"
description: ""
---

# huaweicloud_identity_protection_policy

Manages the account protection policy within HuaweiCloud.

-> **NOTE:**
  You *must* have admin privileges to use this resource.  
  This resource overwrites an existing configuration, make sure one resource per account.  
  During action `terraform destroy` it sets values the same as defaults for this resource.

## Example Usage

### Self-Verification

```hcl
resource "huaweicloud_identity_protection_policy" "test" {
  protection_enabled = true
  self_management {
    access_key = true
    password   = true
    email      = false
    mobile     = false
  }
}
```

### Verification by another person

```hcl
resource "huaweicloud_identity_protection_policy" "verification" {
  protection_enabled = true
  verification_email = "example@email.com"
}
```

## Argument Reference

The following arguments are supported:

* `protection_enabled` - (Required, Bool) Specifies whether to enable operation protection.

* `verification_email` - (Optional, String) Specifies the email address used for verification. An example value is `example@email.com`.

* `verification_mobile` - (Optional, String) Specifies the mobile number used for verification. An example value is *0086-123456789*.

-> If `protection_enabled` is set to true and neither `verification_email` nor `verification_mobile` is specified, IAM users
  perform verification by themselves when performing a critical operation.

* `self_management` - (Optional, List) Specifies the attributes IAM users can modify.
  The [object](#self_management_policy) structure is documented below.

<a name="self_management_policy"></a>
The `self_management` block supports:

* `access_key` - (Optional, Bool) Specifies whether to allow IAM users to manage access keys by themselves.

* `password` - (Optional, Bool) Specifies whether to allow IAM users to change their passwords.

* `email` - (Optional, Bool) Specifies whether to allow IAM users to change their email addresses.

* `mobile` - (Optional, Bool) Specifies whether to allow IAM users to change their mobile numbers.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of account protection policy, which is the same as the account ID.

* `self_verification` - Indicates whether the IAM users perform verification by themselves.

## Import

Identity protection policy can be imported using the account ID or domain ID, e.g.

```bash
$ terraform import huaweicloud_identity_protection_policy.example <your account ID>
```
