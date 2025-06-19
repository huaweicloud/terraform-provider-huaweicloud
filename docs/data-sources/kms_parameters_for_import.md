---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_kms_parameters_for_import"
description: |-
  Use this data source to get parameters required for importing a key.
---

# huaweicloud_kms_parameters_for_import

Use this data source to get parameters required for importing a key, including an import token and a public key.

## Example Usage

```hcl
variable "key_id" {}
variable "wrapping_algorithm" {}

data "huaweicloud_kms_parameters_for_import" "test" {
  key_id             = var.key_id
  wrapping_algorithm = var.wrapping_algorithm
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to obtain the KMS parameters.
  If omitted, the provider-level region will be used.

* `key_id` - (Required, String) Specifies the key ID. It should be `36` bytes and match the regular expression
  **^[0-9a-z]{8}-[0-9a-z]{4}-[0-9a-z]{4}-[0-9a-z]{4}-[0-9a-z]{12}$**.
  For example, **0d0466b0-e727-4d9c-b35d-f84bb474a37f**.

* `wrapping_algorithm` - (Required, String) Specifies the encryption algorithm of key materials.
  The valid values are **RSAES_OAEP_SHA_256** and **SM2_ENCRYPT**.

  -> Some regions do not support **SM2_ENCRYPT** import type.

* `sequence` - (Optional, String) Specifies the `36` bytes sequence number of a request message.
  For Example, **919c82d4-8046-4722-9094-35c3c6524cff**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `import_token` - The key import token.

* `expiration_time` - The import parameter expiration time. The format is 10-digit timestamp in second.

* `public_key` - The public key of the DEK material, in Base64 format.
