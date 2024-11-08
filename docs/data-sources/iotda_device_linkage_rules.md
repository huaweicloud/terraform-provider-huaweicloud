---
subcategory: "IoT Device Access (IoTDA)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_iotda_device_linkage_rules"
description: ""
---

# huaweicloud_iotda_device_linkage_rules

Use this data source to get the list of IoTDA device linkage rules.

-> When accessing an IoTDA **standard** or **enterprise** edition instance, you need to specify the IoTDA service
  endpoint in `provider` block.
  You can login to the IoTDA console, choose the instance **Overview** and click **Access Details**
  to view the HTTPS application access address. An example of the access address might be
  **9bc34xxxxx.st1.iotda-app.ap-southeast-1.myhuaweicloud.com**, then you need to configure the
  `provider` block as follows:

  ```hcl
  provider "huaweicloud" {
    endpoints = {
      iotda = "https://9bc34xxxxx.st1.iotda-app.ap-southeast-1.myhuaweicloud.com"
    }
  }
  ```

## Example Usage

```hcl
variable "rule_id" {}

data "huaweicloud_iotda_device_linkage_rules" "test" {
  rule_id = var.rule_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the device linkage rules.
  If omitted, the provider-level region will be used.

* `rule_id` - (Optional, String) Specifies the ID of the device linkage rule.

* `name` - (Optional, String) Specifies the name of the device linkage rule.

* `type` - (Optional, String) Specifies the type of the device linkage rules.
  The valid values are as follows:
  + **DEVICE_LINKAGE**: Cloud based linkage rule.
  + **DEVICE_SIDE**: Device side rule.

* `status` - (Optional, String) Specifies the current status of the device linkage rule.
  The valid values are as follows:
  + **active**: The device linkage rule is active.
  + **inactive**: The device linkage rule is not enabled.

* `space_id` - (Optional, String) Specifies the ID of the resource space to which the device linkage rules belong.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `rules` - All rules that match the filter parameters.
  The [rules](#iotda_linkage_rules) structure is documented below.

<a name="iotda_linkage_rules"></a>
The `rules` block supports:

* `id` - The ID of the device linkage rule.

* `name` - The name of the device linkage rule.

* `description` - The description of the device linkage rule.

* `type` - The type of the device linkage rule.

* `status` - The current status of the device linkage rule.

* `space_id` - The ID of the resource space to which the device linkage rule belongs.

* `triggers` - The condition list of the device linkage rule.
  The [triggers](#rule_conditions) structure is documented below.

* `actions` - The action list of the device linkage rule.
  The [actions](#rule_actions) structure is documented below.

* `effective_period` - The rule condition triggered validity period of the device linkage rule.
  The [effective_period](#rule_effective_period) structure is documented below.

* `trigger_logic` - The logical relationship between multiple rule conditions of the device linkage rule.

* `updated_at` - The latest update time of the device linkage rule.
  The format is **yyyyMMdd'T'HHmmss'Z'**. e.g. **20151212T121212Z**.

<a name="rule_conditions"></a>
The `triggers` block supports:

* `type` - The rule condition type of the device linkage rule. The valid values are as follows:
  + **DEVICE_DATA**: Triggered upon the property of device.
  + **SIMPLE_TIMER**: Triggered by schedule upon policy.
  + **DAILY_TIMER**: Triggered by schedule upon period.
  + **DEVICE_LINKAGE_STATUS**: Triggered upon the status of device.

* `device_data_condition` - The rule condition triggered upon the property of device.
  The [device_data_condition](#trigger_device_data_condition) structure is documented below.

* `simple_timer_condition` - The rule condition triggered by schedule upon policy.
  The [simple_timer_condition](#trigger_simple_timer_condition) structure is documented below.

* `daily_timer_condition` - The rule condition triggered by schedule upon period.
  The [daily_timer_condition](#trigger_daily_timer_condition) structure is documented below.

* `device_linkage_status_condition` - The rule condition triggered upon the status of device.
  The [device_linkage_status_condition](#trigger_device_linkage_status_condition) structure is documented below.

<a name="trigger_device_data_condition"></a>
The `device_data_condition` block supports:

* `device_id` - The ID of the device which trigger the rule.

* `product_id` - The ID of the product associated with the device.

* `path` - The path of the device property.
  The format is **service_id/DataProperty**. e.g. **DoorWindow/status**

* `operator` - The data comparison operator. The valid values are: **>**, **<**,
  **>=**, **<=**, **=**, **in**, and **between**.

* `value` - The Rvalue of a data comparison expression. When the `operator` is **between**, the Rvalue represents the
  minimum and maximum values, separated by commas, such as **20,30**, which means greater than or equal to `20` and less
  than `30`.

* `in_values` - The Rvalue of a data comparison expression. Only when the `operator` is **in**, this field is valid,
  with a maximum of `20` characters, represents matching within the specified values, e.g. **20,30,40**,

* `trigger_strategy` - The judgment strategy triggered by rule conditions. The valid values are:
  + **pulse**: When the data reported by the device meets the conditions, the rule can be triggered.
  + **reverse**: When the data reported by the device last time does not meet the conditions and the
    data reported this time meets the conditions, the rule can be triggered.

* `data_validatiy_period` - The data validity period, in seconds.

<a name="trigger_simple_timer_condition"></a>
The `simple_timer_condition` block supports:

* `start_time` - The start time for triggering the rule.
  The format is **yyyyMMdd'T'HHmmss'Z'**. e.g. **20151212T121212Z**.

* `repeat_interval` - The repeat time interval for triggering the rule, in minutes.

* `repeat_count` - The repeat times for triggering the rule.
  The valid value is range from `1` to `9,999`.

<a name="trigger_daily_timer_condition"></a>
The `daily_timer_condition` block supports:

* `start_time` - The start time for triggering the rule.
  The format is **HH:MM**. e.g. **10:00**.

* `days_of_week` - The week list of the rule validity period, separated by commas. **1** represents Sunday,
  **2** represents Monday, and so on.

<a name="trigger_device_linkage_status_condition"></a>
The `device_linkage_status_condition` block supports:

* `device_id` - The ID of the device which trigger the rule.

* `product_id` - The ID of the product associated with the device.

* `status_list` - All devices status which trigger the rule. The valid values can be **ONLINE** or **OFFLINE**.

* `duration` - The duration of the device status, in minutes.
  The valid value is range from `0` to `60`.

<a name="rule_actions"></a>
The `actions` block supports:

* `type` - The rule action type of the device linkage rule. The valid values are as follows:
  + **DEVICE_CMD**: Device command deliver messages.
  + **SMN_FORWARDING**: Send SMN messages.
  + **DEVICE_ALARM**: Report device alarm messages.

* `device_command` - The detail of device command.
  The [device_command](#action_device_command) structure is documented below.

* `smn_forwarding` - The detail of SMN notifications.
  The [smn_forwarding](#action_smn_forwarding) structure is documented below.

* `device_alarm` - The detail of device alarm.
  The [device_alarm](#action_device_alarm) structure is documented below.

<a name="action_device_command"></a>
The `device_command` block supports:

* `device_id` - The ID of the device to which the command is delivered.

* `command_name` - The name of the command.

* `command_body` - The command parameters.

* `buffer_timeout` - The cache time of device commands, in seconds.

* `response_timeout` - The effective time of the command response, in seconds.

* `mode` - The issuance mode of device commands, which is only valid when the value of `buffer_timeout` is greater than
  `0`. The valid values are **ACTIVE** and **PASSIVE**.

* `service_id` - The ID of the service to which the command belongs.

<a name="action_smn_forwarding"></a>
The `smn_forwarding` block supports:

* `region` - The region to which the SMN service belongs.

* `topic_name` - The topic name of the SMN.

* `topic_urn` - The topic URN of the SMN.

* `message_title` - The message title.

* `message_content` - The message content.

* `message_template_name` - The template name corresponding to the SMN service.

* `project_id` - The project ID to which the SMN belongs.

<a name="action_device_alarm"></a>
The `device_alarm` block supports:

* `name` - The name of the alarm.

* `type` - The type of the alarm. The valid values are as follows:
  + **fault**: Report alarms.
  + **recovery**: Restore alarms.

* `severity` - The severity level of the alarm.
  The valid values can be **warning**, **minor**, **major**, or **critical**.

* `dimension` - The dimension of the alarm. Combine the alarm name and alarm level to jointly identify an alarm.
  The valid values are **device** and **app**. Defaults to user dimension.

* `description` - The description of the alarm.

<a name="rule_effective_period"></a>
The `effective_period` block supports:

* `start_time` - The start time for triggering the rule.
  The format is **HH:MM**. e.g. **10:00**.

* `end_time` - The end time for triggering the rule.
  The format is **HH:MM**. e.g. **10:00**.

* `days_of_week` - The week list of the rule validity period, separated by commas. **1** represents Sunday,
  **2** represents Monday, and so on.
