---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_workflow_version"
description: |-
  Manages a workflow version resource within HuaweiCloud.
---

# huaweicloud_secmaster_workflow_version

Manages a workflow version resource within HuaweiCloud.

## Example Usage

```hcl
variable "workspace_id" {}
variable "workflow_id" {}
variable "name" {}
variable "taskflow" {}
variable "taskconfig" {}
variable "taskflow_type" {}
variable "aop_type" {}

resource "huaweicloud_secmaster_workflow_version" "test" {
  workspace_id  = var.workspace_id
  workflow_id   = var.workflow_id
  name          = var.name
  taskflow      = var.taskflow
  taskconfig    = var.taskconfig
  taskflow_type = var.taskflow_type
  aop_type      = var.aop_type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `workspace_id` - (Required, String, NonUpdatable) Specifies the ID of the workspace to which the workflow version
  belongs.

* `workflow_id` - (Required, String, NonUpdatable) Specifies the ID of the workflow to which the workflow version
  belongs.

* `name` - (Required, String) Specifies the name of the workflow.

* `taskflow` - (Required, String) Specifies the Base64 encoding of the workflow topology diagram.

* `taskconfig` - (Required, String) Specifies the parameters configuration of the workflow topology diagram.
  In JSON format. e.g. **{\"node_info\":{},\"usertask_info\":{}}**.

* `taskflow_type` - (Required, String) Specifies the taskflow type.
  The value only can be **JSON**.

* `aop_type` - (Required, String) Specifies the type of the workflow.
  The value can be **NORMAL**, **SURVEY**, **HEMOSTASIS** or **EASE**.

* `description` - (Optional, String) Specifies the description of the workflow version.

* `status` - (Optional, String) Specifies the status of the workflow version.
  The valid values are as follows:
  + **pending_submit** (Defaults)
  + **pending_approval**
  + **not_activated**
  + **activated**

  -> 1. This paramater is not supported when creating the resource, only can used during update operation.
    <br/>2. If `status` value is **activated**, the workflow version resource does not support deletion.

-> 1. The parameter `status` and other updatable parameters cannot be update simultaneously.
  <br/>2. If you want to update parameters other than `status`, you need to ensure that `status` is **pending_submit**
  or **pending_approval** status for it to take effect. Once the `status` is other value, it will no longer support
  updating parameters other than `status`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `owner_id` - The owner of the workflow version.

* `creator_id` - The creator of the workflow version.

* `version` - The version number of the workflow version.

* `enabled` - Whether the workflow version is enabled.

## Import

The workflow version can be imported using the `workspace_id` and their `id`, separated by a slash (/), e.g.

```bash
$ terraform import huaweicloud_secmaster_workflow_version.test <workspace_id>/<id>
```
