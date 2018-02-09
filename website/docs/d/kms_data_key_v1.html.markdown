---
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_kms_data_key_v1"
sidebar_current: "docs-huaweicloud-datasource-kms-data-key-v1"
description: |-
  Get information on an HuaweiCloud KMS data encryption key.
---

# huaweicloud\_kms\_data_key_v1

Use this data source to get the plaintext and the ciphertext of an available
HuaweiCloud KMS DEK (data encryption key).

## Example Usage

```hcl

resource "huaweicloud_kms_key_v1" "key1" {
  key_alias       = "key_1"
  pending_days    = "7"
  key_description = "first test key"
}

data "huaweicloud_kms_data_key_v1" "kms_datakey1" {
  key_id         = "${huaweicloud_kms_key_v1.key1.id}"
  datakey_length = "512"
}

```

## Argument Reference

* `key_id` - (Required) The globally unique identifier for the key.
    Changing this gets the new data encryption key.

* `encryption_context` - (Optional) The value of this parameter must be a series of
    "key:value" pairs used to record resource context information. The value of this
    parameter must not contain sensitive information and must be within 8192 characters
    in length. Example: {"Key1":"Value1","Key2":"Value2"}

* `datakey_length` - (Required) Number of bits in the length of a DEK (data encryption keys).
    The maximum number is 512. Changing this gets the new data encryption key.


## Attributes Reference

`id` is set to the date of the found data key. In addition, the following attributes
are exported:

* `plain_text` - The plaintext of a DEK is expressed in hexadecimal format, and two
    characters indicate one byte.
* `cipher_text` - The ciphertext of a DEK is expressed in hexadecimal format, and two
    characters indicate one byte.
