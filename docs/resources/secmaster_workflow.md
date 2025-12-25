---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_workflow"
description: |-
  Manages a workflow resource within HuaweiCloud.
---

# huaweicloud_secmaster_workflow

Manages a workflow resource within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}
variable "name" {}
variable "engine_type" {}
variable "aop_type" {}
variable "dataclass_id" {}

resource "huaweicloud_secmaster_workflow" "test" {
  workspace_id = var.workspace_id
  name         = var.name
  engine_type  = var.engine_type
  aop_type     = var.aop_type
  dataclass_id = var.dataclass_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `workspace_id` - (Required, String, NonUpdatable) Specifies the ID of the workspace to which the workflow belongs.

* `name` - (Required, String) Specifies the name of the workflow.

* `engine_type` - (Required, String, NonUpdatable) Specifies the engine type.
  The value only can be **public_engine**.

* `aop_type` - (Required, String, NonUpdatable) Specifies the workflow type.
  The valid values are as follows:
  + **NORMAL**
  + **SURVEY**
  + **HEMOSTASIS**
  + **EASE**

* `dataclass_id` - (Required, String, NonUpdatable) Specifies the dataclass ID.

* `labels` - (Optional, String) Specifies the workflow entity type.
  The valid values are as follows:
  + **IP**
  + **ACCOUNT**
  + **DOMAIN**

* `description` - (Optional, String) Specifies the description of the workflow.

* `enabled` - (Optional, Bool) Specifies whether to enable workflow version.
  The valid values are as follows:
  + **true**
  + **false**

* `version_id` - (Optional, String) Specifies the ID of the workflow version that has been activated
  under the current workflow.

  -> This parameter is mandatory when `enabled` is set to **true**.

-> 1. The parameters `enabled` and `version_id` are not supported when creating the resource,
  only can used during update operation.
  <br/>2. If you want to edit the parameters `enabled` and `version_id` (enable the workflow), you need to ensure
  the workflow version has been published under the workflow, the workflow version `status` is **not_activated**
  or **activated**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `dataclass_name` - The dataclass name.

* `version` - The current activated workflow version.

* `project_id` - The project ID.

* `owner_id` - The workflow owner ID.

* `creator_id` - The workflow creator ID.

* `creator_name` - The workflow creator name.

* `modifier_id` - The workflow updater ID.

* `modifier_name` - The workflow updater name.

* `create_time` - The workflow creation time.

* `update_time` - The workflow update time.

* `use_role` - The use role.

* `edit_role` - The edit role.

* `approve_role` - The approve role.

* `current_approval_version_id` - The current workflow version ID awaiting approval.

* `current_rejected_version_id` - The current rejected workflow version ID.

## Import

The workflow can be imported using the `workspace_id` and their `id`, separated by a slash (/), e.g.

```bash
$ terraform import huaweicloud_secmaster_workflow.test <workspace_id>/<id>
```
