---
subcategory: "GaussDB(for MySQL)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_mysql_account_privilege"
description: ""
---

# huaweicloud_gaussdb_mysql_account_privilege

Manages a GaussDB MySQL account privilege resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

resource "huaweicloud_gaussdb_mysql_account_privilege" "test" {
  instance_id  = var.instance_id
  account_name = "test_db_name1"
  host         = "10.10.10.10"

  databases {
    name     = "test_db_name"
    readonly = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the GaussDB MySQL instance.

  Changing this parameter will create a new resource.

* `account_name` - (Required, String, ForceNew) Specifies the database username.

  Changing this parameter will create a new resource.

* `host` - (Required, String, ForceNew) Specifies the host IP address which allow database users to connect to the
  database on the current host.

  Changing this parameter will create a new resource.

* `databases` - (Required, List, ForceNew) Specifies the list of the databases. The list contains up to 50 databases.

  Changing this parameter will create a new resource.
The [Database](#GaussDBAccountPrivilege_Database) structure is documented below.

<a name="GaussDBAccountPrivilege_Database"></a>
The `Database` block supports:

* `name` - (Required, String, ForceNew) Specifies the database name.

* `readonly` - (Required, Bool, ForceNew) Specifies whether the database permission is read-only.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID which is formatted `<instance_id>/<account_name>/<host>`.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 10 minutes.
* `update` - Default is 10 minutes.

## Import

The GaussDB MySQL account privilege can be imported using the `instance_id`, `name` and `host` separated by slashes, e.g.

```bash
$ terraform import huaweicloud_gaussdb_mysql_account_privilege.test <instance_id>/<account_name>/<host>
```
