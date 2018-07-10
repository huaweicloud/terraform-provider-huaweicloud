---
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_kms_key_v1"
sidebar_current: "docs-huaweicloud-datasource-kms-key-v1"
description: |-
  Get information on an HuaweiCloud KMS Key.
---

# huaweicloud\_kms\_key_v1

Use this data source to get the ID of an available HuaweiCloud KMS key.

## Example Usage

```hcl

data "huaweicloud_kms_key_v1" "key_1" {
  key_alias        = "test_key"
  key_description  = "test key description"
  key_state        = "2"
  key_id           = "af650527-a0ff-4527-aef3-c493df1f3012"
  realm            = "cn-north-1"
  default_key_flag = "0"
  domain_id        = "b168fe00ff56492495a7d22974df2d0b"
}
```

## Argument Reference

* `key_alias` - (Optional) The alias in which to create the key. It is required when
    we create a new key. Changing this gets the new key.

* `key_description` - (Optional) The description of the key as viewed in Huawei console.
    Changing this gets a new key.

* `realm` - (Optional) Region where a key resides. Changing this gets a new key.

* `key_id` - (Optional) The globally unique identifier for the key. Changing this gets the new key.

* `default_key_flag` - (Optional) Identification of a Master Key. The value "1" indicates a Default
    Master Key, and the value "0" indicates a key. Changing this gets a new key.

* `key_state` - (Optional) The state of a key. "1" indicates that the key is waiting to be activated.
    "2" indicates that the key is enabled. "3" indicates that the key is disabled. "4" indicates that
    the key is scheduled for deletion. Changing this gets a new key.

* `domain_id` - (Optional)  - ID of a user domain for the key. Changing this gets a new key.


## Attributes Reference

`id` is set to the ID of the found key. In addition, the following attributes
are exported:

* `key_alias` - See Argument Reference above.
* `key_description` - See Argument Reference above.
* `realm` - See Argument Reference above.
* `key_id` - See Argument Reference above.
* `default_key_flag` - See Argument Reference above.
* `scheduled_deletion_date` - Scheduled deletion time (time stamp) of a key.
* `domain_id` - See Argument Reference above.
* `expiration_time` - Expiration time.
* `creation_date` - Creation time (time stamp) of a key.
* `key_state` - See Argument Reference above.
