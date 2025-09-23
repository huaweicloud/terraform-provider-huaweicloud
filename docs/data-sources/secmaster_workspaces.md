---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_workspaces"
description: |-
  Use this data source to get the list of SecMaster workspaces.
---

# huaweicloud_secmaster_workspaces

Use this data source to get the list of SecMaster workspaces.

## Example Usage

```hcl
variable "name" {}

data "huaweicloud_secmaster_workspaces" "test" {
  name = var.name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `ids` - (Optional, String) Specifies the workspace IDs, which is separated by commas (,).

* `name` - (Optional, String) Specifies the workspace name.

* `description` - (Optional, String) Specifies the workspace description.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID.

* `is_view` - (Optional, String) Specifies whether to query the view. The value can be **true** or **false**.

* `view_bind_id` - (Optional, String) Specifies the space ID bound to the view.

* `view_bind_name` - (Optional, String) Specifies the space name bound to the view.

* `create_time_start` - (Optional, String) Specifies the creation start time, for example, 2024-04-26T16:08:09Z+0800.

* `create_time_end` - (Optional, String) Specifies the creation end time, for example, 2024-04-2T16:08:09Z+0800.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `workspaces` - The list of workspaces.

  The [workspaces](#workspaces_struct) structure is documented below.

<a name="workspaces_struct"></a>
The `workspaces` block supports:

* `id` - The workspace ID.

* `name` - The workspace name.

* `description` - The workspace description.

* `enterprise_project_id` - The enterprise project ID.

* `enterprise_project_name` - The enterprise project name.

* `is_view` - Whether the view is used.

* `view_bind_id` - The space ID bound to the view.

* `view_bind_name` - The space name bound to the view.

* `created_at` - The creation time.

* `updated_at` - The update time.
