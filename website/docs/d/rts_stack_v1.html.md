---
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rts_stack_v1"
sidebar_current: "docs-huaweicloud-datasource-rts-stack-v1"
description: |-
  Get information on an HuaweiCloud RTS.
---

# huaweicloud_rts_stack_v1

The Huaweicloud `Resource Template Service` Stack data source allows access to stack outputs and other useful data including the template body.

## Example Usage

The following example shows how one might accept a VPC id as a variable and use this data source to obtain the data necessary to create a subnet within it.

```hcl
variable "stack_name" { }

data "huaweicloud_rts_stack_v1" "stacks" {
  name = "${var.stack_name}"
}
```

## Argument Reference
The following arguments are supported:

* `name` - (Required) The name of the stack.

## Attributes Reference

The following attributes are exported:

* `capabilities` - List of stack capabilities for stack.

* `description` - 	Describes the stack.

* `disable_rollback` - Specifies whether to perform a rollback if the update fails.

* `outputs` - A list of stack outputs.

* `parameters` - Specifies the stack parameters.

* `template_body` - Structure containing the template body.

* `timeout_mins` - Specifies the timeout duration.

* `status` - Specifies the stack status.
 
* `name` - Specifies the stack name.
 
* `status_reason` - Specifies the description of the stack operation.

* `notification_topics` - List of notification topics for stack.
