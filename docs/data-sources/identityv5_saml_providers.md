---
subcategory: "IAM"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identityv5_saml_providers"
description: |-
  Use this data source to query the list of SAML identity providers within IAM v5 (new IAM).
---

# huaweicloud_identityv5_saml_providers

Use this data source to query the list of SAML identity providers within IAM v5 (new IAM).

## Example Usage

```hcl
data "huaweicloud_identityv5_saml_providers" "all" {}

data "huaweicloud_identityv5_saml_providers" "filtered" {
  name = "demo-saml-provider"
}

data "huaweicloud_identityv5_saml_providers" "by_id" {
  provider_id = "abc123"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `provider_id` - (Optional, String) Specifies the ID of the SAML identity provider used for filtering.
  If this parameter is set, only the provider with the matched ID is returned.

* `name` - (Optional, String) Specifies the name of the SAML identity provider used for filtering.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `providers` - The list of SAML identity providers that matched the filter parameters.
  The [providers](#identityv5_saml_providers) structure is documented below.

<a name="identityv5_saml_providers"></a>
The `providers` block supports:

* `provider_id` - The ID of the SAML identity provider.

* `name` - The name of the SAML identity provider.

* `description` - The description of the SAML identity provider.

* `saml_metadata_document` - The SAML metadata document of the SAML identity provider.

* `urn` - The uniform resource name of the SAML identity provider.

* `created_at` - The creation time of the SAML identity provider.

* `expires_at` - The expiration time of the SAML identity provider.

* `private_keys` - The list of private keys managed by the SAML identity provider.
  The [private_keys](#private_keys_struct) structure is documented below.

* `assertion_encryption_mode` - The assertion encryption mode of the SAML identity provider.

<a name="private_keys_struct"></a>
The `private_keys` block supports:

* `key_id` - The ID of the private key.

* `timestamp` - The upload time of the private key in ISO 8601 format.
