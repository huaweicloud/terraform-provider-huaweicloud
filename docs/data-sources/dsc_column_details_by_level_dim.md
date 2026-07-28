---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_column_details_by_level_dim"
description: |-
  Use this data source to get the column details of sensitive data by level dimension within HuaweiCloud.
---

# huaweicloud_dsc_column_details_by_level_dim

Use this data source to get the column details of sensitive data by level dimension.

## Example Usage

```hcl
variable "label_id" {}

data "huaweicloud_dsc_column_details_by_level_dim" "test" {
  label_id = var.label_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `label_id` - (Required, String) Specifies the group label ID for filtering.

* `type_id` - (Optional, String) Specifies the type ID for filtering.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `results` - The column details list by level dimension.

  The [results](#results_struct) structure is documented below.

<a name="results_struct"></a>
The `results` block supports:

* `level_name` - The level name.

* `columns` - The column information and match information list.

  The [columns](#columns_struct) structure is documented below.

<a name="columns_struct"></a>
The `columns` block supports:

* `asset_id` - The asset ID.

* `asset_name` - The asset name.

* `column_fqn` - The fully qualified name (FQN) of the column.

* `db_type` - The database type.

* `match_infos` - The match information list.

  The [match_infos](#match_infos_struct) structure is documented below.

<a name="match_infos_struct"></a>
The `match_infos` block supports:

* `classification_id` - The classification ID.

* `classification_name` - The classification name.

* `match_content_cnt` - The match content count.

* `match_rate` - The match rate (percentage).

* `matched_detail` - The matched detail.

* `matched_examples` - The matched examples list.

  The [matched_examples](#matched_examples_struct) structure is documented below.

* `rule_id` - The rule ID.

* `rule_name` - The rule name.

* `security_level_color` - The security level color (RGB value).

* `security_level_id` - The security level ID.

* `security_level_name` - The security level name.

* `template_id` - The template ID.

* `template_name` - The template name.

<a name="matched_examples_struct"></a>
The `matched_examples` block supports:

* `context` - The match context.

* `line_number` - The line number of the match.

* `matched_content` - The matched content.

* `nlp_revised` - Whether the content has been NLP revised.
