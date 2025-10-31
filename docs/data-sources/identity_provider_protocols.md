---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_provider_protocols"
description: |-
  Use this data source to query the list of IAM identity provider protocols within HuaweiCloud.
---

# huaweicloud_identity_provider_protocols

Use this data source to query the list of IAM identity provider protocols within HuaweiCloud.

## Example Usage

```hcl
variable "provider_id" {}

data "huaweicloud_identity_provider_protocols" "test1" {
  provider_id = var.provider_id
}

data "huaweicloud_identity_provider_protocols" "test2" {
  provider_id = var.provider_id
  protocol_id = "saml"
}
```

## Argument Reference

The following arguments are supported:

* `provider_id` - (Required, String) Name of an identity provider.

* `protocol_id` - (Optional, String) Specifies the protocol id.

## Attribute Reference

* `protocols` - Indicates the protocol Information List.
  The [protocols](#IdentityProtocols_Protocols) structure is documented below.

<a name="IdentityProtocols_Protocols"></a>
The `protocols` block contains:

* `id` - The protocol id.

* `mapping_id` - Indicates the mapping id of the protocol.

* `links` - Indicates the links of protocol.
  The [links](#IdentityProtocols_Links) structure is documented below.

<a name="IdentityProtocols_Links"></a>
The `links` block contains:

* `self` - Indicates the protocol resource link.

* `identity_provider` - Indicates the identity provider resource link.
