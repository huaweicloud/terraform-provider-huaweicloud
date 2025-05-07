---
subcategory: "Video on Demand (VOD)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vod_media_asset"
description: ""
---

# huaweicloud_vod_media_asset

Manages a VOD media asset resource within HuaweiCloud.

## Example Usage

### Upload media asset from OBS

```hcl
variable "bucket_name" {}
variable "object_path" {}

resource "huaweicloud_vod_media_asset" "test" {
  name         = "test"
  media_type   = "MP4"
  input_bucket = var.bucket_name
  input_path   = var.object_path
  description  = "test video"
  labels       = "test_label_1,test_lable_2,test_label_3"

  thumbnail {
    type = "time"
    time = 1
  }
}
```

### Upload media asset by URL

```hcl
variable "media_url" {}

resource "huaweicloud_vod_media_asset" "test" {
  name        = "test"
  media_type  = "MP4"
  url         = var.media_url
  description = "test video"
  labels      = "test_label_1,test_lable_2,test_label_3"

  thumbnail {
    type = "time"
    time = 1
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource. If omitted, the
  provider-level region will be used. Changing this creates a new resource.

* `name` - (Required, String) Specifies the media asset name, which contains a maximum of `128` characters.

* `media_type` - (Required, String, ForceNew) Specifies the media type. Valid values are: **MP4**, **TS**, **MOV**,
  **MXF**, **MPG**, **FLV**, **WMV**, **AVI**, **M4V**, **F4V**, **MPEG**, **3GP**, **ASF**, **MKV**, **HLS**,
  **M3U8**, **MP3**, **OGG**, **WAV**, **WMA**, **APE**, **FLAC**, **AAC**, **AC3**, **MMF**, **AMR**, **M4A**,
  **M4R**, **WV**, **MP2**. Changing this creates a new resource.

  -> When `media_type` is set to **HLS**, `storage_mode` must be set to **1** (user bucket), and the `output_path`
  should be set to the same as the `input_path`.

* `description` - (Optional, String) Specifies the media asset description, which contains a maximum of `1,024`
  characters.

* `url` - (Optional, String, ForceNew) Specifies the URL of media source file. Currently only http and https protocols
  are supported. Either this field or `input_bucket` must be specified. Changing this creates a new resource.

* `input_bucket` - (Optional, String, ForceNew) Specifies the OBS bucket name of media source file.
  Either this field or `url` must be specified. Changing this creates a new resource.

* `input_path` - (Optional, String, ForceNew) Specifies the media source file path in the OBS bucket.
  Changing this creates a new resource.

* `storage_mode` - (Optional, Int, ForceNew) Specifies the storage mode. The value can be:
  + **0**: copy the media file to VOD bucket.
  + **1**: save the media file in user bucket.

  Defaults to `0`. Changing this creates a new resource.

* `output_bucket` - (Optional, String, ForceNew) Specifies the output OBS bucket name.
  Changing this creates a new resource.

* `output_path` - (Optional, String, ForceNew) Specifies the output file path in the OBS bucket.
  Changing this creates a new resource.

-> `output_bucket` and `output_path` must be specified when `storage_mode` is set to `1`.

* `category_id` - (Optional, Int) Specifies the category ID of the media asset. Defaults to `-1`, which means the media
  asset will be categorized into the 'Other' category of system presets.

* `labels` - (Optional, String) Specifies the labels of the media asset, which contains a maximum of 16 labels
  separated by commas.

* `template_group_name` - (Optional, String, ForceNew) Specifies the transcoding template group name. If not empty,
  the uploaded media will be transcoded with the specified transcoding template group. Changing this creates a new resource.

* `workflow_name` - (Optional, String, ForceNew) Specifies the workflow name. If not empty, the uploaded media will be
  processed with the specified workflow. Changing this creates a new resource.

* `publish` - (Optional, Bool) Specifies whether to publish the media. Defaults to: **false**.

* `auto_encrypt` - (Optional, Bool, ForceNew) Specifies whether to automatically encrypt the media. If set to **true**,
  `template_group_name` must be specified, and the output format of transcoding is **HLS**. Defaults to: **false**.
  Changing this creates a new resource.

* `auto_preload` - (Optional, Bool, ForceNew) Specifies whether to automatically warm up the media to CDN. Defaults to: **false**.
  Changing this creates a new resource.

* `review_template_id` - (Optional, String, ForceNew) Specifies the review template ID. Changing this creates a new resource.

* `thumbnail` - (Optional, List, ForceNew) Specifies the review thumbnail configurations.
  The [object](#thumbnail_object) structure is documented below. Changing this creates a new resource.

<a name="thumbnail_object"></a>
The `thumbnail` block supports:

* `type` - (Required, String, ForceNew) Specifies the screenshot type. Valid values are: **time** and **dots**.
  Changing this creates a new resource.

* `time` - (Optional, Int, ForceNew) Specifies the screenshot time interval (unit: second). The value range is `1` to `12`.
  Required when `type` is **time**. Changing this creates a new resource.

* `dots` - (Optional, List, ForceNew) Specifies an array of time points of screenshot. Required when `type` is **dots**.
  Changing this creates a new resource.

* `cover_position` - (Optional, Int, ForceNew) Specifies the number of screenshots as the cover. Defaults to `1`.
  Changing this creates a new resource.

* `format` - (Optional, Int, ForceNew) Specifies the screenshot file format. Currently, only `1` (jpg) is supported.
  Defaults to: `1`. Changing this creates a new resource.

* `aspect_ratio` - (Optional, Int, ForceNew) Specifies the screenshot aspect ratio. The value can be:
  + **0**: adaptive (maintain the original aspect ratio).
  + **1**: 16:9.
  
  Defaults to `1`. Changing this creates a new resource.

* `max_length` - (Optional, Int, ForceNew) Specifies the size of the longest side of the screenshot. Unit: pixel.
  The width dimension is calculated by scaling the dimension proportional to the original video pixels.
  Changing this creates a new resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the media asset.

* `media_name` - The name of the media file.

* `media_url` - The URL of original media file.

* `category_name` - The category name of the media asset.

## Timeouts

This resource provides the following timeouts configuration options:

* `delete` - Defaults to `60` seconds.

## Import

The media asset can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_vod_media_asset.test 8754976729b8a2ba745d01036edded2b
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `url`, `input_bucket`,
`input_path`, `output_bucket`, `output_path`, `storage_mode`, `template_group_name`, `workflow_name`, `publish`,
`auto_encrypt`, `auto_preload`, `review_template_id`, `thumbnail`.
It is generally recommended running `terraform plan` after importing a media asset.
You can then decide if changes should be applied to the media asset, or the resource
definition should be updated to align with the media asset. Also you can ignore changes as below.

```hcl
resource "huaweicloud_vod_media_asset" "test" {
    ...
  lifecycle {
    ignore_changes = [
      url, input_bucket, input_path, output_bucket, output_path, storage_mode, template_group_name,
      workflow_name, publish, auto_encrypt, auto_preload, review_template_id, thumbnail,
    ]
  }
}
```
