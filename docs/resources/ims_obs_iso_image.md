---
subcategory: "Image Management Service (IMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ims_obs_iso_image"
description: |-
  Manages an IMS ISO image resource created from external image file in the OBS bucket within HuaweiCloud.
---

# huaweicloud_ims_obs_iso_image

Manages an IMS ISO image resource created from external image file in the OBS bucket within HuaweiCloud.

## Example Usage

### Creating an IMS ISO image from an external image file in the OBS bucket

```hcl
variable "name" {}
variable "image_url" {}
variable "min_disk" {}
variable "os_version" {}

resource "huaweicloud_ims_obs_iso_image" "test" {
  name       = var.name
  image_url  = var.image_url
  min_disk   = var.min_disk
  os_version = var.os_version
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
  storage, and the image file must be ISO format. Changing this parameter will create a new resource.

* `min_disk` - (Required, Int, ForceNew) Specifies the minimum size of the system disk, in GB unit.
  Changing this parameter will create a new resource.
  + When the operating system is Linux, the value ranges from `10` to `1,024`.
  + When the operating system is Windows, the value ranges from `20` to `1,024`.

* `os_version` - (Required, String, ForceNew) Specifies the operating system version of the image.
  Changing this parameter will create a new resource.
  For its values, see [API docs](https://support.huaweicloud.com/intl/en-us/api-ims/ims_03_0910.html).

* `description` - (Optional, String) Specifies the description of the image.

* `is_config` - (Optional, Bool, ForceNew) Specifies whether to automatically configure. The value can be **true** or
  **false**. Defaults to **false**. Changing this parameter will create a new resource.
  About the content of automatic backend configuration, please refer to
  [API docs](https://support.huaweicloud.com/intl/en-us/ims_faq/ims_faq_0020.html).

* `cmk_id` - (Optional, String, ForceNew) Specifies the custom key for creating encrypted image.
  Changing this parameter will create a new resource.

* `architecture` - (Optional, String, ForceNew) Specifies the schema type of the image. The value can be **x86** or
  **arm**. Defaults to **x86**. Changing this parameter will create a new resource.

* `max_ram` - (Optional, Int) Specifies the maximum memory of the image, in MB unit.

* `min_ram` - (Optional, Int) Specifies the minimum memory of the image, in MB unit.

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

* `os_type` - The operating system type of the image. The value can be **Windows** or **Linux**.

* `disk_format` - The image format. The valid value is **iso**.

* `data_origin` - The image source. The format is **file,image_url**.

* `active_at` - The time when the image status changes to active, in RFC3339 format.

* `created_at` - The creation time of the image, in RFC3339 format.

* `updated_at` - The last update time of the image, in RFC3339 format.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 20 minutes.
* `delete` - Default is 10 minutes.

## Import

The IMS OBS ISO image resource can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_ims_obs_iso_image.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `is_config`.
It is generally recommended running `terraform plan` after importing the resource. You can then decide if changes should
be applied to the resource, or the resource definition should be updated to align with the image. Also, you can ignore
changes as below.

```
resource "huaweicloud_ims_obs_iso_image" "test" {
  ...

  lifecycle {
    ignore_changes = [
      is_config,
    ]
  }
}
```
