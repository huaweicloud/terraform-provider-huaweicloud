---
subcategory: "DataArts Studio"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_dataarts_architecture_approvals"
description: |-
  Use this data source to query DataArts Architecture approvals within HuaweiCloud.
---

# huaweicloud_dataarts_architecture_approvals

Use this data source to query DataArts Architecture approvals within HuaweiCloud.

## Example Usage

### Query all approvals under a specified workspace

```hcl
variable "workspace_id" {}

data "huaweicloud_dataarts_architecture_approvals" "test" {
  workspace_id = var.workspace_id
}
```

### Query the approvals under a specified workspace and using approval name to filter

```hcl
variable "workspace_id" {}
variable "approval_name" {}

data "huaweicloud_dataarts_architecture_approvals" "test" {
  workspace_id = var.workspace_id
  name         = var.approval_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region where the approvals are located.  
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the ID of the workspace to which the approvals belong.

* `biz_id` - (Optional, String) Specifies the business ID of the approvals.

* `name` - (Optional, String) Specifies the name of the approvals.  
  Fuzzy search is supported.

* `create_by` - (Optional, String) Specifies the creator of the approvals.

* `approver` - (Optional, String) Specifies the approver of the approvals.

* `approval_status` - (Optional, String) Specifies the approval status of the approvals.  
  The valid values are as follows:
  + **DEVELOPING**
  + **FINISHED**

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `approvals` - The list of approvals that matched filter parameters.  
  The [approvals](#dataarts_architecture_approvals) structure is documented below.

<a name="dataarts_architecture_approvals"></a>
The `approvals` block supports:

* `id` - The ID of the approval, in UUID format.

* `name_ch` - The chinese name of the approval.

* `name_en` - The english name of the approval.

* `biz_id` - The business ID of the approval.

* `biz_type` - The business type of the approval.

* `biz_info` - The business details of the approval.

* `biz_status` - The business status of the approval.

* `approval_status` - The approval status of the approval.

* `approval_type` - The approval type of the approval.

* `submit_time` - The submit time of the approval, in RFC3339 format.

* `create_by` - The creator of the approval.

* `approval_time` - The approval time of the approval, in RFC3339 format.

* `approver` - The approver of the approval.

* `email` - The approver email of the approval.

* `msg` - The approval message of the approval.
