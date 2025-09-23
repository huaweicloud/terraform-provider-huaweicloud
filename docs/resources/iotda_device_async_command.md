---
subcategory: "IoT Device Access (IoTDA)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_iotda_device_async_command"
description: |-
  Manages a device asynchronous command delivery resource within HuaweiCloud.
---

# huaweicloud_iotda_device_async_command

Manages a device asynchronous command delivery resource within HuaweiCloud.

-> 1.This resource is only a one-time action resource for doing API action. Deleting this resource will not clear
  the corresponding request record, but will only remove the resource information from the tfstate file.
  <br>2.Currently, this resource is only supported deliver commands asynchronously to NB-IoT devices.
  <br>3.After the resource is created, please pay attention to the command executed result through `status`,
  you can execute the **terraform plan** command at regular intervals to monitor `status` changes.

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
variable "send_strategy"{}

resource "huaweicloud_iotda_device_async_command" "test" {
  device_id     = var.device_id
  send_strategy = var.send_strategy
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `device_id` - (Required, String, ForceNew) Specifies the ID of the device to which the command is delivered.
  Changing this parameter will create a new resource.

* `send_strategy` - (Required, String, ForceNew) Specifies the delivery policy.
  The valid values are as follows:
  + **immediately**: The command is delivered immediately.
  + **delay**: The command is cached and delivered after the device reports data or goes online.
  
  Changing this parameter will create a new resource.

* `service_id` - (Optional, String, ForceNew) Specifies the ID of the device service to which the device command belongs,
  which is defined in the product model associated with the device.
  This parameter is mandatory if the device requires codecs to parse commands.
  Changing this parameter will create a new resource.

* `name` - (Optional, String, ForceNew) Specifies the command name, which is defined in the product model
  associated with the device.
  This parameter is mandatory if the device requires codecs to parse commands.
  Changing this parameter will create a new resource.

* `paras` - (Optional, Map, ForceNew) Specifies the command executed by the device.
  If `service_id` is specified, each key is the parameter in commands in the product model.
  If `service_id` is left empty, the key can be customized.
  The maximum size of the request object is `256` KB.
  Changing this parameter will create a new resource.

* `expire_time` - (Optional, Int, ForceNew) Specifies the duration of caching commands on the IoT platform.
  This parameter is valid only when `send_strategy` is set to **delay**. The unit is second.
  If `expire_time` is set to **0** or not specified, the command is cached for `24` hours (`86,400` seconds) by default,
  and the maximum cache duration is `2` days.
  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `status` - The status of the command.
  The valid values are as follows:
  + **PENDING**: The command is not delivered and is cached on the platform.
  + **EXPIRED**: The command has expired, the cache time exceeds the value of `expire_time`.
  + **SENT**: The command is being delivered.
  + **DELIVERED**: The command has been delivered.
  + **SUCCESSFUL**: The command has been executed.
  + **FAILED**: The command fails to be executed.
  + **TIMEOUT**: After the command is delivered, no response is received from the device or the response times out.

* `result` - The command execution result.

* `sent_time` - The time of the platform sent the command.
  The format is **yyyyMMdd'T'HHmmss'Z'**, e.g. **20151212T121212Z**.

* `delivered_time` - The time of the device received the command.
  The format is **yyyyMMdd'T'HHmmss'Z'**, e.g. **20151212T121212Z**.

* `response_time` - The time of the device responded to the command.
  The format is **yyyyMMdd'T'HHmmss'Z'**, e.g. **20151212T121212Z**.

* `created_at` - The creation time of the device command.
  The format is **yyyyMMdd'T'HHmmss'Z'**, e.g. **20151212T121212Z**.
