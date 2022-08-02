---
subcategory: "Document Database Service (DDS)"
---

# huaweicloud_dds_database_user

Manages a database user resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "user_name" {}
variable "user_password" {}
variable "owned_role_name" {}
variable "owned_role_db_name" {}

resource "huaweicloud_dds_database_user" "test" {
  instance_id = var.instance_id

  name     = var.user_name
  password = var.user_password
  db_name  = var.db_name

  roles {
    name    = var.owned_role_name
    db_name = var.owned_role_db_name
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the DDS instance is located.
  Changing this parameter will create a new user.

* `instance_id` - (Required, String, ForceNew) Specifies the DDS instance ID to which the user belongs.
  Changing this parameter will create a new user.

* `name` - (Required, String, ForceNew) Specifies the user name.
  The name can contain `1` to `64` characters, only letters, digits, underscores (_), hyphens (-) and dots (.) are
  allowed. And cannot use reserved names: **drsFull** or **drsIncremental**.
  Changing this parameter will create a new user.

* `password` - (Required, String) Specifies the user password.
  The assword content must meet the following conditions:
  + `8` to `32` characters long.
  + Must contains uppercase and lowercase letters, digits, and at least one special character (`~!@#%^*-_=+?`).

* `db_name` - (Required, String, ForceNew) Specifies the database name to which the user belongs.
  The name can contain `1` to `64` characters, only letters, digits and underscores (_) are allowed.
  Changing this parameter will create a new user.

  -> After a DDS instances is created, the default database is **admin**.

* `roles` - (Required, List, ForceNew) Specifies the list of roles owned by the current user.
  The [object](#dds_database_owned_roles) structure is documented below. Changing this parameter will create a new user.

<a name="dds_database_owned_roles"></a>
The `roles` block supports:

* `name` - (Required, String, ForceNew) Specifies the name of role owned by the current user.
  The name can contain `1` to `64` characters, only letters, digits, underscores (_), hyphens (-) and dots (.) are
  allowed. Changing this parameter will create a new user.

* `db_name` - (Required, String, ForceNew) Specifies the database name to which this owned role belongs.
  Changing this parameter will create a new user.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `privileges` - The list of database privileges owned by the current user.
  The [object](#dds_database_privileges) structure is documented below.

* `inherited_privileges` - The list of database privileges owned by the current user, includes all privileges
  inherited by owned roles. The [object](#dds_database_privileges) structure is documented below.

<a name="dds_database_privileges"></a>
The `privileges` and `inherited_privileges` block supports:

* `resources` - The details of the resource to which the privilege belongs.
  The [object](#dds_database_resources) structure is documented below.

* `actions` - The operation permission list.

<a name="dds_database_resources"></a>
The `resources` block supports:

* `collection` - The database collection type.

* `db_name` - The database name.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 2 minute.
* `update` - Default is 2 minute.
* `delete` - Default is 2 minute.

## Import

Database users can be imported using their `id` (combination of `instance_id`, `db_name` and `name`), separated by a
slash (/), e.g.

```
terraform import huaweicloud_dds_database_user.test &ltinstance_id&gt/&ltdb_name&gt/&ltname&gt
```

Due to security reason, the imported state may not be identical to your resource definition (`password` parameter).
It is generally recommended running `terraform plan` after importing a user resource.
You can then decide if changes should be applied to the user, or the resource definition should be updated to align with
the user. Also you can ignore changes as below.

```
resource "huaweicloud_dds_database_user" "test" {
  ...

  lifecycle {
    ignore_changes = [
      password,
    ]
  }
}
```
