---
subcategory: "IoT Device Access (IoTDA)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_iotda_device_linkage_rule"
description: ""
---

# huaweicloud_iotda_device_linkage_rule

Manages an IoTDA device linkage rule within HuaweiCloud.

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
variable "space_id" {}
variable "trigger_device_id" {}
variable "action_device_id" {}


resource "huaweicloud_smn_topic" "topic" {
  name = "iot-demo"
}

resource "huaweicloud_iotda_device_linkage_rule" "test" {
  space_id = var.space_id
  name     = "demoLinkageRule"
  
  triggers {
    type = "SIMPLE_TIMER"
    simple_timer_condition {
      start_time      = "20220622T160000Z"
      repeat_interval = 2
      repeat_count    = 2
    }
  }

  triggers {
    type = "DEVICE_DATA"
    device_data_condition {
      device_id             = var.trigger_device_id
      path                  = "service_id/propertyName_1"
      operator              = "="
      value                 = "5"
      trigger_strategy      = "pulse"
      data_validatiy_period = 300
    }
  }

  triggers {
    type = "DAILY_TIMER"
    daily_timer_condition {
      start_time = "19:02"
    }
  }

  actions {
    type = "SMN_FORWARDING"
    smn_forwarding {
      region          = huaweicloud_smn_topic.topic.region
      topic_name      = huaweicloud_smn_topic.topic.name
      topic_urn       = huaweicloud_smn_topic.topic.topic_urn
      message_title   = "message_title"
      message_content = "message_content"
    }
  }

  actions {
    type = "DEVICE_CMD"
    device_command {
      device_id    = var.action_device_id
      service_id   = "service_id"
      command_name = "cmd_name"
      command_body = "{\"cmd_parameter_1\":\"3\"}"
    }
  }

  effective_period {
    start_time   = "00:00"
    end_time     = "23:59"
    days_of_week = "1,2,3"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the IoTDA device linkage rule
resource. If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the name of the device linkage rule. The name contains a maximum of `128`
characters.

* `space_id` - (Required, String, ForceNew) Specifies the resource space ID to which the device linkage rule belongs.
Changing this parameter will create a new resource.

* `triggers` - (Required, List) Specifies the list of the triggers, at most 10 triggers.
The [triggers](#IoTDA_triggers) structure is documented below.

* `actions` - (Required, List) Specifies the list of the actions, at most 10 actions.
The [actions](#IoTDA_actions) structure is documented below.

* `trigger_logic` - (Optional, String) Specifies the logical relationship between multiple triggers.
The options are as follows:
  + **and**: All of the triggers are met.
  + **or**: Any of the triggers are met.

  Defaults to `and`.

* `description` - (Optional, String) Specifies the description of device linkage rule. The description contains
a maximum of `256` characters.

* `enabled` - (Optional, Bool) Specifies whether to enable the device linkage rule. Defaults to **true**.

* `effective_period` - (Optional, List) Specifies the effective period of the device linkage rule. Always effectives
by default. The [effective_period](#IoTDA_effective_period) structure is documented below.

<a name="IoTDA_triggers"></a>
The `triggers` block supports:

* `type` - (Required, String) Specifies the type of the trigger. The options are as follows:
  + **DEVICE_DATA**: Triggered upon the property of device.
  + **SIMPLE_TIMER**: Triggered by policy.
  + **DAILY_TIMER**: Triggered at specified time every day.

* `device_data_condition` - (Optional, List) Specifies the condition triggered upon the property of device. It is
required when type is `DEVICE_DATA`. The [device_data_condition](#IoTDA_device_data_condition) structure is
documented below.

* `simple_timer_condition` - (Optional, List) Specifies the condition triggered by policy. It is required when type
is `SIMPLE_TIMER`. The [simple_timer_condition](#IoTDA_simple_timer_condition) structure is documented below.

* `daily_timer_condition` - (Optional, List) Specifies the condition triggered at specified time every day. It is
required when type is `DAILY_TIMER`. The [daily_timer_condition](#IoTDA_daily_timer_condition) structure is
documented below.

<a name="IoTDA_device_data_condition"></a>
The `device_data_condition` block supports:

* `device_id` - (Optional, String) Specifies the device id which triggers the rule.
Exactly one of `device_id` or `product_id` must be provided.

* `product_id` - (Optional, String) Specifies the product id, all devices belonging to this product will trigger
the rule. Exactly one of `device_id` or `product_id` must be provided.

* `path` - (Required, String) Specifies the path of the device property, in the format: **service_id/DataProperty**.

* `operator` - (Required, String) Specifies the data comparison operator. The valid values are: **>**, **<**,
**>=**, **<=**, **=** and **between**.

* `value` - (Required, String) Specifies the Rvalue of a data comparison expression. When the `operator` is `between`,
the Rvalue represents the minimum and maximum values, separated by commas, such as "20,30",
which means greater than or equal to 20 and less than 30.

* `trigger_strategy` - (Optional, String) Specifies the trigger strategy. The options are as follows:
  + **pulse**: When the data reported by the device meets the conditions, the rule can be triggered.
  + **reverse**: Repetition suppression. For example, if an alarm is configured to be triggered when the battery level
   is lower than 20%, the alarm will be triggered once the battery initially drops below 20% but will not be triggered
   again each time the battery drops to a lower level.

   Defaults to `pulse`.

* `data_validatiy_period` - (Optional, Int) Specifies data validity period, Unit is `seconds`. Defaults to `300`.
For example, if Data Validity Period is set to 30 minutes, a device generates data at 19:00, and the platform receives
the data at 20:00, the action is not triggered regardless of whether the conditions are met.

<a name="IoTDA_simple_timer_condition"></a>
The `simple_timer_condition` block supports:

* `start_time` - (Required, String) Specifies the start time to trigger the rule, using the UTC time zone,
in the format: yyyyMMdd'T'HHmmss'Z'. For example: `20220622T160000Z`.

* `repeat_interval` - (Required, Int) Specifies the interval of repetition, Unit is `minutes`.

* `repeat_count` - (Required, Int) Specifies total number of repetition.

<a name="IoTDA_daily_timer_condition"></a>
The `daily_timer_condition` block supports:

* `start_time` - (Required, String) Specifies the start time, in the format: `HH:mm`.
For example: `03:00`.

* `days_of_week` - (Optional, String) Specifies a list of days of week, separated by commas. 1 represents Sunday,
2 represents Monday, and so on. Defaults to `1,2,3,4,5,6,7` (every day).

<a name="IoTDA_actions"></a>
The `actions` block supports:

* `type` - (Required, String) Specifies the type of action. The options are as follows:
  + **DEVICE_CMD**: Deliver commands.
  + **SMN_FORWARDING**: Send SMN notifications.
  + **DEVICE_ALARM**: Report alarms or clear alarms.

* `device_command` - (Optional, List) Specifies the detail of device command. It is required when type
is `DEVICE_CMD`. The [device_command](#IoTDA_device_command) structure is documented below.

* `smn_forwarding` - (Optional, List) Specifies the detail of SMN notifications. It is required when type
is `SMN_FORWARDING`. The [smn_forwarding](#IoTDA_smn_forwarding) structure is documented below.

* `device_alarm` - (Optional, List) Specifies the detail of device alarm. It is required when type
is `DEVICE_ALARM`. The [device_alarm](#IoTDA_device_alarm) structure is documented below.

<a name="IoTDA_device_command"></a>
The `device_command` block supports:

* `device_id` - (Required, String) Specifies the device id which executes the command.

* `service_id` - (Required, String) Specifies the service id to which the command belongs.

* `command_name` - (Required, String) Specifies the command name.

* `command_body` - (Required, String) Specifies the command parameters, in json format.
  + Example of device command using LWM2M protocol: `{"value":"1"}`, there are key-value pairs, each key is the
    parameter name of the command in the product model.
  + Example of device command using MQTT protocol: `{"header": {"mode": "ACK","from": "/users/testUser","method":
    "SET_TEMPERATURE_READ_PERIOD","to":"/devices/{device_id }/services/{service_id}"},"body": {"value" : "1"}}`.
      - **mode**: Required, whether the device needs to reply to the confirmation message after receiving the command.
        The default is ACK mode. `ACK` indicates that the confirmation message needs to be replied,
        `NOACK` indicates that the confirmation message does not need to be replied.
      - **from**: Optional, the address of the command sender.
        When the App initiates a request, the format is /users/{userId},
        when the application server initiates a request, the format is /{serviceName},
        and when the IoT platform initiates a request, the format is /cloud/{serviceName}.
      - **to**: optional, the address of the command receiver, the format is /devices/{device_id}/services/{service_id}.
      - **method**: optional, the command name defined in the product model.
      - **body**: optional, the message body of the command, which contains key-value pairs, each key is the parameter
        name of the command in the product model. The specific format requires application and device conventions.

<a name="IoTDA_smn_forwarding"></a>
The `smn_forwarding` block supports:

* `region` - (Required, String) Specifies the region to which the SMN belongs.

* `topic_name` - (Required, String) Specifies the topic name of the SMN.

* `topic_urn` - (Required, String) Specifies the topic URN of the SMN.

* `message_title` - (Required, String) Specifies the message title.

* `message_content` - (Required, String) Specifies the message content.  
  The value can contain a maximum of `256` characters.

* `project_id` - (Optional, String) Specifies the project ID to which the SMN belongs.
If omitted, the default project in the region will be used.

<a name="IoTDA_device_alarm"></a>
The `device_alarm` block supports:

* `name` - (Required, String) Specifies the name of the alarm.

* `type` - (Required, String) Specifies the type of the alarm. The options are as follows:
  + **fault**: Report alarms.
  + **recovery**: Clear alarms.

* `severity` - (Required, String) Specifies the severity level of the alarm.
The valid values are **warning**, **minor**, **major** and **critical**.

* `description` - (Optional, String) Specifies the description of the alarm.  
  The value can contain a maximum of `256` characters.

<a name="IoTDA_effective_period"></a>
The `effective_period` block supports:

* `start_time` - (Required, String) Specifies the start time, in the format: `HH:mm`.
For example: `03:00`.

* `end_time` - (Required, String) Specifies the end time, in the format: `HH:mm`.
For example: `10:00`. If the end time is the same as the start time, the effective period is the whole day.

* `days_of_week` - (Optional, String) Specifies a list of days of week, separated by commas. 1 represents Sunday,
2 represents Monday, and so on. Defaults to `1,2,3,4,5,6,7` (every day).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

## Import

Device linkage rules can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_iotda_device_linkage_rule.test 62b6cc5aa367f403fea86127
```
