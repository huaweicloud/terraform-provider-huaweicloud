---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identityv5_policy_group_attach"
description: |-
  Use this resource to attach an identity policy to an IAM user group within HuaweiCloud.
---

# huaweicloud_identityv5_policy_group_attach

Use this resource to attach an identity policy to an IAM user group within HuaweiCloud.

## Example Usage

```hcl
variable "policy_id" {}
variable "group_id" {}

resource "huaweicloud_identityv5_policy_group_attach" "test" {
  policy_id = var.policy_id
  group_id  = var.group_id
}
```

## Argument Reference

The following arguments are supported:

* `policy_id` - (Required, String, NonUpdatable) Specifies the ID of the identity policy to be attached.

* `group_id` - (Required, String, NonUpdatable) Specifies the ID of the IAM user group associated with the identity policy.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `policy_name` - The name of the identity policy associated with the user group.

* `urn` - The URN of the attached identity policy.

* `attached_at` - The time when the identity policy was attached to the user group, in RFC3339 format.

## Import

The resource can be imported using the `policy_id` and `group_id`, separated by a slash (/), e.g.

```bash
$ terraform import huaweicloud_identityv5_policy_group_attach.test <policy_id>/<group_id>
```
