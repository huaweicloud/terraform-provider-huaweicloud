---
subcategory: "FunctionGraph"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_fgs_application_templates"
description: ""
---

# huaweicloud_fgs_application_templates

Use this data source to get the list of application templates within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_fgs_application_templates" "test" {
  runtime = "Python2.7"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the application templates are located.  
  If omitted, the provider-level region will be used.

* `runtime` - (Optional, String) Specifies the runtime name used to query the application templates.  
  The valid values are as follows:
  + **Node.js6.10**
  + **Node.js8.10**
  + **Node.js10.16**
  + **Node.js12.13**
  + **Node.js14.18**
  + **Node.js16.17**
  + **Node.js18.15**
  + **Python2.7**
  + **Python3.6**
  + **Python3.9**
  + **Python3.10**
  + **Java8**
  + **Java11**
  + **Go1.x**
  + **C#(.NET Core 2.1)**
  + **C#(.NET Core 3.1)**
  + **http**
  + **PHP7.3**
  + **Custom**

* `category` - (Optional, String) Specifies the category used to query the application templates.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `templates` - All application templates that match the filter parameters.  
  The [templates](#fgs_application_templates) structure is documented below.

<a name="fgs_application_templates"></a>
The `templates` block supports:

* `id` - The template ID.

* `name` - The template name.

* `runtime` -  The template runtime.

* `category` - The template category.

* `description` - The description of template.

* `type` - The type of the function application.
