---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_kms_key"
description: ""
---

# huaweicloud_kms_key

Manages a KMS key resource within HuaweiCloud.

## Example Usage

```hcl
resource "huaweicloud_kms_key" "key_1" {
  key_alias       = "key_1"
  pending_days    = "7"
  key_description = "first test key"
  is_enabled      = true
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the KMS key resource. If omitted, the
  provider-level region will be used. Changing this creates a new KMS key resource.

* `key_alias` - (Required, String) The alias in which to create the key. It is required when we create a new key.
  Changing this updates the alias of key.

* `key_description` - (Optional, String) The description of the key as viewed in Huawei console. Changing this updates
  the description of key.

* `key_algorithm` - (Optional, String, ForceNew) The algorithm of the key. Valid values are AES_256, SM4, RSA_2048, RSA_3072,
  RSA_4096, EC_P256, EC_P384, SM2. Changing this creates a new key.

* `pending_days` - (Optional, String) Duration in days after which the key is deleted after destruction of the resource,
  must be between 7 and 1096 days. It doesn't have default value. It only be used when delete a key.

* `is_enabled` - (Optional, Bool) Specifies whether the key is enabled. Defaults to true. This field is not supported
  when creating an external import key for the first time. This field only takes effect when the value
  of `key_state` is **2** or **3**. Changing this updates the state of existing key.

* `rotation_enabled` - (Optional, Bool) Specifies whether the key rotation is enabled. Defaults to false.
  This field is supported when the origin of the key is **kms**.

* `rotation_interval` - (Optional, Int) Specifies the key rotation interval. The valid value is range from 30 to 365,
  defaults to 365. This field is supported when the source of the key is **kms**.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID to which the KMS key belongs.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the kms key.

* `origin` - (Optional, String, ForceNew) Specifies the source of the kms key. Valid values are **kms** and **external**.
  The default value is **kms**. Changing this creates a new key.

* `key_usage` - (Optional, String, ForceNew) Specifies the key usage. The value can be **ENCRYPT_DECRYPT** or **SIGN_VERIFY**.
  For key_usage selection, see the [documentation](https://support.huaweicloud.com/intl/en-us/productdesc-ram/ram_01_0007.html).
  Changing this creates a new key.

* `keystore_id` - (Optional, String, ForceNew) Specifies the keystore ID of the kms key.
  Default value is the KMS default keystore ID. Changing this creates a new key.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies a resource ID in UUID format.
* `key_id` - The globally unique identifier for the key.
* `default_key_flag` - Identification of a Master Key. The value 1 indicates a Default Master Key, and the value 0
  indicates a key.
* `scheduled_deletion_date` - Scheduled deletion time (time stamp) of a key.
* `domain_id` - ID of a user domain for the key.
* `expiration_time` - Expiration time.
* `creation_date` - Creation time (time stamp) of a key.
* `rotation_number` - The total number of key rotations.
* `key_state` - The status of the kms key. The valid values are as follows:
  **1**: To be activated
  **2**: Enabled.
  **3**: Disabled.
  **4**: Pending deletion.
  **5**: Pending import.

## Import

KMS Keys can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_kms_key.key_1 7056d636-ac60-4663-8a6c-82d3c32c1c64
```

Note that the imported state may not be identical to your resource definition,
due to `pending_days` and `is_enabled` are missing from the API response.
It is generally recommended running `terraform plan` after importing a KMS Key.
You can then decide if changes should be applied to the KMS Key, or the resource
definition should be updated to align with the KMS Key. Also you can ignore changes as below.

```hcl
resource "huaweicloud_kms_key" "key_1" {
    ...

  lifecycle {
    ignore_changes = [ pending_days, is_enabled ]
  }
}
```
