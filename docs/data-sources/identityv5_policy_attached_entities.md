---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identityv5_policy_attached_entities"
description: |-
  Use this data source to query the entities attached to an IAM policy.
---

# huaweicloud_identityv5_policy_attached_entities

Use this data source to query the entities attached to an IAM policy.

## Example Usage

```hcl
variable policy_id {}

data "huaweicloud_identityv5_policy_attached_entities" "test" {
  policy_id = var.policy_id
}
```

## Argument Reference

The following arguments are supported:

* `policy_id` - (Required, String) Specifies the ID of the IAM policy.

* `entity_type` - (Optional, String) Specifies the type of the entity, can be `user`, `group` or `agency`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `policy_users` - Indicates the list of IAM users attached to the policy.
  The [policy_users](#IdentityPolicyUsers_List) structure is documented below.

* `policy_groups` - Indicates the list of user groups attached to the policy.
  The [policy_groups](#IdentityPolicyGroups_List) structure is documented below.

* `policy_agencies` - Indicates the list of agencies attached to the policy.
  The [policy_agencies](#IdentityPolicyAgencies_List) structure is documented below.

<a name="IdentityPolicyUsers_List"></a>
The `policy_users` block supports:

* `user_id` - Indicates the ID of the IAM user.

* `attached_at` - Indicates the time when the user was attached to the policy.

<a name="IdentityPolicyGroups_List"></a>
The `policy_groups` block supports:

* `group_id` - Indicates the ID of the user group.

* `attached_at` - Indicates the time when the group was attached to the policy.

<a name="IdentityPolicyAgencies_List"></a>
The `policy_agencies` block supports:

* `agency_id` - Indicates the ID of the agency.

* `attached_at` - Indicates the time when the agency was attached to the policy.
