---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_federation_domains"
description: |-
    Use this data source to query the list of accounts accessible by federated users within HuaweiCloud.
---

# huaweicloud_identity_federation_domains

Use this data source to query the list of accounts accessible by federated users within HuaweiCloud.

## Example Usage

```hcl
variable "federation_token" {}

data "huaweicloud_identity_federation_domains" "test" {
  federation_token = var.federation_token
}
```

## Argument Reference

* `federation_token` - (Required, String) Specifies federated authentication unscoped token.

## Attribute Reference

* `domains` - Indicates the account information list.
  The [domains](#IdentityFederationDomains_Domains) structure is documented below.

<a name="IdentityFederationDomains_Domains"></a>
The `domains` block supports:

* `id` - Indicates the account id.

* `name` - Indicates the account name.

* `description` - Indicates the account description information.

* `enabled` - Indicates whether the account is enabled.
