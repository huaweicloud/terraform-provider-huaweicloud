---
subcategory: "Live"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_live_transcoding"
description: |-
  Manages a Live transcoding resource within HuaweiCloud.
---

# huaweicloud_live_transcoding

Manages a Live transcoding resource within HuaweiCloud.

## Example Usage

### Create a transcoding

```hcl
variable "ingest_domain_name" {}
variable "app_name" {}
variable "video_encoding" {}
variable "template_name" {}

resource "huaweicloud_live_transcoding" "test" {
  domain_name    = var.ingest_domain_name
  app_name       = var.app_name
  video_encoding = var.video_encoding

  templates {
    name    = var.template_name
    width   = 300
    height  = 400
    bitrate = 300
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create this resource.
  If omitted, the provider-level region will be used.
  Changing this parameter will create a new resource.

* `domain_name` - (Required, String, ForceNew) Specifies the ingest domain name.
  Changing this parameter will create a new resource.

* `app_name` - (Required, String, ForceNew) Specifies the application name.
  Changing this parameter will create a new resource.

* `video_encoding` - (Required, String) Specifies the video codec. The valid values are **H264** and **H265**.

* `templates` - (Required, List) Specifies the video quality templates. A maximum of `4` templates can be added.
  The [templates](#transcoding_templates) structure is documented below.  
  For resolution and bitrate settings in the presets,
  please refer to the [document](https://support.huaweicloud.com/intl/en-us/usermanual-live/live01000802.html).

* `trans_type` - (Optional, String) Specifies the transcoding stream trigger mode.
  The valid values are as follows:
  + **play**: Pull stream triggers transcoding.
  + **publish**: Push stream triggers transcoding.

  Defaults to **play**.

* `low_bitrate_hd` - (Optional, Bool) Specifies whether to enable low bitrate HD rates. If enabled
the output media will have a lower bitrate with the same image quality. Defaults to **false**.

<a name="transcoding_templates"></a>
The `templates` block supports:

* `name` - (Required, String) Specifies the template name. The name can contain a maximum of 64 characters, and only
contains letters, digits and hyphens (-).

* `width` - (Required, Int) Specifies video width (unit: pixel).
  + **When the video encoding is H264**, value range: `32` ~ `3,840` and must be a multiple of `2`.
  + **When the video encoding is H265**, value range: `320` ~ `3,840` and must be a multiple of `4`.

* `height` - (Required, Int) Specifies video height (unit: pixel).
  + **When the video encoding is H264**, value range: `32` ~ `2,160` and must be a multiple of `2`.
  + **When the video encoding is H265**, value range: `240` ~ `2,160` and must be a multiple of `4`.

* `bitrate` - (Required, Int) Specifies the bitrate of a transcoded video, in kbit/s. Value range: `40` ~ `30,000`.

* `frame_rate` - (Optional, Int) Specifies the frame rate of the transcoded video, in fps. Value range: `0` ~ `30`.
  Value `0` indicates that the frame rate remains unchanged.

* `protocol` - (Optional, String) Specifies the protocol type supported for transcoding output.
  The valid value is **RTMP**. Defaults to **RTMP**.

* `i_frame_interval` - (Optional, String) Specifies the maximum I-frame interval in frames.
  The value ranges from `0` to `500`, includes `0` and `500`. Defaults to `50`.

  -> If you want to set the i-frame interval through `i_frame_interval`, please set the `gop` to `0` or do not pass the
    `gop` parameter.

* `gop` - (Optional, String) Specifies the interval time for I-frames, in seconds.
  The value ranges from `0` to `10`, includes `0` and `10`. Defaults to `2`.

  -> When `gop` is not `0`, the i-frame interval is set with the `gop` parameter, and the `i_frame_interval` field does
    not take effect.

* `bitrate_adaptive` - (Optional, String) Specifies the adaptive bitrate.
  The valid values are as follows:
  + **off**: Disable rate adaptation and output the target rate according to the set rate.
  + **minimum**: Output the target bitrate based on the minimum value of the set bitrate and source file bitrate.
  + **adaptive**: Adaptive output of target bitrate based on source file bitrate.

  Defaults to **off**.

* `i_frame_policy` - (Optional, String) Specifies the encoding output I-frame strategy.
  The valid values are as follows:
  + **auto**: I-frame output according to the set `gop` duration.
  + **strictSync**: The encoded output I-frame is completely consistent with the source, and the `gop` parameter is
    invalid after setting this value.

  Defaults to **auto**.

  -> In multi bitrate scenarios, it is recommended to enable I-frame random source to ensure alignment of multi bitrate
    I-frames.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in format of **domain_name/app_name**.

## Import

The resource can be imported using the `domain_name` and `app_name`, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_live_transcoding.test <domian_name>/<app_name>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `trans_type`.
It is generally recommended running `terraform plan` after importing the resource.
You can then decide if changes should be applied to the resource, or the resource definition should be updated to align
with the resource. Also, you can ignore changes as below.

```hcl
resource "huaweicloud_live_transcoding" "test" {
  ...

  lifecycle {
    ignore_changes = [
      trans_type,
    ]
  }
}
```
