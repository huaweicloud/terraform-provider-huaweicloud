---
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rts_stack_resource_v1"
sidebar_current: "docs-huaweicloud-datasource-rts-stack-resource-v1"
description: |-
  Provides metadata of an RTS stack resource
---

# Data Source: huaweicloud_rts_stack_resource_v1

The HuaweiCloud RTS Stack Resource data source allows access to stack resource metadata.

## Example Usage

```hcl
variable "stack_name" { }
variable "resource_name" { }

data "huaweicloud_rts_stack_resource_v1" "stackresource" {
  stack_name = "${var.stack_name}"
  resource_name = "${var.resource_name}"
  
}
```

## Argument Reference
The following arguments are supported:

* `stack_name` - (Required) The unique stack name.

* `resource_name` - (Optional) The name of a resource in the stack.

* `physical_resource_id` - (Optional) The physical resource ID.

* `resource_type` - (Optional) The resource type.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `logical_resource_id` - The logical resource ID.

* `resource_status` - The status of the resource.

* `resource_status_reason` - The resource operation reason.
 
* `required_by` - Specifies the resource dependency.



