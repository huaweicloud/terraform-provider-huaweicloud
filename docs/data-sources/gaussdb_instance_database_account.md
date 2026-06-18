---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_instance_database_account"
description: |-
  Manages a GaussDB instance database account resource within HuaweiCloud.
---

# huaweicloud_gaussdb_instance_database_account

Manages a GaussDB instance database account resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "name" {}
variable "password" {}

resource "huaweicloud_gaussdb_instance_database_account" "test" {
  instance_id   = var.instance_id
  name          = var.name
  password      = var.password
  is_login_only = "true"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the GaussDB instance.
  Changing this parameter will create a new resource.

* `name` - (Required, String, NonUpdatable) Specifies the database account name.
  The name must be unique and cannot be the same as any system user names.
  The value can contain 1 to 63 characters, including letters, digits, and underscores.
  It cannot start with **pg** or a digit.
  The system users include **rdsAdmin**, **rdsMetric**, **rdsBackup**, **rdsRepl**, and **root**.
  Changing this parameter will create a new resource.

* `password` - (Required, String) Specifies the database account password.
  The value must meet the following requirements:
  + Cannot be the same as the account name or the reverse of the account name.
  + Cannot be the same as the old password.
  + Must contain 8 to 32 characters.
  + Must contain at least three types of the following characters: uppercase letters, lowercase letters, digits,
      and special characters (**~!@#%^*-_=+?**).

* `is_login_only` - (Optional, String, NonUpdatable) Specifies whether the database account supports login only.
  The valid values are as follows:
  + **true**: The database account has only the login permission.
  + **false**: The database account has the permissions to log in to databases, create databases, and create users.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID. The value is the combination of `instance_id` and `name`, separated by a slash.

* `attribute` - The permission attributes of the database account.
  The [attribute](#gaussdb_instance_database_account_attribute) structure is documented below.

* `memberof` - The default permissions of the database account.

* `lock_status` - Whether the database account is locked.

<a name="gaussdb_instance_database_account_attribute"></a>
The `attribute` block supports:

* `rolsuper` - Whether the user has administrator permissions.

* `rolinherit` - Whether the user automatically inherits the permissions of its roles.

* `rolcreaterole` - Whether the user can create other sub-users.

* `rolcreatedb` - Whether the user can create databases.

* `rolcanlogin` - Whether the user can log in to databases.

* `rolconnlimit` - The maximum number of concurrent connections of the user. The value **-1** indicates no limit.

* `rolreplication` - Whether the user belongs to a replication role.

* `rolbypassrls` - Whether the user bypasses every row-level security policy.

* `rolpassworddeadline` - The password expiration time of the user.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
* `update` - Default is 30 minutes.

## Import

The GaussDB instance database account can be imported using the `instance_id` and `name` separated by a slash, e.g.

```bash
$ terraform import huaweicloud_gaussdb_instance_database_account.test <instance_id>/<name>
```
