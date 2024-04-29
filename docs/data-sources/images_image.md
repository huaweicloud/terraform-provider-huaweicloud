---
subcategory: "Image Management Service (IMS)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_images_image"
description: ""
---

# huaweicloud_images_image

Use this data source to get the ID of an available HuaweiCloud image.

## Example Usage

```hcl
data "huaweicloud_images_image" "ubuntu" {
  name        = "Ubuntu 18.04 server 64bit"
  visibility  = "public"
  most_recent = true
}

data "huaweicloud_images_image" "centos-1" {
  architecture = "x86"
  os_version   = "CentOS 7.4 64bit"
  visibility   = "public"
  most_recent  = true
}

data "huaweicloud_images_image" "centos-2" {
  architecture = "x86"
  name_regex   = "^CentOS 7.4"
  visibility   = "public"
  most_recent  = true
}

data "huaweicloud_images_image" "bms_image" {
  architecture = "x86"
  image_type   = "Ironic"
  os_version   = "CentOS 7.4 64bit"
  visibility   = "public"
  most_recent  = true
}
```

## Argument Reference

* `region` - (Optional, String) The region in which to obtain the images. If omitted, the provider-level region will be
  used.

* `most_recent` - (Optional, Bool) If more than one result is returned, use the latest updated image.

* `name` - (Optional, String) The name of the image. Cannot be used simultaneously with `name_regex`.

* `name_regex` - (Optional, String) The regular expression of the name of the image.
  Cannot be used simultaneously with `name`.

* `visibility` - (Optional, String) The visibility of the image. Must be one of
  **public**, **private**, **market** or **shared**.

* `architecture` - (Optional, String) Specifies the image architecture type. The value can be **x86** and **arm**.

* `os` - (Optional, String) Specifies the image OS type. The value can be **Windows**, **Ubuntu**,
  **RedHat**, **SUSE**, **CentOS**, **Debian**, **OpenSUSE**, **Oracle Linux**, **Fedora**, **Other**,
  **CoreOS**, or **EulerOS**.

* `os_version` - (Optional, String) Specifies the OS version. For example, *CentOS 7.4 64bit* or *Ubuntu 18.04 server 64bit*.
  For all its valid values, see [API docs](https://support.huaweicloud.com/intl/en-us/api-ims/ims_03_0910.html).

* `image_type` - (Optional, String) Specifies the environment where the image is used. For a BMS image, the value is **Ironic**.

* `owner` - (Optional, String) The owner (UUID) of the image.

* `tag` - (Optional, String) Search for images with a specific tag in "Key=Value" format.

* `sort_direction` - (Optional, String) Order the results in either `asc` or `desc`.

* `sort_key` - (Optional, String) Sort images based on a certain key. Must be one of
  "name", "container_format", "disk_format", "status", "id" or "size". Defaults to `name`.

* `enterprise_project_id` - (Optional, String) Specifies the enterprise project ID of the image.

* `flavor_id` - (Optional, String) Specifies the ECS flavor ID used to filter out available images.
  You can specify only one flavor ID and only ECS flavor ID is valid, BMS flavor is not supported.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Specifies a resource ID in UUID format.
* `checksum` - The checksum of the data associated with the image.
* `container_format` - The format of the image's container.
* `disk_format` - The format of the image's disk.
* `file` - the trailing path after the glance endpoint that represent the location of the image or the path to retrieve
  it.
* `metadata` - The metadata associated with the image. Image metadata allow for meaningfully define the image properties
  and tags.
* `min_disk_gb` - The minimum amount of disk space required to use the image.
* `min_ram_mb` - The minimum amount of ram required to use the image.
* `protected` - Whether or not the image is protected.
* `schema` - The path to the JSON-schema that represent the image or image.
* `size_bytes` - The size of the image (in bytes).
* `status` - The status of the image.
* `backup_id` - The backup ID of the whole image in the CBR vault.
* `created_at` - The date when the image was created.
* `updated_at` - The date when the image was last updated.
