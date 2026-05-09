---
subcategory: "AI Development Platform (ModelArts)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_modelartsv2_plugin_templates"
description: |-
  Use this data source to query ModelArts plugin temppates within HuaweiCloud.
---

# huaweicloud_modelartsv2_plugin_templates

Use this data source to query ModelArts plugin temppates within HuaweiCloud.

## Example Usage

```hcl
variable "template_name" {}
variable "pool_name" {}

data "huaweicloud_modelartsv2_plugin_templates" "test" {
  template_name = var.template_name
  pool_name     = var.pool_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the plugin templates are located.  
  If omitted, the provider-level region will be used.

* `template_name` - (Optional, String) Specifies the template name of the plugin templates.

* `pool_name` - (Optional, String) Specifies the pool name of the plugin templates.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `plugin_templates` - The list of plugin templates that matched filter parameters.  
  The [plugin_templates](#modelarts_plugin_templates) structure is documented below.

<a name="modelarts_plugin_templates"></a>
The `plugin_templates` block supports:

* `metadata` - The metadata of the plugin template.  
  The [metadata](#modelarts_plugin_templates_metadata) structure is documented below.

<a name="modelarts_plugin_templates_metadata"></a>
The `metadata` block supports:

* `name` - The metadata name of the plugin template.

* `annotations` - The metadata annotations of the plugin template.

* `spec` - The spec of the plugin template.  
  The [spec](#modelarts_plugin_templates_spec) structure is documented below.

<a name="modelarts_plugin_templates_spec"></a>
The `spec` block supports:

* `optional` - The spec optional of the plugin template.

* `type` - The spec type of the plugin template.

* `logo_url` - The spec logo url of the plugin template.

* `description` - The spec description of the plugin template.

* `versions` - The versions of the plugin template.  
  The [versions](#modelarts_plugin_templates_spec_versions) structure is documented below.

<a name="modelarts_plugin_templates_spec_versions"></a>
The `versions` block supports:

* `version` - The version of the plugin template.

* `creation_timestamp` - The creation timestamp of the plugin template, in RFC3339 format.

* `inputs` - The inputs of the plugin template.

* `translate` - The translate of the plugin template.
