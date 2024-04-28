---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_group_membership"
description: ""
---

# huaweicloud_identity_group_membership

Manages an IAM group membership resource within HuaweiCloud.

-> **NOTE:** You *must* have admin privileges to use this resource.

## Example Usage

```hcl
variable "user_1_password" {}
variable "user_2_password" {}

resource "huaweicloud_identity_group" "group_1" {
  name        = "group1"
  description = "This is a test group"
}

resource "huaweicloud_identity_user" "user_1" {
  name     = "user1"
  enabled  = true
  password = var.user_1_password
}

resource "huaweicloud_identity_user" "user_2" {
  name     = "user2"
  enabled  = true
  password = var.user_2_password
}

resource "huaweicloud_identity_group_membership" "membership_1" {
  group = huaweicloud_identity_group.group_1.id
  users = [
    huaweicloud_identity_user.user_1.id,
    huaweicloud_identity_user.user_2.id
  ]
}
```

## Argument Reference

The following arguments are supported:

* `group` - (Required, String, ForceNew) Specifies the group ID of this membership.

* `users` - (Required, List) Specifies a list of IAM user IDs to associate to the group.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.
