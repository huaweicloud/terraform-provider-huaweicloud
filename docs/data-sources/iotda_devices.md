---
subcategory: "IoT Device Access (IoTDA)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_iotda_devices"
description: |-
  Use this data source to get the list of IoTDA devices within HuaweiCloud.
---

# huaweicloud_iotda_devices

Use this data source to get the list of IoTDA devices within HuaweiCloud.

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
variable "device_id" {}

data "huaweicloud_iotda_devices" "test" {
  device_id = var.device_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the devices.
  If omitted, the provider-level region will be used.

* `space_id` - (Optional, String) Specifies the space ID of the devices to be queried.
  If omitted, query the devices in all spaces under the current instance.

* `product_id` - (Optional, String) Specifies the ID of the product to be queried.
  If omitted, query the devices in all products under the current instance.

* `gateway_id` - (Optional, String) Specifies the gateway ID of the devices to be queried;
  The `gateway_id` is the ID of the parent device to which the devices belong.

* `is_cascade` - (Optional, Bool) Specifies whether to cascade queries, this parameter only takes effect when
  carrying `gateway_id` simultaneously. The default value is **false**.  
  The valid values are as follows:
  + **true**: Represents querying all levels of sub devices under a device with a device ID equal to gateway ID.
  + **false**: Represents the first level child devices under the device with the query device ID equal to gateway ID.
  
* `node_id` - (Optional, String) Specifies the node ID of the device to be queried.

* `node_type` - (Optional, String) Specifies the node type of the devices to be queried.  
  The valid values are as follows:
  + **ENDPOINT**: Non-directly connected devices.
  + **GATEWAY**: Directly connected devices or gateways.
  + **UNKNOWN**: Unknown.

* `status` - (Optional, String) Specifies the status of the devices to be queried.  
  The valid values are as follows:
  + **ONLINE**: Device is online.
  + **OFFLINE**: Device is offline.
  + **ABNORMAL**: Device is abnormal.
  + **INACTIVE**: Device is inactive.
  + **FROZEN**: Device is frozen.

* `device_id` - (Optional, String) Specifies the ID of the device to be queried.

* `name` - (Optional, String) Specifies the name of the device to be queried.

* `start_time` - (Optional, String) Specifies the start time to be queried. The query result shows devices created after
  this time (including devices created at this time). The format is **yyyyMMdd'T'HHmmss'Z**. e.g. **20190528T153000Z**.

* `end_time` - (Optional, String) Specifies the end time to be queried. The query result is for devices created before
  this time (excluding devices created at this time). The format is **yyyyMMdd'T'HHmmss'Z**. e.g. **20190528T153000Z**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `devices` - All devices that match the filter parameters.
  The [devices](#iotda_devices) structure is documented below.

<a name="iotda_devices"></a>
The `devices` block supports:

* `space_id` - The space ID to which the device belongs.

* `space_name` - The space name to which the device belongs.

* `product_id` - The product ID to which the device belongs.

* `product_name` - The product name to which the device belongs.

* `gateway_id` - The ID of the parent device to which the device belongs

* `id` - The device ID.

* `name` - The device name.

* `node_id` - The node ID of the device.

* `node_type` - The node type of the device.

* `description` - The description of the device.

* `status` - The status of the device.

* `fw_version` - The firmware version of the device.

* `sw_version` - The software version of the device.

* `sdk_version` - The SDK information of the device.

* `tags` - The tags of the device, key/value pair format.
