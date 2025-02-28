---
subcategory: "IoT Device Access (IoTDA)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_iotda_device_policy"
description: -|
  Manages an IoTDA device policy resource within HuaweiCloud.
---

# huaweicloud_iotda_device_policy

Manages an IoTDA device policy resource within HuaweiCloud.

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
variable "policy_name" {}
variable "effect" {}
variable "action" {}
variable "resource" {}

resource "huaweicloud_iotda_device_policy" "test" {
  policy_name = var.policy_name

  statement {
    effect    = var.effect
    actions   = [var.action]
    resources = [var.resource]
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `policy_name` - (Required, String) Specifies the device policy name. The length should not exceed `128`, and only
  combinations of letters, numbers, underscores (_), and hyphens (-) are allowed.

* `statement` - (Required, List) Specifies the policy document.  
  The [statement](#device_policy_statement) structure is documented below.

* `space_id` - (Optional, String, NonUpdatable) Specifies the resource space ID to which the device policy belongs.
  If omitted, the created policy will belong to the default resource space.

<a name="device_policy_statement"></a>
The `statement` block supports:

* `effect` - (Required, String) Specifies whether to allow or deny the operation. When there are authorization
  statements that both allow and deny, follow the principle of prioritizing denial.  
  The valid values are as follows:
  + **ALLOW**
  + **DENY**

* `actions` - (Required, List) Specifies the operations allowed or denied by the policy. This value is in string list
  format, the format of a single operation is **service name:resource:operation**.  
  The valid values are as follows:
  + **iotda:devices:publish**: The device uses MQTT protocol to publish messages.
  + **iotda:devices:subscribe**: The device subscribes to messages using the MQTT protocol.

* `resources` - (Required, List) Specifies the resources that allow or deny operations to be performed on them.
  This value is in string list format, the format of a single resource is **resource type:resource name**.
  For example, the resources subscribed to by the device are **topic:/v1/${devices.deviceId}/test/hello**.  
  When using this parameter,
  please refer to the [documentation](https://support.huaweicloud.com/intl/en-us/usermanual-iothub/iot_01_1114.html).

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `create_time` - The creation time of the device policy.
  The format is **yyyyMMdd'T'HHmmss'Z'**. e.g. **20151212T121212Z**.

* `update_time` - The latest update time of the device policy.
  The format is **yyyyMMdd'T'HHmmss'Z'**. e.g. **20151212T121212Z**.

## Import

The device policy can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_iotda_device_policy.test <id>
```
