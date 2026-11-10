---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identityv5_oidc_providers"
description: |-
  Use this data source to query IAM v5 OIDC identity providers within HuaweiCloud.
---

# huaweicloud_identityv5_oidc_providers

Use this data source to query IAM v5 OIDC identity providers within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_identityv5_oidc_providers" "all" {}

data "huaweicloud_identityv5_oidc_providers" "filtered" {
  name = "my-oidc-provider"
}

data "huaweicloud_identityv5_oidc_providers" "by_id" {
  provider_id = "abc123"
}
```

## Argument Reference

The following arguments are supported:

* `provider_id` - (Optional, String) Specifies the ID of the OIDC identity provider used for filtering.
  If this parameter is set, only the provider with the matched ID is returned.

* `name` - (Optional, String) Specifies the name of the OIDC identity provider used for filtering.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `providers` - The list of OIDC identity providers that matched the filter parameters.

The [providers](#identityv5_oidc_providers) structure is documented below.

<a name="identityv5_oidc_providers"></a>
The `providers` block supports:

* `provider_id` - The ID of the OIDC identity provider.

* `name` - The name of the OIDC identity provider.

* `description` - The description of the OIDC identity provider.

* `url` - The URL of the OIDC identity provider.

* `urn` - The uniform resource name of the OIDC identity provider.

* `client_ids` - The list of client IDs registered in the OIDC identity provider.

* `thumbprints` - The list of server certificate fingerprints of the OIDC identity provider.

* `created_at` - The creation time of the OIDC identity provider.

* `expires_at` - The expiration time of the OIDC identity provider.
