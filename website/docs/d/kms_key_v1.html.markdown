---
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_kms_key_v1"
sidebar_current: "docs-huaweicloud-datasource-kms-key-v1"
description: |-
  Get information on an HuaweiCloud KMS Key.
---

# huaweicloud\_kms\_key\_v1

Use this data source to get the ID of an available HuaweiCloud KMS key.

## Example Usage

```hcl

data "huaweicloud_kms_key_v1" "key_1" {
  key_alias       = "test_key"
  key_description = "test key description"
  key_state       = "2"
  key_id          = "${huaweicloud_kms_key_v1.key1.id}"
}
```

## Argument Reference

* `key_alias` - (Optional) The alias in which to create the key. It is required when
    we create a new key. Changing this updates the key's alias.

* `key_description` - (Optional) The description of the key as viewed in Huawei console.
    Changing this creates a new key.

* `realm` - (Optional) Region where a key resides.

* `key_id` - (Optional) The globally unique identifier for the key.

* `default_key_flag` - (Optional) Identification of a Master Key. The value 1 indicates a Default
    Master Key, and the value 0 indicates a key.

* `key_state` - (Optional) The state of a key. "1" indicates that the key is waiting to be activated.
    "2" indicates that the key is enabled. "3" indicates that the key is disabled. "4" indicates that
    the key is scheduled for deletion.


## Attributes Reference

`id` is set to the ID of the found key. In addition, the following attributes
are exported:

* `key_alias` - See Argument Reference above.
* `key_description` - See Argument Reference above.
* `realm` - See Argument Reference above.
* `key_id` - See Argument Reference above.
* `default_key_flag` - See Argument Reference above.
* `origin` - Origin of a key. The default value is kms.
* `scheduled_deletion_date` - Scheduled deletion time (time stamp) of a key.
* `domain_id` - ID of a user domain for the key.
* `expiration_time` - Expiration time.
* `creation_date` - Creation time (time stamp) of a key.
* `key_state` - See Argument Reference above.
