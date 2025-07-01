---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_pg_plugin"
description: ""
---

# huaweicloud_rds_pg_plugin

Manages RDS for PostgreSQL plugin on the databases within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "database_name" {}

resource "huaweicloud_rds_pg_plugin" "test" {
  instance_id   = var.instance_id
  database_name = var.database_name
  name          = "pgaudit"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the PostgreSQL instance ID.

* `name` - (Required, String, NonUpdatable) Specifies the plugin name.

* `database_name` - (Required, String, NonUpdatable) Specifies the database name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID of PostgreSQL plugin which is formatted `<instance_id>/<database_name>/<name>`.

* `version` - The plugin version.

* `shared_preload_libraries` - Dependent preloaded library.

* `description` - The plugin description.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.
* `delete` - Default is 30 minutes.

## Import

The RDS for PostgreSQL plugin can be imported using the `instance_id`, `database_name` and `name` separated by slashs, e.g.:

```bash
$ terraform import huaweicloud_rds_pg_plugin.test <instance_id>/<database_name>/<name>
```
