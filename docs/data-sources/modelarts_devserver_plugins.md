---
subcategory: "AI Development Platform (ModelArts)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_modelarts_devserver_plugins"
description: |-
  Use this data source to get the list of available DevServer plugins of ModelArts within Huaweicloud.
---

# huaweicloud_modelarts_devserver_plugins

Use this data source to get the list of available DevServer plugins of ModelArts within Huaweicloud.

## Example Usage

### Query all DevServer plugins and without any filter

```hcl
data "huaweicloud_modelarts_devserver_plugins" "test" {}
```

### Query the DevServer plugins and using name filter

```hcl
variable "plugin_name" {}

data "huaweicloud_modelarts_devserver_plugins" "test" {
  name = var.plugin_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the DevServer plugins are located.  
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the name of the DevServer plugin to be queried.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `plugins` - The list of DevServer plugins that matched filter parameters.  
  The [plugins](#modelarts_devserver_plugins_attr) structure is documented below.

<a name="modelarts_devserver_plugins_attr"></a>
The `plugins` block supports:

* `name` - The name of the plugin.

* `infos` - The list of plugin details.  
  The [infos](#modelarts_devserver_plugins_infos) structure is documented below.

<a name="modelarts_devserver_plugins_infos"></a>
The `infos` block supports:

* `id` - The ID of the plugin.

* `version` - The version of the plugin.

* `status` - The status of the plugin.

* `url` - The download URL of the plugin.
