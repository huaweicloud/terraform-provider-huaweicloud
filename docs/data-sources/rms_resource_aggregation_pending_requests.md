---
subcategory: "Config"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_rms_resource_aggregation_pending_requests"
description: |-
  Use this data source to get the list of RMS pending resource aggregation requests.
---

# huaweicloud_rms_resource_aggregation_pending_requests

Use this data source to get the list of RMS pending resource aggregation requests.

## Example Usage

```hcl
variable "account_id" {}

"huaweicloud_rms_resource_aggregation_pending_requests" "test" {
  account_id = var.account_id
}
```

## Argument Reference

The following arguments are supported:

* `account_id` - (Optional, String) Specifies the ID of the authorized resource aggregator account.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `pending_aggregation_requests` - The list of pending aggregation requests.

  The [pending_aggregation_requests](#pending_aggregation_requests_struct) structure is documented below.

<a name="pending_aggregation_requests_struct"></a>
The `pending_aggregation_requests` block supports:

* `requester_account_id` - The ID of the account that requests aggregated data.
