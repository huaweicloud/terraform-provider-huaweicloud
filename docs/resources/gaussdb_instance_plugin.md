---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_instance_plugin"
description: |-
  Manages a GaussDB instance kernel plugin resource within HuaweiCloud.
---

# huaweicloud_gaussdb_instance_plugin

Manages a GaussDB instance kernel plugin resource within HuaweiCloud.

-> **NOTE:** This resource is a one-time action resource for installing the GaussDB instance plugin. Deleting this
resource will not clear the corresponding request record, but will only remove the resource information from the
tfstate file.

## Example Usage

```hcl
variable "instance_id" {}
variable "plugin_name" {}
variable "url" {}
variable "sha_256" {}

resource "huaweicloud_gaussdb_instance_plugin" "test" {
  instance_id = var.instance_id
  plugin_name = var.plugin_name
  url         = var.url
  sha_256     = var.sha_256
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the GaussDB instance.

* `plugin_name` - (Required, String, NonUpdatable) Specifies the name of the kernel plugin to install.
  Currently only **postgis** is supported.

* `url` - (Required, String, NonUpdatable) Specifies the OBS sharing URL of the plugin package.

* `sha_256` - (Required, String, NonUpdatable) Specifies the SHA256 value of the plugin package.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in `<instance_id>/<plugin_name>` format.

* `installed` - Whether the plugin is installed.

* `port` - The port number used by the plugin.

* `plugin_version` - The version of the installed plugin.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 30 minutes.

## Import

The GaussDB instance plugin can be imported using the `instance_id` and `plugin_name`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_gaussdb_instance_plugin.test <instance_id>/<plugin_name>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response. The missing attributes include: `url` and `sha_256`. It is generally recommended running `terraform plan`
after importing a GaussDB instance plugin. You can then decide if changes should be applied to the plugin resource, or
the resource definition should be updated to align with the plugin state. Also you can ignore changes as below.

```hcl
resource "huaweicloud_gaussdb_instance_plugin" "test" {
  ...

lifecycle {
  ignore_changes = [ url, sha_256, ]
}
}
```
