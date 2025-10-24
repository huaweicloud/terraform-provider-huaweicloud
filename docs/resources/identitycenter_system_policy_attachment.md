---
subcategory: "IAM Identity Center"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identitycenter_system_policy_attachment"
description: ""
---

# huaweicloud_identitycenter_system_policy_attachment

Manages an Identity Center system policy attachment resource within HuaweiCloud.  

-> **NOTE:** Creating this resource will automatically Provision the Permission Set to apply the corresponding updates
  to all assigned accounts.

## Example Usage

```hcl
variable "permission_set_id" {}
variable "iam_policy_ids" {
  type = list(string)
}

data "huaweicloud_identitycenter_instance" "system" {}

resource "huaweicloud_identitycenter_system_policy_attachment" "test" {
  instance_id       = data.huaweicloud_identitycenter_instance.system.id
  permission_set_id = var.permission_set_id
  policy_ids        = var.iam_policy_ids
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the IAM Identity Center instance.
  Changing this parameter will create a new resource.

* `permission_set_id` - (Required, String, ForceNew) Specifies the ID of the IAM Identity Center permission set.
  Changing this parameter will create a new resource.

* `policy_ids` - (Required, List) Specifies an array of IAM managed system policies/roles to be attached to
  the permission set.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `attached_policies` - All IAM managed system policies/roles attached to the permission set.
  The [object](#AttachedManagedPolicies) structure is documented below.

<a name="AttachedManagedPolicies"></a>
The `attached_policies` block supports:

* `id` - The ID of an IAM system policy/role.

* `name` - The name of an IAM system policy/role.

## Import

The Identity Center system policy attachment can be imported using the `instance_id` and `permission_set_id` separated
by a slash, e.g.

```bash
$ terraform import huaweicloud_identitycenter_system_policy_attachment.test <instance_id>/<permission_set_id>
```
