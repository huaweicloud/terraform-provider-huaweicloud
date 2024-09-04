---
subcategory: "Image Management Service (IMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ims_obs_data_image"
description: |-
  Manages an IMS data image resource created from external image file in the OBS bucket within HuaweiCloud.
---

# huaweicloud_ims_obs_data_image

Manages an IMS data image resource created from external image file in the OBS bucket within HuaweiCloud.

## Example Usage

### Creating an IMS data image from an external image file in the OBS bucket

```hcl
variable "name" {}
variable "image_url" {}
variable "min_disk" {}

resource "huaweicloud_ims_obs_data_image" "test" {
  name      = var.name
  image_url = var.image_url
  min_disk  = var.min_disk
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the name of the image.
  The valid length is limited from `1` to `128` characters.
  The first and last letters of the name cannot be spaces.
  The name can contain uppercase letters, lowercase letters, numbers, spaces, chinese, and special characters (-._).

* `image_url` - (Required, String, ForceNew) Specifies the URL of the external image file in the OBS bucket, the format
  is **OBS bucket name:image file name**. The storage category for OBS bucket and image file must be OBS standard
  storage. Changing this parameter will create a new resource.

* `min_disk` - (Required, Int, ForceNew) Specifies the minimum size of the system disk, in GB unit. The value ranges
  from `40` to `2,048`. Changing this parameter will create a new resource.

* `os_type` - (Optional, String, ForceNew) Specifies the operating system type of the image. The value can be
  **Windows** or **Linux**. Defaults to **Linux**. Changing this parameter will create a new resource.

* `description` - (Optional, String) Specifies the description of the image.

* `cmk_id` - (Optional, String, ForceNew) Specifies the custom key for creating encrypted image.
  Changing this parameter will create a new resource.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the image.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID to which the IMS image belongs.
  For enterprise users, if omitted, default enterprise project will be used.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the image.

* `status` - The status of the image. The value can be **active**, **queued**, **saving**, **deleted**, or **killed*,
  only image with a status of **active** can be used.

* `visibility` - Whether the image is visible to other tenants.

* `image_size` - The size of the image file, in bytes unit.

* `disk_format` - The image format. The value can be **zvhd2**, **vhd**, **zvhd**, **raw**, **qcow2**, or **iso**.

* `data_origin` - The image source. The format is **file,image_url**.

* `active_at` - The time when the image status changes to active, in RFC3339 format.

* `created_at` - The creation time of the image, in RFC3339 format.

* `updated_at` - The last update time of the image, in RFC3339 format.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 20 minutes.
* `delete` - Default is 10 minutes.

## Import

The IMS OBS data image resource can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_ims_obs_data_image.test <id>
```
