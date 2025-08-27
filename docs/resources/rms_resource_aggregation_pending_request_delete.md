---
subcategory: "Config"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rms_resource_aggregation_pending_request_delete"
description: |-
  Manages an RMS pending authorization request delete resource within HuaweiCloud resources.
---

# huaweicloud_rms_resource_aggregation_pending_request_delete

Manages an RMS pending authorization request delete resource within HuaweiCloud resources.

## Example Usage

```hcl
variable "requester_account_id" {}

resource "huaweicloud_rms_resource_aggregation_pending_request_delete" "test" {
  requester_account_id = var.requester_account_id
}
```

## Argument Reference

The following arguments are supported:

* `requester_account_id` - (Required, String, NonUpdatable) Specifies the ID of the account that requests data aggregation.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID which equals to the `requester_account_id`.
