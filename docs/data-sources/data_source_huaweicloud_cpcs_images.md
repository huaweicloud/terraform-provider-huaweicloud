---
subcategory: "Data Encryption Workshop (DEW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cpcs_images"
description: |-
  Use this data source to get the list of CPCS images.
---

# huaweicloud_cpcs_images

Use this data source to get the list of CPCS images.

-> Currently, this data source is valid only in cn-north-9 region.

## Example Usage

```hcl
variable "service_type" {}

data "huaweicloud_cpcs_images" "test" {
  service_type = var.service_type
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `image_name` - (Optional, String) Specifies the name of the image.

* `service_type` - (Optional, String) Specifies the service type of the image.
  Valid values are **ENCRYPT_DECRYPT**, **SIGN_VERIFY**, **KMS**, **TIMESTAMP**, **COLLA_SIGN**, **OTP**,
  **DB_ENCRYPT**, **FILE_ENCRYPT**, **DIGIT_SEAL** and **SSL_VPN**.

* `sort_key` - (Optional, String) Specifies the sort attribute.
  The default value is **create_time**.

* `sort_dir` - (Optional, String) Specifies the sort direction.
  The default value is **DESC**. Valid values are **ASC** and **DESC**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `images` - Indicates the images list.
  The [images](#CPCS_images) structure is documented below.

<a name="CPCS_images"></a>
The `images` block supports:

* `image_id` - The image ID.

* `image_name` - The image name.

* `service_type` - The service type of the image.

* `arch_type` - The system architecture of the image. Valid values are **X86_64** and **ARRCH**.

* `specification_id` - The specification ID.

* `create_time` - The creation time of the image, in UNIX timestamp format.

* `version_type` - The version type.

* `trust_domain` - The domain whitelist.

* `vendor_name` - The vendor name.

* `vendor_image_version` - The vendor image version.

* `ccsp_version_need` - The required platform version.

* `hw_image_version` - The Huawei image version.

* `description` - The description.
