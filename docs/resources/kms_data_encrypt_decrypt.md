---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_kms_data_encrypt_decrypt"
description: -|
  Manages a data encryption or decryption resource within HuaweiCloud.
---

# huaweicloud_kms_data_encrypt_decrypt

Manages a data encryption or decryption resource within HuaweiCloud.

-> Destroying this resource will not change the current status.

-> If you use an asymmetric key to encrypt data, record the selected key ID and encryption algorithm. To decrypt data,
  you need to provider the same key ID and encryption algorithm, or the decryption will fail.

## Example Usage

### encrypt data

```hcl
variable "key_id" {}
variable "plain_text" {}

resource "huaweicloud_kms_data_encrypt_decrypt" "test" {
  key_id     = var.key_id
  action     = "encrypt"
  plain_text = var.plain_text
}
```

### decrypt data

```hcl
variable "cipher_text" {}

resource "huaweicloud_kms_data_encrypt_decrypt" "test" {
  action      = "decrypt"
  cipher_text = var.cipher_text
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `key_id` - (Optional, String, ForceNew) Specifies the key ID.
  Changing this will create a new resource.

  -> This parameter is mandatory for encryption operation.
  <br/>If the ciphertext is encrypted using an asymmetric key, this parameter must be specified
  for decryption operation.

* `action` - (Required, String, ForceNew) Specifies the operation type.
  Changing this will create a new resource.
  The valid values are as follow:
  + **encrypt**
  + **decrypt**

* `encryption_algorithm` - (Optional, String, ForceNew) Specifies the data encryption algorithm.
  Changing this will create a new resource.
  The valid values are as follow, the default value is **SYMMETRIC_DEFAULT**.
  + **SYMMETRIC_DEFAULT**
  + **RSAES_OAEP_SHA_256**
  + **SM2_ENCRYPT**

  -> If the key used is an asymmetric key, this parameter must be specified for encryption or decryption operation.

* `plain_text` - (Optional, String, ForceNew) Specifies the plaintext to be encrypted.
  Changing this will create a new resource.
  It must be `1` to `4,096` bytes long.

  -> This parameter is mandatory for encryption operation.

* `cipher_text` - (Optional, String, ForceNew) Specifies the ciphertext to be decrypted.
  Changing this will create a new resource.
  It must be `188` to `5,648` characters long.

  -> This parameter is mandatory for decryption operation.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `cipher_data` - The encrypted ciphertext, encoding Base64.

* `plain_data` - The decrypted plaintext.

* `plain_text_base64` - The Base64 value in plaintext.
  In asymmetric encryption, if the encrypted plaintext contains invisible characters, the Base64 value will be
  used as the decryption result.
