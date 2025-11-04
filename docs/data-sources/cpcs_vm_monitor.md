---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cpcs_vm_monitor"
description: |-
  Use this data source to query the VM monitoring data from HuaweiCloud CPCS service.
---

# huaweicloud_cpcs_vm_monitor

Use this data source to query the VM monitoring data from HuaweiCloud CPCS service.

## Example Usage

```hcl
data "huaweicloud_cpcs_vm_monitor" "test" {
  namespace   = "ECS"
  metric_name = "mem_util"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `namespace` - (Required, String) Specifies the namespace of the monitoring data.
  Valid values are **ECS** and **DHSM**.

* `metric_name` - (Required, String) Specifies the name of the metric to query.
  Valid values are **cpu_util**, **mem_util**, and **network_outgoing_bytes_rate_inband**.

* `instance_id` - (Optional, String) Specifies the ID of the instance to monitor.

* `vsm_id` - (Optional, String) Specifies the ID of the virtual machine.

* `from` - (Optional, Int) Specifies the start time of the query in milliseconds since epoch.
  Defaults to `0`, which means the query starts from three days ago.

* `to` - (Optional, Int) Specifies the end time of the query in milliseconds since epoch.
  Defaults to `0`, which means the query ends at the current time.

* `period` - (Optional, Int) Specifies the statistical data period.
  Defaults to `0`, which means the default period is `5` minutes.

* `filter` - (Optional, String) Specifies the statistical value type.
  Defaults to **min**, which means querying the minimum value.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `datapoints` - The list of monitoring data points.
  The [datapoints](#datapoints_struct) structure is documented below.

* `metric_name_output` - The name of the metric that was queried.

* `max` - The maximum value of the metric in the time range.

* `average` - The average value of the metric in the time range.

<a name="datapoints_struct"></a>
The `datapoints` block supports:

* `max` - The maximum value of the metric.

* `min` - The minimum value of the metric.

* `average` - The average value of the metric.

* `sum` - The sum of the metric values.

* `variance` - The variance of the metric values.

* `timestamp` - The timestamp of the data point in milliseconds since epoch.

* `unit` - The unit of the metric value.
