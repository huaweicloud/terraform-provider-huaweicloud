---
subcategory: "IoT Device Access (IoTDA)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_iotda_product"
description: -|
  Manages an IoTDA device proxy resource within HuaweiCloud.
---

# huaweicloud_iotda_product

Manages an IoTDA device proxy resource within HuaweiCloud.

-> Currently, device proxy resources are only supported on IoTDA **standard** or **enterprise** edition instance.
  When accessing an IoTDA **standard** or **enterprise** edition instance, you need to specify
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
variable "space_id" {}
variable "name" {}
variable "devices" {}
variable "start_time" {}
variable "end_time" {}

resource "huaweicloud_iotda_device_proxy" "test" {
  space_id = var.space_id
  name     = var.name
  devices  = var.devices

  effective_time_range {
    start_time = var.start_time
    end_time   = var.end_time
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `space_id` - (Required, String, ForceNew) Specifies the resource space ID to which the device proxy belongs.

* `name` - (Required, String) Specifies the device proxy name. The valid length is limited from `1` to `64`.

* `devices` - (Required, List) Specifies the list of proxy device IDs. The number of IDs in the list is limited from
  `2` to `10`. All devices in the list share gateway permissions, which means that any sub device under any gateway in
  the list can go online through any gateway in the group and report data.

* `effective_time_range` - (Required, List) Specifies the validity period of the device proxy rule.
  The [effective_time_range](#IoTDA_effective_time_range) structure is documented below.

<a name="IoTDA_effective_time_range"></a>
The `effective_time_range` block supports:

* `start_time` - (Optional, String) Specifies the effective time of the device proxy, using UTC time zone,
  the format is **yyyyMMdd'T'HHMmmss-Z**. e.g. **20250528T153000Z**.

* `end_time` - (Optional, String) Specifies the device proxy expiration time, must be greater than `start_time`,
  using UTC time zone, the format is **yyyyMMdd'T'HHMmmss-Z**. e.g. **20250528T153000Z**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

## Import

The device proxy can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_iotda_device_proxy.test <id>
```
