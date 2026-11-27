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
  used to verify SAML assertions. The length ranges from 1,000 to 512,000 characters.

* `description` - (Optional, String) Specifies the description of the SAML identity provider.
  The description contains a maximum of 255 characters and cannot contain the following characters:
  @, #, %, &, <, >, ", $, ^ and *.

* `assertion_encryption_mode` - (Optional, String) Specifies the assertion encryption mode of the SAML identity provider.
  The valid values are **Allowed** (accepts encrypted or unencrypted SAML assertions) and
  **Required** (only accepts encrypted SAML assertions).

* `add_private_key` - (Optional, String, Sensitive) Specifies the PEM format RSA private key used to decrypt SAML assertions.
  The length ranges from 1 to 16,384 characters. This parameter is only used for write operations and is not stored
  in the state file.

* `remove_private_key` - (Optional, String) Specifies the ID of the private key to be removed from the SAML identity provider.
  The value contains 22 to 64 characters and can only contain uppercase letters and digits. This parameter is only used
  for write operations and is not stored in the state file.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the SAML identity provider.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which is the provider ID of the SAML identity provider.

* `provider_id` - The ID of the SAML identity provider.

* `urn` - The uniform resource name of the SAML identity provider.

* `created_at` - The creation time of the SAML identity provider.

* `expires_at` - The expiration time of the SAML identity provider.

* `private_keys` - The list of private keys managed by the SAML identity provider.
  The [private_keys](#private_keys_struct) structure is documented below.

<a name="private_keys_struct"></a>
The `private_keys` block supports:

* `key_id` - The ID of the private key.

* `timestamp` - The upload time of the private key in ISO 8601 format.

## Import

The SAML identity provider can be imported using the provider ID, e.g.

```bash
$ terraform import huaweicloud_identityv5_saml_provider.test <provider_id>
```
