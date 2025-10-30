---
subcategory: "Application Operations Management (AOM 2.0)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aom_dashboards"
description: |-
  Use this data source to get the list of AOM dashboards.
---

# huaweicloud_aom_dashboards

Use this data source to get the list of AOM dashboards.

## Example Usage

```hcl
data "huaweicloud_aom_dashboards" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `dashboard_type` - (Optional, String) Specifies the dashboard type.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID to which the dashboard belongs.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `dashboards` - Indicates the dashboards list.
  The [dashboards](#attrblock--dashboards) structure is documented below.

<a name="attrblock--dashboards"></a>
The `dashboards` block supports:

* `id` - Indicates the dashboard ID.

* `dashboard_title` - Indicates the dashboard title.

* `dashboard_title_en` - Indicates the dashboard English title.

* `dashboard_type` - Indicates the dashboard type.

* `display` - Indicates whether the dashboard is displayed.

* `enterprise_project_id` - Indicates the enterprise project ID to which the dashboard belongs.

* `folder_id` - Indicates the folder ID to which the dashboard belongs.

* `folder_title` - Indicates the folder name to which the dashboard belongs.

* `dashboard_tags` - Indicates the dashboard tags.

* `is_favorite` - Indicates whether the dashboard is favorited.

* `version` - Indicates the dashboard version.

* `created_at` - Indicates the create time of the dashboard.

* `created_by` - Indicates the creator of the dashboard.

* `updated_at` - Indicates the update time of the dashboard.

* `updated_by` - Indicates the updator of the dashboard.
