---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_assist_auth_configuration_object_batch_apply"
description: |-
  Manages a Workspace assist auth configuration object batch apply resource within HuaweiCloud.
---

# huaweicloud_workspace_assist_auth_configuration_object_batch_apply

Manages a Workspace assist auth configuration object batch apply resource within HuaweiCloud.

-> This resource is a one-time action resource used to update assist auth configuration apply objects. Deleting this
   resource will not clear the corresponding request record, but will only remove the resource information from
   the tfstate file.

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

resource "huaweicloud_workspace_assist_auth_configuration_object_batch_apply" "test" {
  dynamic "add" {
    for_each = var.apply_users_configuration

    content {
      object_type = "USER"
      object_id   = add.value.id
      object_name = add.value.name
    }
  }

  dynamic "add" {
    for_each = var.apply_user_groups_configuration

    content {
      object_type = "USER_GROUP"
      object_id   = add.value.id
      object_name = add.value.name
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region where the apply objects of assist auth configuration are located.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `add` - (Optional, List) Specifies the list of objects to be added.  
  The [add](#workspace_assist_auth_configuration_object_batch_apply_add) structure is documented below.

* `delete` - (Optional, List) Specifies the list of objects to be removed.  
  The [delete](#workspace_assist_auth_configuration_object_batch_apply_delete) structure is documented below.

<a name="workspace_assist_auth_configuration_object_batch_apply_add"></a>
The `add` block supports:

* `object_type` - (Required, String) Specifies the type of the binding object.  
  Valid values are:
  + **USER** - User
  + **USER_GROUP** - User group
  + **ALL** - All users

* `object_id` - (Required, String) Specifies the ID of the user or user group.

* `object_name` - (Required, String) Specifies the name of the user or user group.

<a name="workspace_assist_auth_configuration_object_batch_apply_delete"></a>
The `delete` block supports:

* `object_type` - (Required, String) Specifies the type of the binding object.  
  Valid values are:
  + **USER** - User
  + **USER_GROUP** - User group
  + **ALL** - All users

* `object_id` - (Required, String) Specifies the ID of the user or user group.

* `object_name` - (Required, String) Specifies the name of the user or user group.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
