---
subcategory: "ServiceStage"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_servicestagev3_runtime_stack_batch_release"
description: |-
  Use this resource to batch release the runtime stacks within HuaweiCloud.
---

# huaweicloud_servicestagev3_runtime_stack_batch_release

Use this resource to batch release the runtime stacks within HuaweiCloud.

-> When deleting resources, all runtime stacks in the list (stored in .tfstate) will be canceled release.

## Example Usage

```hcl
variable "runtime_stack_ids" {
  type = list(string)
}

resource "huaweicloud_servicestagev3_runtime_stack_batch_release" "test" {
  nruntime_stack_idsame = var.runtime_stack_ids
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the runtime stacks are located.  
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `runtime_stack_ids` - (Required, List) Specifies the runtime stack IDs to be released.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, in UUID format.
