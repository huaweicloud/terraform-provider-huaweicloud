---
subcategory: "IoT Device Access (IoTDA)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_iotda_upgrade_packages"
description: |-
  Use this data source to get a list of the upgrade packages.
---

# huaweicloud_iotda_upgrade_packages

Use this data source to get a list of the upgrade packages.

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
variable "type" {}

data "huaweicloud_iotda_upgrade_packages" "test" {
  type = var.type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `type` - (Required, String) Specifies the type of the upgrade package.
  The valid values are as follows:
  + **softwarePackage**
  + **firmwarePackage**

* `space_id` - (Optional, String) Specifies the resource space ID to which the upgrade packages belong.

* `product_id` - (Optional, String) Specifies the product ID associated with the upgrade package.

* `version` - (Optional, String) Specifies the version number of the upgrade package.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `packages` - The list of the upgrade packages.

  The [packages](#packages_struct) structure is documented below.

<a name="packages_struct"></a>
The `packages` block supports:

* `id` - The ID of the upgrade package.

* `space_id` - The resource space ID to which the upgrade package belongs.

* `product_id` - The product ID associated with the upgrade package.

* `type` - The type of the upgrade package.

* `description` - The description of the upgrade package.

* `version` - The version number of the upgrade package.

* `support_source_versions` - The list of source versions that support the upgrade of this version package.

* `custom_info` - The custom information pushed to the device.

* `created_at` - The time when the software and firmware packages are uploaded to the IoT platform.
  The format is **yyyyMMdd'T'HHmmss'Z**. e.g. **20190528T153000Z**.
