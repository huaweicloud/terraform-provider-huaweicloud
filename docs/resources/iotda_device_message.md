---
subcategory: "IoT Device Access (IoTDA)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_iotda_device_message"
description: |-
  Manages a device message delivery resource within HuaweiCloud.
---

# huaweicloud_iotda_device_message

Manages a device message delivery resource within HuaweiCloud.

-> 1.This resource is only a one-time action resource for doing API action. Deleting this resource will not clear
  the corresponding request record, but will only remove the resource information from the tfstate file.
  <br>2.Currently, this resource is only supported deliver message to MQTT devices.
  <br>3.After the resource is created, please pay attention to the message delivery result through `status` attribute,
  you can execute the **terraform plan** command at regular intervals to monitor `status` attribute changes.

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
variable "device_id"{}
variable "message"{}

resource "huaweicloud_iotda_device_message" "test" {
  device_id = var.device_id
  message   = var.message
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `device_id` - (Required, String, ForceNew) Specifies the ID of the device to which the message is delivered.
  Changing this parameter will create a new resource.

* `message` - (Required, String, ForceNew) Specifies message content, supporting string and json formats.
  Changing this parameter will create a new resource.

* `message_id` - (Optional, String, ForceNew) Specifies the message ID, it is recommended to use UUID, which is unique
  within the same device. The length should not exceed `128`, and only combinations of letters, numbers,
  underscores (_), and hyphens (-) are allowed. Changing this parameter will create a new resource.
  + If the message ID filled in is not unique within the device, an error will be returned.
  + If left blank, the system will generate a random ID.

* `name` - (Optional, String, ForceNew) Specifies the message name. The length should not exceed `128`, and only
  Chinese, letters, numbers, and the following characters are allowed: `_?'#().,&%@!-`.
  Changing this parameter will create a new resource.

* `properties` - (Optional, List, ForceNew) Specifies the attribute parameters of the message downstream to the device.
  Changing this parameter will create a new resource.
  The [properties](#iotda_properties) structure is documented below.

* `encoding` - (Optional, String, ForceNew) Specifies the encoding format for message content.
  Changing this parameter will create a new resource.  
  The valid values are as follows:
  + **none**
  + **base64**

  Defaults to **none**.

* `payload_format` - (Optional, String, ForceNew) Specifies the payload format. This parameter is valid only `encoding`
  is set to **none**. Changing this parameter will create a new resource.  
  The valid values are as follows:
  + **standard**: Standard format for platform encapsulation.
  + **raw**: Directly distribute the message content as a payload.

  Defaults to **standard**.

* `topic` - (Optional, String, ForceNew) Specifies the custom topic suffix for message downstream to the device.
  Only topics configured on the tenant product interface can be filled in, otherwise the verification will not pass.
  If the topic is specified, the message will be directed to the device through that topic. If not specified, the
  message will be directed to the device through the system's default topic.
  Changing this parameter will create a new resource.

* `topic_full_name` - (Optional, String, ForceNew) Specifies the complete topic name for the message to be sent to the
  device. When it is necessary to issue a custom topic to the device, this parameter can be used to specify the complete
  topic name. The IoT platform does not verify whether the topic is defined on the platform and directly transmits it to
  the device. The device needs to subscribe to the topic in advance.
  Changing this parameter will create a new resource.

-> Only one of parameters `topic` and `topic_full_name` can be set.

* `ttl` - (Optional, String, ForceNew) Specifies the aging time of the message in the platform cache, in minutes.
  Changing this parameter will create a new resource.
  + The valid value must be a multiple of `5`.
  + When specified as `0`, it means that the message is not cached, and the default maximum caching time is `1,440`.

  Defaults to `1,440`.

<a name="iotda_properties"></a>
The `properties` block supports:

* `correlation_data` - (Optional, String, ForceNew) Specifies relevant data in MQTT 5.0 request and response patterns.
  The length should not exceed `128`, and only combinations of letters, numbers, underscores (_), and hyphens (-) are
  allowed. Changing this parameter will create a new resource.

* `response_topic` - (Optional, String, ForceNew) Specifies response topic in MQTT 5.0 request and response patterns.
  The length should not exceed 128, and only letters, numbers, and the following characters are allowed: `_-?=$#+/`.
  Changing this parameter will create a new resource.

* `user_properties` - (Optional, List, ForceNew) Specifies user-defined attributes. The maximum number that can be
  configured is `20`. Changing this parameter will create a new resource.
  The [user_properties](#iotda_user_properties) structure is documented below.

<a name="iotda_user_properties"></a>
The `user_properties` block supports:

* `prop_key` - (Optional, String, ForceNew) Specifies custom attribute key. The length should not exceed `128`, and only
  combinations of letters, numbers, underscores (_), and hyphens (-) are allowed.
  Changing this parameter will create a new resource.

* `prop_value` - (Optional, String, ForceNew) Specifies custom attribute value. The length should not exceed `128`, and
  only Chinese, letters, numbers, and the following characters are allowed: `_? '#().,&%@!-`.
  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID, same as `message_id`.

* `status` - The status of the message.  
  The valid values are as follows:
  + **PENDING**: The device is not online, the message is cached and will be issued after the device is online.
  + **DELIVERED**: The message sent successfully.
  + **FAILED**: The message sending failed.
  + **TIMEOUT**: The message has not been sent to the device within the default time of the platform (`1` day).

* `error_info` - The message delivery failure details.
  The [error_info](#iotda_error_info) structure is documented below.

* `created_time` - The creation time of the device message.
  The format is **yyyyMMdd'T'HHmmss'Z'**, e.g. **20151212T121212Z**.

* `finished_time` - The end time of the device message. Contains the time for the message to transition to the
  **DELIVERED* and **TIMEOUT** status. The format is **yyyyMMdd'T'HHmmss'Z'**, e.g. **20151212T121212Z**.

<a name="iotda_error_info"></a>
The `error_info` block supports:

* `error_code` - The abnormal information error code.  
  The valid values are as follows:
  + **IOTDA.014016**: Indicates that the device is not online.
  + **IOTDA.014112**: Indicates that the device has not subscribed to the topic.

* `error_msg` - The abnormal information explanation. Includes instructions for devices not online and devices not
  subscribed to the topic.
