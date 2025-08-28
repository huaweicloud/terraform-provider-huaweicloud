---
subcategory: "CodeArts Pipeline"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_codearts_pipeline_plugin_versions"
description: |-
  Use this data source to get a list of CodeArts pipeline plugin versions.
---

# huaweicloud_codearts_pipeline_plugin_versions

Use this data source to get a list of CodeArts pipeline plugin versions.

## Example Usage

```hcl
variable "plugin_name" {}

data "huaweicloud_codearts_pipeline_plugin_versions" "test" {
  plugin_name = var.plugin_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `plugin_name` - (Required, String) Specifies the plugin name.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `versions` - Indicates the version list.

  The [versions](#versions_struct) structure is documented below.

<a name="versions_struct"></a>
The `versions` block supports:

* `plugin_name` - Indicates the extension name.

* `unique_id` - Indicates the unique ID.

* `display_name` - Indicates the display name.

* `version` - Indicates the version.

* `version_description` - Indicates the version description.

* `version_attribution` - Indicates the version attribute.

* `description` - Indicates the description.

* `plugin_attribution` - Indicates the attribute.

* `plugin_composition_type` - Indicates the combination type.

* `icon_url` - Indicates the icon URL.

* `workspace_id` - Indicates the tenant ID.

* `refer_count` - Indicates the number of references.

* `active` - Indicates the activated or not.

* `business_type` - Indicates the service type.

* `business_type_display_name` - Indicates the display name of service type.

* `op_user` - Indicates the operator.

* `op_time` - Indicates the operation time.

* `maintainers` - Indicates the maintenance engineer.

* `usage_count` - Indicates the number of usages.

* `runtime_attribution` - Indicates the runtime attributes.
