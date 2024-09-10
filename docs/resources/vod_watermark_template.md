---
subcategory: "Video on Demand (VOD)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vod_watermark_template"
description: ""
---

# huaweicloud_vod_watermark_template

Manages a VOD watermark template resource within HuaweiCloud.

## Example Usage

```hcl
resource "huaweicloud_vod_watermark_template" "test" {
  name              = "test"
  image_file        = "./test.PNG"
  image_type        = "PNG"
  position          = "TOPLEFT"
  image_process     = "ORIGINAL"
  horizontal_offset = "0.05"
  vertical_offset   = "0.05"
  width             = "0.1"
  height            = "0.1"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource. If omitted, the
  provider-level region will be used. Changing this creates a new resource.

* `name` - (Required, String) Specifies the template name, which contains a maximum of `128` characters.

* `image_file` - (Required, String, ForceNew) Specifies the image file name, e.g. './test.png'.
  Changing this creates a new resource.

* `image_type` - (Required, String, ForceNew) Specifies the image file type. The valid values are **PNG**, **JPG**
  and **JPEG**. Changing this creates a new resource.

* `image_process` - (Optional, String) Specifies the image process. The valid values are:  
  + **TRANSPARENT**: make the background color transparent.
  + **ORIGINAL**: only simple scaling, no other processing.
  + **GRAYED**: make the color image grayed.

  Defaults to: **TRANSPARENT**.

* `horizontal_offset` - (Optional, String) Specifies horizontal offset ratio of the watermark image relative to the
  output video. The value range is [0, 1). It supports 4 decimal places, e.g. 0.9999, the excess will be
  automatically discarded. Defaults to: **0**.

* `vertical_offset` - (Optional, String) Specifies vertical offset ratio of the watermark image relative to the
  output video. The value range is [0, 1). It supports 4 decimal places, e.g. 0.9999, the excess will be
  automatically discarded. Defaults to: **0**.

* `position` - (Optional, String) Specifies the location of the watermark. The valid values are **TOPRIGHT**,
  **TOPLEFT**, **BOTTOMRIGHT** and **BOTTOMLEFT**. Defaults to: **TOPRIGHT**.

* `width` - (Optional, String) Specifies width ratio of the watermark image relative to the output video.
  The value range is (0, 1). It supports 4 decimal places, e.g. 0.9999, the excess will be
  automatically discarded. Defaults to: **0.01**.

* `height` - (Optional, String) Specifies height ratio of the watermark image relative to the output video.
  The value range is (0, 1). It supports 4 decimal places, e.g. 0.9999, the excess will be
  automatically discarded. Defaults to: **0.01**.

* `timeline_start` - (Optional, String) Specifies the watermark start time (Unit: second). The value is a digit
  greater than or equal to **0**. Defaults to: **0**.

* `timeline_duration` - (Optional, String) Specifies the watermark duration (Unit: second). The value is a digit
  greater than or equal to **0**. By default, the watermark lasts until the end of the video.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the template.

* `watermark_type` - The watermark type.

* `image_url` - The watermark image URL.

## Timeouts

This resource provides the following timeouts configuration options:

* `create` - Default is 5 minutes.

## Import

The template can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_vod_watermark_template.test 81ac58796e25842ee2e90a904aa8a719
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `image_file`.
It is generally recommended running `terraform plan` after importing a watermark template.
You can then decide if changes should be applied to the watermark template, or the resource
definition should be updated to align with the watermark template. Also you can ignore changes as below.

```hcl
resource "huaweicloud_vod_watermark_template" "test" {
    ...
  lifecycle {
    ignore_changes = [
      image_file,
    ]
  }
}
```
