---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identityv5_agency_attached_policies"
description: |-
  Use this data source to get policy list attached to an agency within Huaweicloud.
---

# huaweicloud_identityv5_agency_attached_policies

Use this data source to get policy list attached to an agency within Huaweicloud.

## Example Usage

```hcl
variable "agency_id" {}

data "huaweicloud_identityv5_agency_attached_policies" "test" {
  agency_id = var.agency_id
}
```

## Argument Reference

The following arguments are supported:

* `agency_id` - (Required, String) Specifies the ID of the IAM agency.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `attached_policies` - The list of policies attached to the agency.  
  The [attached_policies](#IdentityV5_agency_attached_policies) structure is documented below.

<a name="IdentityV5_agency_attached_policies"></a>
The `attached_policies` block supports:

* `policy_id` - The ID of the policy.

* `policy_name` - The name of the policy.

* `attached_at` - The creation time of the policy.

* `urn` - The Uniform Resource Name (URN) of the policy.
