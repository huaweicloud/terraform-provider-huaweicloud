---
subcategory: "Cloud Eye (CES)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ces_metric_data"
description: |-
  Using the data source to get the list of CES metric data.
---

# huaweicloud_ces_metric_data

Using the data source to get the list of CES metric data.

## Example Usage

```hcl
variable "from" {}
variable "to" {}

data "huaweicloud_ces_metric_data" "test" {
  namespace   = "You.APP"
  metric_name = "cpu_util"
  dim_0       = "platform_id,test_platform_id"
  filter      = "average" 
  period      = 1  
  from        = var.from
  to          = var.to
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `namespace` - (Required, String) Specifies the metric namespace.

* `metric_name` - (Required, String) Specifies the resource metric name.

* `dim_0` - (Required, String) Specifies the level-1 dimension of a metric.
  The dimension value is in the **key,value** format.

* `filter` - (Required, String) Specifies the data aggregation method.
  The valid value can be **max**, **min**, **average**, **sum** or **variance**.

* `period` - (Required, Int) Specifies how often Cloud Eye aggregates data.
  The valid values are as follows:
  + **1**: Cloud Eye performs no aggregation and displays raw data;
  + **300**: Cloud Eye aggregates data every 5 minutes;
  + **1200**: Cloud Eye aggregates data every 20 minutes;
  + **3600**: Cloud Eye aggregates data every hour;
  + **14400**: Cloud Eye aggregates data every 4 hours;
  + **86400**: Cloud Eye aggregates data every 24 hours;

* `from` - (Required, String) Specifies the start time of the query.
  The time is in UTC. The format is **yyyy-MM-dd HH:mm:ss**. The **from** must be earlier than **to**.

* `to` - (Required, String) Specifies the end time of the query.
  The time is in UTC. The format is **yyyy-MM-dd HH:mm:ss**.

* `dim_1` - (Optional, String) Specifies the level-2 dimension of a metric.

* `dim_2` - (Optional, String) Specifies the level-3 dimension of a metric.

* `dim_3` - (Optional, String) Specifies the level-4 dimension of a metric.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `datapoints` - The metric data list.

  The [datapoints](#datapoints_struct) structure is documented below.

<a name="datapoints_struct"></a>
The `datapoints` block supports:

* `max` - The maximum value of metric data within a rollup period.

* `min` - The minimum value of metric data within a rollup period.

* `average` - The average value of metric data within a rollup period.

* `sum` - The sum of metric data within a rollup period.

* `variance` - The variance of metric data within a rollup period.

* `timestamp` - The time when the metric is collected. The time is a UNIX timestamp and the unit is ms.

* `unit` - The metric unit.
