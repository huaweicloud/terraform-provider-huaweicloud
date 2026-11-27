---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identityv5_oidc_provider"
description: |-
  Manages an IAM v5 OIDC identity provider resource within HuaweiCloud.
---

# huaweicloud_identityv5_oidc_provider

Manages an IAM v5 OIDC identity provider resource within HuaweiCloud.

## Example Usage

```hcl
variable "issuer_url" {}

resource "huaweicloud_identityv5_oidc_provider" "test" {
  name        = "my-oidc-provider"
  url         = var.issuer_url
  description = "Created by terraform script"
  client_ids  = ["my-backend-service"]
  thumbprints = ["902ef2deeb3c5b13ea4c3d5193629309e231ae55"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, ForceNew, String) Specifies the name of the OIDC identity provider.  
  The name consists of `1` to `64` characters.  
  Only letters, digits, underscores (_) and hyphens (-) are allowed.

* `url` - (Required, ForceNew, String) Specifies the URL of the OIDC identity provider.  
  The value corresponds to the `issuer` field in the OpenID Connect metadata document
  (`https://{base URL}/.well-known/openid-configuration`) of the identity provider,
  and contains `1` to `255` characters.

* `description` - (Optional, String) Specifies the description of the OIDC identity provider.  
  The maximum length is `255` characters, cannot includes these special characters: `@#$%^&*<>"`.

* `client_ids` - (Required, Set) Specifies the list of client IDs registered in the OIDC identity provider.  
  A maximum of `100` client IDs can be configured, and each client ID contains `1` to `255` characters.  
  When modified, the added or removed client IDs are updated separately.

* `thumbprints` - (Optional, Set) Specifies the list of server certificate fingerprints of the OIDC identity provider.  
  The fingerprint is the SHA-1 hash of the top-intermediate CA certificate of the JWKS endpoint of the identity  
  provider.
  If this parameter is not specified, IAM automatically obtains the CA certificate fingerprint of the JWKS endpoint.  
  When modified, the list is replaced entirely.  

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the OIDC identity provider.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, also the OIDC identity provider ID.

* `urn` - The uniform resource name of the OIDC identity provider.

* `created_at` - The creation time of the OIDC identity provider.

## Import

The IAM v5 OIDC identity provider can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_identityv5_oidc_provider.test <id>
```
