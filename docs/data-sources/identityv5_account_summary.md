---
subcategory: "IAM"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identityv5_account_summary"
description: |-
  Provides account summary information for HuaweiCloud Identity and Access Management V5 service.
---

# huaweicloud_identityv5_account_summary

Provides account summary information for HuaweiCloud Identity and Access Management V5 service.

## Example Usage

```hcl
data "huaweicloud_identityv5_account_summary" "account" {}
```

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates the data source ID.

* `agencies_quota` - Indicates the total quota of agencies and trusted agencies that can be created in this account.

* `attached_policies_per_agency_quota` - Indicates the maximum number of identity policies that can be attached to
  an agency or trusted agency.

* `attached_policies_per_group_quota` - Indicates the maximum number of identity policies that can be attached to
  a user group.

* `attached_policies_per_user_quota` - Indicates the maximum number of identity policies that can be attached to an IAM user.

* `groups_quota` - Indicates the quota of user groups that can be created in this account.

* `policies` - Indicates the current number of custom identity policies created in this account.

* `policies_quota` - Indicates the maximum number of custom identity policies.

* `groups` - Indicates the current number of user groups created in this account.

* `policy_size_quota` - Indicates the maximum number of characters in the policy document of identity policies and
  trusted policies, excluding spaces.

* `root_user_mfa_enabled` - Indicates the number of MFA devices enabled for the root user.

* `users` - Indicates the current number of IAM users created in this account, including the root user.

* `users_quota` - Indicates the quota of IAM users that can be created in this account, including the root user.

* `versions_per_policy_quota` - Indicates the maximum number of versions that can be retained for a custom identity
  policy at the same time.

* `agencies` - Indicates the total number of agencies and trusted agencies created in this account.
