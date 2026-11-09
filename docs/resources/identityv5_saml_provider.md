---
subcategory: "IAM"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identityv5_saml_provider"
description: |-
  Manages a SAML identity provider resource within IAM v5 (new IAM).
---

# huaweicloud_identityv5_saml_provider

Manages a SAML identity provider resource within IAM v5 (new IAM).

## Example Usage

```hcl
variable "saml_metadata_document" {}

resource "huaweicloud_identityv5_saml_provider" "test" {
  name        = "demo-saml-provider"
  description = "Created by terraform"

  saml_metadata_document = var.saml_metadata_document
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the SAML identity provider.
  The name contains 1 to 64 characters and can only contain letters, digits, underscores (_) and hyphens (-).
  Changing this parameter will create a new resource.

* `saml_metadata_document` - (Required, String) Specifies the SAML metadata document obtained from the identity provider.
  The document is an XML string that contains the issuer name, the validity period and the signing certificate
  used to verify SAML assertions.

* `description` - (Optional, String) Specifies the description of the SAML identity provider.
  The description contains a maximum of 255 characters and cannot contain the following characters:
  @, #, %, &, <, >, ", $, ^ and *.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which is the provider ID of the SAML identity provider.

* `provider_id` - The ID of the SAML identity provider.

* `urn` - The uniform resource name of the SAML identity provider.

* `created_at` - The creation time of the SAML identity provider.

* `expires_at` - The expiration time of the SAML identity provider.

* `private_keys` - The list of private keys managed by the SAML identity provider.
  The [private_keys](#private_keys_struct) structure is documented below.

* `assertion_encryption_mode` - The assertion encryption mode of the SAML identity provider.
  The valid values are **Allowed** and **Required**.

### `private_keys` Struct

The `private_keys` block supports the following attributes:

* `key_id` - The ID of the private key.

* `private_key` - The private key content.

* `created_at` - The creation time of the private key.

* `expires_at` - The expiration time of the private key.

## Import

The SAML identity provider can be imported using the provider ID, e.g.

```bash
$ terraform import huaweicloud_identityv5_saml_provider.test <provider_id>
```
