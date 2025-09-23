---
subcategory: "IoT Device Access (IoTDA)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_iotda_device"
description: ""
---

# huaweicloud_iotda_device

Manages an IoTDA device within HuaweiCloud.

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

### Create a directly connected device and an indirectly connected device

```hcl
variable "spaceId" {}
variable "productId" {}
variable "secret" {}

resource "huaweicloud_iotda_device" "device" {
  node_id    = "device_SN_1"
  name       = "device_name"
  space_id   = var.spaceId
  product_id = var.productId
  secret     = var.secret

  tags = {
    foo = "bar"
    key = "value"
  }
}

resource "huaweicloud_iotda_device" "sub_device" {
  node_id    = "device_SN_2"
  name       = "device_name_2"
  space_id   = var.spaceId
  product_id = var.productId
  gateway_id = huaweicloud_iotda_device.device.id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the IoTDA device resource.
If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the device name, which contains `4` to `256` characters. Only letters,
Chinese characters, digits, hyphens (-), underscore (_) and the following special characters are allowed: `?'#().,&%@!`.

* `node_id` - (Required, String, ForceNew) Specifies the node ID, which contains `4` to `256` characters.
The node ID can be IMEI, MAC address, or serial number. Changing this parameter will create a new resource.

* `space_id` - (Required, String, ForceNew) Specifies the resource space ID which the device belongs to.
Changing this parameter will create a new resource.

* `product_id` - (Required, String, ForceNew) Specifies the product ID which the device belongs to.
Changing this parameter will create a new resource.

* `device_id` - (Optional, String, ForceNew) Specifies the device ID, which contains `4` to `256` characters.
Only letters, digits, hyphens (-) and underscore (_) are allowed. If omitted, the platform will automatically allocate
a device ID. Changing this parameter will create a new resource.

* `secret` - (Optional, String) Specifies a primary secret for identity authentication, which contains `8` to `32`
  characters. Only letters, digits, hyphens (-) and underscore (_) are allowed.

* `secondary_secret` - (Optional, String) Specifies a secondary secret for identity authentication.
  When the primary secret verification fails, the secondary secret verification will be enabled, and the secondary
  secret has the same effect as the primary secret; The secondary secret is not effective for devices connected to the
  COAP protocol. Which contains `8` to `32` characters.
  Only letters, digits, hyphens (-) and underscore (_) are allowed.

* `fingerprint` - (Optional, String) Specifies a primary fingerprint of X.509 certificate for identity authentication,
which is a 40-digit or 64-digit hexadecimal string. For more detail, please see
[Registering a Device Authenticated by an X.509 Certificate](https://support.huaweicloud.com/en-us/usermanual-iothub/iot_01_0055.html).

* `secondary_fingerprint` - (Optional, String) Specifies a secondary fingerprint of X.509 certificate for identity
  authentication. When primary fingerprint verification fails, secondary fingerprint verification will be enabled, and
  the secondary fingerprint has the same effectiveness as the primary fingerprint.
  Which is a 40-digit or 64-digit hexadecimal string. For more detail, please see
  [Registering a Device Authenticated by an X.509 Certificate](https://support.huaweicloud.com/en-us/usermanual-iothub/iot_01_0055.html).

-> Only one identity authentication method can be used, either secret or fingerprint. The `secret`, `secondary_secret`
  fields and `fingerprint`, `secondary_fingerprint` fields cannot be set simultaneously.

* `secure_access` - (Optional, Bool) Specifies whether the device is connected through a secure protocol.
  This parameter is only valid when `secret` or `fingerprint` is specified, and suggest setting it to **true**.
  If ignored, it means accessing through insecure protocols, and the device is susceptible to security risks such as
  counterfeiting, please be cautious with this configuration.

* `force_disconnect` - (Optional, Bool) Specifies whether to force device disconnection when resetting secrets or
  fingerprints, currently, only long connections are allowed. The default value is **false**.

* `gateway_id` - (Optional, String) Specifies the gateway ID which is the device ID of the parent device.
The child device is not directly connected to the platform. If omitted, it means to create a device directly connected
to the platform, the `device_id` of the device is the same as the `gateway_id`.

* `description` - (Optional, String) Specifies the description of device. The description contains a maximum of `2,048`
characters. Only letters, Chinese characters, digits, hyphens (-), underscore (_) and the following special characters
are allowed: `?'#().,&%@!`.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the device.

* `frozen` - (Optional, Bool) Specifies whether to freeze the device. Defaults to **false**.

* `extension_info` - (Optional, Map) Specifies the extended information of the device.
  The users can be customized the content. The maximum size of the value is `1` KB.

* `shadow` - (Optional, List) Specifies the initial configuration of the device.
  The [shadow](#device_shadow) structure is documented below.

  -> 1.You can use this parameter to specify the initial configuration for the device. The system compares the property
  values specified by `service_id` and the `desired` section with the default values of the corresponding properties
  in the product. If they are different, the property values specified by the `shadow` parameter are written to the
  device shadow.
  <br>2.The value of `service_id` and the properties in the `desired` section must be defined in the product model.
  <br>3.If you want to config `shadow` data, the `method` value of the properties in the product model must
  contain **W**.

<a name="device_shadow"></a>
The `shadow` block supports:

* `service_id` - (Required, String) Specifies the service ID of the device.
  Which is defined in the product model associated with the device.

* `desired` - (Required, Map) Specifies the initial properties data of the device.
  The each key is a parameter name of a property in the product model.
  If you want to delete the entire `desired`, please enter an empty Map. e.g. **desired = {}**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The device ID in UUID format.

* `status` - The status of device. The valid values are **INACTIVE**, **ONLINE**, **OFFLINE**, **FROZEN**, **ABNORMAL**.

* `auth_type` - The authentication type of device. The options are as follows:
  + **SECRET**: Use a secret for identity authentication.
  + **CERTIFICATES**: Use an x.509 certificate for identity authentication.

* `node_type` - The node type of device. The options are as follows:
  + **GATEWAY**: Directly connected device.
  + **ENDPOINT**: Indirectly connected device.
  + **UNKNOWN**: Unknown type.

## Import

Devices can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_iotda_device.test 10022532f4f94f26b01daa1e424853e1
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason.
The missing attributes include: `force_disconnect`, `extension_info`, `shadow`.
It is generally recommended running `terraform plan` after importing a resource.
You can then decide if changes should be applied to the resource, or the resource definition
should be updated to align with the resource. Also you can ignore changes as below.

```hcl
resource "huaweicloud_iotda_device" "test" { 
  ...
  
  lifecycle {
    ignore_changes = [
      force_disconnect, extension_info, shadow,
    ]
  }
}
```
