---
subcategory: "IAM Identity Center"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identitycenter_service_provider_configuration"
description: |-
  Use this data source to get the Identity Center service provider configuration.
---

# huaweicloud_identitycenter_service_provider_configuration

Use this data source to get the Identity Center service provider configuration.

## Example Usage

```hcl
data "huaweicloud_identitycenter_instance" "test" {}

data "huaweicloud_identitycenter_service_provider_configuration" "test"{
  identity_store_id = data.huaweicloud_identitycenter_instance.test.identity_store_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `identity_store_id` - (Required, String) Specifies the ID of the identity store that associated with IAM Identity
  Center.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `sp_oidc_config` - The oidc configuration of the service provider.
  The [sp_oidc_config](#sp_oidc_config_struct) structure is documented below.

* `sp_saml_config` - The saml configuration of the service provider.
  The [sp_saml_config](#sp_saml_config_struct) structure is documented below.

<a name="sp_oidc_config_struct"></a>
The `sp_oidc_config` block supports:

* `redirect_url` - The redirect url.

<a name="sp_saml_config_struct"></a>
The `sp_saml_config` block supports:

* `acs_url` -  The acs url.

* `issuer` -  The issuer.

* `metadata` - The metadata.
