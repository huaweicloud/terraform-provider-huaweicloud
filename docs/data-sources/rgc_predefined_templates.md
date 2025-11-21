---
subcategory: "Resource Governance Center (RGC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rgc_predefined_templates"
description: |-
  Use this data source to list pre-defined templates in Resource Governance Center.
---

# huaweicloud_rgc_predefined_templates

Use this data source to list pre-defined templates in Resource Governance Center.

## Example Usage

```hcl
data "huaweicloud_rgc_predefined_templates" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

## Attributes Reference

* `id` - The data source ID.

* `templates` - Information about the pre-defined templates list.

The [templates](#templates) structure is documented below.

<a name="templates"></a>
The `templates` block supports:

* `template_name` - The name of the predefined template.

* `template_description` - The description of the predefined template.

* `template_category` - The category of the predefined template.
