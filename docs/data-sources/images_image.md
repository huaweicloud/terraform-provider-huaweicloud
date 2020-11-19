---
subcategory: "Image Management Service (IMS)"
---

# huaweicloud\_images\_image

Use this data source to get the ID of an available HuaweiCloud image.
This is an alternative to `huaweicloud_images_image_v2`

## Example Usage

```hcl
data "huaweicloud_images_image" "ubuntu" {
  name        = "Ubuntu 18.04 server 64bit"
  visibility  = "public"
  most_recent = true
}
```

## Argument Reference

* `region` - (Optional) The region in which to obtain the images. If omitted, the provider-level region will be used.

* `most_recent` - (Optional) If more than one result is returned, use the most
  recent image.

* `name` - (Optional) The name of the image.

* `owner` - (Optional) The owner (UUID) of the image.

* `size_min` - (Optional) The minimum size (in bytes) of the image to return.

* `size_max` - (Optional) The maximum size (in bytes) of the image to return.

* `sort_direction` - (Optional) Order the results in either `asc` or `desc`.

* `sort_key` - (Optional) Sort images based on a certain key. Must be one of
   "name", "container_format", "disk_format", "status", "id" or "size".
   Defaults to `name`.

* `tag` - (Optional) Search for images with a specific tag.

* `visibility` - (Optional) The visibility of the image. Must be one of
   "public", "private", "community", or "shared". Defaults to `private`.


## Attributes Reference

`id` is set to the ID of the found image. In addition, the following attributes
are exported:

* `checksum` - The checksum of the data associated with the image.
* `created_at` - The date the image was created.
* `container_format`: The format of the image's container.
* `disk_format`: The format of the image's disk.
* `file` - the trailing path after the glance endpoint that represent the
   location of the image or the path to retrieve it.
* `metadata` - The metadata associated with the image.
   Image metadata allow for meaningfully define the image properties and tags.
* `min_disk_gb` - The minimum amount of disk space required to use the image.
* `min_ram_mb` - The minimum amount of ram required to use the image.
* `properties` - Freeform information about the image.
* `protected` - Whether or not the image is protected.
* `schema` - The path to the JSON-schema that represent
   the image or image
* `size_bytes` - The size of the image (in bytes).
* `tags` - See Argument Reference above.
* `update_at` - The date the image was last updated.
