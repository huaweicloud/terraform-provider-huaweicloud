---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_unscoped_token_saml"
description: |-
  Manages an unscoped token through IdP-initiated federated authentication within HuaweiCloud.
---

# huaweicloud_identity_unscoped_token_saml

Manages an unscoped token through IdP-initiated federated authentication within HuaweiCloud.

->**Note** The token can not be destroyed. It will be invalid after expiration time.

## Example Usage

```hcl
variable "idp_id" {}
variable "saml_response" {}
variable "with_global_domain" {}

resource "huaweicloud_identity_unscoped_token_saml" "test" {
  idp_id             = var.idp_id
  saml_response      = var.saml_response
  with_global_domain = var.with_global_domain
}
```

## Argument Reference

* `idp_id` - (Required, String, ForceNew) Specifies the identity provider id.

* `saml_response` - (Required, String, ForceNew) Specifies the response body returned after successful IdP authentication.
  You could refer to [workspace](https://support.huaweicloud.com/api-iam/iam_02_0003.html#section3) to get a SAMLResponse.

* `with_global_domain` - (Optional, Bool, ForceNew) Specify whether to use a global domain name to obtain the token.
  If set `with_global_domain=true`, it will call `iam.myhuaweicloud.com.` Otherwise, it will call `iam.{region_id}.myhuaweicloud.com`.
  Default value is `false`.

## Attribute Reference

* `id` - Indicates the resource ID in format `<idp_id>:<username>`.

* `token` - Indicates the token. Validity period is 24 hours.

* `username` - Indicates the user of token.

* `groups` - Indicates the group list of the user.

* `expires_at` - The Time when the token will expire. The value is a UTC time in the YYYY-MM-DDTHH:mm:ss.ssssssZ format.
