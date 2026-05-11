---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_dataservice_approvers"
description: |-
  Use this data source to get the list of approvers for DataArts Studio Data Service within HuaweiCloud.
---

# huaweicloud_dataarts_dataservice_approvers

Use this data source to get the list of approvers for DataArts Studio Data Service within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}

data "huaweicloud_dataarts_dataservice_approvers" "test" {
  workspace_id = var.workspace_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the approvers are located.  
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the ID of the workspace to which the approvers belong.

* `name` - (Optional, String) Specifies the name of the approver to be queried.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `approvers` - The list of approvers that match the filter parameters.  
  The [approvers](#dataarts_dataservice_approvers_attr) structure is documented below.

<a name="dataarts_dataservice_approvers_attr"></a>  
The `approvers` block supports:

* `id` - The ID of the approver.

* `name` - The name of the approver.

* `user_id` - The user ID of the approver.

* `user_name` - The user name of the approver.

* `email` - The email of the approver.

* `phone_number` - The phone number of the approver.

* `department` - The department of the approver.

* `description` - The description of the approver.

* `character` - The character of the approver.

* `create_by` - The creator of the approver.

* `create_time` - The creation time of the approver, in RFC3339 format.

* `app_name` - The application name.

* `topic_urn` - The topic URN.

* `template_id` - The template ID.

* `project_id` - The project ID.

* `approver_type` - The type of the approver.
