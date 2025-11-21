---
subcategory: "IAM Identity Center"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identitycenter_group_memberships"
description: |-
  Use this data source to get the Identity Center group memberships.
---

# huaweicloud_identitycenter_group_memberships

Use this data source to get the Identity Center group memberships.

## Example Usage

```hcl
variable "identity_store_id" {}
variable "group_id" {}

data "huaweicloud_identitycenter_group_memberships" "test"{
  identity_store_id = var.identity_store_id
  group_id          = var.group_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `identity_store_id` - (Required, String) Specifies the ID of the identity store that associated with IAM Identity
  Center.

* `group_id` - (Optional, String) Specifies the ID of the group.

* `user_id` - (Optional, String) Specifies the ID of the user.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `group_memberships` - The list of group membership.
  The [group_memberships](#group_memberships) structure is documented below.

<a name="group_memberships"></a>
The `group_memberships` block supports:

* `group_id` - The ID of the group.

* `member_id` - The ID of the member.
  The [member_id](#member_id) structure is documented below.

* `identity_store_id` - The ID of the identity store.

* `membership_id` - The ID of the membership.

<a name="member_id"></a>
The `member_id` block supports:

* `user_id` - The ID of the user.
