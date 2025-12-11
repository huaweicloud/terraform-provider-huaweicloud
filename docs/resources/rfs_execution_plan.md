---
subcategory: "Resource Formation (RFS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rfs_execution_plan"
description: |-
  Use this resource to manages a execution plan under specified resource stack within HuaweiCloud.
---

# huaweicloud_rfs_execution_plan

Use this resource to manages a execution plan under specified resource stack within HuaweiCloud.

## Example Usage

```hcl
variable "stack_name" {}
variable "execution_plan_name" {}
variable "template_uri" {}
variable "vars_uri" {}

resource "huaweicloud_rfs_execution_plan" "test" {
  stack_name   = var.stack_name
  name         = var.execution_plan_name
  template_uri = var.template_uri
  vars_uri     = var.vars_uri
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the related stack to which the execution plan belongs
  is located. If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `stack_name` - (Required, String, ForceNew) Specifies the name of the resource stack to which the execution plan
  belongs.
  Changing this creates a new resource.
  
* `name` - (Required, String, ForceNew) Specifies the name of the execution plan.
  Changing this creates a new resource.  
  The valid length is limited from `1` to `128`, only Chinese characters, English letters, digits, underscores (_) and
  hyphens (-) are allowed.  
  The name must start with a Chinese character or an English letter, and must be unique.

* `stack_id` - (Optional, String, ForceNew) Specifies the ID of the resource stack to which the execution plan belongs.
  Changing this creates a new resource.

  -> To ensure the correctness of the stack resource being operated (there may be stacks with the same name), it is
     recommended to use `stack_id` for strong matching.

* `description` - (Optional, String, ForceNew) Specifies the description of the execution plan.
  Changing this creates a new resource.

* `template_body` - (Optional, String, ForceNew) Specifies the HCL/JSON template content for deployment resources.
  Changing this creates a new resource.  
  This parameter and `template_uri` are alternative and exactly one of them must be provided.

* `vars_body` - (Optional, String, ForceNew) Specifies the variable content for deployment resources.
  Changing this creates a new resource.  
  This parameter and `vars_uri` parameters cannot be set at the same time.

* `template_uri` - (Optional, String, ForceNew) Specifies the OBS address where the HCL/JSON template archive
  (**.zip** file, which contains all resource **.tf.json** script files to be deployed) or **.tf.json** file is located,
  which describes the target status of the deployment resources.  
  Changing this creates a new resource.

* `vars_uri` - (Optional, String, ForceNew) Specifies the OBS address where the variable (**.tfvars**) file
  corresponding to the HCL/JSON template located, which describes the target status of the deployment resources.  
  Changing this creates a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
