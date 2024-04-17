---
subcategory: "Data Lake Insight (DLI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dli_flink_templates"
description: |-
  Use this data source to get the list of the DLI flink templates.
---

# huaweicloud_dli_flink_templates

Use this data source to get the list of the DLI flink templates.

## Example Usage

```hcl
variable "template_id" {}

data "huaweicloud_dli_flink_templates" "test" {
  template_id = var.template_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the name of the flink template to be queried.

* `template_id` - (Optional, String) Specifies the ID of the flink template to be queried.

* `type` - (Optional, String) Specifies the type of the flink template to be queried.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `templates` - All templates that match the filter parameters.

  The [templates](#template_list_templates_struct) structure is documented below.

<a name="template_list_templates_struct"></a>
The `templates` block supports:

* `id` - The ID of template.

* `name` - The name of template.

* `type` - The type of template.

* `sql` - The stream SQL statement.

* `description` - The description of template.

* `created_at` - The creation time of the template.

* `updated_at` - The latest update time of the template.
