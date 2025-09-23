---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_change_update"
description: |-
  Manages a COC change update resource within HuaweiCloud.
---

# huaweicloud_coc_change_update

Manages a COC change update resource within HuaweiCloud.

~> Deleting change update resource is not supported, it will only be removed from the state.

## Example Usage

```hcl
variable "ticket_id" {}
variable "sub_ticket_id" {}

resource "huaweicloud_coc_change_update" "test" {
  ticket_id = var.ticket_id
  action    = "change_start_change_success"
  sub_tickets {
    ticket_id = var.sub_ticket_id
  }
}
```

## Argument Reference

The following arguments are supported:

* `ticket_id` - (Required, String, NonUpdatable) Specifies the change order number.

* `phase` - (Optional, String, NonUpdatable) Specifies the type of work order operation.
  Values can be as follows:
  + **phase_change_end**: Completed.
  + **phase_change_cancel**: Canceled.
  + **phase_change_draft**: Draft.
  + **phase_change_implement**: Change implementation and verification.
  + **phase_change_apply**: Applicant confirmation.
  + **phase_change_approve**: Approved.
  + **phase_change_close**: Closed.

* `work_flow_status` - (Optional, String, NonUpdatable) Specifies the work order status.

* `action` - (Optional, String, NonUpdatable) Specifies the operation type.
  Values can be as follows:
  + **change_start_change_success**: Change starts.
  + **change_end_change_success**: Change ends.
  + **change_set_change_result_success**: Adds a change result.
  + **change_complete_success**: Closes the order.

* `sub_tickets` - (Optional, List, NonUpdatable) Specifies the change sub-order information.

  The [sub_tickets](#sub_tickets_struct) structure is documented below.

<a name="sub_tickets_struct"></a>
The `sub_tickets` block supports:

* `ticket_id` - (Optional, String, NonUpdatable) Specifies the sub-order ID.

* `change_result` - (Optional, String, NonUpdatable) Specifies the result of the change.

* `is_verified_in_change_time` - (Optional, Bool, NonUpdatable) Specifies whether verification is possible within the
  time window.

* `verified_docs` - (Optional, String, NonUpdatable) Specifies the verification document ID.

* `comment` - (Optional, String, NonUpdatable) Specifies the reason why the change failed.

* `change_fail_type` - (Optional, String, NonUpdatable) Specifies the change failure type.

* `rollback_start_time` - (Optional, Int, NonUpdatable) Specifies the rollback start time.

* `rollback_end_time` - (Optional, Int, NonUpdatable) Specifies the rollback end time.

* `is_rollback_success` - (Optional, Bool, NonUpdatable) Specifies whether the rollback is successful.

* `is_monitor_found` - (Optional, Bool, NonUpdatable) Specifies whether the device is detected by monitoring.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, which equals to `ticket_id`.
