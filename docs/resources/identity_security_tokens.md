---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: identity_security_tokens"
description: ""
---

# identity_security_tokens

Use this data source to create temporary access key by token within HuaweiCloud.

-> **NOTE:** You *must* have IAM read privileges to use this data source.

## Example Usage

```hcl
variable "version" {}
variable "action" {}
variable "effect" {}

resource "huaweicloud_identity_security_tokens" "test" {
  version = var.version
  action = var.action
  effect = var.effect
}
```

## Argument Reference

* `methods` -  (Required, Array of strings) Authentication methods. The value must be ["token"].

* `Version` - (Required, String) Policy version number. Must be set to "1.1" when creating custom policies.

* `Statement` - (Required, Array of objects) Authorization statements that define the custom policy content.

## Attribute Reference

* `expires_at` - The Time when the token will expire. The value is a UTC time in the YYYY-MM-DDTHH:mm:ss.ssssssZ format.

* `access` - The retrieved Access Key (AK).

* `secret ` - The retrieved Secret Key (SK).

* `securitytoken` - An encrypted string containing the obtained AK, SK and other security information.
