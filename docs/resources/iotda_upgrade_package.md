---
subcategory: "IoT Device Access (IoTDA)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_iotda_upgrade_package"
description: -|
  Manages an IoTDA OTA upgrade package within HuaweiCloud.
---

# huaweicloud_iotda_upgrade_package

Manages an IoTDA OTA upgrade package within HuaweiCloud.

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
variable "space_id" {}
variable "type" {}
variable "product_id" {}
variable "version" {}
variable "region" {}
variable "bucket_name" {}
variable "object_key" {}

resource "huaweicloud_iotda_upgrade_package" "test" {
  space_id   = var.space_id
  type       = var.type
  product_id = var.product_id
  version    = var.version
  
  file_location {
    obs_location  { 
      region      = var.region
      bucket_name = var.bucket_name
      object_key  = var.object_key
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the IoTDA upgrade package resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `space_id` - (Required, String, ForceNew) Specifies the resource space ID to which the upgrade package belongs.
  The length does not exceed `36`, and only combinations of letters, numbers, underscores (_), and connectors (-)
  are allowed. Changing this parameter will create a new resource.

* `type` - (Required, String, ForceNew) Specifies the upgrade package type.
  Changing this parameter will create a new resource.  
  The valid value are as follows:
  + `softwarePackage`: Software package.
  + `firmwarePackage`: Firmware package.

* `product_id` - (Required, String, ForceNew) Specifies the product ID associated with the upgrade package. The length
  does not exceed `36`, and only combinations of letters, numbers, underscores (_), and connectors (-) are allowed.
  Changing this parameter will create a new resource.

* `version` - (Required, String, ForceNew) Specifies the version number of the upgrade package. The length does not
  exceed `256`, and only combinations of letters, numbers, underscores (_), connectors (-), and English dots (.)
  are allowed. Changing this parameter will create a new resource.

* `file_location` - (Required, List, ForceNew) Specifies the location of the upgrade package.
  Changing this parameter will create a new resource.  
  The [file_location](#iotda_upgrade_package_file_location) structure is documented below.

* `support_source_versions` - (Optional, List, ForceNew) Specifies a list of device source version numbers supported for
  upgrading this version pack. This is a list of string. The device source version number only allows combinations of
  letters, numbers, underscores (_), connectors (-), and English dots (.).
  Changing this parameter will create a new resource.

* `description` - (Optional, String, ForceNew) Specifies the description of the upgrade package. The length does not
  exceed `1024`. Changing this parameter will create a new resource.

* `custom_info` - (Optional, String, ForceNew) Specifies the custom information to be pushed to the device.
  After creating the upgrade package and completing the upgrade task, when the IoT platform issues an upgrade
  notification to the device, it will send the customized information to the device. The length does not exceed `4096`.
  Changing this parameter will create a new resource.

<a name="iotda_upgrade_package_file_location"></a>
The `file_location` block supports:

* `obs_location` - (Optional, List, ForceNew) Specifies the location of the OBS object associated with the upgrade
  package. Changing this parameter will create a new resource.  
  The [obs_location](#iotda_upgrade_package_obs_location) structure is documented below.

<a name="iotda_upgrade_package_obs_location"></a>
The `obs_location` block supports:

* `region` - (Required, String, ForceNew) Specifies the region where OBS is located.
  Changing this parameter will create a new resource.

* `bucket_name` - (Required, String, ForceNew) Specifies the name of the OBS bucket where the upgrade package is located.
  Changing this parameter will create a new resource.

* `object_key` - (Required, String, ForceNew) Specifies the name of the OBS object where the upgrade package is located,
  including the folder path. The maximum size of OBS objects is **1GB**, and only supports files in **.bin**, **.dav**,
  **.tar**, **.gz**, **.zip**, **.gzip**, **.apk**, **.ta.gz**, **.tar.xz**, **.pack**, **.exe**, **.bat** and **.img**
  formats. The valid length is limited from `1` to `1024`.
  Changing this parameter will create a new resource.

* `sign` - (Optional, String, ForceNew) Specifies the signature value of the upgrade package calculated by SHA256
  algorithm. After added the upgrade package and created the upgrade task, when the IoT platform issues an upgrade
  notification to the device, it will send the signature to the device.
  The valid length is `64`, only letters `a(A)` to `f(F)` and digits are allowed.

  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `created_at` - The time for uploading the upgrade package to the IoT platform.
  The format is **yyyyMMdd'T'HHmmss'Z**. e.g. **20190528T153000Z**.

## Import

The OTA upgrade package can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_iotda_upgrade_package.test <id>
```
