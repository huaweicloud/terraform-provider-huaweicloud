---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identityv5_access_key"
description: |-
  Use this data source to get the first available access key under the specified user within HuaweiCloud.
---

# huaweicloud_identityv5_access_key

Use this data source to get the first available access key under the specified user within HuaweiCloud.

## Example Usage

```hcl
variable "user_id" {}

data "huaweicloud_identityv5_access_key" "test" {
  user_id = var.user_id
}
```

## Argument Reference

The following arguments are supported:

* `user_id` - (Required, String) Specifies the ID of the IAM user.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the data source.

* `access_key_id` - The ID of the access key.

* `status` - The status of the access key.  
  + **active**
  + **inactive**

* `created_at` - The creation time of the access key.

* `last_used_at` - The time when the access key was last used.
