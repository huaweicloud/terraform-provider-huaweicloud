---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_kms_encrypt_datakey"
description: |-
  Manages a KMS encrypt datakey resource within HuaweiCloud.
---

# huaweicloud_kms_encrypt_datakey

Manages a KMS encrypt datakey resource within HuaweiCloud.

-> The current resource is a one-time resource, and destroying this resource will not recover the encrypted datakey,
but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "key_id" {}
variable "plain_text" {}
variable "datakey_plain_length" {}

resource "huaweicloud_kms_encrypt_datakey" "test" {
  key_id               = var.key_id
  plain_text           = var.plain_text
  datakey_plain_length = var.datakey_plain_length
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

* `plain_text` - (Required, String, NonUpdatable) Specifies the plaintext of data encryption key.
  When CMK is AES, the plaintext of DEK and SHA256 hash value (`32` bytes) of the plaintext; When CMK is SM4,
  the plaintext of DEK and SM3 hash value (`32` bytes) of the plaintext, both represented as hexadecimal strings.

* `datakey_plain_length` - (Required, String, NonUpdatable) Specifies the byte length of the DEK plaintext.
  The vaild value ranges from `1` to `1024`, with a commonly used value of `64`.

* `sequence` - (Optional, String, NonUpdatable) Specifies the sequence number of the request message, `36` bytes.
  For example: **919c82d4-8046-4722-9094-35c3c6524cff**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `cipher_text` - The DEK ciphertext in hexadecimal, two characters represent `1` byte.
