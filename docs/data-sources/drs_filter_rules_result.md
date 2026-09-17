---
subcategory: "DRS"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_drs_filter_rules_result"
description: |-
  Use this data source to query the data filter rules result of a DRS job within HuaweiCloud.
---

# huaweicloud_drs_filter_rules_result

Use this data source to query the data filter rules result of a DRS job within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
variable "job_id" {}

data "huaweicloud_drs_filter_rules_result" "test" {
  job_id = var.job_id
}
```

### Query with Request ID

```hcl
variable "job_id" {}
variable "query_id" {}

data "huaweicloud_drs_filter_rules_result" "test" {
  job_id   = var.job_id
  query_id = var.query_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the filter rules result.  
  If omitted, the provider-level region will be used.

* `job_id` - (Required, String) Specifies the ID of the DRS job to query filter rules result.

* `query_id` - (Optional, String) Specifies the request ID used to associate the previous filter rules query request.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `status` - The query progress of the filter rules.

* `total_count` - The total number of the filter rules.

* `rules` - The list of the data filter rules.  
  The [rules](#drs_filter_rules_result_rules) structure is documented below.

<a name="drs_filter_rules_result_rules"></a>
The `rules` block supports:

* `batch_id` - The batch ID of the data filter rule.

* `filter_object_list` - The list of the filter objects.  
  The [filter_object_list](#drs_filter_rules_result_filter_object_list) structure is documented below.

* `filter_rules_resp` - The filter condition response of the rule.  
  The [filter_rules_resp](#drs_filter_rules_result_filter_rules_resp) structure is documented below.

* `advanced_setting_resp` - The advanced setting response of the rule.  
  The [advanced_setting_resp](#drs_filter_rules_result_advanced_setting_resp) structure is documented below.

<a name="drs_filter_rules_result_filter_object_list"></a>
The `filter_object_list` block supports:

* `batch_id` - The batch ID of the data filter object.

* `id` - The unique identifier of the data filter object.

* `parent_id` - The parent object identifier of the data filter object.

* `object_name` - The name of the data filter object.

* `object_alias_name` - The alias name of the data filter object.

* `object_type` - The type of the data filter object.

* `db_name` - The database name of the data filter object.

* `schema_name` - The schema name of the data filter object.

* `table_name` - The table name of the data filter object.

* `is_data_filter` - Whether the data filter is enabled for the object.

* `data_filter_type` - The data filter type of the object.

* `filter_condition` - The filter condition of the object.

* `is_synchronized` - Whether the data filter condition is synchronized.

<a name="drs_filter_rules_result_filter_rules_resp"></a>
The `filter_rules_resp` block supports:

* `source_filter_condition` - The source database filter condition SQL expression.

* `target_filter_condition` - The target database filter condition SQL expression.

* `condition_separated` - Whether the source and target database filter conditions are configured separately.

<a name="drs_filter_rules_result_advanced_setting_resp"></a>
The `advanced_setting_resp` block supports:

* `db_name` - The database name.

* `table_name` - The table name.

* `col_names` - The column names.

* `prim_key_or_indexes` - The primary key or unique index.

* `indexes` - The indexes used for optimized queries.

* `values` - The filter condition.
