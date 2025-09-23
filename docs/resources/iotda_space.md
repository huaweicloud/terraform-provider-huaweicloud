---
subcategory: "IoT Device Access (IoTDA)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_iotda_space"
description: -|
  Manages an IoTDA resource space within HuaweiCloud.
---

# huaweicloud_iotda_space

Manages an IoTDA resource space within HuaweiCloud.

A resource space is the basic unit of service management and provides independent device management and platform
configuration capabilities at the service layer. Resources (such as products and devices) must be created on
a resource space.

-> The **basic** edition instance does not support updating the resource.

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
variable "name" {}

resource "huaweicloud_iotda_space" "test" {
  name = var.name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the IoTDA resource space resource.
If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the space name. The name contains a maximum of `64` characters.
Only letters, digits, hyphens (-), underscore (_) and the following special characters are allowed: `?'#().,&%@!`.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

* `is_default` - Whether it is the default resource space. The IoT platform automatically creates and assigns
a default resource space (undeletable) to your account.

## Import

The resource can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_iotda_space.test <id>
```
