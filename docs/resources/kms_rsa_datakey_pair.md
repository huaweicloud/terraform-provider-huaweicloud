---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_kms_rsa_datakey_pair"
description: |-
  Manages a resource to create RSA datakey pair within HuaweiCloud.
---

# huaweicloud_kms_rsa_datakey_pair

Manages a resource to create RSA datakey pair within HuaweiCloud.

-> This resource is a one-time action resource. Deleting this resource will not clear the corresponding request record,
  but will only remove the resource information from the tf state file.

## Example Usage

```hcl
variable "key_id" {}
variable "key_spec" {}

resource "huaweicloud_kms_rsa_datakey_pair" "test" {
  key_id   = var.key_id
  key_spec = var.key_spec
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `key_id` - (Required, String, NonUpdatable) Specifies the key ID.

* `key_spec` - (Required, String, NonUpdatable) Specifies the algorithm type.
  The valid values are as follows:
  + **RSA_2048**
  + **RSA_3072**
  + **RSA_4096**

* `with_plain_text` - (Optional, Bool, NonUpdatable) Specifies whether to return the plaintext private key.
  The valid values are as follows:
  + **true** (Default value)
  + **false**

* `additional_authenticated_data` - (Optional, String, NonUpdatable) Specifies the authenticate and encrypt additional
  information.
  Please do not fill in sensitive information.

* `sequence` - (Optional, String, NonUpdatable) Specifies the sequence number of the request message, `36` bytes.
  For example: **919c82d4-8046-4722-9094-35c3c6524cff**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `public_key` - The plaintext public key information.

* `private_key_cipher_text` - The ciphertext private key.

* `private_key_plain_text` - The plaintext private key.

* `wrapped_private_key` - The ciphertext private key encrypted by a custom private key.

* `ciphertext_recipient` - The ciphertext private key encrypted by the Qingtian public key information.
