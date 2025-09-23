---
subcategory: "IoT Device Access (IoTDA)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_iotda_device_binding_groups"
description: |-
  Use this data source to get the list of IoTDA device groups bound to a specified device within HuaweiCloud.
---

# huaweicloud_iotda_device_binding_groups

Use this data source to get the list of IoTDA device groups bound to a specified device within HuaweiCloud.

-> Currently, this data source is only supported on IoTDA **standard** or **enterprise** edition instance.
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
variable "device_id" {}

data "huaweicloud_iotda_device_binding_groups" "test" {
  device_id = var.device_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the device groups.
  If omitted, the provider-level region will be used.

* `device_id` - (Required, String) Specifies the ID of the device.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `groups` - The list of device groups bound to specified device.
  The [groups](#iotda_device_groups) structure is documented below.

<a name="iotda_device_groups"></a>
The `groups` block supports:

* `id` - The ID of the device group.

* `name` - The name of the device group.

* `description` - The description of the device group.

* `parent_group_id` - The ID of the parent device group to which the device group belongs.

* `type` - The type of the device group. The valid values are **STATIC** and **DYNAMIC**.
