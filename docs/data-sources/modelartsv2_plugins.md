---
subcategory: "AI Development Platform (ModelArts)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_modelartsv2_plugins"
description: |-
  Use this data source to query the plugins of a ModelArts resource pool within HuaweiCloud.
---

# huaweicloud_modelartsv2_plugins

Use this data source to query the plugins of a ModelArts resource pool within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
variable "pool_id" {}

data "huaweicloud_modelartsv2_plugins" "test" {
  pool_id = var.pool_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the plugins are located.  
  If omitted, the provider-level region will be used.

* `pool_id` - (Required, String) Specifies the ID of the resource pool.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `plugins` - The list of the plugins that matched filter parameters.  
  The [plugins](#modelartsv2_plugins_plugins) structure is documented below.

<a name="modelartsv2_plugins_plugins"></a>
The `plugins` block supports:

* `api_version` - The API version of the plugin.

* `kind` - The type of the plugin instance.

* `metadata` - The metadata information of the plugin.  
  The [metadata](#modelartsv2_plugins_metadata) structure is documented below.

* `spec` - The spec information of the plugin.  
  The [spec](#modelartsv2_plugins_spec) structure is documented below.

* `status` - The status information of the plugin.  
  The [status](#modelartsv2_plugins_status) structure is documented below.

<a name="modelartsv2_plugins_metadata"></a>
The `metadata` block supports:

* `name` - The name of the plugin instance.

* `created_at` - The creation time of the plugin instance.

<a name="modelartsv2_plugins_spec"></a>
The `spec` block supports:

* `template` - The template information of the plugin.  
  The [template](#modelartsv2_plugins_spec_template) structure is documented below.

<a name="modelartsv2_plugins_spec_template"></a>
The `template` block supports:

* `name` - The name of the plugin template.

* `version` - The version of the plugin template.

* `inputs` - The installation parameters of the plugin template.

<a name="modelartsv2_plugins_status"></a>
The `status` block supports:

* `phase` - The status of the plugin instance.

* `version` - The version of the plugin instance.

* `reason` - The reason for the plugin instance installation failure.

* `values` - The installation parameters of the plugin instance.

* `resources` - The resources occupied by the plugin instance.  
  The [resources](#modelartsv2_plugins_status_resources) structure is documented below.

<a name="modelartsv2_plugins_status_resources"></a>
The `resources` block supports:

* `involved_object` - The referenced resource object of the plugin.  
  The [involved_object](#modelartsv2_plugins_status_resources_involved_object) structure is documented below.

* `replicas` - The number of replicas of the resource object.

* `limits` - The resource limits of the plugin.

* `requests` - The resource requests of the plugin.

<a name="modelartsv2_plugins_status_resources_involved_object"></a>
The `involved_object` block supports:

* `kind` - The API type of the resource object.

* `api_version` - The API version of the resource object.

* `namespace` - The namespace of the resource object.

* `name` - The name of the resource object.

* `uid` - The unique identifier of the resource object.

* `resource_version` - The current version of the resource object.
