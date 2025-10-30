---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_kms_verify_mac"
description: |-
  Manages a resource to verify a MAC within HuaweiCloud.
---

# huaweicloud_kms_verify_mac

Manages a resource to verify a MAC within HuaweiCloud.

-> This resource is a one-time action resource. Deleting this resource will not clear the corresponding request record,
  but will only remove the resource information from the tf state file.

-> This resource only supports KMS keys with `key_usage` set to **GENERATE_VERIFY_MAC**.

## Example Usage

```hcl
variable "key_id" {}
variable "mac_algorithm" {}
variable "message" {}
variable "mac" {}

resource "huaweicloud_kms_verify_mac" "test" {
  key_id        = var.key_id
  mac_algorithm = var.mac_algorithm
  message       = var.message
  mac           = var.mac
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `key_id` - (Required, String, NonUpdatable) Specifies the key ID.

* `mac_algorithm` - (Required, String, NonUpdatable) Specifies the MAC algorithm.
  The valid values are as follows:
  + **HMAC_SHA_256**
  + **HMAC_SHA_384**
  + **HMAC_SHA_512**
  + **HMAC_SM3** (Only supported in the Chinese region)

  -> The `mac_algorithm` value should be consistent with the value of KMS key `key_algorithm`.

* `message` - (Required, String, NonUpdatable) Specifies the message to be processed.
  The valid length from `1` to `4,096` characters and should be encoded using Base64.

* `mac` - (Required, String, NonUpdatable) Specifies the MAC to be verified.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `mac_valid` - The MAC verification result.
