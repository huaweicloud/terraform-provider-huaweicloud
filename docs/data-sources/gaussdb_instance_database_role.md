---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_instance_database_role"
description: |-
  Manages a GaussDB instance database role resource within HuaweiCloud.
---

# huaweicloud_gaussdb_instance_database_role

Manages a GaussDB instance database role resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_gaussdb_instance_database_role" "test" {
  instance_id = var.instance_id
  name        = "test_role"
  password    = "Test@123456"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the database role.  
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String, NonUpdatable) Specifies the instance ID of the GaussDB instance.  
  This parameter is the unique identifier of the instance created by the user.

* `name` - (Required, String, NonUpdatable) Specifies the database role name.  
  The name must be unique and cannot be the same as existing roles or system users/roles.  
  System users/roles include: rdsAdmin, rdsMetric, rdsBackup, rdsRepl, root.  
  The name must be 1 to 63 characters long and can contain letters, digits, and underscores.  
  It cannot start with "pg" or a digit, and cannot contain other special characters.

* `password` - (Required, String, NonUpdatable) Specifies the database role password.  
  The password must be 8 to 32 characters long and contain at least three types of the following characters:  
  uppercase letters, lowercase letters, digits, and special characters (~!@#%^*-_=+?).  
  The password cannot be the same as the role name or its reverse.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, composed of `instance_id` and `name`, separated by a slash (`/`).

* `attribute` - The permission attributes of the database role.  
  The [attribute](#gaussdb_instance_database_role_attribute) structure is documented below.

* `memberof` - The default permissions of the database role.

* `lock_status` - Whether the database role is locked.

<a name="gaussdb_instance_database_role_attribute"></a>
The `attribute` block supports:

* `rolsuper` - Whether the role has administrator privileges.

* `rolinherit` - Whether the role automatically inherits the permissions of its member roles.

* `rolcreaterole` - Whether the role can create other sub-roles.

* `rolcreatedb` - Whether the role can create databases.

* `rolcanlogin` - Whether the role can log in to the database.

* `rolconnlimit` - The maximum number of concurrent connections for the role. `-1` means no limit.

* `rolreplication` - Whether the role is a replication role.

* `rolbypassrls` - Whether the role bypasses every row-level security policy.

* `rolpassworddeadline` - The password expiration time of the role.

## Import

The GaussDB instance database role resource can be imported using `instance_id` and `name` separated by a slash (`/`),
e.g.

```bash
$ terraform import huaweicloud_gaussdb_instance_database_role.test <instance_id>/<name>
```
