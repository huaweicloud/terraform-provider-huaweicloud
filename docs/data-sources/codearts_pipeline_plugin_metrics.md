---
subcategory: "CodeArts Pipeline"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_codearts_pipeline_plugin_metrics"
description: |-
  Use this data source to get a list of CodeArts pipeline plugin mertics.
---

# huaweicloud_codearts_pipeline_plugin_metrics

Use this data source to get a list of CodeArts pipeline plugin mertics.

## Example Usage

```hcl
variable "plugin_name" {}
variable "version" {}

data "huaweicloud_codearts_pipeline_plugin_metrics" "test" {
  plugin_name = var.plugin_name
  version     = var.version
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `plugin_name` - (Required, String) Specifies the plugin name.

* `display_name` - (Optional, String) Specifies the display name.

* `plugin_attribution` - (Optional, String) Specifies the extension attribute. The value can be **official** or **custom**.

* `version` - (Optional, String) Specifies the version.

* `version_attribution` - (Optional, String) Specifies the version attribute. The value can be **draft** or **formal**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `metrics` - Indicates the plugin list.
  The [metrics](#attrblock--metrics) structure is documented below.

<a name="attrblock--metrics"></a>
The `metrics` block supports:

* `version` - Indicates the version.

* `output_key` - Indicates the output key.

* `output_value` - Indicates the output value.
