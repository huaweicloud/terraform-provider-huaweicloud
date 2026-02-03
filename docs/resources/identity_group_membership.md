---
subcategory: "Identity and Access Management (IAM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_identity_group_membership"
description: |-
  Manages an IAM group membership resource within HuaweiCloud.
---

# huaweicloud_identity_group_membership

Manages an IAM group membership resource within HuaweiCloud.

-> **NOTE:** You **must** have admin privileges to use this resource.

## Example Usage

```hcl
variable "group_id" {}
variable "user_ids" {
  type = list(string)
}

resource "huaweicloud_identity_group_membership" "test" {
  group = var.group_id
  users = var.user_ids
}
```

## Argument Reference

The following arguments are supported:

* `group` - (Required, String, NonUpdatable) Specifies the ID of the group to which the users belong.

* `users` - (Required, List) Specifies the list of user IDs associated with the group.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

## Import

The memberships of group can be imported using group ID (`id`), e.g.

```bash
$ terraform import huaweicloud_identity_agency.test <id>
```

~> During the import process, all memberships managed remotely will be synchronized to the `terraform.tfstate` file.
