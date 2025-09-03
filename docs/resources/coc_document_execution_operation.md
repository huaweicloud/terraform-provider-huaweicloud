---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_document_execution_operation"
description: |-
  Manages a COC document execution operation resource within HuaweiCloud.
---

# huaweicloud_coc_document_execution_operation

Manages a COC document execution operation resource within HuaweiCloud.

~> Deleting document execution operation resource is not supported, it will only be removed from the state.

## Example Usage

```hcl
variable "execution_id" {}

resource "huaweicloud_coc_document_execution_operation" "test" {
  execution_id = var.execution_id
  operate_type = "TERMINATE"
}
```

## Argument Reference

The following arguments are supported:

* `execution_id` - (Required, String, NonUpdatable) Specifies the work order ID.

* `operate_type` - (Required, String, NonUpdatable) Specifies the operation type.
  Values can be as follows:
  + **RESUME**: Resubmit.
  + **TERMINATE**: Force termination.
  + **RETRY**: Retry.
  + **SKIP_BATCH**: Skip batching.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which equals to `execution_id`.
