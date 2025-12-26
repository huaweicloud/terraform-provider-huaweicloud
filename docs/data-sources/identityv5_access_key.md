---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identityv5_access_key"
description: |-
  Use this data source to get the list of access key in the Identity and Access Management V5 service.
---

# huaweicloud_identityv5_access_key

Use this data source to get the list of access key in the Identity and Access Management V5 service.

## Example Usage

```hcl
resource "huaweicloud_identityv5_user" "user_1" {
  name = "Test_accessKey"
}

resource "huaweicloud_identityv5_access_key" "key_1" {
  user_id = huaweicloud_identityv5_user.user_1.id
  status  = "inactive"
}

data "huaweicloud_identityv5_access_key" "test" {
  user_id = huaweicloud_identityv5_access_key.key_1.user_id
}
```

## Argument Reference

The following arguments are supported:

* `user_id` - (Required, String) Specifies the ID of the IAM user.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `access_key_id` - Indicates the ID of the access key.

* `status` - Indicates the status of the access key. The value can be `active` or `inactive`.

* `created_at` - Indicates the time when the access key was created.

* `last_used_at` - Indicates the time when the access key was last used.
