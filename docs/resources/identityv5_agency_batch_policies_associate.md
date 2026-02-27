---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identityv5_agency_batch_policies_associate"
description: |-
  Using this resource to associate IAM v5 identity policies to an agency within HuaweiCloud.
---

# huaweicloud_identityv5_agency_batch_policies_associate

Using this resource to associate IAM v5 identity policies to an agency within HuaweiCloud.

-> You **must** have admin privileges to use this resource.

## Example Usage

```hcl
variable "agency_id" {}
variable "policy_ids" {
  type = list(string)
}

resource "huaweicloud_identityv5_agency_batch_policies_associate" "test" {
  agency_id = huaweicloud_identity_agency.test.id

  dynamic "policies" {
    for_each = 

    content {
      id = policies.value
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `agency_id` - (Required, String, NonUpdatable) Specifies the ID of the IAM agency to which the policies will be attached.

* `policies` - (Required, List) Specifies the list of policies to be attached to the agency.  
  The [policies](#iam_v5_agency_batch_associated_policies) structure is documented below.

<a name="iam_v5_agency_batch_associated_policies"></a>
The `policies` block supports:

* `id` - (Required, String) Specifies the ID of the identity policy.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID (also agency ID).

* `policies` - The list of policies attached to the agency.  
  The [policies](#iam_v5_agency_batch_associated_policies_attr) structure is documented below.

<a name="iam_v5_agency_batch_associated_policies_attr"></a>
The `policies` block supports:

* `name` - The name of the identity policy.

* `urn` - The URN of the identity policy.

* `attached_at` - The attached time of the identity policy, in RFC3339 format.

## Import

Associated policies under agency can be imported using the `agency_id`, e.g.

```bash
$ terraform import huaweicloud_identityv5_agency_batch_policies_associate.test <agency_id>
```

~> During the import process, all associated policies managed remotely will be synchronized to the
   `terraform.tfstate` file.
