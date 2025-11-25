---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_app_server_metric_data"
description: |-
  Use this data source to query the monitoring metric data of the Workspace APP server within HuaweiCloud.
---

# huaweicloud_workspace_app_server_metric_data

Use this data source to query the monitoring metric data of the Workspace APP server within HuaweiCloud.

## Example Usage

```hcl
variable "server_id" {}
variable "start_time" {}
variable "end_time" {}

data "huaweicloud_workspace_app_server_metric_data" "test" {
  server_id   = var.server_id
  namespace   = "SYS.ECS"
  metric_name = "cpu_util"
  from        = var.start_time
  to          = var.end_time
  period      = 1
  filter      = "average"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the server metric data are located.  
  If omitted, the provider-level region will be used.

* `server_id` - (Required, String) Specifies the ID of the server.

* `namespace` - (Required, String) Specifies the namespace of the service.  
  The valid values are as follows:
  + **SYS.ECS**: Basic monitoring metrics of Elastic Cloud Server (ECS).
  + **AGT.ECS**: Operating system monitoring metrics of ECS (GPU metrics).

* `metric_name` - (Required, String) Specifies the name of the monitoring metric.
  + For **SYS.ECS** namespace, please refer to the [documentation](https://support.huaweicloud.com/intl/en-us/usermanual-ecs/ecs_03_1002.html).
  + For **AGT.ECS** namespace, please refer to the [documentation](https://support.huaweicloud.com/intl/en-us/usermanual-ecs/ecs_03_1003.html#section11).

* `from` - (Required, String) Specifies the start time of the query data, in RFC3339 format.

* `to` - (Required, String) Specifies the end time of the query data, in RFC3339 format.

* `period` - (Required, Int) Specifies the granularity of monitoring data.  
  The valid values are as follows:
  + **1**: Real-time data.
  + **300**: `5` minute granularity.
  + **1200**: `20` minute granularity.
  + **3600**: `1` hour granularity.
  + **14400**: `4` hour granularity.
  + **86400**: `1` day granularity.

* `filter` - (Required, String) Specifies the data aggregation method.  
  The valid values are as follows:
  + **average**: The average value of the metric data within the aggregation period.
  + **max**: The maximum value of the metric data within the aggregation period.
  + **min**: The minimum value of the metric data within the aggregation period.
  + **sum**: The sum value of the metric data within the aggregation period.
  + **variance**: The variance of the metric data within the aggregation period.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `metrics` - The list of server metric data.  
  The [metrics](#app_server_metric_data) structure is documented below.

<a name="app_server_metric_data"></a>
The `metrics` block supports:

* `metric_name` - The name of the monitoring metric.

* `dimension_value` - The dimension value.  
  This field is valid only when querying GPU monitoring information.

* `datapoints` - The list of metric data points.  
  The [datapoints](#app_server_metric_data_datapoints) structure is documented below.

<a name="app_server_metric_data_datapoints"></a>
The `datapoints` block supports:

* `average` - The average value of the metric data within the aggregation period.  
  This field is valid only when the `filter` is set to `average`.

* `max` - The maximum value of the metric data within the aggregation period.  
  This field is valid only when the `filter` is set to `max`.

* `min` - The minimum value of the metric data within the aggregation period.  
  This field is valid only when the `filter` is set to `min`.

* `sum` - The sum value of the metric data within the aggregation period.  
  This field is valid only when the `filter` is set to `sum`.

* `variance` - The variance of the metric data within the aggregation period.  
  This field is valid only when the `filter` is set to `variance`.

* `collection_time` - The collection time of the metric, in RFC3339 format.

* `unit` - The unit of the metric.
