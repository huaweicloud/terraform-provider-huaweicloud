---
subcategory: "SecMaster"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_secmaster_workflow_version_validation"
description: |-
  Manages a resource to validate workflow version within HuaweiCloud.
---

# huaweicloud_secmaster_workflow_version_validation

Manages a resource to validate workflow version within HuaweiCloud.

-> This resource is a one-time action resource. Deleting this resource will not clear the corresponding request record,
  but will only remove the resource information from the tf state file.

## Example Usage

```hcl
variable "workspace_id" {}
variable "aopworkflow_id" {}
variable "mode" {}
variable "taskconfig" {}
variable "taskflow" {}

resource "huaweicloud_secmaster_workflow_version_validation" "test" {
  workspace_id   = var.workspace_id
  aopworkflow_id = var.aopworkflow_id
  mode           = var.mode
  taskconfig     = var.taskconfig
  taskflow       = var.taskflow
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `workspace_id` - (Required, String, NonUpdatable) Specifies the workspace ID to which the workflow belongs.

* `aopworkflow_id` - (Required, String, NonUpdatable) Specifies the workflow ID.

* `mode` - (Required, String, NonUpdatable) Specifies the workflow verify type.
  The valid values are as follows:
  + **BASIC**
  + **CIRCLE**
  + **APP_PARAMS**

* `taskconfig` - (Required, String, NonUpdatable) Specifies the parameters configuration of the workflow topology diagram.
  In JSON format. e.g. **{\"node_info\":{},\"usertask_info\":{}}**.

* `taskflow` - (Required, String, NonUpdatable) Specifies the Base64 encoding of the workflow topology diagram.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `result` - The verify result.

* `detail` - The verify result description.
