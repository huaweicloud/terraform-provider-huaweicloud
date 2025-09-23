---
subcategory: "Log Tank Service (LTS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_lts_structuring_custom_templates"
description: |-
  Use this data source to get the list of LTS structuring custom templates.
---

# huaweicloud_lts_structuring_custom_templates

Use this data source to get the list of LTS structuring custom templates.

## Example Usage

```hcl
data "huaweicloud_lts_structuring_custom_templates" "test" {
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `template_id` - (Optional, String) Specifies the custom template ID to be queried.

* `name` - (Optional, String) Specifies the custom template name to be queried.

* `type` - (Optional, String) Specifies the custom template type to be queried. Valid values are: **regex**, **json**,
  **split** and **nginx**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `templates` - The list of LTS structuring custom templates.
  The [templates](#CustomTemplates_templates) structure is documented below.

<a name="CustomTemplates_templates"></a>
The `templates` block supports:

* `id` - The structuring custom template ID.

* `name` - The structuring custom template name.

* `type` - The structuring custom template type.

* `demo_log` - The sample log event.
