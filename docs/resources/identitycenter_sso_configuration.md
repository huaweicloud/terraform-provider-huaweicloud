---
subcategory: "IAM Identity Center"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identitycenter_sso_configuration"
description: |-
  Manages an Identity Center sso configuration resource within HuaweiCloud.
---

# huaweicloud_identitycenter_sso_configuration

Manages an Identity Center sso configuration resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "mfa_mode" {}
variable "allowed_mfa_types" {}
variable "no_mfa_signin_behavior" {}
variable "no_password_signin_behavior" {}
variable "max_authentication_age" {}
variable "configuration_type" {}

resource "huaweicloud_identitycenter_sso_configuration" "test" {
  instance_id                 = var.instance_id
  mfa_mode                    = var.mfa_mode
  allowed_mfa_types           = var.allowed_mfa_types
  no_mfa_signin_behavior      = var.no_mfa_signin_behavior
  no_password_signin_behavior = var.no_password_signin_behavior
  max_authentication_age      = var.max_authentication_age
  configuration_type          = var.configuration_type
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the Identity Center instance.

* `configuration_type` - (Required, String) Specifies the sso configuration type.
  Value Options: **APP_AUTHENTICATION_CONFIGURATION**, **SESSION**.

* `mfa_mode` - (Required, String) Specifies the mfa mode.
  Value Options: **CONTEXT_AWARE**, **DISABLED**, **ALWAYS_ON**.

* `allowed_mfa_types` - (Required, List) Specifies the mfa types.
  Value Options: **TOTP**, **WEBAUTHN**, **WEBAUTHN_SECURITY_KEY**.

* `no_mfa_signin_behavior` - (Required, String) Specifies the sign in behavior if not have mfa.
  Value Options: **ALLOWED_WITH_ENROLLMENT**, **ALLOWED**, **EMAIL_OTP**, **BLOCKED**.

* `no_password_signin_behavior` - (Required, String) Specifies the sign in behavior if not have password.
  Value Options: **BLOCKED**, **EMAIL_OTP**.

* `max_authentication_age` - (Required, String) Specifies the max authentication age.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
