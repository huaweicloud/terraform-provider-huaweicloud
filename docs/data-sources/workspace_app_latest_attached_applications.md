---
subcategory: "Workspace"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_workspace_app_latest_attached_applications"
description: |-
  Use this data source to get the latest application list attached to the Workspace APP image server within HuaweiCloud.
---

# huaweicloud_workspace_app_latest_attached_applications

Use this data source to get the latest application list attached to the Workspace APP image server within HuaweiCloud.

## Example Usage

```hcl
variable "server_id" {}

data "huaweicloud_workspace_app_latest_attached_applications" "test" {
  server_id = var.server_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.  
  If omitted, the provider-level region will be used.

* `server_id` - (Required, String) Specifies the ID of the image server instance.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `applications` - The list of latest attached applications.
  The [applications](#app_latest_attached_applications) structure is documented below.

<a name="app_latest_attached_applications"></a>
The `applications` block supports:

* `app_id` - The ID of the attached application.

* `record_id` - The record ID of the attached application.
