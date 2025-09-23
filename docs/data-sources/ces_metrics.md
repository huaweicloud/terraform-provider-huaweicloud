---
subcategory: "Cloud Eye (CES)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ces_metrics"
description: |-
  Use this data source to get the list of CES metrics.
---

# huaweicloud_ces_metrics

Use this data source to get the list of CES metrics.

## Example Usage

```hcl
variable "namespace" {}
variable "metric_name" {}

data "huaweicloud_ces_metrics" "test" {
  namespace   = var.namespace
  metric_name = var.metric_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `dim_0` - (Optional, String) The first metric dimension.

* `dim_1` - (Optional, String) The second metric dimension.

* `dim_2` - (Optional, String) The third metric dimension.

* `metric_name` - (Optional, String) The metric name.

* `namespace` - (Optional, String) The metric namespace.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `metrics` - The metric information list.

  The [metrics](#metrics_struct) structure is documented below.

<a name="metrics_struct"></a>
The `metrics` block supports:

* `namespace` - The metric namespace.

* `unit` - The metric unit.

* `dimensions` - The metric dimension list.

  The [dimensions](#metrics_dimensions_struct) structure is documented below.

* `metric_name` - The metric name.

<a name="metrics_dimensions_struct"></a>
The `dimensions` block supports:

* `name` - The resource dimension name.

* `value` - The resource dimension value.
