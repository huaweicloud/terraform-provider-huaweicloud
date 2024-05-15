---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_api_associated_plugins"
description: |-
  Use this data source to query the plugins associated with the specified API within HuaweiCloud.
---

# huaweicloud_apig_api_associated_plugins

Use this data source to query the plugins associated with the specified API within HuaweiCloud.

## Example Usage

### Query the contents of all plugins bound to the current API

```hcl
variable "instance_id" {}
variable "associated_api_id" {}

data "huaweicloud_apig_api_associated_plugins" "test" {
  instance_id = var.instance_id
  api_id      = var.associated_api_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the associated plugins.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the dedicated instance to which the plugins belong.

* `api_id` - (Required, String) Specifies the ID of the API bound to the plugin.

* `plugin_id` - (Optional, String) Specifies the ID of the plugin.

* `name` - (Optional, String) Specifies the name of the plugin.

* `type` - (Optional, String) Specifies the type of the plugin.

* `env_id` - (Optional, String) Specifies the ID of the environment where the API is published.

* `env_name` - (Optional, String) Specifies the name of the environment where the API is published.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `plugins` - All plugins that match the filter parameters.
  The [plugins](#api_associated_plugins) structure is documented below.

<a name="api_associated_plugins"></a>
The `plugins` block supports:

* `id` - The ID of the plugin.

* `name` - The name of the plugin.

* `type` - The type of the plugin.

* `description` - The description of the plugin.

* `content` - The configuration details for the plugin.

* `env_id` - The ID of the environment where the API is published.

* `env_name` - The name of the environment where the API is published.

* `bind_id` - The bind ID.

* `bind_time` - The time that the plugin is bound to the API.
