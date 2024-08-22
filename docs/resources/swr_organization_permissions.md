---
subcategory: "Software Repository for Container (SWR)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_swr_organization_permissions"
description: ""
---

# huaweicloud_swr_organization_permissions

Manages user permissions for the SWR organization resource within HuaweiCloud.

## Example Usage

```hcl
variable "organization_name" {}
variable "user_1" {}
variable "user_2" {}

resource "huaweicloud_swr_organization_permissions" "test" {
  organization = var.organization_name

  users {
    user_name  = var.user_1.name
    user_id    = var.user_1.id
    permission = "Read"
  }

  users {
    user_name  = var.user_2.name
    user_id    = var.user_2.id
    permission = "Read"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource. If omitted, the
  provider-level region will be used. Changing this creates a new resource.

* `organization` - (Required, String, ForceNew) Specifies the name of the organization (namespace) to be accessed.
  Changing this creates a new resource.

* `users` - (Required, List) Specifies the users to access to the organization (namespace).
  Structure is documented below.

The `users` block supports:

* `user_id` - (Required, String) Specifies the ID of the existing HuaweiCloud user.

* `user_name` - (Optional, String) Specifies the name of the existing HuaweiCloud user.

* `permission` - (Required, String) Specifies the permission of the existing HuaweiCloud user.
  The values can be **Manage**, **Write** and **Read**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the permissions. The value is the name of the organization.

* `creator` - The creator user name of the organization.

* `self_permission` - The permission information of current user.

The `self_permission` block supports:

* `user_name` - The name of current user.

* `user_id` - The ID of current user.

* `permission` - The permission of current user.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minutes.
* `delete` - Default is 5 minutes.

## Import

Organization Permissions can be imported using the `id` (organization name), e.g.

```bash
$ terraform import huaweicloud_swr_organization_permissions.test terraform-test
```
