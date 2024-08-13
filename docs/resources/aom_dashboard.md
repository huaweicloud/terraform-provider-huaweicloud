---
subcategory: "Application Operations Management (AOM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aom_dashboard"
description:  |-
  Manages an AOM dashboard resource within HuaweiCloud.
---

# huaweicloud_aom_dashboard

Manages an AOM dashboard resource within HuaweiCloud.

## Example Usage

```hcl
variable "dashboard_title " {}
variable "folder_title " {}
variable "dashboard_type " {}

resource "huaweicloud_aom_dashboard" "test" {
  dashboard_title = var.dashboard_title
  folder_title    = var.folder_title
  dashboard_type  = var.dashboard_type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `dashboard_title` - (Required, String) Specifies the dashboard title.

* `folder_title` - (Required, String) Specifies the folder title.

* `dashboard_type` - (Required, String) Specifies the dashboard type. It's customized by user.

* `is_favorite` - (Optional, Bool) Specifies whether to favorite the dashboard. Defaults to **false**.

* `dashboard_tags` - (Optional, List) Specifies the dashboard tags. It's an array of map.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project ID to which the dashboards
  belongs. Defaults to **0**. Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates the resource ID.

## Import

The AOM dashboard resource can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_aom_dashboard.test <id>
```
