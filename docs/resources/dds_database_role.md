---
subcategory: "Document Database Service (DDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dds_database_role"
description: ""
---

# huaweicloud_dds_database_role

Manages a DDS database role resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "role_name" {}
variable "db_name" {}
variable "owned_role_name" {}
variable "owned_role_db_name" {}

resource "huaweicloud_dds_database_role" "test" {
  instance_id = var.instance_id

  name    = var.role_name
  db_name = var.db_name

  roles {
    name    = var.owned_role_name
    db_name = var.owned_role_db_name
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the DDS instance is located.
  Changing this parameter will create a new role.

* `instance_id` - (Required, String, ForceNew) Specifies the DDS instance ID to which the role belongs.
  Changing this parameter will create a new role.

* `name` - (Required, String, ForceNew) Specifies the role name.
  The name can contain `1` to `64` characters, only letters, digits, underscores (_), hyphens (-) and dots (.) are
  allowed. Changing this parameter will create a new role.

* `db_name` - (Required, String, ForceNew) Specifies the database name to which the role belongs.
  The name can contain `1` to `64` characters, only letters, digits and underscores (_) are allowed.
  Changing this parameter will create a new role.

  -> After a DDS instances is created, the default database is **admin**.

* `roles` - (Optional, List, ForceNew) Specifies the list of roles owned by the current role.
  The [roles](#dds_database_owned_roles) structure is documented below.
  Changing this parameter will create a new role.

<a name="dds_database_owned_roles"></a>
The `roles` block supports:

* `name` - (Required, String, ForceNew) Specifies the name of role owned by the current role.
  The name can contain `1` to `64` characters, only letters, digits, underscores (_), hyphens (-) and dots (.) are
  allowed. Changing this parameter will create a new role.

* `db_name` - (Required, String, ForceNew) Specifies the database name to which this owned role belongs.
  Changing this parameter will create a new role.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in format of `<instance_id>/<db_name>/<name>`.

* `privileges` - The list of database privileges owned by the current role.
  The [privileges](#dds_database_privileges) structure is documented below.

* `inherited_privileges` - The list of database privileges owned by the current role, includes all privileges
  inherited by owned roles. The [inherited_privileges](#dds_database_privileges) structure is documented below.

<a name="dds_database_privileges"></a>
The `privileges` and `inherited_privileges` block supports:

* `resources` - The details of the resource to which the privilege belongs.
  The [resources](#dds_database_resources) structure is documented below.

* `actions` - The operation permission list.

<a name="dds_database_resources"></a>
The `resources` block supports:

* `collection` - The database collection type.

* `db_name` - The database name.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 60 minutes.
* `delete` - Default is 60 minutes.

## Import

DDS database roles can be imported using the `instance_id`, `db_name` and `name` separated by slashes (/), e.g.

```bash
terraform import huaweicloud_dds_database_role.test <instance_id>/<db_name>/<name>
```
