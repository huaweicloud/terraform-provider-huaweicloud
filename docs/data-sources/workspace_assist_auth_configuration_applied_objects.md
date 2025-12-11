---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_assist_auth_configuration_applied_objects"
description: |-
  Use this data source to get the applied object list of assist auth configuration within HuaweiCloud.
---

# huaweicloud_workspace_assist_auth_configuration_applied_objects

Use this data source to get the applied object list of assist auth configuration within HuaweiCloud.

## Example Usage

### Query all applied objects of assist auth configuration

```hcl
data "huaweicloud_workspace_assist_auth_configuration_applied_objects" "test" {}
```

### Query all applied users of assist auth configuration

```hcl
data "huaweicloud_workspace_assist_auth_configuration_applied_objects" "test" {
  object_type = "USER"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) The region where the applied objects of assist auth configuration are located.  
  If omitted, the provider-level region will be used.

* `object_type` - (Optional, String) Specifies the type of the applied objects.  
  The valid values are as follows:
  + **USER**
  + **USER_GROUP**
  + **ALL**

  Defaults to **ALL**.

* `object_name` - (Optional, String) Specifies the name of the object.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `objects` - The list of assist auth configuration applied objects that matched filter parameters.  
  The [objects](#workspace_assist_auth_configuration_applied_objects_attr) structure is documented below.

<a name="workspace_assist_auth_configuration_applied_objects_attr"></a>
The `objects` block supports:

* `id` - The ID of the user or user group.

* `type` - The type of the binding object.
  + **USER** - User
  + **USER_GROUP** - User group
  + **ALL** - All users

* `name` - The name of the user or user group.
