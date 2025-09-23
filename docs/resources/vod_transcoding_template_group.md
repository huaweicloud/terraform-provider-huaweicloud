---
subcategory: "Video on Demand (VOD)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_vod_transcoding_template_group"
description: ""
---

# huaweicloud_vod_transcoding_template_group

Manages a VOD transcoding template group resource within HuaweiCloud.

## Example Usage

```hcl
resource "huaweicloud_vod_transcoding_template_group" "test" {
  name                 = "test"
  description          = "test group"
  audio_codec          = "AAC"
  hls_segment_duration = 5
  low_bitrate_hd       = false
  video_codec          = "H264"

  quality_info {
    output_format = "HLS"

    audio {
      bitrate     = 0
      channels    = 2
      sample_rate = 1
    }

    video {
      bitrate    = 1000
      frame_rate = 1
      height     = 720
      quality    = "HD"
      width      = 1280
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource. If omitted, the
  provider-level region will be used. Changing this creates a new resource.

* `name` - (Required, String) Specifies the name of the template group. The value can be a string of `1` to `128`
  characters that can consist of letters, digits and underscores (_).

* `quality_info` - (Required, List) Specifies the quality info list of the template group.
  The [object](#quality_info_object) structure is documented below.

* `low_bitrate_hd` - (Optional, Bool) Specifies whether to enable low bitrate HD. Defaults to: **false**.

* `audio_codec` - (Optional, String) Specifies the audio codec. The value can be: **AAC** and **HEAAC1**.
  Defaults to: **AAC**.

* `video_codec` - (Optional, String) Specifies the video codec. The value can be: **H264** and **H265**.
  Defaults to: **H264**.

* `hls_segment_duration` - (Optional, Int) Specifies the HLS segment duration. The value can be: `2`, `3`, `5`
  and `10`. Defaults to: `5`. This parameter is used only when `output_format` is set to **HLS** or **DASH_HLS**.

* `description` - (Optional, String) Specifies the description of the template group.

* `auto_encrypt` - (Optional, Bool) Specifies whether to automatically encrypt. Before enabling, you need to configure
  the HLS encryption key URL. When `auto_encrypt` is **true**, the `output_format` must be **HLS**.
  Defaults to: **false**.

* `is_default` - (Optional, Bool) Specifies whether to use this group as default group. Defaults to: **false**.

* `watermark_template_ids` - (Optional, List) Specifies the list of the watermark template IDs.

<a name="quality_info_object"></a>
The `quality_info` block supports:

* `output_format` - (Required, String) Specifies the output format. The value can be: **HLS**, **MP4**, **DASH**,
  **DASH_HLS**, **MP3** and **ADTS**.

* `video` - (Optional, List) Specifies the video configurations.
  The [object](#video_object) structure is documented below.

* `audio` - (Optional, List) Specifies the audio configurations.
  The [object](#audio_object) structure is documented below.

<a name="video_object"></a>
The `video` block supports:

* `quality` - (Required, String) Specifies the video quality.
  The value can be: **4K**, **2K**, **FHD**, **SD**, **LD** and **HD**.

* `width` - (Optional, Int) Specifies the video width. The value can be `0` or range from `128` to `3,840`.
  Defaults to `0`. If set to `0`, the system will automatically adjust the `width` according to the `height`.

* `height` - (Optional, Int) Specifies the video height. The value can be `0` or range from `128` to `2,160`.
  Defaults to `0`. If set to `0`, the system will automatically adjust the `height` according to the `width`.

  -> If the quality of the original file is **2K** or **4K**, and the `video_codec` is **H264**, the `width` and
  `height` of the cannot be set to `0`, otherwise the transcoding will fail.

* `bitrate` - (Optional, Int) Specifies the video bitrate. The value can be `0` or range from `700` to `3,000`.
  Defaults to `0`. If set to `0`, the output video will be produced at the recommended bitrate.

* `frame_rate` - (Optional, Int) Specifies the video frame rate. The value ranges from `1` to `75`.
  Defaults to `1`. If set to `1`, the frame rate of the transcoded video is the same as that of the untransocded video.

<a name="audio_object"></a>
The `audio` block supports:

* `sample_rate` - (Required, Int) Specifies the audio sample rate. The value can be:
  + **1**: AUTO
  + **2**: 22,050 Hz
  + **3**: 32,000 Hz
  + **4**: 44,100 Hz
  + **5**: 48,000 Hz
  + **6**: 96,000 Hz
  
* `channels` - (Required, Int) Specifies the audio channels. The value can be:
  + **1**: Mono
  + **2**: Stereo

* `bitrate` - (Optional, Int) Specifies the audio bitrate. The value can be `0` or range from
  `8` to `1,000`. Defaults to `0`. If set to `0`, the output audio will be produced at the recommended bitrate.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Indicates the transcoding template group ID.

* `type` - Indicates the type of the template group.

## Import

VOD transcoding template groups can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_vod_transcoding_template_group.test 589e49809bb84447a759f6fa9aa19949
```
