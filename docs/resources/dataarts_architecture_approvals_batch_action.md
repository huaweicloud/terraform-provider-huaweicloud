---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_architecture_approvals_batch_action"
description: |-
  Use this resource to operate a DataArts Architecture approvals batch action within HuaweiCloud.
---

# huaweicloud_dataarts_architecture_approvals_batch_action

Use this resource to operate a DataArts Architecture approvals batch action within HuaweiCloud.

-> This resource is only a one-time action resource for operating approvals. Deleting this resource will not
   clear the corresponding request record, but will only remove the resource information from the tfstate file.

## Example Usage

```hcl
variable "workspace_id" {}
variable "approval_ids" {
  type = list(string)
}

resource "huaweicloud_dataarts_architecture_approvals_batch_action" "test" {
  workspace_id = var.workspace_id
  approval_ids = join(",", var.approval_ids)
  message      = "approve"
  action       = "resolve"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the approval is located.  
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `workspace_id` - (Required, String, NonUpdatable) Specifies the ID of the workspace to which the approval belongs.

* `approval_ids` - (Required, String, NonUpdatable) Specifies the IDs of the approvals.  
  Multiple approval IDs separated by commas (,).

* `action` - (Required, String, NonUpdatable) Specifies the action of the approval status to be approved.  
  The valid values are as follows:
  + **resolve**
  + **reject**
  + **recall**

* `message` - (Optional, String, NonUpdatable) Specifies the approval message is required for the approval action.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
