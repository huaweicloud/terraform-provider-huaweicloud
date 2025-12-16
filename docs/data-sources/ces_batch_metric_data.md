---
subcategory: "Cloud Eye (CES)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ces_batch_metric_data"
description: |-
  Using the data source to get the list of metric data in batches.
---

# huaweicloud_ces_batch_metric_data

Using the data source to get the list of metric data in batches.

## Example Usage

```hcl
data "huaweicloud_ces_batch_metric_data" "test" {
  namespace        = "SYS.ECS"
  metric_name      = "cpu_util"
  metric_dimension = "instance_id"
  from             = "2025-12-15 06:36:46"
  to               = "2025-12-15 06:40:46"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `namespace` - (Required, String) Specifies the namespace of a service.

* `metric_name` - (Required, String) Specifies the metric name of a resource.

* `metric_dimension` - (Required, String) Specifies the metric dimension. Multiple dimensions are separated by commas (,).

* `from` - (Optional, String) Specifies the start time of the query.
  The time is in UTC. The format is **yyyy-MM-dd HH:mm:ss**. The `from` must be earlier than `to`.

* `to` - (Optional, String) Specifies the end time of the query.
  The time is in UTC. The format is **yyyy-MM-dd HH:mm:ss**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data_points` - Indicates the metric data list.
  The [data_points](#data_points_struct) structure is documented below.

<a name="data_points_struct"></a>
The `data_points` block supports:

* `dimensions` - Indicates the dimension information.
  The [dimensions](#dimensions_struct) structure is documented below.

* `timestamp` - Indicates the metric collection time.

* `value` - Indicates the metric value.

* `unit` - Indicates the data unit.

<a name="dimensions_struct"></a>
The `dimensions` block supports:

* `name` - Indicates the metric dimension name.

* `value` - Indicates the metric dimension value.
