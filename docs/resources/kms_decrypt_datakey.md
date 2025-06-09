---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_kms_decrypt_datakey"
description: |-
  Manages a KMS decrypt datakey resource within HuaweiCloud.
---

# huaweicloud_kms_decrypt_datakey

Manages a KMS decrypt datakey resource within HuaweiCloud.

-> The current resource is a one-time resource, and destroying this resource will not recover the decrypted datakey,
but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "key_id" {}
variable "cipher_text" {}
variable "datakey_cipher_length" {}

resource "huaweicloud_kms_decrypt_datakey" "test" {
  key_id                = var.key_id
  cipher_text           = var.cipher_text
  datakey_cipher_length = var.datakey_cipher_length
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `key_id` - (Required, String, NonUpdatable) Specifies the key ID.
  The valid length is `36` bytes, meeting regular match **^[0-9a-z]{8}-[0-9a-z]{4}-[0-9a-z]{4}-[0-9a-z]{4}-[0-9a-z]{12}$**.
  For example: **0d0466b0-e727-4d9c-b35d-f84bb474a37f**.

* `cipher_text` - (Required, String, NonUpdatable) Specifies the DEK ciphertext and metadata in hexadecimal string.
  The value is the `cipher_text` in the result of encrypting the data key.

* `datakey_cipher_length` - (Required, String, NonUpdatable) Specifies the byte length of the DEK ciphertext.
  The valid value ranges from `1` to `1024`, with a commonly used value of `64`.

* `sequence` - (Optional, String, NonUpdatable) Specifies the sequence number of the request message, `36` bytes.
  For example: **919c82d4-8046-4722-9094-35c3c6524cff**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `data_key` - The DEK plaintext in hexadecimal string.

* `datakey_length` - The byte length of DEK plaintext.

* `datakey_dgst` - The SHA256 value of the DEK plaintext in hexadecimal string.
