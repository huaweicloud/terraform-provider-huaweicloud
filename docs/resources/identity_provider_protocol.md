---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_provider_protocol"
description: |-
  Manages an IAM provider protocol resource within HuaweiCloud.
---

# huaweicloud_identity_provider_protocol

Manages an IAM provider protocol resource within HuaweiCloud.

->**Note:** When you create identity provider through console or terraform resource `huaweicloud_identity_provider`,
it will create a default protocol with a default mapping. If the mapping of the existed protocol is different from the
mapping you specified in `huaweicloud_identity_provider_protocol`, it will try to update the protocol.

## Example Usage

```hcl
variable "provider_name" {}

# Without protocol setting
resource "huaweicloud_identity_provider" "test" {
  name = var.provider_name
}

resource "huaweicloud_identity_provider_conversion" "test" {
  provider_id = huaweicloud_identity_provider.test.id

  conversion_rules {
    ...
  }
}

resource "huaweicloud_identity_provider_protocol" "test" {
  provider_id = huaweicloud_identity_provider.test.id
  protocol_id = "saml"
  mapping_id  = huaweicloud_identity_provider_conversion.conversion.id
}
```

## Argument Reference

The following arguments are supported:

* `provider_id` - (Required, String) Specifies the ID of the identity provider used to manage the protocol.

* `protocol_id` - (Required, String) Specifies the identity protocol of the identity provider.  
  The valid values are as follows:
  + **saml**: SAML protocol.
  + **oidc**: OpenID Connect protocol.

* `mapping_id` - (Optional, String) Specifies the mapping_id for the protocol. When the identity provider type is
  `iam_user_sso`, there is no need to bind a mapping ID, and this field does not need to be passed; otherwise,
  this field is mandatory.

## Attribute Reference

* `links` - Indicates the links of protocol.  
  The [links](#identity_provider_protocol_links) structure is documented below.

<a name="identity_provider_protocol_links"></a>
The `links` block supports:

* `self` - Indicates the resource link.

* `identity_provider` - Indicates the identity provider resource link.

## Import

Identity provider protocol can be imported using the `<provider_id>:<protocol_id>`. For example,
if you have provider_id `provider_test` and protocol `saml`, you should use `provider_test:saml` to import.

```bash
$ terraform import huaweicloud_identity_provider_protocol.test <provider_id>:<protocol_id>
```
