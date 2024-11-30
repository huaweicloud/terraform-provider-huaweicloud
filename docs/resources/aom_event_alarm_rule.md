---
subcategory: "Application Operations Management (AOM)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aom_event_alarm_rule"
description: ""
---

# huaweicloud_aom_event_alarm_rule

Manages an AOM event alarm rule resource within HuaweiCloud.

## Example Usage

variable "action_rule_name" {}

```hcl
resource "huaweicloud_aom_event_alarm_rule" "test" {
  name                = "test_rule"
  description         = "terraform test"
  alarm_type          = "notification"
  action_rule         = var.action_rule_name
  enabled             = true
  trigger_type        = "accumulative"
  period              = "300"
  comparison_operator = ">="
  trigger_count       = 2
  alarm_source        = "AOM"

  select_object = {
    "event_type"     ="alarm",
    "event_severity" = "Critical"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of the rule.

  Changing this parameter will create a new resource.

* `alarm_type` - (Required, String) Specifies the alarm type of the rule.
  The value can be **notification** and **denoising**.

* `alarm_source` - (Required, String) Specifies the alarm source of the rule.

* `select_object` - (Required, Map) Specifies the select object of the rule.

* `trigger_type` - (Required, String) Specifies the trigger type.
  The value can be **accumulative** and **immediately**.

* `description` - (Optional, String) Specifies the description of the rule.

* `action_rule` - (Optional, String) Specifies the action rule name.

* `grouping_rule` - (Optional, String) Specifies the route grouping rule name.

* `enabled` - (Optional, Bool) Specifies whether the rule is enabled. Defaults to **true**.

* `trigger_count` - (Optional, Int) Specifies the accumulated times to trigger the alarm.

* `comparison_operator` - (Optional, String) Specifies the comparison condition of alarm.
  The value can be **>** and **>=**.

* `period` - (Optional, Int) Specifies the monitoring period in seconds.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID. The value is the rule name.

* `created_at` - The creation time.

* `updated_at` - The last updated time.

## Import

The application operations management can be imported using the `id` (name), e.g.

```bash
$ terraform import huaweicloud_aom_event_alarm_rule.test test_rule
```
