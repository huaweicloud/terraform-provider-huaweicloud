---
subcategory: "Live"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_live_channels"
description: |-
  Use this data source to get the list of Live channels within HuaweiCloud.
---

# huaweicloud_live_channels

Use this data source to get the list of Live channels within HuaweiCloud.

## Example Usage

```hcl
data "huaweicloud_live_channels" "test" {
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String) Specifies the region in which to query the resource.
  If omitted, the provider-level region will be used.

* `domain_name` - (Optional, String) Specifies the channel streaming domain name.

* `app_name` - (Optional, String) Specifies the group name or application name.

* `channel_id` - (Optional, String) Specifies the channel ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.

* `channels` - The channel information.

  The [channels](#channels_struct) structure is documented below.

<a name="channels_struct"></a>
The `channels` block supports:

* `domain_name` - The channel streaming domain name.

* `app_name` - The group name or application name.

* `state` - The channel status. Valid values are:
  + **ON**: After a channel is delivered, functions such as stream pull, transcoding, and recording are automatically enabled.
  + **OFF**: Only the channel information is saved but the channel is not started.

* `input` - The channel input information.
  The [input](#LiveChannel_Input) structure is documented below.

* `record_settings` - The configuration for replaying a recording.
  The [record_settings](#LiveChannel_RecordSettings) structure is documented below.

* `endpoints` - The channel outflow information.
  The [endpoints](#LiveChannel_Endpoints) structure is documented below.

* `encoder_settings_expand` - The audio output configuration.
  The [encoder_settings_expand](#LiveChannel_EncoderSettingsExpand) structure is documented below.

* `encoder_settings` - The transcoding template configuration.
  The [encoder_settings](#LiveChannel_EncoderSettings) structure is documented below.

* `name` - The channel name. The name can be duplicated.

* `id` - The channel ID.

<a name="LiveChannel_Input"></a>
The `input` block supports:

* `input_protocol` - The channel input protocol. Valid values are:
  + **FLV_PULL**.
  + **RTMP_PUSH**.
  + **HLS_PULL**.
  + **SRT_PULL**.
  + **SRT_PUSH**.

* `sources` - The channel main source stream information.
  The [sources](#LiveChannel_Sources) structure is documented below.

* `secondary_sources` - The prepared stream array.
  The [secondary_sources](#LiveChannel_SecondarySources) structure is documented below.

* `failover_conditions` - The configuration of switching between primary and backup audio and video stream URLs.
  The [failover_conditions](#LiveChannel_FailoverConditions) structure is documented below.

* `max_bandwidth_limit` - The maximum bandwidth that needs to be configured when the inbound protocol is **HLS_PULL**.
  The unit is **bps**.

* `ip_port_mode` - The IP port mode.

* `ip_whitelist` - The IP whitelist when protocol is **SRT_PUSH**.

* `scte35_source` - The advertisement scte35 signal source.

* `ad_triggers` - The ad trigger configuration list. Valid Values are:
  + **Splice insert**.
  + **Provider advertisement**.
  + **Distributor advertisement**.
  + **Provider placement opportunity**.
  + **Distributor placement opportunity**.

* `audio_selectors` - The audio selector configuration.
  The [audio_selectors](#LiveChannel_AudioSelectors) structure is documented below.

<a name="LiveChannel_Sources"></a>
The `sources` block supports:

* `url` - The channel source stream URL, used for external streaming.

* `bitrate` - The bitrate. The unit is **bps**.

* `width` - The resolution corresponds to the width value.

* `height` - The resolution corresponds to the high value.

* `enable_snapshot` - Whether to use this stream to take screenshots.

* `bitrate_for3u8` - Whether to use bitrate to fix the bitrate.

* `passphrase` - The encrypted information when the protocol is **SRT_PUSH**.

* `backup_urls` - The list of backup stream addresses.

* `stream_id` - The stream ID of the stream pull address when the channel type is **SRT_PULL**.

* `latency` - The streaming delay when the channel type is **SRT_PULL**.

<a name="LiveChannel_SecondarySources"></a>
The `secondary_sources` block supports:

* `url` - The channel source stream URL, used for external streaming.

* `bitrate` - The bitrate. The unit is **bps**.

* `width` - The resolution corresponds to the width value.

* `height` - The resolution corresponds to the high value.

* `bitrate_for3u8` - Whether to use bitrate to fix the bitrate.

* `passphrase` - The encrypted information when the protocol is **SRT_PUSH**.

* `backup_urls` - The list of backup stream addresses.

* `stream_id` - The stream ID of the stream pull address when the channel type is **SRT_PULL**.

* `latency` - The streaming delay when the channel type is **SRT_PULL**.

<a name="LiveChannel_FailoverConditions"></a>
The `failover_conditions` block supports:

* `input_loss_threshold_msec` - The duration threshold of inflow stop. The unit is millisecond.

* `input_preference` - The input preference type. Valid values are:
  + **PRIMARY**: The main incoming URL is the first priority.
  + **EQUAL**: Equal switching between primary and backup URLs.

<a name="LiveChannel_AudioSelectors"></a>
The `audio_selectors` block supports:

* `name` - The name of the audio selector.

* `selector_settings` - The audio selector configuration.
  The [selector_settings](#LiveChannel_SelectorSettings) structure is documented below.

<a name="LiveChannel_SelectorSettings"></a>
The `selector_settings` block supports:

* `audio_language_selection` - The language selector configuration.
  The [audio_language_selection](#LiveChannel_AudioLanguageSelection) structure is documented below.

* `audio_pid_selection` - The PID selector configuration.
  The [audio_pid_selection](#LiveChannel_AudioPidSelection) structure is documented below.

* `audio_hls_selection` - The HLS selector configuration.
  The [audio_hls_selection](#LiveChannel_AudioHlsSelection) structure is documented below.

<a name="LiveChannel_AudioLanguageSelection"></a>
The `audio_language_selection` block supports:

* `language_code` - The language abbreviation. Supports `2` or `3` lowercase letter language codes.

* `language_selection_policy` - The language output strategy. Valid values are:
  + **LOOSE**: Loose matching. For example, "eng" will prioritize matching tracks with English as the language in the
    source stream. If no match is found, the track with the smallest PID will be selected.
  + **STRICT**: Strict matching. For example, "eng" will strictly match the audio track in the source stream whose
    language is English. If no match is found, the media live broadcast service will automatically fill in a silent
    segment. When the terminal uses this audio selector to play the video, it will be played silently.

<a name="LiveChannel_AudioPidSelection"></a>
The `audio_pid_selection` block supports:

* `pid` - The value of PID.

<a name="LiveChannel_AudioHlsSelection"></a>
The `audio_hls_selection` block supports:

* `name` - The HLS audio selector name.

* `group_id` - The HLS audio selector gid.

<a name="LiveChannel_RecordSettings"></a>
The `record_settings` block supports:

* `rollingbuffer_duration` - The maximum playback recording time. During this time period, the recording will continue.
  The unit is second.

<a name="LiveChannel_Endpoints"></a>
The `endpoints` block supports:

* `hls_package` - The HLS packaging information.
  The [hls_package](#LiveChannel_HlsPackage) structure is documented below.

* `dash_package` - The DASH packaging information.
  The [dash_package](#LiveChannel_DashPackage) structure is documented below.

* `mss_package` - The MSS packaging information.
  The [mss_package](#LiveChannel_MssPackage) structure is documented below.

<a name="LiveChannel_HlsPackage"></a>
The `hls_package` block supports:

* `url` - The customer-defined streaming address, including method, domain name, and path.

* `stream_selection` - The stream selection. Filter out the specified range of streams from the full stream.
  The [stream_selection](#LiveChannel_StreamSelection) structure is documented below.

* `hls_version` - The HLS version.

* `segment_duration_seconds` - The duration of the channel output segment. The unit is second.

* `playlist_window_seconds` - The window length of the channel live broadcast return shard. The unit is second.

* `encryption` - The encrypted information.
  The [encryption](#LiveChannel_Encryption) structure is documented below.

* `request_args` - The play related configuration.
  The [request_args](#LiveChannel_RequestArgs) structure is documented below.

* `ad_marker` - The advertising marker.

<a name="LiveChannel_DashPackage"></a>
The `dash_package` block supports:

* `url` - The customer-defined streaming address, including method, domain name, and path.

* `stream_selection` - The stream selection. Filter out the specified range of streams from the full stream.
  The [stream_selection](#LiveChannel_StreamSelection) structure is documented below.

* `segment_duration_seconds` - The duration of the channel output segment. The unit is second.

* `playlist_window_seconds` - The window length of the channel live broadcast return shard. The unit is second.

* `encryption` - The encrypted information.
  The [encryption](#LiveChannel_Encryption) structure is documented below.

* `request_args` - The play related configuration.
  The [request_args](#LiveChannel_RequestArgs) structure is documented below.

* `ad_marker` - The advertising marker.

<a name="LiveChannel_MssPackage"></a>
The `mss_package` block supports:

* `url` - The customer-defined streaming address, including method, domain name, and path.

* `stream_selection` - The stream selection. Filter out the specified range of streams from the full stream.
  The [stream_selection](#LiveChannel_StreamSelection) structure is documented below.

* `segment_duration_seconds` - The duration of the channel output segment. The unit is second.

* `playlist_window_seconds` - The window length of the channel live broadcast return shard. The unit is second.

* `encryption` - The encrypted information.
  The [encryption](#LiveChannel_Encryption) structure is documented below.

* `delay_segment` - The delayed playback time. The unit is second.

* `request_args` - The play related configuration.
  The [request_args](#LiveChannel_RequestArgs) structure is documented below.

<a name="LiveChannel_StreamSelection"></a>
The `stream_selection` block supports:

* `key` - The key used for bitrate filtering in streaming URLs.

* `max_bandwidth` - The maximum code rate. The unit is bps.

* `min_bandwidth` - The minimum code rate. The unit is bps.

<a name="LiveChannel_Encryption"></a>
The `encryption` block supports:

* `level` - The level. Valid values are:
  + **content**: One channel corresponds to one key.
  + **profile**: One code rate corresponds to one key.

* `resource_id` - The customer-generated DRM content ID.

* `system_ids` - The system ID enumeration values. Valid values are **FairPlay** (HLS),
  **Widevine** (DASH), **PlayReady** (DASH), and **PlayReady** (MSS).

* `url` - The DRM address of the key.

* `speke_version` - The DRM spec version number.

* `request_mode` - The request mode. Valid values are:
  + **direct_http**: HTTP(S) direct access to DRM.
  + **functiongraph_proxy**: FunctionGraph proxy access to DRM.

* `key_rotation_interval_seconds` - The key rotation interval seconds.

* `encryption_method` - The encryption method.

* `http_headers` - The authentication information that needs to be added to the DRM request header.
  The [http_headers](#LiveChannel_HttpHeader) structure is documented below.

* `urn` - The URN of the function graph.

<a name="LiveChannel_HttpHeader"></a>
The `http_headers` block supports:

* `key` - The key field name in the request header.

* `value` - The value corresponding to the key in the request header.

<a name="LiveChannel_RequestArgs"></a>
The `request_args` block supports:

* `record` - The recording and playback related configuration.
  The [record](#LiveChannel_RequestArgsRecord) structure is documented below.

* `timeshift` - The time-shift playback configuration.
  The [timeshift](#LiveChannel_RequestArgsTimeShift) structure is documented below.

* `live` - The live broadcast configuration.
  The [live](#LiveChannel_RequestArgsLive) structure is documented below.

<a name="LiveChannel_RequestArgsRecord"></a>
The `record` block supports:

* `start_time` - The start time.

* `end_time` - The end time.

* `format` - The format.

* `unit` - The unit.

<a name="LiveChannel_RequestArgsTimeShift"></a>
The `timeshift` block supports:

* `back_time` - The time shift duration field name.

* `unit` - The unit.

<a name="LiveChannel_RequestArgsLive"></a>
The `live` block supports:

* `delay` - The delay field.

* `unit` - The unit.

<a name="LiveChannel_EncoderSettingsExpand"></a>
The `encoder_settings_expand` block supports:

* `audio_descriptions` - The description of the audio output configuration.
  The [audio_descriptions](#LiveChannel_AudioDescriptions) structure is documented below.

<a name="LiveChannel_AudioDescriptions"></a>
The `audio_descriptions` block supports:

* `name` - The name of the audio output configuration.

* `audio_selector_name` - The audio selector name.

* `language_code_control` - The language code control configuration. Valid values are:
  + **FOLLOW_INPUT**: If the output audio corresponding to the selected audio selector has a language, it will be
    consistent with it, otherwise it will be backed up by the language code and stream name configured here.
    The current option is recommended and is the default value.
  + **USE_CONFIGURED**: Users can customize the language and stream name of the output audio based on actual conditions.

* `language_code` - The language code.

* `stream_name` - The stream name.

<a name="LiveChannel_EncoderSettings"></a>
The `encoder_settings` block supports:

* `template_id` - The transcoding template ID.
