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

* `is_overwrite` - (Optional, Bool, NonUpdatable) Specifies whether to overwrite an existing alarm template with the same
  template name when creating a template. Default to **false**.

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
  + When `type` is **0**, metric alarm template value options: **>**, **<**, **=**, **>=**, **<=**, **!=**,
    **cycle_decrease**, **cycle_increase**, **cycle_wave**.
  + When `type` is **2**, event alarm template value options: **>**, **<**, **=**, **>=**, **<=**, **!=**.

* `count` - (Required, Int) Specifies the number of consecutive alarm triggering times.
  + For event alarms, the value ranges from **1** to **180**.
  + For metric and website alarms, the value can be **1**, **2**, **3**, **4**, **5**, **10**, **15**, **30**, **60**,
    **90**, **120**, **180**.

* `suppress_duration` - (Required, Int) Specifies the alarm suppression cycle. Unit: second.
  Only one alarm is sent when the alarm suppression period is **0**.
  Value options: **0**, **300**, **600**, **900**, **1800**, **3600**, **10800**, **21600**,
  **43200**, **86400**.

* `value` - (Optional, Int) Specifies the alarm threshold.

* `hierarchical_value` - (Optional, List) Specifies the multiple levels of alarm thresholds.
  The [hierarchical_value](#CesAlarmTemplate_Policy_hierarchical_value) structure is documented below.

-> When `hierarchical_value` and `value` are used at the same time, `hierarchical_value` takes precedence.

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

<a name="CesAlarmTemplate_Policy_hierarchical_value"></a>
The `hierarchical_value` block supports:

* `critical` - (Optional, Float) Specifies the threshold for the critical level.

* `major` - (Optional, Float) Specifies the threshold for the major level.

* `minor` - (Optional, Float) Specifies the threshold for the minor level.

* `info` - (Optional, Float) Specifies the threshold for the info level.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `association_alarm_total` - Indicates the total num of the alarm that associated with the alarm template.

## Import

The ces alarm template can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_ces_alarm_template.test <template_id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `is_overwrite`.
It is generally recommended running `terraform plan` after importing a alarm template.
You can then decide if changes should be applied to the alarm template, or the resource definition should be updated to
align with the alarm template. Also you can ignore changes as below.

```hcl
resource "huaweicloud_ces_alarm_template" "test" {
  ...

  lifecycle {
    ignore_changes = [
      is_overwrite,
    ]
  }
}
```
