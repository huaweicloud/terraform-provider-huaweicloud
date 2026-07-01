---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_instance_account_permission"
description: |-
  Manages a GaussDB instance database account permission resource within HuaweiCloud.
---

# huaweicloud_gaussdb_instance_account_permission

Manages a GaussDB instance database account permission resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "db_name" {}
variable "name" {}
variable "schema_name" {}

resource "huaweicloud_gaussdb_instance_account_permission" "test" {
  instance_id = var.instance_id
  db_name     = var.db_name

  users {
    name        = var.name
    readonly    = "true"
    schema_name = var.schema_name
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to set the database account permission.  
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String, NonUpdatable) Specifies the instance ID of the GaussDB instance.  
  This parameter is the unique identifier of the instance created by the user.

* `db_name` - (Required, String, NonUpdatable) Specifies the database name.  
  Template databases such as `postgres`, `template0`, and `template1` cannot be used. The database must already exist.

* `users` - (Required, List) Specifies the list of database accounts associated with the database.  
  A maximum of 50 elements can be configured in a single request.  
  The [users](#gaussdb_instance_account_permission_users) structure is documented below.

<a name="gaussdb_instance_account_permission_users"></a>
The `users` block supports:

* `name` - (Required, String, NonUpdatable) Specifies the database account name.  
  System users such as `rdsAdmin`, `rdsMetric`, `rdsBackup`, `rdsRepl`, and `root` cannot be used.
  The account name must already exist.

* `readonly` - (Required, String) Specifies the database account permission.  
  The valid values are as follows:
  + **true**: Read-only permission.
  + **false**: Read and write permission.

* `schema_name` - (Required, String) Specifies the schema name.  
  The schemas `public` and `information_schema` cannot be used. The schema name must already exist.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID which is formatted `<instance_id>/<db_name>`.
