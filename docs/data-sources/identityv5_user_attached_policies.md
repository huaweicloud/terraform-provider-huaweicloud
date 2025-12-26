---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identityv5_user_attached_policies"
description: |-
  Use this data source to get attached policies for a specified IAM V5 user.
---

# huaweicloud_identityv5_user_attached_policies

Use this data source to get attached policies for a specified IAM V5 user.

## Example Usage

```hcl
variable user_id {}

data "huaweicloud_identityv5_user_attached_policies" "test" {
  user_id = var.user_id
}
```

## Argument Reference

The following arguments are supported:

* `user_id` - (Required, String) Specifies the ID of the IAM user.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `attached_policies` - Indicate the attached policies of the IAM user.
  The [attached_policies](#IdentityV5Attached_policies) structure is documented below.

<a name="IdentityV5Attached_policies"></a>
The `attached_policies` block contains:

* `attached_at` - Indicate the time when the policy was attached.

* `policy_id` - Indicate the ID of the policy.

* `policy_name` - Indicate the name of the policy.

* `urn` - Indicate the Uniform Resource Name (URN) of the policy.
