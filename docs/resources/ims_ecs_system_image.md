---
subcategory: "Image Management Service (IMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_ims_ecs_system_image"
description: |-
  Manages an IMS system image resource created from ECS instance within HuaweiCloud.
---

# huaweicloud_ims_ecs_system_image

Manages an IMS system image resource created from ECS instance within HuaweiCloud.

## Example Usage

### Creating an IMS system image from an existing ECS instance

```hcl
variable "name" {}
variable "instance_id" {}

resource "huaweicloud_ims_ecs_system_image" "test" {
  name        = var.name
  instance_id = var.instance_id
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

* `instance_id` - (Required, String, ForceNew) Specifies the source ECS instance ID used to create the image.
  Changing this parameter will create a new resource.

* `description` - (Optional, String) Specifies the description of the image.

* `max_ram` - (Optional, Int) Specifies the maximum memory of the image, in MB unit.
  The default value is `0`, indicating that the memory is not restricted.

* `min_ram` - (Optional, Int) Specifies the minimum memory of the image, in MB unit.
  The default value is `0`, indicating that the memory is not restricted.

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

* `min_disk` - The minimum disk space required to run an image, in GB unit.

* `disk_format` - The image format. The value can be **zvhd2**, **vhd**, **zvhd**, **raw**, or **qcow2**.

* `data_origin` - The image resource. The format is **instance,instance_id**.

* `os_version` - The operating system version of the image.

* `active_at` - The time when the image status changes to active, in RFC3339 format.

* `created_at` - The creation time of the image, in RFC3339 format.

* `updated_at` - The last update time of the image, in RFC3339 format.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 20 minutes.
* `delete` - Default is 10 minutes.

## Import

The IMS ECS system image resource can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_ims_ecs_system_image.test <id>
```
