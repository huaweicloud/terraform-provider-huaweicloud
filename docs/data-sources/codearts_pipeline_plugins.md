---
subcategory: "CodeArts Pipeline"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_codearts_pipeline_plugins"
description: |-
  Use this data source to get a list of CodeArts pipeline plugins.
---

# huaweicloud_codearts_pipeline_plugins

Use this data source to get a list of CodeArts pipeline plugins.

## Example Usage

```hcl
data "huaweicloud_codearts_pipeline_plugins" "test" {
  business_type = ["Build", "Gate", "Deploy", "Test", "Normal"]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `business_type` - (Optional, List) Specifies the service type.
  Valid values are **Build**, **Gate**, **Deploy**, **Test** and **Normal**.

* `maintainer` - (Optional, String) Specifies the maintenance engineer.

* `plugin_attribution` - (Optional, String) Specifies the extension attribute.
  Valid values are **official** and **custom**.

* `plugin_name` - (Optional, String) Specifies the plugin name.

* `regex_name` - (Optional, String) Specifies the match name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `plugins` - Indicates the plugin list.
  The [plugins](#attrblock--plugins) structure is documented below.

<a name="attrblock--plugins"></a>
The `plugins` block supports:

* `unique_id` - Indicates the unique ID.

* `plugin_name` - Indicates the plugin name.

* `active` - Indicates whether the plugin is activate or not.

* `business_type` - Indicates the service type.

* `business_type_display_name` - Indicates the display name of service type.

* `description` - Indicates the description.

* `display_name` - Indicates the display name.

* `icon_url` - Indicates the icon URL.

* `maintainers` - Indicates the maintenance engineer.

* `op_time` - Indicates the operation time.

* `op_user` - Indicates the operator.

* `plugin_attribution` - Indicates the attribute.

* `plugin_composition_type` - Indicates the combination type.

* `refer_count` - Indicates the number of references.

* `runtime_attribution` - Indicates the runtime attribution.

* `usage_count` - Indicates the number of usages.

* `version` - Indicates the version.

* `version_attribution` - Indicates the version attribution.

* `version_description` - Indicates the version description.

* `workspace_id` - Indicates the tenant ID.
