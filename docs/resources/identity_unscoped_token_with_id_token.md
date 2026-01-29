---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_unscoped_token_with_id_token"
description: |-
  Manages a federated authentication unscoped token via the OpenID Connect ID token method within HuaweiCloud.
---

# huaweicloud_identity_unscoped_token_with_id_token

Manages a federated authentication unscoped token via the OpenID Connect ID token method within HuaweiCloud.

->**Note** The token can not be destroyed. It will be invalid after expiration time.

## Example Usage

```hcl
variable "idp_id" {}
variable "protocol_id" {}
variable "id_token" {}

resource "huaweicloud_identity_unscoped_token_with_id_token" "test" {
  idp_id      = var.idp_id
  protocol_id = var.protocol_id
  id_token    = var.id_token
}
```

## Argument Reference

* `idp_id` - (Required, String, ForceNew) Specifies the identity provider id.

* `protocol_id` - (Required, String, ForceNew) Specifies the protocol id.  
  The valid value is **oidc**.

* `id_token` - (Required, String, ForceNew) The security token of the OpenID Connect Identity Provider,
  format is Bearer {ID Token}.

## Attribute Reference

* `id` - Indicates the resource ID in format `<idp_id>:<username>`.

* `token` - Indicates the details of the obtained token. The valid period is `24` hours.

* `username` - Indicates the user of token.

* `groups` - Indicates the group list of the user.

* `expires_at` - The Time when the token will expire. The value is a UTC time in the YYYY-MM-DDTHH:mm:ss.ssssssZ format.
