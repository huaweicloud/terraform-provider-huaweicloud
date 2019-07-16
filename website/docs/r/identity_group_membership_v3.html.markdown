---
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_group_membership_v3"
sidebar_current: "docs-huaweicloud-resource-identity-group-membership-v3"
description: |-
  Manages the membership combine User Group resource and User resource  within
  HuaweiCloud IAM service.
---

# huaweicloud\_identity\_group_membership_v3

Manages a User Group Membership resource within HuaweiCloud IAM service.

Note: You _must_ have admin privileges in your HuaweiCloud cloud to use
this resource.

## Example Usage

```hcl
resource "huaweicloud_identity_group_v3" "group_1" {
  name        = "group1"
  description = "This is a test group"
}

resource "huaweicloud_identity_user_v3" "user_1" {
  name     = "user1"
  enabled  = true
  password = "password12345!"
}

resource "huaweicloud_identity_user_v3" "user_2" {
  name     = "user2"
  enabled  = true
  password = "password12345!"
}

resource "huaweicloud_identity_group_membership_v3" "membership_1" {
  group = "${huaweicloud_identity_group_v3.group_1.id}"
  users = ["${huaweicloud_identity_user_v3.user_1.id}",
    "${huaweicloud_identity_user_v3.user_2.id}"
  ]
}
```

## Argument Reference

The following arguments are supported:

* `group` - (Required) The group ID of this membership. 

* `users` - (Required) A List of user IDs to associate to the group.

## Attributes Reference

The following attributes are exported:

* `group` - See Argument Reference above.

* `users` - See Argument Reference above.

