---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_instance_plugin_extensions_config"
description: |-
  Manages a GaussDB instance plugin extension configuration resource within HuaweiCloud.
---

# huaweicloud_gaussdb_instance_plugin_extensions_config

Manages a GaussDB instance plugin extension configuration resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "db_name" {}
variable "plugin_name" {}
variable "extension_name" {}

resource "huaweicloud_gaussdb_instance_plugin_extensions_config" "test" {
  instance_id    = var.instance_id
  db_name        = var.db_name
  plugin_name    = var.plugin_name
  extension_name = var.extension_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the GaussDB instance.

* `db_name` - (Required, String, NonUpdatable) Specifies the name of the database.

* `plugin_name` - (Required, String, NonUpdatable) Specifies the name of the plugin.

* `extension_name` - (Required, String, NonUpdatable) Specifies the name of the extension.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in `<instance_id>/<db_name>/<plugin_name>/<extension_name>` format.

* `status` - The status of the plugin extension.

## Import

The GaussDB instance plugin extension configuration can be imported using the `instance_id`, `db_name`, `plugin_name`
and `extension_name`, separated by slashes, e.g.

```bash
$ terraform import huaweicloud_gaussdb_instance_plugin_extensions_config.test \
  <instance_id>/<db_name>/<plugin_name>/<extension_name>
```
