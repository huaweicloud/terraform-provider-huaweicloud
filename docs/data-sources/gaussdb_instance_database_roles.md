---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_instance_database_roles"
description: |-
  Use this data source to get the list of GaussDB instance database roles.
---

# huaweicloud_gaussdb_instance_database_roles

Use this data source to get the list of GaussDB instance database roles.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_gaussdb_instance_database_roles" "test" {
  instance_id = var.instance_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the GaussDB instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `roles` - The list of database roles.
  The [roles](#roles) structure is documented below.

<a name="roles"></a>
The `roles` block supports:

* `name` - The name of the user/role.

* `memberof` - The default permissions of the user/role.

* `lock_status` - Whether the user/role is locked.

* `attribute` - The permission attributes of the user/role.
  The [attribute](#attribute) structure is documented below.

<a name="attribute"></a>
The `attribute` block supports:

* `rolsuper` - Whether the user/role has administrator privileges.

* `rolinherit` - Whether the user/role automatically inherits the permissions of its parent role.

* `rolcreaterole` - Whether the user/role can create other child users.

* `rolcreatedb` - Whether the user/role can create databases.

* `rolcanlogin` - Whether the user/role can log in to the database.

* `rolconnlimit` - The maximum number of concurrent connections for the user/role.
  The value **-1** indicates no limit.

* `rolreplication` - Whether the user/role is a replication role.

* `rolbypassrls` - Whether the user/role bypasses each row-level security policy.

* `rolpassworddeadline` - The password expiration time of the user/role.
