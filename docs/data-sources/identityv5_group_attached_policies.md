---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identityv5_group_attached_policies"
description: |-
  Use this data source to get policy list attached to a group within Huaweicloud.
---
 
# huaweicloud_identityv5_group_attached_policies

Use this data source to get policy list attached to a group within Huaweicloud.

## Example Usage

```hcl
variable "group_id" {}

data "huaweicloud_identityv5_group_attached_policies" "test" {
  group_id = var.group_id
}
```

## Argument Reference

The following arguments are supported:

* `group_id` - (Required, String) Specifies the ID of the IAM user group.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:、

* `id` - The data source ID.

* `attached_policies` - Indicate the list of policies attached to the user group.  
  The [attached_policies](#IdentityV5Attached_policies) structure is documented below.

<a name="IdentityV5Attached_policies"></a>
The `attached_policies` block supports:

* `policy_id` - The ID of the policy.

* `policy_name` - The name of the policy.

* `attached_at` - The creation time of the policy.

* `urn` - The Uniform Resource Name (URN) of the policy.
