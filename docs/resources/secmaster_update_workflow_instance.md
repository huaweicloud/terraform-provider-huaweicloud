---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_update_workflow_instance"
description: |-
  Manages a resource to update workflow instance within HuaweiCloud.
---

# huaweicloud_secmaster_update_workflow_instance

Manages a resource to update workflow instance within HuaweiCloud.

-> This resource is a one-time action resource. Deleting this resource will not clear the corresponding request record,
  but will only remove the resource information from the tf state file.

## Example Usage

```hcl
variable "workspace_id" {}
variable "instance_id" {}

resource "huaweicloud_secmaster_update_workflow_instance" "test" {
  workspace_id = var.workspace_id
  instance_id  = var.instance_id
  command_type = "ActionInstanceTerminateCommand"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this creates a new resource.

* `workspace_id` - (Required, String, NonUpdatable) Specifies the workspace ID to which the workflow belongs.

* `instance_id` - (Required, String, NonUpdatable) Specifies the workflow instance ID.

* `command_type` - (Required, String, NonUpdatable) Specifies the command type. Valid values are:
  + **ActionInstanceTerminateCommand**: Terminate process instance.
  + **ActionInstanceRetryCommand**: Retry Process Example.
  + **ActionInstanceDebugCommand**: Update the debugging results of the process instance.

  -> Fields `task_id` and `inputdataobject` are required parameters when `command_type` is set to
     **ActionInstanceDebugCommand**.

* `task_id` - (Optional, String, NonUpdatable) Specifies the update process debugging instance node ID.

* `input_dataobject` - (Optional, String, NonUpdatable) Specifies the update process debug instance node parameters.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID (also `instance_id`).
