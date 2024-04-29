---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_pg_accounts"
description: ""
---

# huaweicloud_rds_pg_accounts

Use this data source to get the list of RDS PostgreSQL accounts.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_rds_pg_accounts" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the PostgreSQL instance ID.

* `user_name` - (Optional, String) Specifies the username of the DB account.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `users` - Indicates the user list.
  The [users](#PgAccounts_Users) structure is documented below.

<a name="PgAccounts_Users"></a>
The `users` block supports:

* `name` - Indicates the username of the DB account.

* `attributes` - Indicates the permission attributes of a user.
  The [attributes](#PgAccounts_Attributes) structure is documented below.

* `memberof` - Indicates the default rights of a user.

* `description` - Indicates the remarks of the DB account.

<a name="PgAccounts_Attributes"></a>
The `attributes` block supports:

* `rolsuper` - Indicates whether a user has the super user permission. The value is **false**.

* `rolinherit` - Indicates whether a user automatically inherits the permissions of the role to which the user belongs.
  The value can be **true** or **false**.

* `rolcreaterole` - Indicates whether a user can create other sub-users. The value can be **true** or **false**.

* `rolcreatedb` - Indicates whether a user can create a database. The value can be **true** or **false**.

* `rolcanlogin` - Indicates whether a user can log in to the database. The value can be **true** or **false**.

* `rolconnlimit` - Indicates the maximum number of concurrent connections to a DB instance. The value **-1** indicates
  that there are no limitations on the number of concurrent connections.

* `rolreplication` - Indicates whether the user is a replication role. The value can be **true** or **false**.

* `rolbypassrls` - Indicates whether a user bypasses each row-level security policy. The value can be **true** or **false**.
