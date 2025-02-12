---
subcategory: "Cloud Eye (CES)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ces_multiple_metrics_data"
description: |-
  Use this data source to get the list of CES multiple metrics data.
---

# huaweicloud_ces_multiple_metrics_data

Use this data source to get the list of CES multiple metrics data.

## Example Usage

```hcl
variable "from" {}
variable "to" {}

data "huaweicloud_ces_multiple_metrics_data" "test" {
  metrics {
    namespace   = "YOU.APP"
    metric_name = "cpu_util"

    dimensions {
      name  = "platform_id"
      value = "test_platform_id"
    }

    dimensions {
      name  = "instance_id"
      value = "test_instance_id"
    }

    dimensions {
      name  = "cpu_type"
      value = "test_cpu_type"
    }
  }

  metrics {
    namespace   = "MINE.APP"
    metric_name = "mem_util"

    dimensions {
      name  = "instance_id"
      value = "test_instance_id"
    }

    dimensions {
      name  = "memory_type"
      value = "test_memory_type"
    }
  }

  from   = var.from
  to     = var.to
  period = 1
  filter = "average"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `metrics` - (Required, List) Specifies the metric data. Up to 500 metrics can be specified at a time.
  The [metrics](#Metrics) structure is documented below.

* `from` - (Required, String) Specifies the start time of the query.
  The time is in UTC. The format is **yyyy-MM-dd HH:mm:ss**. The **from** must be earlier than **to**.

* `to` - (Required, String) Specifies the end time of the query.
  The time is in UTC. The format is **yyyy-MM-dd HH:mm:ss**.

* `period` - (Required, String) Specifies how often Cloud Eye aggregates data.
  The valid values are as follows:
  + **1**: Cloud Eye performs no aggregation and displays raw data;
  + **300**: Cloud Eye aggregates data every 5 minutes;
  + **1200**: Cloud Eye aggregates data every 20 minutes;
  + **3600**: Cloud Eye aggregates data every hour;
  + **14400**: Cloud Eye aggregates data every 4 hours;
  + **86400**: Cloud Eye aggregates data every 24 hours;

* `filter` - (Required, String) Specifies the data rollup method.
  The valid value can be **max**, **min**, **average**, **sum** or **variance**.
  The field does not affect the query result of raw data. (The period is **1**.)

<a name="Metrics"></a>
The `metrics` block supports:

* `namespace` - (Required, String) Specifies the namespace of a service.

* `metric_name` - (Required, String) Specifies the metric ID.

* `dimensions` - (Required, List) Specifies metric dimensions.
  The [dimensions](#MetricsDimensions) structure is documented below.

<a name="MetricsDimensions"></a>
The `dimensions` block supports:

* `name` - (Required, String) Specifies the dimension.

* `value` - (Required, String) Specifies the dimension value.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data` - The metric data.
  The [data](#DataAttr) structure is documented below.

<a name="DataAttr"></a>
The `data` block supports:

* `namespace` - The namespace of a service.

* `metric_name` - The metric ID.

* `dimensions` - The metric dimensions.
  The [dimensions](#DataDimensionsAttr) structure is documented below.

* `datapoints` - The metric data list. Up to 3000 data points can be returned.
  The [datapoints](#DataDataPointsAttr) structure is documented below.

* `unit` - The metric unit.

<a name="DataDimensionsAttr"></a>
The `dimensions` block supports:

* `name` - The dimension.

* `value` - The dimension value.

<a name="DataDataPointsAttr"></a>
The `datapoints` block supports:

* `average` - The average value of metric data within a rollup period.

* `max` - The maximum value of metric data within a rollup period.

* `min` - The minimum value of metric data within a rollup period.

* `sum` - The sum of metric data within a rollup period.

* `variance` - The variance of metric data within a rollup period.

* `timestamp` - The time when the metric is collected. The time is a UNIX timestamp and the unit is ms.
