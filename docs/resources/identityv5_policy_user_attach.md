---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identityv5_policy_user_attach"
description: |-
  Use this resource to attach an identity policy to an IAM user within HuaweiCloud.
---

# huaweicloud_identityv5_policy_user_attach

Use this resource to attach an identity policy to an IAM user within HuaweiCloud.

## Example Usage

```hcl
variable "policy_id" {}
variable "user_id" {}

resource "huaweicloud_identityv5_policy_user_attach" "test" {
  policy_id = var.policy_id
  user_id   = var.user_id
}
```

## Argument Reference

The following arguments are supported:

* `policy_id` - (Required, String, NonUpdatable) Specifies the ID of the identity policy to be attached.

* `user_id` - (Required, String, NonUpdatable) Specifies the ID of the IAM user associated with the identity policy.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `policy_name` - The name of the identity policy associated with the user.

* `urn` - The URN of the attached identity policy.

* `attached_at` - The time when the identity policy was attached to the user, in RFC3339 format.

## Import

The resource can be imported using the `policy_id` and `user_id`, separated by a slash (/), e.g.

```bash
$ terraform import huaweicloud_identityv5_policy_user_attach.test <policy_id>/<user_id>
```
