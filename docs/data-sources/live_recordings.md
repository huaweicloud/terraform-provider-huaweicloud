---
subcategory: "Live"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_live_recordings"
description: |-
  Use this data source to get the list of recording rules.
---

# huaweicloud_live_recordings

Use this datasource to get the list of recording rules.

## Example Usage

```hcl
data "huaweicloud_live_recordings" "test" {}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `domain_name` - (Optional, String) Specifies the ingest domain name to which the recording rules belong.

* `app_name` - (Optional, String) Specifies the application name of the recording rule.

* `stream_name` - (Optional, String) Specifies the stream name of the recording rule.

* `type` - (Optional, String) Specifies the recording type of the recording rule.
  The valid values are as follows:
  + **CONTINUOUS_RECORD**: Indicates continuous recording.
  + **COMMAND_RECORD**: Indicates command recording.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `rules` - The list of the recording rules.

  The [rules](#rules_struct) structure is documented below.

<a name="rules_struct"></a>
The `rules` block supports:

* `id` - The recording rule ID.

* `domain_name` - The ingest domain name to which the recording rule belongs.

* `app_name` - The application name of the recording rule.

* `stream_name` - The stream name of the recording rule.

* `type` - The recording type of the recording rule.

* `default_record_config` - The default recording configuration rule.

  The [default_record_config](#rules_default_record_config_struct) structure is documented below.

* `created_at` - The creation time of the recording rule.
  The format is **yyyy-mm-ddThh:mm:ssZ**. e.g. **2024-09-01T15:30:20Z**.

* `updated_at` - The lasted update time of the recording rule.
  The format is **yyyy-mm-ddThh:mm:ssZ**. e.g. **2024-09-01T15:30:20Z**.

<a name="rules_default_record_config_struct"></a>
The `default_record_config` block supports:

* `record_format` - The recording format.
  The valid values are **HLS**, **FLV** and **MP4**.

* `obs` - The OBS bucket information for storing recordings.

  The [obs](#default_record_config_obs_struct) structure is documented below.

* `hls` - The HLS configuration rule.

  The [hls](#default_record_config_hls_struct) structure is documented below.

* `flv` - The FLV configuration rule.

  The [flv](#default_record_config_flv_struct) structure is documented below.

* `mp4` - The MP4 configuration rule.

  The [mp4](#default_record_config_mp4_struct) structure is documented below.

<a name="default_record_config_obs_struct"></a>
The `obs` block supports:

* `bucket` - The OBS bucket name.

* `region` - The region to which the OBS bucket belongs.

* `object` - The OBS object storage path.

<a name="default_record_config_hls_struct"></a>
The `hls` block supports:

* `recording_length` - The periodic recording duration, in seconds.

* `file_naming` - The file path and file name prefix of the recorded M3U8 file.

* `ts_file_naming` - The file name prefix of recorded TS file.

* `record_slice_duration` - The TS slicing duration during HLS recording, in seconds.

* `max_stream_pause_length` - The recording HLS file concatenation duration, in seconds.

<a name="default_record_config_flv_struct"></a>
The `flv` block supports:

* `recording_length` - The periodic recording duration, in seconds.

* `file_naming` - The file path and file name prefix of the recorded FLV file.

* `max_stream_pause_length` - The recording FLV file concatenation duration, in seconds.

<a name="default_record_config_mp4_struct"></a>
The `mp4` block supports:

* `recording_length` - The periodic recording duration, in seconds.

* `file_naming` - The file path and file name prefix of the recorded MP4 file.

* `max_stream_pause_length` - The recording MP4 file concatenation duration, in seconds.
