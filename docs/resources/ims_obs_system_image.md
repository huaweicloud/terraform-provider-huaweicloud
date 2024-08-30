---
subcategory: "Image Management Service (IMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ims_obs_system_image"
description: |-
  Manages an IMS system image resource created from external image file in the OBS bucket within HuaweiCloud.
---

# huaweicloud_ims_obs_system_image

Manages an IMS system image resource created from external image file in the OBS bucket within HuaweiCloud.

## Example Usage

### Creating an IMS system image from an external image file in the OBS bucket

```hcl
variable "name" {}
variable "image_url" {}
variable "min_disk" {}

resource "huaweicloud_ims_obs_system_image" "test" {
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
  is **OBS bucket name:image file name**. Changing this parameter will create a new resource.

* `min_disk` - (Required, Int, ForceNew) Specifies the minimum size of the system disk, in GB unit.
  Changing this parameter will create a new resource.
  + When the operating system is Linux, the value ranges from `10` to `1,024`.
  + When the operating system is Windows, the value ranges from `20` to `1,024`.

* `description` - (Optional, String) Specifies the description of the image.

* `type` - (Optional, String, ForceNew) Specifies the image type. The value can be **ECS**, **FusionCompute**, **BMS**,
  or **Ironic**. Defaults to **ECS**. Changing this parameter will create a new resource.
  + Set to **ECS** or **FusionCompute** represent the creation of ECS server image.
  + Set to **BMS** or **Ironic** represent the creation of BMS server image.

* `is_config` - (Optional, Bool, ForceNew) Specifies whether to automatically configure. The value can be **true** or
  **false**. Defaults to **false**. Changing this parameter will create a new resource.
  About the content of automatic backend configuration, please refer to
  [API docs](https://support.huaweicloud.com/intl/en-us/ims_faq/ims_faq_0020.html).

* `cmk_id` - (Optional, String, ForceNew) Specifies the custom key for creating encrypted image.
  Changing this parameter will create a new resource.

* `is_quick_import` - (Optional, Bool, ForceNew) Specifies whether to use the image file quick import method to create
  an image. The value can be **true** or **false**. Changing this parameter will create a new resource.
  For constraints and limitations on fast import of image files,
  please refer to [API docs](https://support.huaweicloud.com/intl/en-us/api-ims/ims_03_0605.html).

-> 1. When the `is_quick_import` set to **true**, IMS will not parse the specified external image file, so the
  `os_type`, `os_version`, and `architecture` parameters is based on the specified value.
  <br/>2. When ignoring the `is_quick_import` or set to **false** , IMS will parse the external image file and confirm
  the `os_type`, `os_version`, and `architecture` of the image, if parsing fails, the specified value shall prevail.

* `os_type` - (Optional, String, ForceNew) Specifies the operating system type of the image. The value can be
  **Windows** or **Linux**. Changing this parameter will create a new resource.

* `os_version` - (Optional, String, ForceNew) Specifies the operating system version of the image. This field is
  required when `is_quick_import` set to **true**. Changing this parameter will create a new resource.
  For its values, see [API docs](https://support.huaweicloud.com/intl/en-us/api-ims/ims_03_0910.html).

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

* `disk_format` - The image format. The value can be **zvhd2**, **vhd**, **zvhd**, **raw**, or **qcow2**.

* `data_origin` - The image source. The format is **file,image_url**.

* `active_at` - The time when the image status changes to active, in RFC3339 format.

* `created_at` - The creation time of the image, in RFC3339 format.

* `updated_at` - The last update time of the image, in RFC3339 format.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 20 minutes.
* `delete` - Default is 10 minutes.

## Import

The IMS OBS system image resource can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_ims_obs_system_image.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `type`, `is_config`, `is_quick_import`.
It is generally recommended running `terraform plan` after importing the resource. You can then decide if changes should
be applied to the resource, or the resource definition should be updated to align with the image. Also, you can ignore
changes as below.

```
resource "huaweicloud_ims_obs_system_image" "test" {
    ...

  lifecycle {
    ignore_changes = [
      type, is_config, is_quick_import,
    ]
  }
}
```
