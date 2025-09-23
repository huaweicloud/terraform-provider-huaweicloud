---
subcategory: "Cloud Eye (CES)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ces_dashboard_widget"
description: |-
  Manages a CES dashboard widget resource within HuaweiCloud.
---

# huaweicloud_ces_dashboard_widget

 Manages a CES dashboard widget resource within HuaweiCloud.

## Example Usage

```hcl
variable "title" {}
variable "dashboard_id" {}
variable "instance_id" {}
variable "left" {}
variable "top" {}
variable "width" {}
variable "height" {}

resource "huaweicloud_ces_dashboard_widget" "test" {
  dashboard_id        = var.dashboard_id
  title               = var.title
  view                = "line"
  metric_display_mode = "single"

  metrics {
    metric_name = "cpu_util"
    namespace   = "SYS.ECS"

    dimensions  {
      name        = "instance_id"
      filter_type = "specific_instances"
      values      = [var.instance_id]
    }
  }

  location {
    left   = var.left
    top    = var.top
    width  = var.width
    height = var.height
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `dashboard_id` - (Required, String, NonUpdatable) Specifies the dashboard ID.

* `title` - (Required, String) Specifies the dashboard widget title.

* `metrics` - (Required, List) Specifies the metric list.

  The [metrics](#Metrics) structure is documented below.

* `view` - (Required, String, NonUpdatable) Specifies the graph type.
  The valid values are as follows:
  + **bar**: bar chart.
  + **line**: line chart.
  + **bar_chart**: histogram.
  + **table**: table.
  + **circular_bar**: pie chart.
  + **area_chart**: area chart.

* `metric_display_mode` - (Required, String) Specifies how many metrics will be displayed on one widget.
  The valid values are as follows:
  + **single**: one metric.
  + **multiple**: multiple metrics.

* `location` - (Required, List) Specifies the dashboard widget coordinates.
  
  The [location](#Location) structure is documented below.

* `properties` - (Optional, List) Specifies additional information.
  
  The [properties](#Properties) structure is documented below.

* `unit` - (Optional, String) Specifies the metric unit.

<a name="Metrics"></a>
The `metrics` block supports:

* `namespace` - (Required, String) Specifies the cloud service dimension.

* `dimensions` - (Required, List) Specifies the dimension list.
  
  The [dimensions](#MetricsDimensions) structure is documented below.

* `metric_name` - (Required, String) Specifies the metric name.

* `alias` - (Optional, List) Specifies the alias list of metrics.

<a name="Location"></a>
The `location` block supports:

* `top` - (Required, Int) Specifies the grids between the widget and the top of the dashboard.

* `left` - (Required, Int) Specifies the grids between the widget and the left side of the dashboard.

* `width` - (Required, Int) Specifies the dashboard widget width.

* `height` - (Required, Int) Specifies the dashboard widget height.

<a name="Properties"></a>
The `properties` block supports:

* `top_n` - (Required, Int) Specifies the top n resources sorted by a metric.

* `filter` - (Optional, String) Specifies how metric data is aggregated.
  The value can only be **topN**.

* `order` - (Optional, String) Specifies how top n resources by a metric are sorted on a dashboard widget.
  The value can be **asc** or **desc**.

<a name="MetricsDimensions"></a>
The `dimensions` block supports:

* `filter_type` - (Required, String) Specifies the resource type.
  The value can be **all_instances** (all resources) or **specific_instances** (specified resources).

* `name` - (Required, String) Specifies the dimension name.

* `values` - (Optional, List) Specifies the dimension value list.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `created_at` - When the dashboard widget was created.

## Import

The dashboard widget can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_ces_dashboard_widget.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason.The missing attributes is `dashboard_id`.
It is generally recommended running `terraform plan` after importing the resource.
You can then decide if changes should be applied to the dashboard widget, or the resource definition should be updated to
align with the dashboard widget. Also you can ignore changes as below.

```hcl
resource "huaweicloud_ces_dashboard_widget" "test" {
    ...

  lifecycle {
    ignore_changes = [
      dashboard_id,
    ]
  }
}
```
