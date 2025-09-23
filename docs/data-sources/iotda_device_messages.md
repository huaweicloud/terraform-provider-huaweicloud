---
subcategory: "IoT Device Access (IoTDA)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_iotda_device_messages"
description: |-
  Use this data source to get the list of IoTDA device messages within HuaweiCloud.
---

# huaweicloud_iotda_device_messages

Use this data source to get the list of IoTDA device messages within HuaweiCloud.

-> When accessing an IoTDA **standard** or **enterprise** edition instance, you need to specify
  the IoTDA service endpoint in `provider` block.
  You can login to the IoTDA console, choose the instance **Overview** and click **Access Details**
  to view the HTTPS application access address. An example of the access address might be
  *9bc34xxxxx.st1.iotda-app.ap-southeast-1.myhuaweicloud.com*, then you need to configure the
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
variable "device_id" {}

data "huaweicloud_iotda_device_messages" "test" {
  device_id = var.device_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the device messages.
  If omitted, the provider-level region will be used.

* `device_id` - (Required, String) Specifies the device ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `messages` - All device messages that match the filter parameters.
  The [messages](#iotda_messages) structure is documented below.

<a name="iotda_messages"></a>
The `messages` block supports:

* `id` - The message ID.

* `name` - The message name.

* `message` - The message content.

* `encoding` - The encoding format for message content. The value can be **none** or **base64**.

* `payload_format` - The payload format. The value can be **standard** or **raw**.

* `topic` - The message topic.

* `properties` - The attribute parameters of the message downstream to the device.
  The [properties](#iotda_properties) structure is documented below.

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

<a name="iotda_properties"></a>
The `properties` block supports:

* `correlation_data` - The relevant data in MQTT 5.0 request and response patterns.

* `response_topic` - The response topic in MQTT 5.0 request and response patterns.

* `user_properties` - The user-defined attributes.
  The [user_properties](#iotda_user_properties) structure is documented below.

<a name="iotda_user_properties"></a>
The `user_properties` block supports:

* `prop_key` - The custom attribute key.

* `prop_value` - The custom attribute value.

<a name="iotda_error_info"></a>
The `error_info` block supports:

* `error_code` - The abnormal information error code.  
  The valid values are as follows:
  + **IOTDA.014016**: Indicates that the device is not online.
  + **IOTDA.014112**: Indicates that the device has not subscribed to the topic.

* `error_msg` - The abnormal information explanation. Includes instructions for devices not online and devices not
  subscribed to the topic.
