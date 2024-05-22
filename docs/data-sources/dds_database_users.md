---
subcategory: "Document Database Service (DDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dds_database_users"
description: |-
  Use this data source to get the list of DDS instance database users.
---

# huaweicloud_dds_database_users

Use this data source to get the list of DDS instance database users.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_dds_database_users" "test"{
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the instance ID.

* `name` - (Optional, String) Specifies the username.

* `db_name` - (Optional, String) Specifies the database name to which the user belongs.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `users` - Indicates the backup list.
  The [users](#users_struct) structure is documented below.

<a name="users_struct"></a>
The `users` block supports:

* `name` - Indicates the username.

* `db_name` - Indicates the database name to which the user belongs.

* `roles` - Indicates the list of roles owned by the current user.
  The [roles](#dds_database_owned_roles) structure is documented below.

* `privileges` - The list of database privileges owned by the current user.
  The [privileges](#dds_database_privileges) structure is documented below.

* `inherited_privileges` - The list of database privileges owned by the current user, includes all privileges
  inherited by owned roles.
  The [inherited_privileges](#dds_database_privileges) structure is documented below.

<a name="dds_database_privileges"></a>
The `privileges` and `inherited_privileges` block supports:

* `resources` - The details of the resource to which the privilege belongs.
  The [resources](#dds_database_resources) structure is documented below.

* `actions` - The operation permission list.

<a name="dds_database_resources"></a>
The `resources` block supports:

* `collection` - The database collection type.

* `db_name` - The database name.

<a name="dds_database_owned_roles"></a>
The `roles` block supports:

* `name` - Indicates the name of role owned by the current user.

* `db_name` - Indicates the database name to which this owned role belongs.
