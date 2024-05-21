---
subcategory: "Cloud Search Service (CSS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_css_logstash_templates"
description: |-
  Use this data source to get the list of CSS logstash templates.
---

# huaweicloud_css_logstash_templates

Use this data source to get the list of CSS logstash templates.

## Example Usage

```hcl
variable "type" {}
variable "name" {}

data "huaweicloud_css_logstash_templates" "test" {
  type = var.type
  name = var.name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `type` - (Optional, String) Specifies the type of the CSS logstash configuration template.
  The values can be **custom** and **system**.

* `name` - (Optional, String) Specifies the name of the template.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `system_templates` - The list of the system templates.

  The [system_templates](#system_templates_struct) structure is documented below.

* `custom_templates` - The list of the custom templates.

  The [custom_templates](#custom_templates_struct) structure is documented below.

<a name="system_templates_struct"></a>
The `system_templates` block supports:

* `id` - The ID of the system template.

* `name` - The name of the system template.

* `conf_content` - The configuration file content of the system template.

* `desc` - The description of the system template.

<a name="custom_templates_struct"></a>
The `custom_templates` block supports:

* `id` - The ID of the custom template.

* `name` - The name of the custom template.

* `conf_content` - The configuration file content of the custom template.

* `desc` - The description of the custom template.
