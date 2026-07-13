---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_column_details_by_classification"
description: |-
  Use this data source to get the column details by classification dimension within HuaweiCloud.
---

# huaweicloud_dsc_column_details_by_classification

Use this data source to get the column details by classification dimension within HuaweiCloud.

-> Either `label_id` or `type_id` must be specified as a query parameter.

## Example Usage

```hcl
variable "type_id" {}

data "huaweicloud_dsc_column_details_by_classification" "test" {
  type_id = var.type_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `label_id` - (Optional, String) Specifies the group label ID for filtering.
  Either `label_id` or `type_id` must be specified.

* `type_id` - (Optional, String) Specifies the type ID for filtering.
  Either `label_id` or `type_id` must be specified.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `results` - The column details list by classification dimension.
  The [results](#dsc_results_struct) structure is documented below.

<a name="dsc_results_struct"></a>
The `results` block supports:

* `template_name` - The template name.

* `count` - The match count.

* `classifications` - The classification information list.
  The [classifications](#dsc_classifications_struct) structure is documented below.

<a name="dsc_classifications_struct"></a>
The `classifications` block supports:

* `classification_name` - The classification name.

* `count` - The match count.

* `columns` - The column information and match information list.
  The [columns](#dsc_columns_struct) structure is documented below.

<a name="dsc_columns_struct"></a>
The `columns` block supports:

* `asset_id` - The asset ID.

* `asset_name` - The asset name.

* `column_fqn` - The column fully qualified name.

* `db_type` - The database type.

* `match_infos` - The match information list.
  The [match_infos](#dsc_match_infos_struct) structure is documented below.

<a name="dsc_match_infos_struct"></a>
The `match_infos` block supports:

* `classification_id` - The classification ID.

* `classification_name` - The classification name.

* `match_content_cnt` - The matched content count.

* `match_rate` - The match rate (percentage).

* `matched_detail` - The match detail.

* `matched_examples` - The matched example list.
  The [matched_examples](#dsc_matched_examples_struct) structure is documented below.

* `rule_id` - The rule ID.

* `rule_name` - The rule name.

* `security_level_color` - The security level color.

* `security_level_id` - The security level ID.

* `security_level_name` - The security level name.

* `template_id` - The template ID.

* `template_name` - The template name.

<a name="dsc_matched_examples_struct"></a>
The `matched_examples` block supports:

* `context` - The match context.

* `line_number` - The line number of the match.

* `matched_content` - The matched content.

* `nlp_revised` - Whether it has been NLP revised.
