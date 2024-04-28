---
subcategory: "IoT Device Access (IoTDA)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_iotda_device_groups"
description: ""
---

# huaweicloud_iotda_device_groups

Use this data source to get the list of the IoTDA device groups.

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
variable "device_group_id" {}
data "huaweicloud_iotda_device_groups" "test" {
  group_id = var.device_group_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the device groups.
  If omitted, the provider-level region will be used.

* `group_id` - (Optional, String) Specifies the ID of the device group.

* `name` - (Optional, String) Specifies the name of the device group.

* `type` - (Optional, String) Specifies the type of the device groups.
  The valid values are as follows:
  + **STATIC**: The device group is a static group.
  + **DYNAMIC**: The device group is a dynamical group.

* `parent_group_id` - (Optional, String) Specifies the ID of the parent device group to which the device group belongs.

* `space_id` - (Optional, String) Specifies the ID of the resource space to which the device groups belong.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `groups` - The list of the device groups.
  The [groups](#iotda_device_groups) structure is documented below.

<a name="iotda_device_groups"></a>
The `groups` block supports:

* `id` - The ID of the device group.

* `name` - The name of the device group.

* `description` - The description of the device group.

* `parent_group_id` - The ID of the parent device group to which the device group belongs.

* `type` - The type of the device group.
