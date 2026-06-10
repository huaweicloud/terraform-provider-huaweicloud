---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_instance_plugin_extensions"
description: |-
  Use this data source to query the plugin extensions of a GaussDB instance within HuaweiCloud.
---

# huaweicloud_gaussdb_instance_plugin_extensions

Use this data source to query the plugin extensions of a GaussDB instance within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "plugin_name" {}
variable "db_name" {}

data "huaweicloud_gaussdb_instance_plugin_extensions" "test" {
  instance_id = var.instance_id
  plugin_name = var.plugin_name
  db_name     = var.db_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the GaussDB instance.

* `plugin_name` - (Required, String) Specifies the name of the plugin.

* `db_name` - (Required, String) Specifies the name of the database.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `extensions` - The list of plugin extensions.
  The [extensions](#gaussdb_instance_plugin_extensions_struct) structure is documented below.

<a name="gaussdb_instance_plugin_extensions_struct"></a>
The `extensions` block supports:

* `extension_name` - The name of the extension.

* `status` - The status of the extension.
