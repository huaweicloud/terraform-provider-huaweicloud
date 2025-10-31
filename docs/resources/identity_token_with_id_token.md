---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_token_with_id_token"
description: |-
  Manages a federated authentication token via the OpenID Connect ID token method within HuaweiCloud.
---

# huaweicloud_identity_token_with_id_token

Manages a federated authentication token via the OpenID Connect ID token method within HuaweiCloud.

->**Note** The token can not be destroyed. It will be invalid after expiration time.

## Example Usage

```hcl
variable "idp_id" {}
variable "id_token" {}
variable "domain_name" {}
variable "project_name" {}

# Get unscoped token
resource "huaweicloud_identity_token_with_id_token" "test1" {
  idp_id   = var.idp_id
  id_token = var.id_token
}

# Get token with domain scope
resource "huaweicloud_identity_token_with_id_token" "test2" {
  idp_id      = var.idp_id
  id_token    = var.id_token
  domain_name = var.domain_name
}

# Get token with project scope
resource "huaweicloud_identity_token_with_id_token" "test3" {
  idp_id       = var.idp_id
  id_token     = var.id_token
  project_name = var.project_name
}
```

## Argument Reference

* `idp_id` - (Required, String, ForceNew) Specifies the identity provider id.

* `id_token` - (Required, String, ForceNew) Specify ID Token of the OpenID Connect Identity Provider.

* `domain_name` - (Optional, String, ForceNew) Specify the domain/account name to get a domain scoped token.

* `project_name` - (Optional, String, ForceNew) Specify the project name to get a project scoped token.

## Attribute Reference

* `id` - Indicates the resource ID in format `<idp_id>:<username>`.

* `token` - Indicates the token. Validity period is 24 hours.

* `username` - Indicates the user of token.

* `groups` - Indicates the group list of the user.

* `expires_at` - The Time when the token will expire. The value is a UTC time in the YYYY-MM-DDTHH:mm:ss.ssssssZ format.
