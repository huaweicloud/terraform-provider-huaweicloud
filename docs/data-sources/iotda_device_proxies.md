---
subcategory: "IoT Device Access (IoTDA)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_iotda_device_proxies"
description: |-
  Use this data source to get the list of IoTDA device proxies within HuaweiCloud.
---

# huaweicloud_iotda_device_proxies

Use this data source to get the list of IoTDA device proxies within HuaweiCloud.

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
variable "name" {}

data "huaweicloud_iotda_device_proxies" "test" {
  name = var.name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the device proxies.
  If omitted, the provider-level region will be used.

* `space_id` - (Optional, String) Specifies the space ID to which the device proxies belong.
  If omitted, query all device proxies under the current instance.

* `name` - (Optional, String) Specifies the name of the device proxy.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID in UUID format.

* `proxies` - All device proxies that match the filter parameters.
  The [proxies](#iotda_proxies) structure is documented below.

<a name="iotda_proxies"></a>
The `proxies` block supports:

* `id` - The device proxy ID.

* `name` - The device proxy name.

* `space_id` - The space ID to which the device proxy belongs.

* `effective_time_range` - The validity period of the device proxy rule.
  The [effective_time_range](#IoTDA_effective_time_range) structure is documented below.

<a name="IoTDA_effective_time_range"></a>
The `effective_time_range` block supports:

* `start_time` - The effective time of the device proxy, using UTC time zone,
  the format is **yyyyMMdd'T'HHMmmss-Z**. e.g. **20250528T153000Z**.

* `end_time` - The device proxy expiration time, using UTC time zone,
  the format is **yyyyMMdd'T'HHMmmss-Z**. e.g. **20250528T153000Z**.
