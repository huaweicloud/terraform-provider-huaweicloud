---
subcategory: "Cloud Eye (CES)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ces_dashboard_widgets"
description: |-
  Use this data source to get the list of CES dashboard widgets.
---

# huaweicloud_ces_dashboard_widgets

Use this data source to get the list of CES dashboard widgets.

## Example Usage

```hcl
variable "dashboard_id" {}

data "huaweicloud_ces_dashboard_widgets" "test" {
  dashboard_id = var.dashboard_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `dashboard_id` - (Required, String) Specifies the dashboard ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `widgets` - The dashboard widget list.

  The [widgets](#widgets_struct) structure is documented below.

<a name="widgets_struct"></a>
The `widgets` block supports:

* `widget_id` - The dashboard widget ID.

* `metric_display_mode` - How many metrics will be displayed on one widget.

* `properties` - The additional information.

  The [properties](#widgets_properties_struct) structure is documented below.

* `view` - The graph type.

* `title` - The dashboard widget title.

* `threshold` - The threshold of metrics on the graph.

* `threshold_enabled` - Whether to display the threshold of metrics.

* `location` - The dashboard widget coordinates.

  The [location](#widgets_location_struct) structure is documented below.

* `unit` - The metric unit.

* `metrics` - The metric list.

  The [metrics](#widgets_metrics_struct) structure is documented below.

* `created_at` - When the dashboard widget was created.

<a name="widgets_properties_struct"></a>
The `properties` block supports:

* `filter` - How metric data is aggregated.

* `top_n` - The top n resources sorted by a metric.

* `order` - How top n resources by a metric are sorted on a widget.

<a name="widgets_location_struct"></a>
The `location` block supports:

* `top` - The grids between the widget and the top of the dashboard.

* `left` - The grids between the widget and the left side of the dashboard.

* `width` - The dashboard widget width.

* `height` - The dashboard widget height.

<a name="widgets_metrics_struct"></a>
The `metrics` block supports:

* `alias` - The alias list of metrics on the dashboard widget.

* `namespace` - The cloud service dimension.

* `dimensions` - The dimension list.

  The [dimensions](#metrics_dimensions_struct) structure is documented below.

* `metric_name` - The metric name.

<a name="metrics_dimensions_struct"></a>
The `dimensions` block supports:

* `values` - The dimension value list.

* `name` - The dimension name.

* `filter_type` - The resource type.
