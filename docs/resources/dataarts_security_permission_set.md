---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_security_permission_set"
description: ""
---

# huaweicloud_dataarts_security_permission_set

Manages DataArts Security permission set resource within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}
variable "name" {}
variable "manager_id" {}

resource "huaweicloud_dataarts_security_permission_set" "test" {
  workspace_id = var.workspace_id
  name         = var.name
  parent_id    = "0"
  manager_id   = var.manager_id
}

resource "huaweicloud_dataarts_security_permission_set" "sub_test" {
  workspace_id = var.workspace_id
  name         = var.name
  parent_id    = huaweicloud_dataarts_security_permission_set.test.id
  manager_id   = var.manager_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the permission set resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `workspace_id` - (Required, String, ForceNew) Specifies the workspace ID in which the permission set.
  Changing this creates a new permission set.

* `name` - (Required, String) Specifies the name of the permission set. The name can contain `1` to `128` characters.
  Only letters, digits, and underscores (_) are allowed.

* `parent_id` - (Required, String, ForceNew) Specifies the parent ID of the permission set.
  The parent ID can contain `1` to `128` characters. The value of parent_id is `0`
  when we want to create a workspace permission set.

* `manager_id` - (Required, String) Specifies the manager ID of the permission set. The manager can choose from
  member management under the workspace. The manager ID cancontain
  `1` to `128` characters.

* `manager_name` - (Optional, String) Specifies the manager name of the permission set. The manager name can
  contain `1` to `128` characters.

* `manager_type` - (Optional, String) Specifies the manager type of the permission set. The valid
  values are **USER** and **USER_GROUP**.

* `description` - (Optional, String) Specifies the description of the permission set. The description can contain
  `0` to `10240` characters.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, also the ID of the permission set.

* `datasource_type` - The data source type managed by permission sets.

* `instance_id` - The ID of the instance to which the permission set belongs.

* `type` - The type of the permission set.

* `created_at` - The create time of the permission set.

* `created_by` - The creator of the permission set.

* `updated_at` - The update time of the permission set.

* `updated_by` - The updator of the permission set.

## Import

The DataArts Security permission set can be imported using the `workspace_id` and `id` separated by a slash, e.g.

```bash
$ terraform import huaweicloud_dataarts_security_permission_set.test <workspace_id>/<id>
```
