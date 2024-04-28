---
subcategory: "AI Development Platform (ModelArts)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_modelarts_notebook_images"
description: ""
---

# huaweicloud_modelarts_notebook_images

Use this data source to get a list of ModelArts notebook images.

## Example Usage

```hcl
data "huaweicloud_modelarts_notebook_images" "test" {
  type = "BUILD_IN"
}

output "image_id" {
  value = data.huaweicloud_modelarts_notebook_images.test.images[0].id
}
```

## Argument Reference

The arguments of this data source act as filters for querying the available images in the current region.
 All images that meet the filter criteria will be exported as attributes.

* `region` - (Optional, String) Specifies the region in which to obtain images. If omitted, the provider-level region
 will be used.

* `name` - (Optional, String) Specifies the name of image.

* `organization` - (Optional, String) Specifies the name of the organization (namespace) which image belongs to.

* `type` - (Optional, String) Specifies the type of image. The options are:
  + `BUILD_IN`: The system built-in image.
  + `DEDICATED`: User-saved images.

 The default value is `BUILD_IN`.

* `cpu_arch` - (Optional, String) Specifies the CPU architecture of image. The value can be **x86_64** and **aarch64**.

* `workspace_id` - (Optional, String) Specifies the workspace ID which image belongs to.

## Attribute Reference

The following attributes are exported:

* `id` - Indicates a data source ID.

* `images` - Indicates a list of all images found. Structure is documented below.

The `images` block contains:

* `id` - The ID of the image.
* `name` - The name of the image.
* `swr_path` - The path the image in HuaweiCloud SWR service (SoftWare Repository for Container).
* `type` - The type of the image.
* `cpu_arch` - The CPU architecture of the image. The value can be **x86_64** and **aarch64**.
* `description` - The description of the image.
