---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_instance_role_permission"
description: |-
  Manages a GaussDB instance database role permission resource within HuaweiCloud.
---

# huaweicloud_gaussdb_instance_role_permission

Manages a GaussDB instance database role permission resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "db_name" {}
variable "role_name" {}
variable "schema" {}

resource "huaweicloud_gaussdb_instance_role_permission" "test" {
  instance_id = var.instance_id
  db_name     = var.db_name

  user {
    name                      = var.role_name
    readonly                  = "false"
    schema                    = var.schema
    default_privilege_grantee = ""
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the GaussDB instance.

* `db_name` - (Required, String, NonUpdatable) Specifies the database name.
  The value cannot be a template database name.
  The template databases include **postgres**, **template0**, and **template1**.

* `user` - (Required, List) Specifies the role permission information.
  The [user](#gaussdb_instance_role_permission_user) structure is documented below.

<a name="gaussdb_instance_role_permission_user"></a>
The `user` block supports:

* `name` - (Required, String, NonUpdatable) Specifies the database role name.
  The value cannot be a system user or role name.
  The system users/roles include **rdsAdmin**, **rdsMetric**, **rdsBackup**, **rdsRepl**, and **root**.

* `readonly` - (Required, String) Specifies the database permission.
  The valid values are as follows:
  + **true**: Read-only.
  + **false**: Readable and writable.

* `schema` - (Required, String) Specifies the schema name.
  The value cannot be **public** or **information_schema**.

* `default_privilege_grantee` - (Optional, String) Specifies the database user/role name.
  This field is used to grant the permissions of this user/role to the role specified by `name`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID which is formatted `<instance_id>/<db_name>`.
