---
subcategory: "Data Lake Insight (DLI)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dli_sql_templates"
description: |-
  Use this data source to get the list of the DLI SQL templates.
---

# huaweicloud_dli_sql_templates

Use this data source to get the list of the DLI SQL templates.

## Example Usage

```hcl
variable "template_name" {}

data "huaweicloud_dli_sql_templates" "test" {
  name = var.template_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the name of the SQL template to be queried.

* `template_id` - (Optional, String) Specifies the ID of the SQL template to be queried.

* `group` - (Optional, String) Specifies the group name to which the SQL templates belong.

* `owner` - (Optional, String) Specifies user ID of owner to be queried.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `templates` - All templates that match the filter parameters.

  The [templates](#templates_struct) structure is documented below.

<a name="templates_struct"></a>
The `templates` block supports:

* `id` - The ID of SQL template.

* `name` - The name of SQL template.

* `sql` - The SQL statement of SQL template.

* `group` - The group name to which the SQL template belongs.

* `owner` - The user ID of owner.

* `description` - The description of the SQL template.
