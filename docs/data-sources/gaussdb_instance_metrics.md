---
subcategory: "GaussDB"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_gaussdb_instance_metrics"
description: |-
  Use this data source to query the metric data of a specified GaussDB instance within HuaweiCloud.
---

# huaweicloud_gaussdb_instance_metrics

Use this data source to query the metric data of a specified GaussDB instance within HuaweiCloud.

## Example Usage

### Basic Usage

```hcl
variable "instance_id" {}
variable "metric" {
  type = list(string)
}
variable "node_ids" {
  type = list(string)
}

data "huaweicloud_gaussdb_instance_metrics" "test" {
  instance_id = var.instance_id
  start_time  = "1756971683303"
  end_time    = "1756975283303"
  metric      = var.metric
  node_id     = var.node_id
}
```

### Query with Component ID

```hcl
variable "instance_id" {}
variable "metric" {
  type = list(string)
}
variable "node_ids" {
  type = list(string)
}
variable "component_id" {
  type = list(string)
}

data "huaweicloud_gaussdb_instance_metrics" "test" {
  instance_id  = var.instance_id
  start_time   = "1756971683303"
  end_time     = "1756975283303"
  metric       = var.metric
  node_id      = var.node_id
  component_id = var.component_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the instance metrics.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the GaussDB instance.

* `start_time` - (Required, String) Specifies the start time of the query, in timestamp format
  (e.g., `1756971683303`).

* `end_time` - (Required, String) Specifies the end time of the query, in timestamp format
  (e.g., `1756975283303`).

* `metric` - (Required, List) Specifies the metric IDs to query.
  You can obtain the metric IDs through the
  [Query Metric Names API](https://support.huaweicloud.com/intl/zh-cn/api-gaussdb/gaussdb_api_214.html).
  For example, `rds001_cpu_util` for CPU utilization.

* `node_id` - (Required, List) Specifies the node IDs to query.

* `component_id` - (Optional, List) Specifies the component IDs to query (e.g., `dn_6001`).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `metrics` - The metric data collection.
  The [metrics](#instance_metrics_metrics_attr) structure is documented below.

<a name="instance_metrics_metrics_attr"></a>
The `metrics` block supports:

* `metric` - The metric ID.

* `type` - The metric type.
  The valid values are as follows:
  + **INSTANCE**: Instance type.
  + **NODE**: Node type.
  + **COMPONENT**: Component type.

* `unit` - The metric unit.

* `datapoints` - The metric dimension and metric values.
  The [datapoints](#instance_metrics_datapoints_attr) structure is documented below.

* `timestamps` - The timestamps.

<a name="instance_metrics_datapoints_attr"></a>
The `datapoints` block supports:

* `datapoint_name` - The metric item name. For instance metrics, it is the instance ID;
  for node metrics, it is the node name; for component metrics, it is the component name.

* `datapoint_values` - The metric value collection.
