---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_metrics"
description: |-
  Use this data source to get the list of SecMaster metrics within HuaweiCloud.
---

# huaweicloud_secmaster_metrics

Use this data source to get the list of SecMaster metrics within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}

data "huaweicloud_secmaster_metrics" "test" {
  workspace_id = var.workspace_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data` - The result data list.

  The [data](#data_struct) structure is documented below.

<a name="data_struct"></a>
The `data` block supports:

* `name` - The metric name.

* `id` - The metric ID.

* `metric_type` - The metric type.

* `data_type` - The data type.

* `metric_dimension` - The metric result dimension.

* `cache_ttl` - The cache TTL, in seconds.

* `report_period` - The report period, in seconds.

* `is_built_in` - Whether the metric is built-in.

* `effective_column` - The effective column.

* `max_query_range` - The maximum query range, in days.

* `derived_metrics` - The derived metrics list.

  The [derived_metrics](#derived_metrics_struct) structure is documented below.

* `compound_expression` - The compound expression.

* `metric_format` - The metric format.

  The [metric_format](#metric_format_struct) structure is documented below.

* `metric_expand_dim` - The metric dimension expand parameter.

  The [metric_expand_dim](#metric_expand_dim_struct) structure is documented below.

* `version` - The SecMaster version.

<a name="derived_metrics_struct"></a>
The `derived_metrics` block supports:

* `metric_dimension` - The derived metric result dimension.

* `max_query_range` - The maximum query range, in days.

* `date_start` - The relative start time of the metric query range.

* `date_end` - The relative end time of the metric query range.

* `date_format` - The time format.

* `query_type` - The query type.

* `query_function` - The query function.

<a name="metric_format_struct"></a>
The `metric_format` block supports:

* `data` - The data format.

* `display` - The display format.

* `display_param` - The display parameters.

* `data_param` - The data parameters.

<a name="metric_expand_dim_struct"></a>
The `metric_expand_dim` block supports:

* `labels` - The dimension expand labels.

* `functions` - The dimension expand functions.
