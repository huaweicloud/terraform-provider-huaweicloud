---
subcategory: "Live"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_live_recording"
description: |-
  Manages a recording template within HuaweiCloud Live.
---

# huaweicloud_live_recording

Manages a recording template within HuaweiCloud Live.

## Example Usage

### Create a recording template for an ingest domain name

```hcl
variable "ingest_domain_name" {}
variable "bucket_region" {}
variable "bucket_name" {}

resource "huaweicloud_live_domain" "ingestDomain" {
  name = var.ingest_domain_name
  type = "push"
}

resource "huaweicloud_live_recording" "recording" {
  domain_name = huaweicloud_live_domain.ingestDomain.name
  app_name    = "live"
  stream_name = "stream_name"
  type        = "CONTINUOUS_RECORD"

  obs {
    region = var.bucket_region
    bucket = var.bucket_name
  }

  hls {
    recording_length = 15
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `domain_name` - (Required, String) Specifies the ingest domain name.

* `app_name` - (Required, String) Specifies the application name. To match all names, use an asterisk (*).

* `stream_name` - (Required, String) Specifies the stream name. To match all names, use an asterisk (*).

* `type` - (Optional, String, ForceNew) Specifies the types of recording notifications. The options are as follows:
  + **CONTINUOUS_RECORD**: continuous recording. Recording is triggered once streams are pushed to the recording system.
  + **COMMAND_RECORD**: command-based recording. Tenants need to run commands to start and stop recording after streams
    are pushed to the recording system.

  Defaults to `CONTINUOUS_RECORD`. Changing this parameter will create a new resource.

* `obs` - (Required, List) Specifies the obs for storing recordings.
  The [obs](#recording_obs) structure is documented below.

* `hls` - (Optional, List) Specifies the HLS configuration rule for storing recording as HLS.
  The [hls](#recording_HLS) structure is documented below.

* `flv` - (Optional, List) Specifies the FLV configuration rule for storing recording as FLV.
  The [flv](#recording_FLV_MP4) structure is documented below.

* `mp4` - (Optional, List) Specifies the MP4 configuration rule for storing recording as MP4.
  The [mp4](#recording_FLV_MP4) structure is documented below.

-> At least one of `hls`, `flv`, `mp4` must be specified.

<a name="recording_obs"></a>
The `obs` block supports:

* `region` - (Required, String) Specifies the region of OBS.

* `bucket` - (Required, String) Specifies OBS bucket.

* `object` - (Optional, String) Specifies OBS object path. If omitted, recordings will be saved to the root directory.

<a name="recording_HLS"></a>
The `hls` block supports:

* `recording_length` - (Required, Int) Specifies the recording length. Value range: `15` ~ `720`, unit: `minute`.
  A stream exceeding the recording length will generate a new recording.

* `file_naming` - (Optional, String) Specifies the path and file name prefix of an M3U8 file. The default value is
  `Record/{publish_domain}/{app}/{record_type}/{record_format}/{stream}_{file_start_time}/{stream}_{file_start_time}`.

* `ts_file_naming` - (Optional, String) Specifies TS file name prefix.
  The default value is `{file_start_time_unix}_{file_end_time_unix}_{ts_sequence_number}`.

* `record_slice_duration` - (Optional, Int) Specifies the TS slice duration for HLS recording.
  Value range: `2` ~ `60`, unit: `second`. Defaults to `10`.

* `max_stream_pause_length` - (Optional, Int) Specifies the interval threshold for combining HLS chunks. If the stream
  pause length exceeds the value of this parameter, a new recording is generated.
  Value range: `-1` ~ `300`, unit: `second`. Defaults to `0`.

  -> 1. If the value is set to `0`, a new file will be generated once the stream is interrupted.
    <br/>2. If the value is set to `-1`, the HLS chunks will be combined to the previous file generated within `30` days
    after the same stream is recovered.

<a name="recording_FLV_MP4"></a>
The `flv` and `mp4` block support:

* `recording_length` - (Required, Int) Specifies the recording length. Value range: `15` ~ `180`, unit: `minute`.
  A stream exceeding the recording length will generate a new recording.

* `file_naming` - (Optional, String) Specifies the path and file name prefix of a recording file. The default value is
  `Record/{publish_domain}/{app}/{record_type}/{record_format}/{stream}_{file_start_time}/{file_start_time}`.

* `max_stream_pause_length` - (Optional, Int) Specifies the interval threshold for combining recording chunks. If the
  stream pause length exceeds the value of this parameter, a new recording is generated.
  Value range: `0` ~ `300`, unit: `second`. Defaults to `0`.
  If the value is set to `0`, a new file will be generated once the stream is interrupted.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID in UUID format.

## Import

Recording templates can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_live_recording.test <id>
```
