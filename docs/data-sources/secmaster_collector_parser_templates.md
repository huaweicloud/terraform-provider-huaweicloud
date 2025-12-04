---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_collector_parser_templates"
description: |-
  Use this data source to query the SecMaster collector parser templates within HuaweiCloud.
---

# huaweicloud_secmaster_collector_parser_templates

Use this data source to query the SecMaster collector parser templates within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}

data "huaweicloud_secmaster_collector_parser_templates" "test" {
  workspace_id = var.workspace_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the parser templates.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the ID of the workspace to query parser templates.

* `title` - (Optional, String) Specifies the title of the parser template to filter the results.

* `description` - (Optional, String) Specifies the description to filter the parser templates.

* `sort_key` - (Optional, String) Specifies the key for sorting the results.

* `sort_dir` - (Optional, String) Specifies the sort direction.
  Valid values are `asc` (ascending) and `desc` (descending).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `records` - The list of parser templates that match the query criteria.

  The [records](#records_struct) structure is documented below.

<a name="records_struct"></a>
The `records` block supports:

* `title` - The title of the parser template.

* `description` - The description of the parser template.

* `parser_id` - The unique identifier of the parser template.
