---
subcategory: "CodeArts Pipeline""
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_codearts_pipeline_available_plugins"
description: |-
  Use this data source to get a list of CodeArts pipelines available plugins.
---

# huaweicloud_codearts_pipeline_available_plugins

Use this data source to get a list of CodeArts pipelines available plugins.

## Example Usage

```hcl
data "huaweicloud_codearts_pipeline_available_plugins" "test" {
  use_condition = "pipeline"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `business_type` - (Optional, String) Specifies the service type.

* `regex_name` - (Optional, String) Specifies the regex name.

* `use_condition` - (Optional, String) Specifies the use condition.
  Options: **pipeline** or **template**.

* `input_repo_type` - (Optional, String) Specifies the source code repository type.
  Options: **codehub**, **gitee**, **github**, **gitcode**, **gitlab**.

* `input_source_type` - (Optional, String) Specifies the input source type, whether a pipeline has one source or
  multiple sources.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data` - Indicates the result set.

  The [data](#data_struct) structure is documented below.

<a name="data_struct"></a>
The `data` block supports:

* `unique_id` - Indicates the unique ID.

* `display_name` - Indicates the display name.

* `business_type` - Indicates the service type.

* `conditions` - Indicates the conditions.

* `removable` - Indicates whether it is removable.

* `cloneable` - Indicates whether it is replicable.

* `disabled` - Indicates whether it is disabled.

* `editable` - Indicates whether it is editable.

* `plugins_list` - Indicates the extension list.

  The [plugins_list](#data_plugins_list_struct) structure is documented below.

<a name="data_plugins_list_struct"></a>
The `plugins_list` block supports:

* `unique_id` - Indicates the unique ID.

* `display_name` - Indicates the display name.

* `plugin_name` - Indicates the extension name.

* `disabled` - Indicates whether it is disabled.

* `group_name` - Indicates the group name.

* `group_type` - Indicates the group type.

* `plugin_attribution` - Indicates the attribute.

* `plugin_composition_type` - Indicates the combination extension.

* `runtime_attribution` - Indicates the runtime attributes.

* `description` - Indicates the description.

* `version_attribution` - Indicates the version attribute.

* `icon_url` - Indicates the icon URL.

* `multi_step_editable` - Indicates whether it is editable.

* `location` - Indicates the address.

* `publisher_unique_id` - Indicates the publisher ID.

* `manifest_version` - Indicates the version.

* `all_steps` - Indicates the basic extension list.

  The [all_steps](#plugins_list_all_steps_struct) structure is documented below.

<a name="plugins_list_all_steps_struct"></a>
The `all_steps` block supports:

* `plugin_name` - Indicates the extension name.

* `display_name` - Indicates the display name.

* `version` - Indicates the version.
