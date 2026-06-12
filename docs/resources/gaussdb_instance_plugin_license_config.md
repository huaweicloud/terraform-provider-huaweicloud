---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_instance_plugin_license_config"
description: |-
  Manages a GaussDB instance plugin license configuration resource within HuaweiCloud.
---

# huaweicloud_gaussdb_instance_plugin_license_config

Manages a GaussDB instance plugin license configuration resource within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "license_str" {}

resource "huaweicloud_gaussdb_instance_plugin_license_config" "test" {
  instance_id = var.instance_id
  license_str = var.license_str
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `instance_id` - (Required, String, NonUpdatable) Specifies the ID of the GaussDB instance. Changing this parameter will
  create a new resource.

* `license_str` - (Required, String, NonUpdatable) Specifies the license string for the plugin. Changing this parameter
  will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.
