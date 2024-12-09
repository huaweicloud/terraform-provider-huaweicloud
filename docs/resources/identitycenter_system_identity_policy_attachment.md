---
subcategory: "IAM Identity Center"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identitycenter_system_identity_policy_attachment"
description: |-
  Manages an Identity Center system identity policy attachment resource within HuaweiCloud.
---

# huaweicloud_identitycenter_system_identity_policy_attachment

Manages an Identity Center system identity policy attachment resource within HuaweiCloud.

-> **NOTE:** Creating this resource will automatically provision the Permission Set to apply the corresponding updates
  to all assigned accounts.

## Example Usage

```hcl
variable "instance_id" {}
variable "permission_set_id" {}
variable "iam_identity_policy_ids" {}

resource "huaweicloud_identitycenter_system_identity_policy_attachment" "test" {
  instance_id       = var.instance_id
  permission_set_id = var.permission_set_id
  policy_ids        = var.iam_identity_policy_ids
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the IAM Identity Center instance.
  Changing this creates a new resource.

* `permission_set_id` - (Required, String, ForceNew) Specifies the ID of the IAM Identity Center permission set.
  Changing this creates a new resource.

* `policy_ids` - (Required, List) Specifies an array of IAM managed system identity policies to be attached to
  the permission set.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `attached_policies` - All IAM managed system identity policies attached to the permission set.
  The [attached_policies](#attached_policies) structure is documented below.

<a name="attached_policies"></a>
The `attached_policies` block supports:

* `id` - The ID of an IAM system identity policy.

* `name` - The name of an IAM system identity policy.

## Import

The Identity Center system identity policy attachment can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_identitycenter_system_identity_policy_attachment.test <id>
```
