---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identityv5_agency_attached_policies"
description: |-
  Use this data source to get attached policies for a specified IAM V5 agency.
---

# huaweicloud_identityv5_agency_attached_policies

Use this data source to get attached policies for a specified IAM V5 agency.

## Example Usage

```hcl
variable agency_id {}

data "huaweicloud_identityv5_agency_attached_policies" "test" {
  agency_id = var.agency_id
}
```

## Argument Reference

The following arguments are supported:

* `agency_id` - (Required, String) Specifies the ID of the IAM agency.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `attached_policies` - Indicate the attached policies list.
  The [attached_policies](#IdentityV5_agency_attached_policies) structure is documented below.

<a name="IdentityV5_agency_attached_policies"></a>
The `attached_policies` block contains:

* `attached_at` - Indicate the time when the policy was attached.

* `policy_id` - Indicate the ID of the policy.

* `policy_name` - Indicate the name of the policy.

* `urn` - Indicate the Uniform Resource Name (URN) of the policy.
