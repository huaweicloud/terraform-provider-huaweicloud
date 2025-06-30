---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_playbook_approvals"
description: |-
  Use this data source to get the list of SecMaster playbook review results.
---

# huaweicloud_secmaster_playbook_approvals

Use this data source to get the list of SecMaster playbook review results.

## Example Usage

```hcl
variable "workspace_id" {}
variable "resource_id" {}
variable "approve_type" {}

data "huaweicloud_secmaster_playbook_approvals" "test" {
  workspace_id = var.workspace_id
  resource_id  = var.resource_id
  approve_type = var.approve_type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `workspace_id` - (Required, String) Specifies the workspace ID.

* `resource_id` - (Required, String) Specifies the resource ID.

* `approve_type` - (Required, String) Specifies the review type.
  The valid values are as follows:
  + **PLAYBOOK**: Indicates playbook.
  + **AOP_WORKFLOW**: Indicates workflow.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `data` - The list of playbook review result.

  The [data](#data_struct) structure is documented below.

<a name="data_struct"></a>
The `data` block supports:

* `id` - The approval ID.

* `result` - The review result.
  The valid values are as follows:
  + **PASS**: Indicates review pass.
  + **UN_PASS**: Indicates review not pass.

* `content` - The review content.

* `type` - The resource type.

* `resource_id` - The resource ID.

* `user_id` - The reviewer ID.

* `create_time` - The creation time of the playbook review.

* `update_time` - The update time of the playbook review.
