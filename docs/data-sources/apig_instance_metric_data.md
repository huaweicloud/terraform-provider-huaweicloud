---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_instance_metric_data"
description: |-
  Use this data source to query the metric data of the dedicated instance within HuaweiCloud.
---

# huaweicloud_apig_instance_metric_data

Use this data source to query the metric data of the dedicated instance within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}
variable "from_unix_timestamp" {}
variable "to_unix_timestamp" {}

data "huaweicloud_apig_instance_metric_data" "test" {
  instance_id = var.instance_id
  dim         = "inbound_eip"
  metric_name = "upstream_bandwidth"
  from        = var.from_unix_timestamp
  to          = var.to_unix_timestamp
  period      = 300
  filter      = "average"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the dedicated instance is located.  
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the dedicated instance.

* `dim` - (Required, String) Specifies the dimension of the metric data.  
  The valid values are as follows:
  + **inbound_eip**: The inbound public network bandwidth, only supported by ELB type instances.
  + **outbound_eip**: The outbound public network bandwidth.

* `metric_name` - (Required, String) Specifies name of the metric data.  
  The valid values are as follows:
  + **upstream_bandwidth**: The outbound bandwidth.
  + **downstream_bandwidth**: The inbound bandwidth.
  + **upstream_bandwidth_usage**: The outbound bandwidth usage rate.
  + **downstream_bandwidth_usage**: The inbound bandwidth usage rate.
  + **up_stream**: The outbound traffic.
  + **down_stream**: The inbound traffic.

* `from` - (Required, String) Specifies the start time of the metric data, UNIX timestamp in milliseconds.

* `to` - (Required, String) Specifies the end time of the metric data, UNIX timestamp in milliseconds.  
  The value of `from` must be less than the value of `to`.

* `period` - (Required, Int) Specifies the granularity of the metric data.  
  The valid values are as follows:
  + **1**: Real-time data.
  + **300**: `5` minutes granularity.
  + **1,200**: `20` minutes granularity.
  + **3,600**: `1` hour granularity.
  + **14,400**: `4` hours granularity.
  + **86,400**: `1` day granularity.

* `filter` - (Required, String) Specifies the data aggregation method of the metric data.  
  The valid values are as follows:
  + **average**: The average value of the metric data within the aggregation period.
  + **max**: The maximum value of the metric data within the aggregation period.
  + **min**: The minimum value of the metric data within the aggregation period.
  + **sum**: The sum value of the metric data within the aggregation period.
  + **variance**: The variance value of the metric data within the aggregation period.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `datapoints` - The list of the metric data points that matched the filter parameters.  
  The [datapoints](#apig_metric_data_datapoints_attr) structure is documented below.

<a name="apig_metric_data_datapoints_attr"></a>
The `datapoints` block supports:

* `average` - The average value of the metric data within the aggregation period.  
  Required if the `filter` parameter is set to `average`, this field is available.

* `max` - The maximum value of the metric data within the aggregation period.  
  Required if the `filter` parameter is set to `max`, this field is available.

* `min` - The minimum value of the metric data within the aggregation period.  
  Required if the `filter` parameter is set to `min`, this field is available.

* `sum` - The sum value of the metric data within the aggregation period.  
  Required if the `filter` parameter is set to `sum`, this field is available.

* `variance` - The variance value of the metric data within the aggregation period.  
  Required if the `filter` parameter is set to `variance`, this field is available.

* `timestamp` - The collection time of the metric data, UNIX timestamp in milliseconds.

* `unit` - The unit of the metric.
