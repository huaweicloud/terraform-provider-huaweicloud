---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_studio_workspace_user_roles"
description: |-
  Use this data source to query workspace roles of DataArts Studio within HuaweiCloud.
---

# huaweicloud_dataarts_studio_workspace_user_roles

Use this data source to query workspace user roles of DataArts Studio within HuaweiCloud.

## Example Usage

### Query the user roles under the full instance

```hcl
variable "instance_id" {}

data "huaweicloud_dataarts_studio_workspace_user_roles" "test" {
  instance_id = var.instance_id
}
```

### Query the user roles under the specified workspace

```hcl
variable "workspace_id" {}

data "huaweicloud_dataarts_studio_workspace_user_roles" "test" {
  workspace_id = var.workspace_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the workspace user roles are located.  
  If omitted, the provider-level region will be used.

* `instance_id` - (Optional, String) Specifies the instance ID to query user roles.

* `workspace_id` - (Optional, String) Specifies the workspace ID to query user roles.

  -> At least one of `instance_id` and `workspace_id` must be set.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `roles` - The list of user roles that matched filter parameters.
  The [roles](#dataarts_studio_workspace_user_roles_attr) structure is documented below.

<a name="dataarts_studio_workspace_user_roles_attr"></a>
The `roles` block supports:

* `id` - The role ID.

* `code` - The role code.

* `name` - The role name.

* `description` - The role description.
