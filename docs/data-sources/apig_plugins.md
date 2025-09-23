---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_plugins"
description: |-
  Use this data source to get plugin list under the specified dedicated instance within HuaweiCloud.
---

# huaweicloud_apig_plugins

Use this data source to get plugin list under the specified dedicated instance within HuaweiCloud.

## Example Usage

### Query all plugins under specified instance

```hcl
variable "instance_id" {}

data "huaweicloud_apig_plugins" "test" {
  instance_id = var.instance_id
}
```

### Query the plugins by the specified plugin ID

```hcl
variable "instance_id" {}
variable "plugin_id" {}

data "huaweicloud_apig_plugins" "test" {
  instance_id = var.instance_id
  plugin_id   = var.plugin_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the dedicated instance to which the plugin belongs.

* `plugin_id` - (Optional, String) Specifies the ID of the plugin.

* `name` - (Optional, String) Specifies the name of the plugin. Fuzzy search is supported.

* `type` - (Optional, String) Specifies the type of the plugin.  
  The valid values are as follows:
  + **cors**
  + **set_resp_headers**
  + **rate_limit**
  + **kafka_log**
  + **breaker**
  + **third_auth**
  + **proxy_cache**
  + **proxy_mirror**

* `plugin_scope` - (Optional, String) Specifies the scope of the plugin.  
  The valid values are as follows:
  + **global**

* `precise_search` - (Optional, String) Specifies the name of the parameter to be matched exactly.  
  The valid values are as follows:
  + **name**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `plugins` - All plugins that match the filter parameters.

  The [plugins](#data_plugins_struct) structure is documented below.

<a name="data_plugins_struct"></a>
The `plugins` block supports:

* `id` - The ID of the plugin.

* `name` - The name of the plugin.

* `type` - The type of the plugin.

* `content` - The content of the plugin.

* `plugin_scope` - The scope of the plugin.

* `description` - The description of the plugin.

* `created_at` - The creation time of the plugin, in RFC3339 format.

* `updated_at` - The latest update time of the plugin, in RFC3339 format.
