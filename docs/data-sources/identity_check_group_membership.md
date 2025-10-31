---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_check_group_membership"
description: |-
    Use this data source to query whether an IAM user is in a group within HuaweiCloud.
---

# huaweicloud_identity_check_group_membership

Use this data source to query whether an IAM user is in a group within HuaweiCloud.

-> **NOTE:** You *must* have IAM read privileges to use this data source.

## Example Usage

```hcl
variable "group_id" {}
variable "user_id" {}

data "huaweicloud_identity_check_group_membership" "test" {
  group_id = var.group_id
  user_id  = var.user_id
}
```

## Argument Reference

* `group_id` - (Required, String) Specifies the group id.

* `user_id` - (Required, String) Specifies the user id.

## Attribute Reference

* `result` - Indicates whether the user is in the group.
