---
subcategory: "IoT Device Access (IoTDA)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_iotda_device_certificate"
description: ""
---

# huaweicloud_iotda_device_certificate

Manages an IoTDA device CA certificate within HuaweiCloud.

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
variable "certificateContent" {}

resource "huaweicloud_iotda_device_certificate" "test" {
  content  = var.certificateContent
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the IoTDA device CA certificate
resource. If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `content` - (Required, String, ForceNew) Specifies the content of device CA certificate.
Changing this parameter will create a new resource. Create a private CA certificate,
please following [reference](https://support.huaweicloud.com/usermanual-iothub/iot_01_0104.html)

* `space_id` - (Optional, String, ForceNew) Specifies the resource space ID to which the device CA certificate belongs.
If omitted, the certificate will belong to the default resource space.
Changing this parameter will create a new resource.

* `verify_content` - (Optional, String) Specifies the content of verification certificate. Can only be used to verify
the validity of the device CA certificate after creation. Get the verification certificate,
please following [reference](https://support.huaweicloud.com/usermanual-iothub/iot_01_0106.html)

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `cn` - The CN name of the device CA certificate.

* `owner` - The owner of the device CA certificate.

* `status` - The status of the device CA certificate. The valid values are **Unverified** and **Verified**.

* `verify_code` - The verify code of the device CA certificate.

* `effective_date` - The effective date of the device CA certificate.
The format is: **yyyyMMdd'T'HHmmss'Z'**, e.g., **20151212T121212Z**.

* `expiry_date` - The expiry date of the device CA certificate.
The format is: **yyyyMMdd'T'HHmmss'Z'**, e.g., **20151212T121212Z**.

## Import

Device CA certificates can be imported by `id`, e.g.

```bash
$ terraform import huaweicloud_iotda_device_certificate.test 62b3cec5558d4b703f064534
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `content`, `space_id`, `verify_content`.
It is generally recommended running `terraform plan` after importing the resource. You can then decide if changes should
be applied to the resource, or the resource definition should be updated to align with the group. Also you can ignore
changes as below.

```hcl
resource "huaweicloud_iotda_device_certificate" "test" {
    ...

  lifecycle {
    ignore_changes = [
      content, space_id, verify_content
    ]
  }
}
```
