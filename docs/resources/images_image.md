---
subcategory: "Deprecated"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_images_image"
description: |-
  Manages an IMS image resource within HuaweiCloud.
---

# huaweicloud_images_image

!> **WARNING:** It has been deprecated, please select the corresponding resource replacement based on the image type and
creation method, please use resources named in `huaweicloud_ims_xxx_xxx_image` format instead.

Manages an IMS image resource within HuaweiCloud.

## Example Usage

### Creating a system image from an existing ECS instance

```hcl
variable "name" {}
variable "instance_id" {}

resource "huaweicloud_images_image" "test" {
  name        = var.name
  instance_id = var.instance_id

  tags = {
    foo = "bar"
    key = "value"
  }
}
```

### Creating a system image from OBS bucket

```hcl
variable "name" {}
variable "image_url" {}
variable "min_disk" {}

resource "huaweicloud_images_image" "test" {
  name      = var.name
  image_url = var.image_url
  min_disk  = var.min_disk
}
```

### Creating a whole image from an existing ECS instance

```hcl
variable "name" {}
variable "instance_id" {}
variable "vault_id" {}

resource "huaweicloud_images_image" "test" {
  name        = var.name
  instance_id = var.instance_id
  vault_id    = var.vault_id
}
```

### Creating a whole image from CBR backup

```hcl
variable "name" {}
variable "backup_id" {}

resource "huaweicloud_images_image" "test" {
  name      = var.name
  backup_id = var.backup_id
}
```

### Creating a data image from the data disk bound to the ECS instance

```hcl
variable "name" {}
variable "volume_id" {}

resource "huaweicloud_images_image" "test" {
  name      = var.name
  volume_id = var.volume_id
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the name of the image.

* `instance_id` - (Optional, String, ForceNew) Specifies the ID of the ECS that needs to be converted into an image.
  This parameter is valid and mandatory only when you create a private system image or a private whole image from an
  ECS instance. Changing this parameter will create a new resource.
  + If the value of `vault_id` is empty, then a private system image will be created.
  + If the value of `vault_id` is not empty, then a private whole image will be created.

-> Exactly one of `instance_id`, `backup_id`, `volume_id` or `image_url` must be set.

* `vault_id` - (Optional, String, ForceNew) Specifies the ID of the vault to which an ECS instance is to be added or has
  been added. This parameter can only be used with `instance_id`. Changing this parameter will create a new resource.

* `backup_id` - (Optional, String, ForceNew) Specifies the ID of the CBR backup that needs to be converted into an
  image. This parameter is valid and mandatory only when you create a private whole image from a CBR backup.
  Changing this parameter will create a new resource.

* `volume_id` - (Optional, String, ForceNew) Specifies the ID of the data disk. This parameter is valid and mandatory
  when you create a private data image from an ECS instance, and the data disk must be bound to the ECS instance.
  Changing this parameter will create a new resource.

* `image_url` - (Optional, String, ForceNew) Specifies the URL of the external image file in the OBS bucket, the format
  is **OBS bucket name:Image file name**, e.g. **obs_bucket_name:image_test.vhd**. The storage category for OBS bucket
  and image file must be OBS standard storage. This parameter is valid and mandatory when you create a private system
  image from an external file uploaded to an OBS bucket, and this parameter can only be used with `min_disk`.
  Changing this parameter will create a new resource.

* `min_disk` - (Optional, Int, ForceNew) Specifies the minimum size of the system disk in the unit of GB. This parameter
  is valid and mandatory when you create a private system image from an external file uploaded to an OBS bucket.
  Changing this parameter will create a new resource.
  + When the operating system is Linux, the value ranges from `10` to `1,024`.
  + When the operating system is Windows, the value ranges from `20` to `1,024`.

* `os_version` - (Optional, String, ForceNew) Specifies the OS version.
  Changing this parameter will create a new resource.
  For its values, see [API docs](https://support.huaweicloud.com/intl/en-us/api-ims/ims_03_0910.html).

* `is_config` - (Optional, Bool, ForceNew) Specifies whether to automatically configure. If automatic backend
  configuration is required, set the value to **true**, Otherwise, set it to **false**. The default value is **false**.
  Changing this parameter will create a new resource.

* `cmk_id` - (Optional, String, ForceNew) Specifies the master key used for encrypting an image.
  Changing this parameter will create a new resource.

* `type` - (Optional, String, ForceNew) Specifies the image type. The value can be **ECS**, **FusionCompute**, **BMS**,
  or **Ironic**. Changing this parameter will create a new resource.

-> The `os_version`, `is_config`, `cmk_id`, and `type` parameters are valid only when you create a private system image
   from an external file uploaded to an OBS bucket.

* `description` - (Optional, String) Specifies the description of the image.

* `max_ram` - (Optional, Int) Specifies the maximum memory of the image in the unit of MB.

* `min_ram` - (Optional, Int) Specifies the minimum memory of the image in the unit of MB. The default value is `0`,
  indicating that the memory is not restricted.

-> When creating a private data image using `volume_id`, the `min_ram` and `max_ram` parameters do not take effect,
   please ignore them when creating. You can update them after the image is successfully created.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the image.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project ID to which the IMS image
  belongs. For enterprise users, if omitted, default enterprise project will be used.
  Changing this parameter will create a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - A unique ID assigned by IMS.

* `visibility` - Whether the image is visible to other tenants.

* `data_origin` - The image resource. The pattern can be **server_backup,backup_id**, **instance,instance_id**,
  **file,image_url**, or **volume,volume_id**.

* `disk_format` - The image file format. The value can be **vhd**, **zvhd**, **raw**, **zvhd2**, or **qcow2**.

* `image_size` - The size(bytes) of the image file format.

* `checksum` - The checksum of the data associated with the image.

* `status` - The status of the image.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 20 minutes.
* `delete` - Default is 10 minutes.

## Import

Image can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_images_image.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response. The missing attributes include: `vault_id`. It is generally recommended running `terraform plan` after
importing the image. You can then decide if changes should be applied to the image, or the resource
definition should be updated to align with the image. Also, you can ignore changes as below.

```hcl
resource "huaweicloud_images_image" "test" {
  ...

  lifecycle {
    ignore_changes = [
      vault_id,
    ]
  }
}
```
