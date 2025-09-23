---
subcategory: "Cloud Eye (CES)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ces_dashboard"
description: |-
  Manages a CES dashboard resource within HuaweiCloud.
---

# huaweicloud_ces_dashboard

Manages a CES dashboard resource within HuaweiCloud.

## Example Usage

### Basic example

```hcl
variable "name" {}
variable "row_widget_num" {}

resource "huaweicloud_ces_dashboard" "test" {
  name           = var.name
  row_widget_num = var.row_widget_num
}
```

### Copy dashboard example

```hcl
variable "name" {}
variable "dashboard_id" {}
variable "row_widget_num" {}

resource "huaweicloud_ces_dashboard" "test" {
  name           = var.name
  dashboard_id   = var.dashboard_id
  row_widget_num = var.row_widget_num
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `name` - (Required, String) Specifies the dashboard name.

* `row_widget_num` - (Required, Int) Specifies the monitoring view display mode.
  The valid values are as follows:
  + **0**: custom coordinates;
  + **1**: one per row;
  + **2**: two per row;
  + **3**: three per row.

* `dashboard_id` - (Optional, String, NonUpdatable) Specifies the copied dashboard ID.

* `enterprise_project_id` - (Optional, String, NonUpdatable) Specifies the enterprise project ID of the dashboard.

* `is_favorite` - (Optional, Bool) Specifies whether the dashboard is favorite.

* `extend_info` - (Optional, List) Specifies the information about the extension.
  The [extend_info](#extend_info_struct) structure is documented below.

<a name="extend_info_struct"></a>
The `extend_info` block supports:

* `filter` - (Optional, String) Specifies the metric aggregation method.
  Values can be as follows:
  + **average**: Average value.
  + **min**: Minimum value.
  + **max**: Maximum value.
  + **sum**: Sum.

* `period` - (Optional, String) Specifies the metric aggregation period.
  Values can be as follows:
  + **1**: Original value.
  + **60**: One minute.
  + **300**: Five minutes.
  + **1200**: Twenty minutes.
  + **3600**: One hour.
  + **14400**: Four hours.
  + **86400**: One day.

* `display_time` - (Optional, Int) Specifies the display time.
  Values can be as follows:
  + **0**: Using custom time display.
  + **5**: Five minutes.
  + **15**: Fifteen minutes.
  + **30**: Thirty minutes.
  + **60**: One hour.
  + **120**: Two hours.
  + **180**: Three hours.
  + **720**: Twelve hours.
  + **1440**: One day.
  + **10080**: Seven days.
  + **43200**: Thirty days.

* `refresh_time` - (Optional, Int) Specifies the refresh time.
  Values can be as follows:
  + **0**: No refresh.
  + **10**: Ten seconds.
  + **60**: One minute.
  + **300**: Five minutes.
  + **1200**: Twenty minutes.

* `from` - (Optional, Int) Specifies the start time.

* `to` - (Optional, Int) Specifies the end time.

* `screen_color` - (Optional, String) Specifies the monitoring screen background color.

* `enable_screen_auto_play` - (Optional, Bool) Specifies whether the monitoring screen switches automatically.

* `time_interval` - (Optional, Int) Specifies the automatic switching time interval of the monitoring screen.
  Values can be as follows:
  + **10000**: Ten seconds.
  + **30000**: Thirty seconds.
  + **60000**: One minute.

* `enable_legend` - (Optional, Bool) Specifies whether to enable the legend.

* `full_screen_widget_num` - (Optional, Int) Specifies the number of large screen display views.
  Values can be **1**, **4**, **9**, **16** and **25**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `created_at` - Indicates the creation time of the dashboard.

* `creator_name` - Indicates the creator name of the dashboard.

* `namespace` - Indicates the namespace.

* `sub_product` - Indicates the sub-product ID.

* `dashboard_template_id` - Indicates the monitoring disk template ID.

* `widgets_num` - Indicates the total number of views under the board.

## Import

The dashboard can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_ces_dashboard.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attribute is `dashboard_id`.
It is generally recommended running `terraform plan` after importing the resource.
You can then decide if changes should be applied to the dashboard, or the resource definition should be updated to
align with the dashboard. Also you can ignore changes as below.

```hcl
resource "huaweicloud_ces_dashboard" "test" {
    ...

  lifecycle {
    ignore_changes = [
      dashboard_id,
    ]
  }
}
```
