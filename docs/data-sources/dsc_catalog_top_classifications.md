---
subcategory: "Data Security Center (DSC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dsc_catalog_top_classifications"
description: |-
  Use this data source to query the top 5 classifications in the asset catalog within HuaweiCloud.
---

# huaweicloud_dsc_catalog_top_classifications

Use this data source to query the top 5 classifications in the asset catalog within HuaweiCloud.

## Example Usage

```hcl
variable "type_id" {}

data "huaweicloud_dsc_catalog_top_classifications" "test" {
  type_id = var.type_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the catalog top classifications are located.
  If omitted, the provider-level region will be used.

* `label_id` - (Optional, String) Specifies the ID of the group label.

* `type_id` - (Optional, String) Specifies the ID of the data type.

-> Exactly one of the `label_id` and `type_id` parameters must be specified.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `classifications` - The list of classifications.  
  The [classifications](#catalog_top_classifications_attr) structure is documented below.

<a name="catalog_top_classifications_attr"></a>
The `classifications` block supports:

* `classification_name` - The name of the classification.

* `hit_number` - The number of matched records.

* `column_details` - The column detail list of the classification.  
  The [column_details](#catalog_top_classifications_column_details) structure is documented below.

<a name="catalog_top_classifications_column_details"></a>
The `column_details` block supports:

* `asset_id` - The ID of the asset corresponding to the database.

* `asset_name` - The name of the asset corresponding to the database.

* `column_fqn` - The fully qualified name of the column.

* `db_type` - The type of the database.

* `match_infos` - The match information list.  
  The [match_infos](#catalog_top_classifications_column_details_match_infos) structure is documented below.

<a name="catalog_top_classifications_column_details_match_infos"></a>
The `match_infos` block supports:

* `classification_id` - The ID of the classification.

* `classification_name` - The name of the classification.

* `match_content_cnt` - The matched content count.

* `match_rate` - The match rate (percentage).

* `matched_detail` - The match detail.

* `matched_examples` - The matched example list.  
  The [matched_examples](#catalog_top_classifications_column_details_match_infos_matched_examples) structure is
  documented below.

* `rule_id` - The ID of the rule.

* `rule_name` - The name of the rule.

* `security_level_id` - The ID of the security level.

* `security_level_name` - The name of the security level.

* `security_level_color` - The color corresponding to the security level.

* `template_id` - The ID of the template.

* `template_name` - The name of the template.

<a name="catalog_top_classifications_column_details_match_infos_matched_examples"></a>
The `matched_examples` block supports:

* `context` - The match context.

* `line_number` - The line number of the match.

* `matched_content` - The matched content.

* `nlp_revised` - Whether it has been NLP revised.
