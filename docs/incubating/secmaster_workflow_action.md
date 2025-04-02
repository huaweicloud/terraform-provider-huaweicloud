---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_workflow_action"
description: |-
  Manages a SecMaster workflow action resource within HuaweiCloud.
---

# huaweicloud_secmaster_workflow_action

Manages a SecMaster workflow action resource within HuaweiCloud.

-> Destroying this resource will not change the status of the workflow action resource.

## Example Usage

### Basic Example

```hcl
variable "workspace_id" {}
variable "workflow_id" {}

resource "huaweicloud_secmaster_workflow_action" "test" {
  workspace_id = var.workspace_id
  workflow_id  = var.workflow_id
  command_type = "ActionInstanceRunCommand"
  action_type  = "workflow"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `workspace_id` - (Required, String, NonUpdatable) Sepcifies the workspace ID.

* `workflow_id` - (Required, String, NonUpdatable) Sepcifies the workflow ID.

* `command_type` - (Required, String, NonUpdatable) Sepcifies the command type.
  The value can be: **ActionInstanceRunCommand**, **ActionInstanceDebugCommand**, **ActionInstanceTerminateCommand**,
  **ActionInstanceRetryCommand**.

* `action_type` - (Required, String, NonUpdatable) Sepcifies the action type, e.g. **workflow**.

* `action_instance_id` - (Optional, String, NonUpdatable) Sepcifies the action instance ID.

* `playbook_context` - (Optional, String, NonUpdatable) Sepcifies the playbook context.

* `simulation_context` - (Optional, String, NonUpdatable) Sepcifies the simulation context.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
