---
subcategory: "Cloud Eye (CES)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ces_alarm_template"
description: ""
---

# huaweicloud_ces_alarm_template

Manages a CES alarm template resource within HuaweiCloud.

## Example Usage

### Create a metric alarm template

```hcl
variable "name" {}

resource "huaweicloud_ces_alarm_template" "test"{
  name = var.name

  policies {
    namespace           = "SYS.APIG"
    dimension_name      = "api_id"
    metric_name         = "req_count_2xx"
    period              = 1
    filter              = "average"
    comparison_operator = ">="
    value               = "10"
    unit                = "times/minute"
    count               = 3
    alarm_level         = 2
    suppress_duration   = 43200
  }
}
```

### Create an event alarm template

```hcl
variable "name" {}

resource "huaweicloud_ces_alarm_template" "test"{
  name = var.name
  type = 2

  policies {
    namespace           = "SYS.VPC"
    metric_name         = "modifyVpc"
    period              = 0
    filter              = "average"
    comparison_operator = ">="
    value               = "1"
    unit                = "count"
    count               = 1
    alarm_level         = 2
    suppress_duration   = 0
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the name of the CES alarm template.
  An alarm template name starts with a letter or Chinese, consists of `1` to `128` characters,
  and can contain only letters, Chinese characters, digits, hyphens (-) and hyphens (-).

* `policies` - (Required, List) Specifies the policy list of the CES alarm template.
The [Policy](#CesAlarmTemplate_Policy) structure is documented below.

* `type` - (Optional, Int, NonUpdatable) Specifies the type of the CES alarm template.
  Default to `0`. The valid values are as follows:
  + **0**: metric alarm template.
  + **2**: event alarm template.

* `delete_associate_alarm` - (Optional, Bool) Specifies whether delete the alarm rule which the alarm
  template associated with. Default to **false**.

* `description` - (Optional, String) Specifies the description of the CES alarm template.
  The description can contain a maximum of `256` characters.

<a name="CesAlarmTemplate_Policy"></a>
The `Policy` block supports:

* `namespace` - (Required, String) Specifies the namespace of the service.

* `metric_name` - (Required, String) Specifies the alarm metric name.

* `period` - (Required, Int) Specifies the judgment period of alarm condition.
  Value options: **0**, **1**, **300**, **1200**, **3600**, **14400**, **86400**.

* `filter` - (Required, String) Specifies the data rollup methods.
  Value options: **average**, **variance**, **min**, **max**, **sum**.

* `comparison_operator` - (Required, String) Specifies the comparison conditions for alarm threshold.
  Value options: **>**, **<**, **=**, **>=**, **<=**.

* `value` - (Required, Int) Specifies the alarm threshold.

* `count` - (Required, Int) Specifies the number of consecutive triggering of alarms. The value ranges from `1` to `5`.

* `suppress_duration` - (Required, Int) Specifies the alarm suppression cycle. Unit: second.
  Only one alarm is sent when the alarm suppression period is **0**.
  Value options: **0**, **300**, **600**, **900**, **1800**, **3600**, **10800**, **21600**,
  **43200**, **86400**.

* `alarm_level` - (Optional, Int) Specifies the alarm level. It means no level if not set.
  The valid values are as follows:
  + **1**: critical.
  + **2**: major.
  + **3**: minor.
  + **4**: informational.

* `unit` - (Optional, String) Specifies the unit string of the alarm threshold.
  The unit can contain a maximum of `32` characters.

* `dimension_name` - (Optional, String) Specifies the resource dimension.
  The name starts with a letter and separated by commas(,) for multiple dimensions,
  can contain only letters, digits, hyphens (-) and hyphens (-),
  and contain a maximum of `32` characters for each dimension.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `association_alarm_total` - Indicates the total num of the alarm that associated with the alarm template.

## Import

The ces alarm template can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_ces_alarm_template.test <template_id>
```
