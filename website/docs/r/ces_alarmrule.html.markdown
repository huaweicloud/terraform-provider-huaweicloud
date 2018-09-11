---
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ces-alarmrule"
sidebar_current: "docs-huaweicloud-resource-ces-alarmrule"
description: |-
  Manages a V2 topic resource within HuaweiCloud.
---

# huaweicloud\_ces\_alarmrule

Manages a V2 topic resource within HuaweiCloud.

## Example Usage

```hcl
resource "huaweicloud_ces_alarmrule" "alarm_rule" {
  "alarm_name" = "alarm_rule"
  "metric" {
    "namespace" = "SYS.ECS"
    "metric_name" = "network_outgoing_bytes_rate_inband"
    "dimensions" {
        "name" = "instance_id"
        "value" = "${huaweicloud_compute_instance_v2.webserver.id}"
    }
  }
  "condition"  {
    "period" = 300
    "filter" = "average"
    "comparison_operator" = ">"
    "value" = 6
    "unit" = "B/s"
    "count" = 1
  }
  "alarm_actions" {
    "type" = "notification"
    "notification_list" = [
      "${huaweicloud_smn_topic_v2.topic.id}"
    ]
  }
}
```

## Argument Reference

The following arguments are supported:

* `alarm_name` - (Required) Specifies the name of an alarm rule. The value can
    be a string of 1 to 128 characters that can consist of numbers, lowercase letters,
    uppercase letters, underscores (_), or hyphens (-).

* `alarm_description` - (Optional) The value can be a string of 0 to 256 characters.

* `metric` - (Required) Specifies the alarm metrics. The structure is described
    below.

* `condition` - (Required) Specifies the alarm triggering condition. The structure
    is described below.

* `alarm_actions` - (Optional) Specifies the action triggered by an alarm. The
    structure is described below.

* `insufficientdata_actions` - (Optional) Specifies the action triggered by
    data insufficiency. The structure is described below.

* `ok_actions` - (Optional) Specifies the action triggered by the clearing of
    an alarm. The structure is described below.

* `alarm_enabled` - (Optional) Specifies whether to enable the alarm. The default
    value is true.

* `alarm_action_enabled` - (Optional) Specifies whether to enable the action
    to be triggered by an alarm. The default value is true.
    Note: If alarm_action_enabled is set to true, at least one of the following
    parameters alarm_actions, insufficientdata_actions, and ok_actions cannot
    be empty. If alarm_actions, insufficientdata_actions, and ok_actions coexist,
    their corresponding notification_list must be of the same value.

The `metric` block supports:

* `namespace` - (Required) Specifies the namespace in service.item format. service.item
    can be a string of 3 to 32 characters that must start with a letter and can
    consists of uppercase letters, lowercase letters, numbers, or underscores (_).

* `metric_name` - (Required) Specifies the metric name. The value can be a string
    of 1 to 64 characters that must start with a letter and can consists of uppercase
    letters, lowercase letters, numbers, or underscores (_).

* `dimensions` - (Required) Specifies the list of metric dimensions. Currently,
    the maximum length of the dimesion list that are supported is 3. The structure
    is described below.

The `dimensions` block supports:

* `name` - (Required) Specifies the dimension name. The value can be a string
    of 1 to 32 characters that must start with a letter and can consists of uppercase
    letters, lowercase letters, numbers, underscores (_), or hyphens (-).

* `value` - (Required) Specifies the dimension value. The value can be a string
    of 1 to 64 characters that must start with a letter or a number and can consists
    of uppercase letters, lowercase letters, numbers, underscores (_), or hyphens
    (-).

The `condition` block supports:

* `period` - (Required) Specifies the alarm checking period in seconds. The
    value can be 1, 300, 1200, 3600, 14400, and 86400.
    Note: If period is set to 1, the raw metric data is used to determine
    whether to generate an alarm.

* `filter` - (Required) Specifies the data rollup methods. The value can be
    max, min, average, sum, and vaiance.

* `comparison_operator` - (Required) Specifies the comparison condition of alarm
    thresholds. The value can be >, =, <, >=, or <=.

* `value` - (Required) Specifies the alarm threshold. The value ranges from
    0 to Number of 1.7976931348623157e+308.

* `unit` - (Optional) Specifies the data unit.

* `count` - (Required) Specifies the number of consecutive occurrence times.
    The value ranges from 1 to 5.

the `alarm_actions` block supports:

* `type` - (Optional) specifies the type of action triggered by an alarm. the
    value can be notification or autoscaling.
    notification: indicates that a notification will be sent to the user.
    autoscaling: indicates that a scaling action will be triggered.

* `notification_list` - (Optional) specifies the topic urn list of the target
    notification objects. the maximum length is 5. the topic urn list can be
    obtained from simple message notification (smn) and in the following format:
    urn: smn:([a-z]|[a-z]|[0-9]|\-){1,32}:([a-z]|[a-z]|[0-9]){32}:([a-z]|[a-z]|[0-9]|\-|\_){1,256}.
    if type is set to notification, the value of notification_list cannot be
    empty. if type is set to autoscaling, the value of notification_list must
    be [] and the value of namespace must be sys.as.
    Note: to enable the as alarm rules take effect, you must bind scaling
    policies. for details, see the auto scaling api reference.

the `insufficientdata_actions` block supports:

* `type` - (Optional) specifies the type of action triggered by an alarm. the
    value is notification.
    notification: indicates that a notification will be sent to the user.

* `notification_list` - (Optional) indicates the list of objects to be notified
    if the alarm status changes. the maximum length is 5.

the `ok_actions` block supports:

* `type` - (Optional) specifies the type of action triggered by an alarm. the
    value is notification.
    notification: indicates that a notification will be sent to the user.

* `notification_list` - (Optional) indicates the list of objects to be notified
    if the alarm status changes. the maximum length is 5.

## Attributes Reference

The following attributes are exported:

* `alarm_name` - See Argument Reference above.
* `alarm_description` - See Argument Reference above.
* `metric` - See Argument Reference above.
* `condition` - See Argument Reference above.
* `alarm_actions` - See Argument Reference above.
* `insufficientdata_actions` - See Argument Reference above.
* `ok_actions` - See Argument Reference above.
* `alarm_enabled` - See Argument Reference above.
* `alarm_action_enabled` - See Argument Reference above.
* `id` - Specifies the alarm rule ID.
* `update_time` - Specifies the time when the alarm status changed. The value
    is a UNIX timestamp and the unit is ms.
* `alarm_state` - Specifies the alarm status. The value can be:
    ok: The alarm status is normal,
    alarm: An alarm is generated,
    insufficient_data: The required data is insufficient.

