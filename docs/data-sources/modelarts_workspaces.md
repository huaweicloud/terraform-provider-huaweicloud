---
subcategory: "AI Development Platform (ModelArts)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_modelarts_workspaces"
description: ""
---

# huaweicloud_modelarts_workspaces

Use this data source to get workspaces of ModelArts.

## Example Usage

```hcl
variable "workspace_name" {}

data "huaweicloud_modelarts_workspaces" "test" {
  name = var.workspace_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `name` - (Optional, String) Workspace name. Fuzzy match is supported.  

* `enterprise_project_id` - (Optional, String) The enterprise project ID to which the workspace belongs.  

* `filter_accessible` - (Optional, Bool) Whether to filter that the current user does not have permission to access.  
  Defaults to **false**, query all workspaces.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `workspaces` - The list of workspaces.
  The [workspaces](#Workspaces_Workspaces) structure is documented below.

<a name="Workspaces_Workspaces"></a>
The `workspaces` block supports:

* `id` - Workspace ID.

* `name` - Workspace name.

* `auth_type` - Authorization type.  
  Value options are as follows:
    + **public**: public access within the tenant.
    + **private**: Only the creator and master account can access.
    + **internal**: Accessible to the creator, main account, and specified IAM sub-accounts.

* `description` - The description of the workspace.  

* `owner` - Account name of the owner.

* `enterprise_project_id` - The enterprise project ID to which the workspace belongs.  

* `status` - Workspace status.  
  Valid values are **CREATE_FAILED**, **NORMALL**, **DELETING** and **DELETE_FAILED**.

* `status_info` - Status details.  
  If the deletion fails, you can check the reason through this field.
