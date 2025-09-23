---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_metric_results"
description: |-
  Use this data source to get the list of SecMaster metric results.
---

# huaweicloud_secmaster_metric_results

Use this data source to get the list of SecMaster metric results.

## Example Usage

```hcl
variable "workspace_id" {}
variable "metric_ids" {}

data "huaweicloud_secmaster_metric_results" "test" {
  workspace_id = var.workspace_id
  metric_ids   = var.metric_ids
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID.

* `metric_ids` - (Required, List) Specifies the metrics IDs.

* `timespan` - (Optional, String) Specifies the time range for querying metrics.
  The format is ISO8601, for example, **2007-03-01T13:00:00Z/2008-05-11T15:30:00Z**,
  **2007-03-01T13:00:00Z/P1Y2M10DT2H30M**, or **P1Y2M10DT2H30M/2008-05-11T15:30:00Z**.

* `cache` - (Optional, String) Specifies whether the cache is enabled.

* `field_ids` - (Optional, List) Specifies the indicator card IDs.

* `params` - (Optional, List) Specifies the parameter list of the metric.
  The number of elements must be the same as that of the metric_ids list.
  For details, see [About Metrics](https://support.huaweicloud.com/intl/en-us/api-secmaster/secmaster_03_0028.html)

* `interactive_params` - (Optional, List) Specifies the interactive parameters.
  For details, see [About Metrics](https://support.huaweicloud.com/intl/en-us/api-secmaster/secmaster_03_0028.html)

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `metric_results` - The metric results.
  The [metric_results](#metric_results) structure is documented below.

<a name="metric_results"></a>
The `metric_results` block supports:

* `id` - The metric ID.

* `labels` - The statistical labels of the metric.
  The value in the label corresponds to the value in a piece of `data_row` one-to-one.

* `data_rows` - All statistical data of the metric.
  The [data_rows](#metric_results_data_rows) structure is documented below.

<a name="metric_results_data_rows"></a>
The `data_rows` block supports:

* `data_row` - A piece of data in the metric results.
