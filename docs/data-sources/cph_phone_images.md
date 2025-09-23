---
subcategory: "Cloud Phone (CPH)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cph_phone_images"
description: |-
  Use this data source to get available images of CPH phone.
---

# huaweicloud_cph_phone_images

Use this data source to get available images of CPH phone.

## Example Usage

```hcl
data "huaweicloud_cph_phone_images" "images" {
  is_public = 1
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the data source.
  If omitted, the provider-level region will be used.

* `is_public` - (Optional, Int) The image type.  
  The options are as follows:
  + **1**: Public image.
  + **2**: Private image.

* `image_label` - (Optional, String) The label of image.  
  The valid values are **cloud_phone**, **cloud_game**, **qemu_phone**, **cloud_phone_1620**, and **cloud_game_1620**.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `images` - The list of images detail.
  The [Images](#phoneImages_Images) structure is documented below.

<a name="phoneImages_Images"></a>
The `Images` block supports:

* `id` - The ID of the image.

* `name` - The name of the image.

* `os_type` - The os type of the image.

* `os_name` - The os name of the image.

* `is_public` - The image type.  
  The options are as follows:
  + **1**: Public image.
  + **2**: Private image.

* `image_label` - The label of the image.  
  The valid values are **cloud_phone**, **cloud_game**, **qemu_phone**, **cloud_phone_1620**, and **cloud_game_1620**.
