---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_metric"
description: |-
  Manages a metric resource within HuaweiCloud.
---

# huaweicloud_secmaster_metric

Manages a metric resource within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}
variable "metric_name" {}

resource "huaweicloud_secmaster_metric" "test" {
  workspace_id     = var.workspace_id
  name             = var.metric_name
  metric_type      = "LOGGING"
  data_type        = "STATISTICS"
  metric_dimension = 1
  cache_ttl        = 10
  report_period    = 0
  is_built_in      = false
  max_query_range  = 5
  version          = "0.1"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `workspace_id` - (Required, String, NonUpdatable) Specifies the ID of the workspace to which the metric belongs.

* `name` - (Required, String) Specifies the name of the metric.

* `metric_type` - (Required, String, NonUpdatable) Specifies the type of the metric.
  Currently, only **LOGGING** (log type) is supported for creation.

* `data_type` - (Required, String, NonUpdatable) Specifies the data type of the metric.
  Currently, only **STATISTICS** is supported for creation.

* `cache_ttl` - (Required, Int) Specifies the cache lifecycle in seconds.

* `metric_dimension` - (Optional, Int) Specifies the metric result dimension.
  **0** means a single number, **2** means a chart or table, **3+** means a multi-label chart.
  Required when `metric_type` is **DERIVED**; must be **0** when `metric_type` is **COMPOUND**.

* `report_period` - (Optional, Int) Specifies the report period in seconds.
  Required for buried point metrics.

* `is_built_in` - (Optional, Bool, NonUpdatable) Specifies whether the metric is a system metric. Defaults to **false**.

* `effective_column` - (Optional, String) Specifies the effective column.
  When this parameter is present, the specified column is used as the metric data result.

* `max_query_range` - (Optional, Int) Specifies the maximum query range of the metric in days.
  For compound metrics, the value is the minimum of all elements in the `derived_metrics` list.

* `derived_metrics` - (Optional, List) Specifies the derived metrics list.
  For non-compound metrics, there is only one element; for compound metrics, it contains the definitions
  of each derived metric.
  The [derived_metrics](#derived_metrics_block) structure is documented below.

* `compound_expression` - (Optional, String) Specifies the compound expression.
  Required when `metric_type` is **DERIVED**.

* `metric_format` - (Optional, List) Specifies the metric format list.
  The [metric_format](#metric_format_block) structure is documented below.

* `metric_expand_dim` - (Optional, List) Specifies the metric dimension expansion parameters.
  The [metric_expand_dim](#metric_expand_dim_block) structure is documented below.

* `version` - (Optional, String, NonUpdatable) Specifies the SecMaster version.

<a name="derived_metrics_block"></a>
The `derived_metrics` block supports:

* `metric_dimension` - (Required, Int) Specifies the derived metric result dimension.
  **0** means a single number, **2** means a chart or table, **3+** means a multi-label chart.

* `max_query_range` - (Optional, Int) Specifies the maximum query range of the metric in days.

* `date_start` - (Optional, String) Specifies the relative start time of the metric query range
  (datemath expression).

* `date_end` - (Optional, String) Specifies the relative end time of the metric query range
  (datemath expression).

* `date_format` - (Optional, String) Specifies the time format.
  Valid values: **epoch_millis**, **epoch_second**, **yyyy-MM-dd'T'HH:mm:ss.SSSZ**.

* `query_type` - (Required, String) Specifies the method to obtain metric results.
  Valid values: **cbsl**, **api**, **dsl**, **sql**.

* `query_function` - (Required, String) Specifies the function to obtain metric results, escaped as a string.

<a name="metric_format_block"></a>
The `metric_format` block supports:

* `data` - (Optional, String) Specifies the data format.

* `display` - (Optional, String) Specifies the display format.

* `display_param` - (Optional, Map) Specifies the display parameters.

* `data_param` - (Optional, Map) Specifies the data parameters.

<a name="metric_expand_dim_block"></a>
The `metric_expand_dim` block supports:

* `labels` - (Required, List) Specifies the dimension expansion labels.

* `functions` - (Required, List) Specifies the dimension expansion functions.
  Fill in the built-in methods of the metric data plane, with parameter index starting from 1.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which is the metric ID.

## Import

The metric can be imported using the `workspace_id` and their `id`, separated by a slash (/), e.g.

```bash
$ terraform import huaweicloud_secmaster_metric.test <workspace_id>/<metric_id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response. The missing attributes include `version`. It is generally recommended running `terraform plan`
after importing a resource. You can then decide if changes should be applied to the resource,
or the resource definition should be updated to align with the resource. Also you can ignore changes as below.

```hcl
resource "huaweicloud_secmaster_metric" "test" {
  ...

  lifecycle {
    ignore_changes = [
      version,
    ]
  }
}
```
