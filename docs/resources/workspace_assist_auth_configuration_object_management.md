---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_assist_auth_configuration_object_management"
description: |-
  Use this resource to manage users and user groups to associate the assist auth configuration within HuaweiCloud.
---

# huaweicloud_workspace_assist_auth_configuration_object_management

Use this resource to manage users and user groups to associate the assist auth configuration within HuaweiCloud.

## Example Usage

```hcl
variable "apply_users_configuration" {
  type = list(object({
    id   = string
    name = string
  }))
}
variable "apply_user_groups_configuration" {
  type = list(object({
    id   = string
    name = string
  }))
}

resource "huaweicloud_workspace_assist_auth_configuration_object_management" "test" {
  dynamic "objects" {
    for_each = var.apply_users_configuration

    content {
      type = "USER"
      id   = objects.value.id
      name = objects.value.name
    }
  }

  dynamic "objects" {
    for_each = var.apply_user_groups_configuration

    content {
      type = "USER_GROUP"
      id   = objects.value.id
      name = objects.value.name
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region where the apply objects of assist auth configuration are located.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `objects` - (Required, List) Specifies the list of objects to be managed with the assist auth configuration.  
  The [objects](#workspace_assist_auth_configuration_object_management_objects) structure is documented below.

<a name="workspace_assist_auth_configuration_object_management_objects"></a>
The `objects` block supports:

* `type` - (Required, String) Specifies the type of the binding object.  
  The valid values are as follows:
  + **USER** - User
  + **USER_GROUP** - User group
  + **ALL** - All users

* `id` - (Required, String) Specifies the ID of the binding object.

* `name` - (Required, String) Specifies the name of the binding object.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

Managed objects can be imported using a random UUID, e.g.

```bash
$ terraform import huaweicloud_workspace_assist_auth_configuration_object_management.test <id>
```

~> During import, all associated objects from the remote service will be included in tfstate file.
