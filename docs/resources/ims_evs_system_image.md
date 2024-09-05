---
subcategory: "Image Management Service (IMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ims_evs_system_image"
description: |-
  Manages an IMS system image resource created from EVS volume within HuaweiCloud.
---

# huaweicloud_ims_evs_system_image

Manages an IMS system image resource created from EVS volume within HuaweiCloud.

## Example Usage

### Creating an IMS system image from EVS volume

```hcl
variable "name" {}
variable "volume_id" {}
variable "os_version" {}

resource "huaweicloud_ims_evs_system_image" "test" {
  name       = var.name
  volume_id  = var.volume_id
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

* `volume_id` - (Required, String, ForceNew) Specifies the EVS volume ID used to create the image. When creating an ECS
  system image, the volume must be bound to the ECS instance. Changing this parameter will create a new resource.

* `os_version` - (Required, String, ForceNew) Specifies the operating system version of the image. Changing this
  parameter will create a new resource.
  For its values, see [API docs](https://support.huaweicloud.com/intl/en-us/api-ims/ims_03_0910.html).

  -> 1. Please strictly follow the requirements in the API docs to fill in this value, otherwise it may cause the system
     image created to be unusable.<br/>2. During the process of creating a system image, if the system can obtain the
     operating system in the volume, the operating system version in the volume will be used as the standard. At this
     time, set `os_version` is invalid. If the system is unable to retrieve the operating system from the volume, the
     system will use the input `os_version` value as the reference.

* `type` - (Optional, String, ForceNew) Specifies the image type. The value can be **ECS**, **FusionCompute**, **BMS**,
  or **Ironic**. Defaults to **ECS**. Changing this parameter will create a new resource.
  + Set to **ECS** or **FusionCompute** represent the creation of ECS server image.
  + Set to **BMS** or **Ironic** represent the creation of BMS server image.

* `description` - (Optional, String) Specifies the description of the image.

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

* `os_type` - The operating system type of the image. The value can be **Linux** or **Windows**.

* `min_disk` - The minimum disk space required to run an image, in GB unit.

* `disk_format` - The image format. The value can be **zvhd2**, **vhd**, **zvhd**, **raw**, **qcow2**, or **iso**.

* `data_origin` - The image source. The format is **volume,volume_id**.

* `active_at` - The time when the image status changes to active, in RFC3339 format.

* `created_at` - The creation time of the image, in RFC3339 format.

* `updated_at` - The last update time of the image, in RFC3339 format.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 20 minutes.
* `delete` - Default is 10 minutes.

## Import

The IMS EVS system image resource can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_ims_evs_system_image.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `type`.
It is generally recommended running `terraform plan` after importing the resource. You can then decide if changes should
be applied to the resource, or the resource definition should be updated to align with the image. Also, you can ignore
changes as below.

```
resource "huaweicloud_ims_evs_system_image" "test" {
    ...

  lifecycle {
    ignore_changes = [
      type,
    ]
  }
}
```
