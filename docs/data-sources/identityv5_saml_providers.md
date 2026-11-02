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
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

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

* `create_time` - The creation time of the SAML identity provider.

* `expire_time` - The expiration time of the SAML identity provider.
