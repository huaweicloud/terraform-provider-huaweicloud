---
subcategory: "Live"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_live_transcodings"
description: |-
  Use this data source to get the list of the transcoding templates.
---

# huaweicloud_live_transcodings

Use this data source to get the list of the transcoding templates.

## Example Usage

```hcl
variable "domain_name" {}

data "huaweicloud_live_transcodings" "test" {
  domain_name = var.domain_name
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `domain_name` - (Required, String) Specifies the ingest domain name to which the transcoding templates blong.

* `app_name` - (Optional, String) Specifies the application name of the transcoding template.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `templates` - The list of the transcoding templates.

  The [templates](#templates_struct) structure is documented below.

<a name="templates_struct"></a>
The `templates` block supports:

* `app_name` - The application name of the transcoding template.

* `quality_info` - The video quality information.

  The [quality_info](#templates_quality_info_struct) structure is documented below.

<a name="templates_quality_info_struct"></a>
The `quality_info` block supports:

* `name` - The transcoding template name.

* `quality` - The video quality.
  The valid values are as follows:
  + **lud**: Indicates ultra high definition.
  + **lhd**: Indicates high definition.
  + **lsd**: Indicates standard definition.
  + **lld**: Indicates smooth.
  + **userdefine**: Indicates customization of video quality.

* `video_encoding` - The video encoding format.
  The value can be **H264** or **H265**.

* `width` - The video long edge (width of horizontal screen, height of vertical screen), in pixels.

* `height` - The video short edge (horizontal screen height, vertical screen width), in pixels.

* `bitrate` - The bitrate of the transcoding video, in Kbps.

* `frame_rate` - The frame rate of transcoding video, in fps.

* `protocol` - The protocol type of transcoding output.
  The value can be **RTMP**.

* `low_bitrate_hd` - Whether to enable high-definition and low bitrate.
  The value can be **on** or **off**.

* `gop` - The I frame interval, in seconds.

* `bitrate_adaptive` - The adaptive bitrate.
  The valid values are as follows:
  + **off**: Turn off rate adaptation and output the target rate at the set rate.
  + **minimum**: The target bitrate is output at the minimum value of the set bitrate and the source file bitrate.
  + **adaptive**: The target bitrate is adaptively output based on source file bitrate.

* `i_frame_interval` - The maximum I frame interval, in frame.

* `i_frame_policy` - The encoding output I frame policy.
  The valid values are as follows:
  + **auto**: The I frame output according to the set `gop` duration.
  + **strictSync**: The encoding output I frame is completely consistent with the source.
