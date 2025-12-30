---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_collector_module_templates"
description: |-
  Use this data source to get the list of collector module templates.
---

# huaweicloud_secmaster_collector_module_templates

Use this data source to get the list of collector module templates.

## Example Usage

```hcl
variable "workspace_id" {}
variable "parser_type" {}

data "huaweicloud_secmaster_collector_module_templates" "test" {
  workspace_id = var.workspace_id
  parser_type  = var.parser_type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID.

* `parser_type` - (Optional, String) Specifies the parser type.
  The value can be **FILTER**, **INPUT** or **OUTPUT**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `common` - The commonly used parsing templates.

  The [common](#common_templates_struct) structure is documented below.

* `list` - List the parsing template.

  The [list](#list_templates_struct) structure is documented below.

<a name="common_templates_struct"></a>
The `common` block supports:

* `template_id` - The template ID.

* `name` - The template name.

* `description` - The template description.

* `title` - The template title.

<a name="list_templates_struct"></a>
The `list` block supports:

* `template_id` - The template ID.

* `name` - The template name.

* `description` - The template description.

* `title` - The template title.
