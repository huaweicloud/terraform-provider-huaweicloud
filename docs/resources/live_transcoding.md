---
subcategory: "Live"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_live_transcoding"
description: ""
---

# huaweicloud_live_transcoding

Manages a Live transcoding within HuaweiCloud.

## Example Usage

### Create a transcoding

```hcl
variable "ingest_domain_name" {}

resource "huaweicloud_live_domain" "ingestDomain" {
  name = var.ingest_domain_name
  type = "push"
}

resource "huaweicloud_live_transcoding" "test" {
  domain_name    = huaweicloud_live_domain.ingestDomain.name
  app_name       = "live"
  video_encoding = "H264"

  templates {
    name    = "L"
    width   = 300
    height  = 400
    bitrate = 300
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create this resource. If omitted,
the provider-level region will be used. Changing this parameter will create a new resource.

* `domain_name` - (Required, String, ForceNew) Specifies the ingest domain name.
Changing this parameter will create a new resource.

* `app_name` - (Required, String, ForceNew) Specifies the application name.
Changing this parameter will create a new resource.

* `video_encoding` - (Required, String) Specifies the video codec. The valid values are **H264** and **H265**.

* `templates` - (Required, List) Specifies the video quality templates.
The [object](#templates_resource) structure is documented below. A maximum of `4` templates can be added.
For resolution and bitrate settings in the presets,
please refer to the [document](https://support.huaweicloud.com/intl/en-us/usermanual-live/live01000802.html).

* `low_bitrate_hd` - (Optional, Bool) Specifies whether to enable low bitrate HD rates. If enabled
the output media will have a lower bitrate with the same image quality. Defaults to **false**.

<a name="templates_resource"></a>
The `templates` block supports:

* `name` - (Required, String) Specifies the template name. The name can contain a maximum of 64 characters, and only
contains letters, digits and hyphens (-).

* `width` - (Required, Int) Specifies video width (unit: pixel).
  + **When the video encoding is H264**, value range: `32` ~ `3,840` and must be a multiple of `2`.
  + **When the video encoding is H265**, value range: `320` ~ `3,840` and must be a multiple of `4`.

* `height` - (Required, Int) Specifies video height (unit: pixel).
  + **When the video encoding is H264**, value range: 32 ~ 2160 and must be a multiple of `2`.
  + **When the video encoding is H265**, value range: 240 ~ 2160 and must be a multiple of `4`.

* `bitrate` - (Required, Int) Specifies the bitrate of a transcoded video, in kbit/s. Value range: `40` ~ `30,000`.

* `frame_rate` - (Optional, Int) Specifies the frame rate of the transcoded video, in fps. Value range: `0` ~ `30`.
Value 0 indicates that the frame rate remains unchanged.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in format of **domain_name/app_name**. It is composed of domain name and the application name,
separated by a slash.

## Import

Transcodings can be imported using the `domain_name` and `app_name`, separated by a slash. e.g.

```bash
$ terraform import huaweicloud_live_transcoding.test play.example.demo.com/live
```
