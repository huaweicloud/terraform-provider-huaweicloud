---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_script_order_operation"
description: |-
  Manages a COC script order operation resource within HuaweiCloud.
---

# huaweicloud_coc_script_order_operation

Manages a COC script order operation resource within HuaweiCloud.

~> Deleting script order operation resource is not supported, it will only be removed from the state.

## Example Usage

```hcl
variable "execute_uuid" {}
variable "operation_type" {}

resource "huaweicloud_coc_script_order_operation" "test" {
  execute_uuid   = var.execute_uuid
  operation_type = var.operation_type
}
```

## Argument Reference

The following arguments are supported:

* `execute_uuid` - (Required, String, NonUpdatable) Specifies the script order ID.

* `operation_type` - (Required, String, NonUpdatable) Specifies the operation type.
  Values can be as follows:
  + **CANCEL_INSTANCE**: Cancel an instance.
  + **SKIP_BATCH**: Skip batch.
  + **CANCEL_ORDER**: Cancel the script order.
  + **PAUSE_ORDER**: Pause the script order.
  + **CONTINUE_ORDER**: Continue the script order.

* `batch_index` - (Optional, Int, NonUpdatable) Specifies the batch index.

* `instance_id` - (Optional, Int, NonUpdatable) Specifies the instance ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The script order ID, which equals to `execute_uuid`.
