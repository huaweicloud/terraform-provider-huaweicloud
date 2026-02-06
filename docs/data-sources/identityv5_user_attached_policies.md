---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identityv5_user_attached_policies"
description: |-
  Use this data source to get policy list attached to a user within Huaweicloud.
---

# huaweicloud_identityv5_user_attached_policies

Use this data source to get policy list attached to a user within Huaweicloud.

## Example Usage

```hcl
variable "user_id" {}

data "huaweicloud_identityv5_user_attached_policies" "test" {
  user_id = var.user_id
}
```

## Argument Reference

The following arguments are supported:

* `user_id` - (Required, String) Specifies the ID of the IAM user.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `attached_policies` - The list of policies attached to the user.  
  The [attached_policies](#v5_user_attached_policies) structure is documented below.

<a name="v5_user_attached_policies"></a>
The `attached_policies` block supports:

* `policy_id` - The ID of the policy.

* `policy_name` - The name of the policy.

* `attached_at` - The creation time of the policy.

* `urn` - The Uniform Resource Name (URN) of the policy.
