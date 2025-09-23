---
subcategory: "Cloud Eye (CES)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ces_dashboards"
description: |-
  Use this data source to get the list of CES dashboards.
---

# huaweicloud_ces_dashboards

Use this data source to get the list of CES dashboards.

## Example Usage

```hcl
data "huaweicloud_ces_dashboards" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Specifies the dashboard name.

* `dashboard_id` - (Optional, String) Specifies the dashboard ID.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.

* `is_favorite` - (Optional, Bool) Specifies whether a dashboard in an enterprise project is added to favorites.
  The value can be **true** (added to favorites) or **false** (not added to favorites).
  If this parameter is specified, **enterprise_project_id** is mandatory.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `dashboards` - The dashboard list.

  The [dashboards](#dashboards_struct) structure is documented below.

<a name="dashboards_struct"></a>
The `dashboards` block supports:

* `dashboard_id` - The dashboard ID.

* `name` - The name of the dashboard.

* `enterprise_project_id` - The enterprise project ID.

* `creator_name` - The creator of the dashboard.

* `created_at` - The creation time of the dashboard.

* `row_widget_num` - The monitoring view display mode.

* `is_favorite` - Whether a dashboard is added to favorites.
