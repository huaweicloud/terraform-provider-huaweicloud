---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identityv5_policy_agency_attach"
description: |-
  Use this resource to attach an IAM V5 policy to an IAM agency within HuaweiCloud.
---

# huaweicloud_identityv5_policy_agency_attach

Use this resource to attach an IAM V5 policy to an IAM agency within HuaweiCloud.

## Example Usage

```hcl
variable "policy_id" {}
variable "agency_id" {}

resource "huaweicloud_identityv5_policy_agency_attach" "test" {
  policy_id = var.policy_id
  agency_id = var.agency_id
}
```

## Argument Reference

The following arguments are supported:

* `policy_id` - (Required, String, NonUpdatable) Specifies the ID of the IAM V5 policy to be attached.

* `agency_id` - (Required, String, NonUpdatable) Specifies the ID of the IAM agency associated with the IAM V5 policy.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `policy_name` - The name of the IAM V5 policy associated with the agency.

* `policy_urn` - The URN of the attached IAM V5 policy.

* `attached_at` - The time when the IAM V5 policy was attached to the agency, in RFC3339 format.

## Import

The resource can be imported using the `policy_id` and `agency_id`, separated by a slash (/), e.g.

```bash
$ terraform import huaweicloud_identityv5_policy_agency_attach.test <policy_id>/<agency_id>
```
