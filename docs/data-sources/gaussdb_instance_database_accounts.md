---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_instance_database_accounts"
description: |-
  Use this data source to query the database users of a GaussDB instance within HuaweiCloud.
---

# huaweicloud_gaussdb_instance_database_accounts

Use this data source to query the database users of a GaussDB instance within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
data "huaweicloud_gaussdb_instance_database_accounts" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the database accounts.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the GaussDB instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `users` - The list of database users.
  The [users](#users) structure is documented below.

<a name="users"></a>
The `users` block supports:

* `name` - The user name.

* `attribute` - The user permission attributes.
  The [attribute](#users_attribute) structure is documented below.

* `memberof` - The default permissions of the user.

* `lock_status` - Whether the user is locked.

<a name="users_attribute"></a>
The `attribute` block supports:

* `rolsuper` - Whether the user has administrator privileges.

* `rolinherit` - Whether the user automatically inherits the permissions of its roles.

* `rolcreaterole` - Whether the user can create other sub-users.

* `rolcreatedb` - Whether the user can create databases.

* `rolcanlogin` - Whether the user can log in to the database.

* `rolconnlimit` - The maximum number of concurrent connections for the user. `-1` means no limit.

* `rolreplication` - Whether the user belongs to a replication role.

* `rolbypassrls` - Whether the user bypasses every row-level security policy.

* `rolpassworddeadline` - The user password expiration time.
