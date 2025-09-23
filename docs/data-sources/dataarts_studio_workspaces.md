---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_studio_workspaces"
description: ""
---

# huaweicloud_dataarts_studio_workspaces

Use this data source to get a list of DataArts Studio workspaces within HuaweiCloud.

## Example Usage

```hcl
variable "instance_id" {}

data "huaweicloud_dataarts_studio_workspaces" "all"{
  instance_id = var.instance_id
}
```

## Argument Reference

* `region` - (Optional, String) Specifies the region in which to obtain the workspaces.
  If omitted, the provider-level region will be used.

* `instance_id` - (Required, String) Specifies the ID of the instance which the workspaces in.

* `name` - (Optional, String) Specifies the workspace name used to filter results.

* `workspace_id` - (Optional, String) Specifies the workspace ID used to filter results.

* `created_by` - (Optional, String) Specifies the user creating workspaces used to filter results.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID used to filter results.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `workspaces` - The filtered workspaces.The [workspaces](#block-workspaces) structure is documented below.

<a name="block-workspaces"></a>
The `workspaces` block supports:

* `name` - The workspace name.

* `id` - The workspace ID.

* `created_by` - The user creating the workspace.

* `enterprise_project_id` - The enterprise project ID of workspace.

* `description` - The description of workspace.

* `bad_record_location_name` - The bad record location name of workspace.

* `job_log_location_name` - The job log location name of workspace.

* `member_num` - The member num of the workspace.

* `is_default` - Indicates the workspace is default sapce or not.
  + **0** means private space.
  + **1** means default space.
  + **2** means public space.

* `created_at` - The create time of the workspace.

* `updated_at` - The update time of the workspace.

* `created_by` - The user creating the workspace.

* `updated_by` - The user updating the workspace.
