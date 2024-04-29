---
subcategory: "Relational Database Service (RDS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rds_pg_plugins"
description: ""
---

# huaweicloud_rds_pg_plugins

Use this data source to get the list of RDS PostgreSQL plugins.

## Example Usage

```hcl
variable "instance_id" {}
variable "database_name" {}

data "huaweicloud_rds_pg_plugins" "plugins" {
  instance_id   = var.instance_id
  database_name = var.database_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of a PostgreSQL instance.

* `database_name` - (Required, String) Specifies the database name of a PostgreSQL instance.

* `name` - (Optional, String) Specifies the plugin name.

* `version` - (Optional, String) Specifies the plugin version.

* `created` - (Optional, Bool) Specifies whether the plugin has been created. Defaults to: **false**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `plugins` - Indicates the plugin list.
  The [plugins](#PgPlugins_Plugin) structure is documented below.

<a name="PgPlugins_Plugin"></a>
The `plugins` block supports:

* `name` - Indicates the plugin name.

* `version` - Indicates the plugin version.

* `shared_preload_libraries` - Indicates the dependent preloaded library.

* `created` - Indicates whether the plugin has been created.

* `description` - Indicates the plugin description.
