---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_kms_verify_sign"
description: |-
  Manages a KMS verify resource within HuaweiCloud.
---

# huaweicloud_kms_verify_sign

Manages a KMS verify resource within HuaweiCloud.

-> The current resource is a one-time resource, and destroying this resource will not recover verify,
  but will only remove the resource information from the tfstate file.

-> 1. This resource only supports signature verification operations for asymmetric keys with `key_usage` set to **SIGN_VERIFY**.
  <br/>2. When using SM2 keys for signature, only can be used to sign message digests.

## Example Usage

```hcl
variable "key_id" {}
variable "message" {}
variable "signature" {}
variable "signing_algorithm" {}

resource "huaweicloud_kms_verify_sign" "test" {
  key_id            = var.key_id
  message           = var.message
  signature         = var.signature
  signing_algorithm = var.signing_algorithm
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

* `message` - (Required, String, NonUpdatable) Specifies the message digest or message to be signed.
  The message length must be less than `4,096` bytes and should be encoded using Base64.

* `signature` - (Required, String, NonUpdatable) Specifies the signature value to be verified, encoded using Base64.

* `signing_algorithm` - (Required, String, NonUpdatable) Specifies the signature algorithm.
  The valid values are as follows:
  + **SM2DSA_SM3**
  + **RSASSA_PSS_SHA_256**
  + **RSASSA_PSS_SHA_384**
  + **RSASSA_PSS_SHA_512**
  + **RSASSA_PKCS1_V1_5_SHA_256**
  + **RSASSA_PKCS1_V1_5_SHA_384**
  + **RSASSA_PKCS1_V1_5_SHA_512**
  + **ECDSA_SHA_256**
  + **ECDSA_SHA_384**
  + **ECDSA_SHA_512**

* `message_type` - (Optional, String, NonUpdatable) Specifies the type of the message.
  The valid values are as follows, the default value is **DIGEST**.
  + **DIGEST**: Message digest.
  + **RAW**: Original message.

* `sequence` - (Optional, String, NonUpdatable) Specifies the sequence number of the request message, `36` bytes.
  For example: **919c82d4-8046-4722-9094-35c3c6524cff**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `signature_valid` - The validity of signature verification.
  The valid values are as follows:
  + **true**: The signature is legal.
  + **false**: The signature is illegal.
