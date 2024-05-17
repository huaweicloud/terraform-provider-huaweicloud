---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_providers"
description: |-
  Use this data source to get the list of IAM identity providers.
---

# huaweicloud_identity_providers

Use this data source to get the list of IAM identity providers.

## Example Usage

```hcl
variable "provider_name" {}

data "huaweicloud_identity_providers" "test" {
  name = var.provider_name
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional, String) Specifies the name of the identity provider.

* `sso_type` - (Optional, String) Specifies the single sign-on type of the identity provider.

* `status` - (Optional, String) Specifies the status of the identity provider. The value can be **true** or **false**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `identity_providers` - The list of identity providers.

  The [identity_providers](#identity_providers_struct) structure is documented below.

<a name="identity_providers_struct"></a>
The `identity_providers` block supports:

* `id` - The identity provider ID which equals the identity provider name.

* `description` - The description of the identity provider.

* `sso_type` - The single sign-on type of the identity provider.

* `status` - The enabled status for the identity provider.

* `remote_ids` - The list of federated user IDs configured for the identity provider.

* `links` - The links of identity provider.

  The [links](#identity_providers_links_struct) structure is documented below.

<a name="identity_providers_links_struct"></a>
The `links` block supports:

* `self` - The identity provider resource link.

* `protocols` - The protocol resource link.
