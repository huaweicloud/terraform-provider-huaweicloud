---
subcategory: "Image Management Service (IMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_images_image_copy"
description: |-
  Manages an IMS image copy resource within HuaweiCloud.
---

# huaweicloud_images_image_copy

Manages an IMS image copy resource within HuaweiCloud.

## Example Usage

### Copy image within region

```hcl
variable "source_image_id" {}
variable "name" {}
variable "kms_key_id" {}

resource "huaweicloud_images_image_copy" "test" {
  source_image_id = var.source_image_id
  name            = var.name
  kms_key_id      = var.kms_key_id
}
```

### Copy image cross region

```hcl
variable "source_image_id" {}
variable "name" {}
variable "target_region" {}
variable "agency_name" {}

resource "huaweicloud_images_image_copy" "test" {
  source_image_id = var.source_image_id
  name            = var.name
  target_region   = var.target_region
  agency_name     = var.agency_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region to which the source image belongs.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `source_image_id` - (Required, String, ForceNew) Specifies the ID of the copied image.
  Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the name of the copy image. The name can contain `1` to `128` characters,
  only Chinese and English letters, digits, underscore (_), hyphens (-), dots (.) and space are allowed, but it cannot
  start or end with a space.

* `target_region` - (Optional, String, ForceNew) Specifies the target region name.
  If specified, it means cross-region replication.

* `kms_key_id` - (Optional, String, ForceNew) Specifies the master key used for encrypting an image.
  Only copying scene within a region is supported. Changing this parameter will create a new resource.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID of the image.
  Only copying scene within a region is supported, for enterprise users, if omitted, default enterprise project will
  be used.

* `agency_name` - (Optional, String, ForceNew) Specifies the agency name. It is required in the cross-region scene.
  Changing this parameter will create a new resource.

* `vault_id` - (Optional, String, ForceNew) Specifies the ID of the vault. It is used in the cross-region scene, it is
  mandatory if you are replicating a full-ECS image, and the region to which the vault belongs must be consistent with
  the value of `target_region`.
  Changing this parameter will create a new resource.

* `min_ram` - (Optional, Int) Specifies the minimum memory of the copy image in the unit of MB. The default value is
  `0`, indicating that the memory is not restricted.

* `max_ram` - (Optional, Int) Specifies the maximum memory of the copy image in the unit of MB.

* `description` - (Optional, String) Specifies the description of the copy image.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the copy image.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `instance_id` - Indicates the ID of the ECS that needs to be converted into an image.

* `os_version` - Indicates the OS version.

* `visibility` - Indicates whether the image is visible to other tenants.

* `min_disk` - The minimum disk space required to run an image, in GB unit.

* `data_origin` - Indicates the image source.
  The format is **image,region,source_image_id**, e.g. **image,cn-north-4,xxxxxx**.

* `disk_format` - Indicates the image file format.
  The value can be **zvhd2**, **vhd**, **zvhd**, **raw**, **qcow2**, or **iso**.

* `image_size` - Indicates the size(bytes) of the image file.

* `status` - Indicates the status of the image. The value can be **active**, **queued**, **saving**, **deleted**,
  or **killed*, only image with a status of **active** can be used.

* `active_at` - The time when the image status changes to active, in RFC3339 format.

* `created_at` - The creation time of the image, in RFC3339 format.

* `updated_at` - The last update time of the image, in RFC3339 format.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 60 minutes.
* `delete` - Default is 3 minutes.
