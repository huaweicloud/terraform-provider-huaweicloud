---
subcategory: "Config"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rms_assignment_package_templates"
description: ""
---

# huaweicloud_rms_assignment_package_templates

Use this data source to get the list of RMS assignment package templates.

## Example Usage

```hcl
data "huaweicloud_rms_assignment_package_templates" "test" {
  template_key = "test_template_key.json"
  description  = "test_template_description"
}
```

## Argument Reference

The following arguments are supported:

* `template_key` - (Optional, String) Specifies the name of a built-in assignment package template.

* `description` - (Optional, String) Specifies the description for a built-in assignment package template.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `templates` - Indicates the list of RMS assignment package templates.
  The [templates](#Templates_Template) structure is documented below.

<a name="Templates_Template"></a>
The `templates` block supports:

* `id` - Indicates the ID of a built-in assignment package template.

* `template_key` - Indicates the name of a built-in assignment package template.

* `description` - Indicates the description for a built-in assignment package template.

* `template_body` - Indicates the content of a built-in assignment package template.

* `parameters` - Indicates the parameters for a built-in assignment package template.
  The [parameters](#Templates_TemplateParameter) structure is documented below.

<a name="Templates_TemplateParameter"></a>
The `parameters` block supports:

* `name` - Indicates the name of a parameter for a built-in assignment package template.

* `description` - Indicates the description of a parameter for a built-in assignment package template.

* `default_value` - Indicates the default value of a parameter for a built-in assignment package template.

* `type` - Indicates the type of a parameter for a built-in assignment package template.
