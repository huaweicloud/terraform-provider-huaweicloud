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

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `created_at` - The creation time of the dashboard.

* `creator_name` - The creator name of the dashboard.

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
