---
subcategory: "Cloud Eye (CES)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ces_one_click_alarm_rule_action"
description: |-
  Manages an CES batch enable or disable alarm rules for one service in one-click monitoring resource within HuaweiCloud.
---

# huaweicloud_ces_one_click_alarm_rule_action

Manages an CES batch enable or disable alarm rules for one service in one-click monitoring resource within HuaweiCloud.

-> Deleting the batch enable or disable alarm rules resource is not supported. The batch enable or disable alarm rules
  for one service in one-click monitoring resource is only removed from the state.

## Example Usage

```hcl
variable "one_click_alarm_id" {}
variable "alarm_id" {}

resource "huaweicloud_ces_one_click_alarm_rule_action" "test" {
  one_click_alarm_id = var.one_click_alarm_id
  alarm_ids          = [var.alarm_id]
  alarm_enabled      = false
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `one_click_alarm_id` - (Required, String, NonUpdatable) Specifies the one-click monitoring ID for a service.

* `alarm_ids` - (Required, List, NonUpdatable) Specifies IDs of alarm rules to be enabled or disabled in batches.

* `alarm_enabled` - (Required, Bool, NonUpdatable) Specifies whether to generate alarms when the alarm triggering
  conditions are met.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `success_alarm_ids` - Indicates IDs of alarm rules that were enabled or disabled.
