---
subcategory: "Image Management Service (IMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ims_image_export"
description: |-
  Manages an IMS image export resource within HuaweiCloud.
---

# huaweicloud_ims_image_export

Manages an IMS image export resource within HuaweiCloud.

-> 1. The whole image, ISO image, private images created by public images of Windows, SUSE, Red Hat, Ubuntu, and
   Oracle Linux, and private images created by market images are not allowed to be exported.
   <br/>2. The image size must be less than `1` TB. The images larger than `128` GB only support fast export, the
   maximum image size supported for export in some regions may be greater than `128` GB, please refer to the actual
   operation prompts on the console for accuracy.<br/>3. Destroying resource does not change the current action of the
   image export resource.

## Example Usage

### Ordinary export image

```hcl
variable "image_id" {}
variable "bucket_url" {}
variable "file_format" {}

resource "huaweicloud_ims_image_export" "test" {
  image_id    = var.image_id
  bucket_url  = var.bucket_url
  file_format = var.file_format
}
```

### Quickly export image

```hcl
variable "image_id" {}
variable "bucket_url" {}

resource "huaweicloud_ims_image_export" "test" {
  image_id        = var.image_id
  bucket_url      = var.bucket_url
  is_quick_export = true
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `image_id` - (Required, String, ForceNew) Specifies the image ID to export.
  Changing this parameter will create a new resource.

* `bucket_url` - (Required, String, ForceNew) Specifies the URL of the image file to be exported to the OBS bucket, the
  format is **OBS bucket name:image file name**, e.g. **test_bucket:test_image_file**. The storage category of the OBS
  bucket and image file here must be OBS standard storage. Changing this parameter will create a new resource.

* `file_format` - (Optional, String, ForceNew) Specifies the format of the image file to be exported. The valid values
  are **qcow2**, **vhd**, **zvhd**, or **vmdk**. Changing this parameter will create a new resource.

* `is_quick_export` - (Optional, Bool, ForceNew) Specifies whether to use quick export. The valid value is **true** or
  **false**. Changing this parameter will create a new resource.

-> 1. When the `is_quick_export` parameter is ignored or set to **false**, the `file_format` parameter is required.
   <br/>2. When the `is_quick_export` parameter is set to **true**, the `file_format` parameter must be ignored, and
   the exported image file format is **zvhd2** at this time.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID, same as `image_id`.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 20 minutes.
