---
subcategory: "Image Management Service (IMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_images_image"
description: ""
---

# huaweicloud_images_image

Manages an Image resource within HuaweiCloud IMS.

## Example Usage

### Creating an image from an existing ECS

```hcl
variable "instance_name" {}
variable "image_name" {}

data "huaweicloud_compute_instance" "test" {
  name = var.instance_name
}

resource "huaweicloud_images_image" "test" {
  name        = var.image_name
  instance_id = data.huaweicloud_compute_instance.test.id
  description = "created by Terraform"

  tags = {
    foo = "bar"
    key = "value"
  }
}
```

### Creating an image from OBS bucket

```hcl
resource "huaweicloud_images_image" "ims_test_file" {
  name        = "ims_test_file"
  image_url   = "ims-image:centos70.qcow2"
  min_disk    = 40
  description = "Create an image from the OBS bucket."

  tags = {
    foo = "bar1"
    key = "value"
  }
}
```

### Creating a whole image from an existing ECS

```hcl
variable "vault_id" {}
variable "instance_id" {}

resource "huaweicloud_images_image" "test" {
  name        = "test_whole_image"
  instance_id = var.instance_id
  vault_id    = var.vault_id

  tags = {
    foo = "bar2"
    key = "value"
  }
}
```

### Creating a whole image from CBR backup

```hcl
variable "backup_id" {}

resource "huaweicloud_images_image" "test" {
  name      = "test_whole_image"
  backup_id = var.backup_id

  tags = {
    foo = "bar1"
    key = "value"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, String) The name of the image.

* `instance_id` - (Optional, String, ForceNew) The ID of the ECS that needs to be converted into an image. This
  parameter is mandatory when you create a private image or a private whole image from an ECS.
  If the value of `vault_id` is not empty, then a whole image will be created.

* `backup_id` - (Optional, String, ForceNew) The ID of the CBR backup that needs to be converted into an image. This
  parameter is mandatory when you create a private whole image from a CBR backup.

* `image_url` - (Optional, String, ForceNew) The URL of the external image file in the OBS bucket. This parameter is
  mandatory when you create a private image from an external file uploaded to an OBS bucket. The format is *OBS bucket
  name:Image file name*.

* `min_ram` - (Optional, Int) The minimum memory of the image in the unit of MB. The default value is 0,
  indicating that the memory is not restricted.

* `max_ram` - (Optional, Int) The maximum memory of the image in the unit of MB.

* `description` - (Optional, String) A description of the image.

* `tags` - (Optional, Map) The tags of the image.

* `min_disk` - (Optional, Int, ForceNew) The minimum size of the system disk in the unit of GB. This parameter is
  mandatory when you create a private image from an external file uploaded to an OBS bucket. The value ranges from 1 GB
  to 1024 GB.

* `os_version` - (Optional, String, ForceNew) The OS version. This parameter is valid when you create a private image
  from an external file uploaded to an OBS bucket.
  For its values, see [API docs](https://support.huaweicloud.com/intl/en-us/api-ims/ims_03_0910.html).

* `is_config` - (Optional, Bool, ForceNew) If automatic configuration is required, set the value to true. Otherwise, set
  the value to false.

* `cmk_id` - (Optional, String, ForceNew) The master key used for encrypting an image.

* `type` - (Optional, String, ForceNew) The image type. Must be one of `ECS`, `FusionCompute`, `BMS`, or `Ironic`.

* `enterprise_project_id` - (Optional, String, ForceNew) The enterprise project id of the image. Changing this creates a
  new image.

* `vault_id` - (Optional, String, ForceNew) The ID of the vault to which an ECS is to be added or has been added.
  This parameter is mandatory when you create a private whole image from an ECS.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - A unique ID assigned by IMS.

* `visibility` - Whether the image is visible to other tenants.

* `data_origin` - The image resource. The pattern can be 'instance,*instance_id*', 'file,*image_url*'
  or 'server_backup,*backup_id*'.

* `disk_format` - The image file format. The value can be `vhd`, `zvhd`, `raw`, `zvhd2`, or `qcow2`.

* `image_size` - The size(bytes) of the image file format.

* `checksum` - The checksum of the data associated with the image.

* `status` - The status of the image.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 20 minutes.
* `delete` - Default is 10 minutes.

## Import

Images can be imported using the `id`, e.g.

```bash
terraform import huaweicloud_images_image.my_image <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response. The missing attributes include: `vault_id`. It is generally recommended running `terraform plan` after
importing the image. You can then decide if changes should be applied to the image, or the resource
definition should be updated to align with the image. Also you can ignore changes as below.

```
resource "huaweicloud_images_image" "test" {
  ...

  lifecycle {
    ignore_changes = [
      vault_id,
    ]
  }
}
```
