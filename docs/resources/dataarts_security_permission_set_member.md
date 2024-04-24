---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_security_permission_set_member"
description: |-
  Manages DataArts Security permission set member resource within HuaweiCloud.
---

# huaweicloud_dataarts_security_permission_set_member

Manages DataArts Security permission set member resource within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}
variable "permission_set_id" {}
variable "object_name" {}
variable "object_id" {}

resource "huaweicloud_dataarts_security_permission_set_member" "test" {
  workspace_id      = var.workspace_id
  permission_set_id = var.permission_set_id
  object_id         = var.object_id
  name              = var.object_name
  type              = "USER"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `workspace_id` - (Required, String, ForceNew) Specifies the workspace ID to which the permission set and member belongs.
  Changing this creates a new resource.

* `permission_set_id` - (Required, String, ForceNew) Specifies the permission set ID to which the member belongs.
  Changing this creates a new resource.

* `object_id` - (Required, String, ForceNew) Specifies the ID of the member object. The valid value ranges from `1` to `128`.
  Changing this creates a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the member object.
  Changing this creates a new resource.

* `type` - (Required, String, ForceNew) Specifies the type of the member object.
  The valid values are **USER**and **USER_GROUP**.
  Changing this creates a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `member_id` - The ID of the member.

## Import

The permission set member can be imported using `workspace_id`, `permission_set_id` and `object_id`, separated by
slashes (/), e.g.

```bash
$ terraform import huaweicloud_dataarts_security_permission_set_member.test <workspace_id>/<permission_set_id>/<object_id>
```
