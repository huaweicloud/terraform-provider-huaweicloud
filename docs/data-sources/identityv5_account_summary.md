---
subcategory: "IAM"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identityv5_account_summary"
description: |-
  Use this data source to get quota information for resources under the IAM account within HuaweiCloud.
---

# huaweicloud_identityv5_account_summary

Use this data source to get quota information for resources under the IAM account within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_identityv5_account_summary" "test" {}
```

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `agencies_quota` - The total quota of agencies and trusted agencies that can be created in this account.

* `attached_policies_per_agency_quota` - The maximum number of identity policies that can be attached to
  an agency or trusted agency.

* `attached_policies_per_group_quota` - The maximum number of identity policies that can be attached to a user group.

* `attached_policies_per_user_quota` - The maximum number of identity policies that can be attached to an IAM user.

* `groups_quota` - The quota of user groups that can be created in this account.

* `policies` - The current number of custom identity policies created in this account.

* `policies_quota` - The maximum number of custom identity policies.

* `groups` - The current number of user groups created in this account.

* `policy_size_quota` - The maximum number of characters in the policy document of identity policies and
  trusted policies, excluding spaces.

* `root_user_mfa_enabled` - The number of MFA devices enabled for the root user.

* `users` - The current number of IAM users created in this account, including the root user.

* `users_quota` - The quota of IAM users that can be created in this account, including the root user.

* `versions_per_policy_quota` - The maximum number of versions that can be retained for a custom identity
  policy at the same time.

* `agencies` - The total number of agencies and trusted agencies created in this account.
