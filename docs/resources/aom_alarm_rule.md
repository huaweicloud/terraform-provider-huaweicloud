---
subcategory: "Application Operations Management (AOM 1.0)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_aom_alarm_rule"
description: ""
---

# huaweicloud_aom_alarm_rule

Manages an AOM alarm rule resource within HuaweiCloud.

## Example Usage

```hcl
resource "huaweicloud_aom_alarm_rule" "alarm_rule" {  
  name        = "test-rule"
  alarm_level = 3
  description = "test rule"

  namespace   = "PAAS.NODE"
  metric_name = "cupUsage"

  dimensions {
    name  = "hostID"
    value = var.instance_id
  }

  comparison_operator = ">="
  period              = 60000
  statistic           = "average"
  threshold           = 3
  unit                = "Percent"
  evaluation_periods  = 2
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the alarm rule resource. If omitted, the
  provider-level region will be used. Changing this creates a new resource.

* `name` - (Required, String, ForceNew) Specifies the name of an alarm rule. The value can be a string of `1` to `100`
  characters that can consist of letters, digits, underscores (_), hyphens (-) and chinese characters,
  and it must start and end with letters, digits or chinese characters. Changing this creates a new resource.

* `metric_name` - (Required, String, ForceNew) Specifies the alarm metric name. Changing this creates a new resource.

* `namespace` - (Required, String, ForceNew) Specifies the alarm namespace. Changing this creates a new resource.

* `dimensions` - (Required, List, ForceNew) Specifies the list of metric dimensions. The structure is described below.
  Changing this creates a new resource.

* `period` - (Required, Int) Specifies the alarm checking period in milliseconds.
  The value can be **60,000**, **300,000**, **900,000** and **3,600,000**.

* `statistic` - (Required, String, ForceNew) Specifies the data rollup methods. The value can be **maximum**,
  **minimum**, **average**, **sum** and **sampleCount**. Changing this creates a new resource.

* `comparison_operator` - (Required, String) Specifies the comparison condition of alarm thresholds.
  The value can be **>**, **=**, **<**, **>=** or **<=**.

* `threshold` - (Required, String) Specifies the alarm threshold.

* `unit` - (Required, String, ForceNew) Specifies the data unit.  
  The valid value is range from `1` to `32`.  
  Changing this creates a new resource.

* `evaluation_periods` - (Required, Int) Specifies the alarm checking evaluation periods.
  The value can be `1`, `2`, `3`, `4` and `5`.

* `description` - (Optional, String) Specifies the description of the alarm rule.
 The value can be a string of `0` to `1,000` characters.

* `alarm_level` - (Optional, Int) Specifies the alarm severity. The value can be `1`, `2`, `3` or `4`,
  which indicates *critical*, *major*, *minor*, and *informational*, respectively.
  The default value is `2`.

* `alarm_actions` - (Optional, List, ForceNew) Specifies the action triggered by an alarm. This is a list of strings.
  Changing this creates a new resource.

* `alarm_action_enabled` - (Optional, Bool, ForceNew) Specifies whether to enable the action to be triggered by an alarm.
  The default value is true. Changing this creates a new resource.

* `ok_actions` - (Optional, List, ForceNew) Specifies the action triggered by the clearing of an alarm.
  This is a list of strings. Changing this creates a new resource.

* `insufficient_data_actions` - (Optional, List, ForceNew) Specifies the action triggered when the data is not enough.
  This is a list of strings. Changing this creates a new resource.

The `dimensions` block supports:

* `name` - (Required, String, ForceNew) Specifies the dimension name. Changing this creates a new resource.

* `value` - (Required, String, ForceNew) Specifies the dimension value. Changing this creates a new resource.

-> **NOTE:** You can get more information about `metric_name`, `namespace`, `unit` and `dimensions`
  from [Metric Overview](https://support.huaweicloud.com/intl/en-us/productdesc-aom/aom_06_0014.html).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates the alarm rule ID.

* `alarm_enabled` - Indicates whether the alarm rule is enabled.

* `state_value` - Indicates the alarm status.

* `state_reason` - Indicates the reason of alarm status.

## Import

AOM alarm rules can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_aom_alarm_rule.alarm_rule <id>
```
