---
subcategory: "IoT Device Access (IoTDA)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_iotda_device_certificates"
description: ""
---

# huaweicloud_iotda_device_certificates

Use this data source to get the list of IoTDA device CA certificates within HuaweiCloud.

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
variable "certificate_id" {}

data "huaweicloud_iotda_device_certificates" "test" {
  certificate_id = var.certificate_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the certificates.
  If omitted, the provider-level region will be used.

* `space_id` - (Optional, String) Specifies the space ID of the certificates to be queried.
  If omitted, query the certificates for all spaces under the current instance.

* `certificate_id` - (Optional, String) Specifies the ID of the certificate to be queried.

* `cn` - (Optional, String) Specifies the CN name of the certificates to be queried.

* `status` - (Optional, String) Specifies the verification status of the certificates to be queried.  
  The value can be **Verified** or **Unverified**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `certificates` - All certificates that match the filter parameters.
  The [certificates](#iotda_device_certificates) structure is documented below.

<a name="iotda_device_certificates"></a>
The `certificates` block supports:

* `id` - The certificate ID.

* `cn` - The certificate CN name.

* `owner` - The certificate owner.

* `status` - The certificate verification status.

* `verify_code` - The certificate verification code.

* `created_at` - The creation time of the certificate.

* `effective_date` - The effective time of the certificate.

* `expiry_date` - The expiration time of the certificate.
