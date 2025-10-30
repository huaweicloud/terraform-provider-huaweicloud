---
subcategory: "Application Operations Management (AOM 2.0)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aom_dashboards_folders"
description: |-
  Use this data source to get the list of AOM dashboards folders.
---

# huaweicloud_aom_dashboards_folders

Use this data source to get the list of AOM dashboards folders.

## Example Usage

```hcl
data "huaweicloud_aom_dashboards_folders" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID to which the folder belongs.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `folders` - Indicates the dashboards folders list.
  The [folders](#folders_struct) structure is documented below.

<a name="folders_struct"></a>
The `folders` block supports:

* `id` - Indicates the folder ID.

* `folder_title` - Indicates the folder title.

* `folder_title_en` - Indicates the folder English title.

* `display` - Indicates whether to display the folder.

* `enterprise_project_id` - Indicates the enterprise project ID to which the folder belongs.

* `dashboard_ids` - Indicates the dashboard IDs under the folder.

* `is_template` - Indicates whether the folder is default.

* `created_by` - Indicates the creator of the folder.

* `created_at` - Indicates the create time of the folder.

* `updated_at` - Indicates the update time of the folder.
