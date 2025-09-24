---
subcategory: "Cloud Operations Center (COC)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_coc_alarm_clear"
description: |-
  Manages a COC alarm clear resource within HuaweiCloud.
---

# huaweicloud_coc_alarm_clear

Manages a COC alarm clear resource within HuaweiCloud.

~> Deleting alarm clear resource is not supported, it will only be removed from the state.

## Example Usage

```hcl
variable "alarm_id" {}

resource "huaweicloud_coc_alarm_clear" "test" {
  alarm_ids            = var.alarm_id
  remarks              = "this is remark"
  is_service_interrupt = false
}
```

## Argument Reference

The following arguments are supported:

* `alarm_ids` - (Required, String, NonUpdatable) Specifies the list of alarm IDs, separated by commas.

* `remarks` - (Optional, String, NonUpdatable) Specifies the alarm remarks.

* `is_service_interrupt` - (Optional, Bool, NonUpdatable) Specifies whether to interrupt.

* `start_time` - (Optional, Int, NonUpdatable) Specifies the time when the fault occurred.

* `fault_recovery_time` - (Optional, Int, NonUpdatable) Specifies the fault recovery time.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.
