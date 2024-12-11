---
subcategory: "Live"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_live_channel"
description: |-
  Manages a Live channel resource within HuaweiCloud.
---

# huaweicloud_live_channel

Manages a Live channel resource within HuaweiCloud.

## Example Usage

### Create a Live channel with RTMP_PUSH protocol

```hcl
variable "domain_name" {}
variable "template_id" {}
variable "hls_package_url" {}
variable "hls_package_encryption_url" {}
variable "hls_package_encryption_urn" {}
variable "mss_package_url" {}

resource "huaweicloud_live_channel" "test" {
  app_name    = "live"
  domain_name = var.domain_name
  name        = "test-name"
  state       = "ON"

  encoder_settings {
    template_id = var.template_id
  }

  endpoints {
    hls_package {
      hls_version              = "v3"
      playlist_window_seconds  = 24
      segment_duration_seconds = 4
      url                      = var.hls_package_url

      encryption {
        level         = "profile"
        request_mode  = "functiongraph_proxy"
        resource_id   = "test-resource-id"
        speke_version = "1.0"
        url           = var.hls_package_encryption_url
        urn           = var.hls_package_encryption_urn

        system_ids = [
          "FairPlay",
        ]
      }

      request_args {
        record {
          end_time   = "end"
          format     = "timestamp"
          start_time = "begin"
          unit       = "second"
        }

        timeshift {
          back_time = "delay"
          unit      = "second"
        }
      }
    }

    mss_package {
      playlist_window_seconds  = 42
      segment_duration_seconds = 4
      url                      = var.mss_package_url

      encryption {
        level         = "content"
        request_mode  = "direct_http"
        resource_id   = "test-resource-id"
        speke_version = "1.0"
        url           = "https://test-url.cd"

        system_ids = [
          "PlayReady",
        ]
      }

      request_args {
        record {
          end_time   = "end"
          format     = "timestamp"
          start_time = "begin"
          unit       = "second"
        }

        timeshift {
          back_time = "delay"
          unit      = "second"
        }
      }
    }
  }

  input {
    input_protocol = "RTMP_PUSH"

    sources {
      bitrate = 100
    }
  }

  record_settings {
    rollingbuffer_duration = 20
  }
}
```

### Create a Live channel with FLV_PULL protocol

```hcl
variable "domain_name" {}
variable "template_id" {}
variable "dash_package_url" {}
variable "hls_package_url" {}

resource "huaweicloud_live_channel" "test" {
  app_name    = "live"
  domain_name = var.domain_name
  name        = "test-name"
  state       = "ON"

  encoder_settings {
    template_id = var.template_id
  }

  endpoints {
    dash_package {
      playlist_window_seconds  = 24
      segment_duration_seconds = 4
      url                      = var.dash_package_url

      request_args {
        record {
          end_time   = "end"
          format     = "timestamp"
          start_time = "begin"
          unit       = "second"
        }

        timeshift {
          back_time = "delay"
          unit      = "second"
        }
      }
    }

    hls_package {
      hls_version              = "v3"
      playlist_window_seconds  = 8
      segment_duration_seconds = 4
      url                      = var.hls_package_url

      encryption {
        level         = "content"
        request_mode  = "direct_http"
        resource_id   = "test-resource-id"
        speke_version = "1.0"
        url           = "http://xxx.sp"

        system_ids = [
          "FairPlay",
        ]

        http_headers {
          key   = "aaa"
          value = "sss"
        }

        http_headers {
          key   = "www"
          value = "qqq"
        }
      }

      request_args {
        record {
          end_time   = "end"
          format     = "timestamp"
          start_time = "begin"
          unit       = "second"
        }

        timeshift {
          back_time = "delay"
          unit      = "second"
        }
      }
    }
  }

  input {
    input_protocol = "FLV_PULL"

    failover_conditions {
      input_loss_threshold_msec = 4000
      input_preference          = "EQUAL"
    }

    secondary_sources {
      bitrate = 100
      url     = "https://hgf.vv"
    }

    sources {
      bitrate = 100
      url     = "https://qwe.cc"
    }
  }

  record_settings {
    rollingbuffer_duration = 12
  }
}
```

### Create a Live channel with HLS_PULL protocol

```hcl
variable "domain_name" {}
variable "template_id" {}
variable "hls_package_url" {}
variable "hls_package_encryption_url" {}
variable "hls_package_encryption_urn" {}
variable "mss_package_url" {}

resource "huaweicloud_live_channel" "test" {
  app_name    = "live"
  domain_name = var.domain_name
  name        = "test-name"
  state       = "ON"

  encoder_settings {
    template_id = var.template_id
  }

  endpoints {
    hls_package {
      hls_version              = "v3"
      playlist_window_seconds  = 4
      segment_duration_seconds = 2
      url                      = var.hls_package_url

      encryption {
        level         = "content"
        request_mode  = "functiongraph_proxy"
        resource_id   = "test-resource-id"
        speke_version = "1.0"
        url           = var.hls_package_encryption_url
        urn           = var.hls_package_encryption_urn

        system_ids = [
          "FairPlay",
        ]
      }

      request_args {
        record {
          end_time   = "end"
          format     = "timestamp"
          start_time = "begin"
          unit       = "second"
        }

        timeshift {
          back_time = "delay"
          unit      = "second"
        }
      }
    }

    mss_package {
      playlist_window_seconds  = 8
      segment_duration_seconds = 2
      url                      = var.mss_package_url

      encryption {
        level         = "content"
        request_mode  = "direct_http"
        resource_id   = "dfge"
        speke_version = "1.0"
        url           = "https://ssc.cd"

        system_ids = [
          "PlayReady",
        ]

        http_headers {
          key   = "aa"
          value = "ss"
        }

        http_headers {
          key   = "gg"
          value = "ff"
        }
      }

      request_args {
        record {
          end_time   = "end"
          format     = "timestamp"
          start_time = "begin"
          unit       = "second"
        }

        timeshift {
          back_time = "delay"
          unit      = "second"
        }
      }
    }
  }

  input {
    input_protocol      = "HLS_PULL"
    max_bandwidth_limit = 200

    failover_conditions {
      input_loss_threshold_msec = 2000
      input_preference          = "PRIMARY"
    }

    secondary_sources {
      bitrate = 100
      url     = "https://qqwe.dd"
    }

    sources {
      bitrate = 100
      url     = "https://ssa.qw"
    }
  }

  record_settings {
    rollingbuffer_duration = 3
  }
}
```

### Create a Live channel with SRT_PUSH protocol

```hcl
variable "domain_name" {}
variable "template_id" {}
variable "hls_package_url" {}
variable "mss_package_url" {}

resource "huaweicloud_live_channel" "test" {
  app_name    = "live"
  domain_name = var.domain_name
  name        = "test-name"
  state       = "ON"

  encoder_settings {
    template_id = var.template_id
  }

  endpoints {
    hls_package {
      hls_version              = "v3"
      playlist_window_seconds  = 4
      segment_duration_seconds = 2
      url                      = var.hls_package_url

      encryption {
        level         = "content"
        request_mode  = "direct_http"
        resource_id   = "test-resource-id"
        speke_version = "1.0"
        url           = "https://qqq.co"

        system_ids = [
          "FairPlay",
        ]

        http_headers {
          key   = "aa"
          value = "sss"
        }

        http_headers {
          key   = "dd"
          value = "sss"
        }
      }

      request_args {
        record {
          end_time   = "end"
          format     = "timestamp"
          start_time = "begin"
          unit       = "second"
        }

        timeshift {
          back_time = "delay"
          unit      = "second"
        }
      }
    }

    mss_package {
      playlist_window_seconds  = 8
      segment_duration_seconds = 2
      url                      = var.mss_package_url

      request_args {
        record {
          end_time   = "end"
          format     = "timestamp"
          start_time = "begin"
          unit       = "second"
        }

        timeshift {
          back_time = "delay"
          unit      = "second"
        }
      }
    }
  }

  input {
    input_protocol = "SRT_PUSH"
    ip_whitelist   = "192.168.0.1/16,192.168.1.1/16,192.168.2.1/16"

    audio_selectors {
      name = "test-audio-selectors1"

      selector_settings {
        audio_pid_selection {
          pid = 2
        }
      }
    }

    audio_selectors {
      name = "test-audio-selectors2"

      selector_settings {
        audio_language_selection {
          language_code             = "ch"
          language_selection_policy = "LOOSE"
        }
      }
    }

    audio_selectors {
      name = "test-audio-selectors3"

      selector_settings {
        audio_pid_selection {
          pid = 0
        }
      }
    }

    sources {
      bitrate = 100
    }
  }

  record_settings {
    rollingbuffer_duration = 2
  }
}
```

### Create a Live channel with SRT_PULL protocol

```hcl
variable "domain_name" {}
variable "template_id" {}
variable "hls_package_url" {}
variable "mss_package_url" {}

resource "huaweicloud_live_channel" "test" {
  app_name    = "live"
  domain_name = var.domain_name
  name        = "test-name"
  state       = "ON"

  encoder_settings {
    template_id = var.template_id
  }

  endpoints {
    hls_package {
      hls_version              = "v3"
      playlist_window_seconds  = 16
      segment_duration_seconds = 4
      url                      = var.hls_package_url

      encryption {
        level         = "content"
        request_mode  = "direct_http"
        resource_id   = "test-resource-id"
        speke_version = "1.0"
        url           = "https://sss.cc"

        system_ids = [
          "FairPlay",
        ]

        http_headers {
          key   = "aa"
          value = "ss"
        }

        http_headers {
          key   = "ff"
          value = "dd"
        }
      }

      request_args {
        record {
          end_time   = "end"
          format     = "timestamp"
          start_time = "begin"
          unit       = "second"
        }

        timeshift {
          back_time = "delay"
          unit      = "second"
        }
      }
    }

    mss_package {
      playlist_window_seconds  = 24
      segment_duration_seconds = 4
      url                      = var.mss_package_url

      request_args {
        record {
          end_time   = "end"
          format     = "timestamp"
          start_time = "begin"
          unit       = "second"
        }

        timeshift {
          back_time = "delay"
          unit      = "second"
        }
      }
    }
  }

  input {
    input_protocol = "SRT_PULL"

    audio_selectors {
      name = "test-audio-selectors1"

      selector_settings {
        audio_language_selection {
          language_code             = "dfg"
          language_selection_policy = "LOOSE"
        }
      }
    }

    audio_selectors {
      name = "test-audio-selectors2"

      selector_settings {
        audio_pid_selection {
          pid = 13
        }
      }
    }

    audio_selectors {
      name = "test-audio-selectors3"

      selector_settings {
        audio_pid_selection {
          pid = 0
        }
      }
    }

    failover_conditions {
      input_loss_threshold_msec = 2000
      input_preference          = "EQUAL"
    }

    secondary_sources {
      bitrate   = 100
      latency   = 1000
      stream_id = "vcbeer"
      url       = "srt://192.168.1.215:9001"
    }

    sources {
      bitrate   = 100
      latency   = 2000
      stream_id = "dfawerw"
      url       = "srt://192.168.1.216:9001"
    }
  }

  record_settings {
    rollingbuffer_duration = 4
  }
}
```

### Create a Live channel with custom channel ID

```hcl
variable "domain_name" {}
variable "channel_id" {}
variable "template_id" {}
variable "hls_package_url" {}

resource "huaweicloud_live_channel" "test" {
  app_name    = "live"
  domain_name = var.domain_name
  state       = "OFF"
  channel_id  = var.channel_id

  encoder_settings {
    template_id = var.template_id
  }

  endpoints {
    hls_package {
      hls_version              = "v3"
      playlist_window_seconds  = 40
      segment_duration_seconds = 4
      url                      = var.hls_package_url

      request_args {
        record {
          end_time   = "end"
          format     = "timestamp"
          start_time = "begin"
          unit       = "second"
        }

        timeshift {
          back_time = "delay"
          unit      = "second"
        }
      }
    }
  }

  input {
    input_protocol = "RTMP_PUSH"

    sources {
      bitrate = 100
    }
  }

  record_settings {
    rollingbuffer_duration = 0
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `domain_name` - (Required, String, ForceNew) Specifies the channel streaming domain name.

  Changing this parameter will create a new resource.

* `app_name` - (Required, String, ForceNew) Specifies the group name or application name. Currently, only supports **live**.

  Changing this parameter will create a new resource.

* `state` - (Required, String) Specifies the channel status. Valid values are:
  + **ON**: After a channel is delivered, functions such as stream pull, transcoding, and recording are automatically enabled.
  + **OFF**: Only the channel information is saved but the channel is not started.

* `input` - (Required, List) Specifies the channel input information.
  The [input](#LiveChannel_Input) structure is documented below.

* `record_settings` - (Required, List) Specifies the configuration for replaying a recording.
  The [record_settings](#LiveChannel_RecordSettings) structure is documented below.

* `endpoints` - (Required, List) Specifies the channel outflow information.
  The [endpoints](#LiveChannel_Endpoints) structure is documented below.

* `encoder_settings_expand` - (Optional, List) Specifies the audio output configuration.
  The [encoder_settings_expand](#LiveChannel_EncoderSettingsExpand) structure is documented below.

* `encoder_settings` - (Optional, List) Specifies the transcoding template configuration.
  The [encoder_settings](#LiveChannel_EncoderSettings) structure is documented below.

* `name` - (Optional, String) Specifies the channel name. The name can be duplicated.

* `channel_id` - (Optional, String, ForceNew) Specifies the unique channel ID.

  Changing this parameter will create a new resource.

  -> If this field is configured, it will be used as the unique channel ID. Otherwise, a randomly generated string
  will be used as the unique ID.

<a name="LiveChannel_Input"></a>
The `input` block supports:

* `input_protocol` - (Required, String, ForceNew) Specifies the channel input protocol. Valid values are:
  + **FLV_PULL**.
  + **RTMP_PUSH**.
  + **HLS_PULL**.
  + **SRT_PULL**.
  + **SRT_PUSH**.

  Changing this parameter will create a new resource.

* `sources` - (Optional, List) Specifies the channel main source stream information. This parameter is optional
  when the stream input protocol is **RTMP_PUSH** or **SRT_PUSH**. In other cases, this parameter is mandatory.
  The [sources](#LiveChannel_Sources) structure is documented below.

* `secondary_sources` - (Optional, List) Specifies the prepared stream array. If this parameter is configured, ensure
  that the number of channels, codec, and resolution of the primary and standby input streams are the same.
  This field does not need to be configured when the stream input protocol is **RTMP_PUSH**.
  The [secondary_sources](#LiveChannel_SecondarySources) structure is documented below.

* `failover_conditions` - (Optional, List) Specifies the configuration of switching between primary and backup audio
  and video stream URLs.
  The [failover_conditions](#LiveChannel_FailoverConditions) structure is documented below.

* `max_bandwidth_limit` - (Optional, Int) Specifies the maximum bandwidth that needs to be configured when the inbound
  protocol is **HLS_PULL**. The unit is **bps**.

  -> In the streaming URL provided by the user, the bandwidth parameter "BANDWIDTH" will be carried for audio and video
  with different bit rates.<br/>1. If the maximum bandwidth is configured here, when the media live broadcast service
  pulls the stream from the URL, it will select the audio and video stream with a smaller bandwidth and the highest bit
  rate and push it to the source station.<br/>2. If the maximum bandwidth is not configured here, when the media live
  broadcast service pulls the stream from the URL, it will select the audio and video stream with the highest "BANDWIDTH"
  by default and push the stream to the source station.

* `ip_port_mode` - (Optional, Bool) Specifies the IP port mode.

  -> When the stream push protocol is **SRT_PUSH** and streams are pushed to the origin server, set this parameter
  to **true**.

* `ip_whitelist` - (Optional, String) Specifies the IP whitelist when protocol is **SRT_PUSH**.

* `scte35_source` - (Optional, String) Specifies the advertisement scte35 signal source. This configuration is only
  supported for **HLS_PULL** channels, and currently only supports **SEGMENTS**.

* `ad_triggers` - (Optional, List) Specifies the ad trigger configuration list. Valid Values are:
  + **Splice insert**.
  + **Provider advertisement**.
  + **Distributor advertisement**.
  + **Provider placement opportunity**.
  + **Distributor placement opportunity**.

* `audio_selectors` - (Optional, List) Specifies the audio selector configuration. Set up to `8` audio selectors.
  The [audio_selectors](#LiveChannel_AudioSelectors) structure is documented below.

<a name="LiveChannel_Sources"></a>
The `sources` block supports:

* `url` - (Optional, String) Specifies the channel source stream URL, used for external streaming.

* `bitrate` - (Optional, Int) Specifies the bitrate. This parameter is required when live transcoding is not required.
  The unit is **bps**. Value ranges from `0` to `104,857,600`.

* `width` - (Optional, Int) Specifies the resolution corresponds to the width value. Value ranges from `0` to `4,096`.

* `height` - (Optional, Int) Specifies the resolution corresponds to the high value. Value ranges from `0` to `2,160`.

* `enable_snapshot` - (Optional, Bool) Specifies whether to use this stream to take screenshots.

* `bitrate_for3u8` - (Optional, Bool) Specifies whether to use bitrate to fix the bitrate. Defaults to **false**.

* `passphrase` - (Optional, String) Specifies the encrypted information when the protocol is **SRT_PUSH**.

* `backup_urls` - (Optional, List) Specifies the list of backup stream addresses.

* `stream_id` - (Optional, String) Specifies the stream ID of the stream pull address when the channel type is **SRT_PULL**.

* `latency` - (Optional, Int) Specifies the streaming delay when the channel type is **SRT_PULL**.

<a name="LiveChannel_SecondarySources"></a>
The `secondary_sources` block supports:

* `url` - (Optional, String) Specifies the channel source stream URL, used for external streaming.

* `bitrate` - (Optional, Int) Specifies the bitrate. This parameter is required when live transcoding is not required.
  The unit is **bps**. Value ranges from `0` to `104,857,600`.

* `width` - (Optional, Int) Specifies the resolution corresponds to the width value. Value ranges from `0` to `4,096`.

* `height` - (Optional, Int) Specifies the resolution corresponds to the high value. Value ranges from `0` to `2,160`.

* `bitrate_for3u8` - (Optional, Bool) Specifies whether to use bitrate to fix the bitrate. Defaults to **false**.

* `passphrase` - (Optional, String) Specifies the encrypted information when the protocol is **SRT_PUSH**.

* `backup_urls` - (Optional, List) Specifies the list of backup stream addresses.

* `stream_id` - (Optional, String) Specifies the stream ID of the stream pull address when the channel type is **SRT_PULL**.

* `latency` - (Optional, Int) Specifies the streaming delay when the channel type is **SRT_PULL**.

<a name="LiveChannel_FailoverConditions"></a>
The `failover_conditions` block supports:

* `input_loss_threshold_msec` - (Optional, Int) Specifies the duration threshold of inflow stop.
  When this threshold is reached, the active/standby switchover is automatically triggered. The unit is millisecond.
  Value ranges from `0` to `3,600,000`. Defaults to `2,000` ms.

* `input_preference` - (Optional, String) Specifies the input preference type. Valid values are:
  + **PRIMARY**: The main incoming URL is the first priority.
  + **EQUAL**: Equal switching between primary and backup URLs.

  Defaults to **EQUAL**.

  -> If equal switching is used and the backup URL is used, it will not automatically switch to the primary URL.

<a name="LiveChannel_AudioSelectors"></a>
The `audio_selectors` block supports:

* `name` - (Required, String) Specifies the name of the audio selector.

* `selector_settings` - (Optional, List) Specifies the audio selector configuration.
  The [selector_settings](#LiveChannel_SelectorSettings) structure is documented below.

<a name="LiveChannel_SelectorSettings"></a>
The `selector_settings` block supports:

* `audio_language_selection` - (Optional, List) Specifies the language selector configuration.
  The [audio_language_selection](#LiveChannel_AudioLanguageSelection) structure is documented below.

* `audio_pid_selection` - (Optional, List) Specifies the PID selector configuration.
  The [audio_pid_selection](#LiveChannel_AudioPidSelection) structure is documented below.

* `audio_hls_selection` - (Optional, List) Specifies the HLS selector configuration.
  The [audio_hls_selection](#LiveChannel_AudioHlsSelection) structure is documented below.

<a name="LiveChannel_AudioLanguageSelection"></a>
The `audio_language_selection` block supports:

* `language_code` - (Required, String) Specifies the language abbreviation. Supports `2` or `3` lowercase letter language
  codes.

* `language_selection_policy` - (Optional, String) Specifies the language output strategy. Valid values are:
  + **LOOSE**: Loose matching. For example, "eng" will prioritize matching tracks with English as the language in the
    source stream. If no match is found, the track with the smallest PID will be selected.
  + **STRICT**: Strict matching. For example, "eng" will strictly match the audio track in the source stream whose
    language is English. If no match is found, the media live broadcast service will automatically fill in a silent
    segment. When the terminal uses this audio selector to play the video, it will be played silently.

<a name="LiveChannel_AudioPidSelection"></a>
The `audio_pid_selection` block supports:

* `pid` - (Required, Int) Specifies the value of PID.

<a name="LiveChannel_AudioHlsSelection"></a>
The `audio_hls_selection` block supports:

* `name` - (Required, String) Specifies the HLS audio selector name.

* `group_id` - (Required, String) Specifies the HLS audio selector gid.

<a name="LiveChannel_RecordSettings"></a>
The `record_settings` block supports:

* `rollingbuffer_duration` - (Required, Int) Specifies the maximum playback recording time. During this time period,
  the recording will continue. The unit is second.
  When the value is `0`, it means that recording is not supported. The maximum supported recording period is `14` days.

<a name="LiveChannel_Endpoints"></a>
The `endpoints` block supports:

* `hls_package` - (Optional, List) Specifies the HLS packaging information.
  The [hls_package](#LiveChannel_HlsPackage) structure is documented below.

* `dash_package` - (Optional, List) Specifies the DASH packaging information.
  The [dash_package](#LiveChannel_DashPackage) structure is documented below.

* `mss_package` - (Optional, List) Specifies the MSS packaging information.
  The [mss_package](#LiveChannel_MssPackage) structure is documented below.

<a name="LiveChannel_HlsPackage"></a>
The `hls_package` block supports:

* `url` - (Required, String) Specifies the customer-defined streaming address, including method, domain name, and path.

* `stream_selection` - (Optional, List) Specifies the stream selection. Filter out the specified range of streams from
  the full stream.
  The [stream_selection](#LiveChannel_StreamSelection) structure is documented below.

* `hls_version` - (Optional, String) Specifies the HLS version.

* `segment_duration_seconds` - (Required, Int) Specifies the duration of the channel output segment. The unit is second.
  Value ranges from `1` to `10`.

  -> Modifying the segment duration will affect the time-shift and playback services of the recorded content, so please
  modify with caution!

* `playlist_window_seconds` - (Optional, Int) Specifies the window length of the channel live broadcast return shard.
  The value is the output segment duration multiplied by the number of segments. There are at least three returned segments.
  The unit is second. Value ranges from `0` to `86,400`.

* `encryption` - (Optional, List) Specifies the encrypted information.
  The [encryption](#LiveChannel_Encryption) structure is documented below.

* `request_args` - (Optional, List) Specifies the play related configuration.
  The [request_args](#LiveChannel_RequestArgs) structure is documented below.

* `ad_marker` - (Optional, List) Specifies the advertising marker. The HLS value is **ENHANCED_SCTE35**.

<a name="LiveChannel_DashPackage"></a>
The `dash_package` block supports:

* `url` - (Required, String) Specifies the customer-defined streaming address, including method, domain name, and path.

* `stream_selection` - (Optional, List) Specifies the stream selection. Filter out the specified range of streams from
  the full stream.
  The [stream_selection](#LiveChannel_StreamSelection) structure is documented below.

* `segment_duration_seconds` - (Required, Int) Specifies the duration of the channel output segment. The unit is second.
  Value ranges from `1` to `10`.

  -> Modifying the segment duration will affect the time-shift and playback services of the recorded content, so please
  modify with caution!

* `playlist_window_seconds` - (Optional, Int) Specifies the window length of the channel live broadcast return shard.
  The value is the output segment duration multiplied by the number of segments. There are at least three returned segments.
  The unit is second. Value ranges from `0` to `86,400`.

* `encryption` - (Optional, List) Specifies the encrypted information.
  The [encryption](#LiveChannel_Encryption) structure is documented below.

* `request_args` - (Optional, List) Specifies the play related configuration.
  The [request_args](#LiveChannel_RequestArgs) structure is documented below.

* `ad_marker` - (Optional, String) Specifies the advertising marker. The DASH value is **xml+bin**.

<a name="LiveChannel_MssPackage"></a>
The `mss_package` block supports:

* `url` - (Required, String) Specifies the customer-defined streaming address, including method, domain name, and path.

* `stream_selection` - (Optional, List) Specifies the stream selection. Filter out the specified range of streams from
  the full stream.
  The [stream_selection](#LiveChannel_StreamSelection) structure is documented below.

* `segment_duration_seconds` - (Required, Int) Specifies the duration of the channel output segment. The unit is second.
  Value ranges from `1` to `10`.

  -> Modifying the segment duration will affect the time-shift and playback services of the recorded content, so please
  modify with caution!

* `playlist_window_seconds` - (Optional, Int) Specifies the window length of the channel live broadcast return shard.
  The value is the output segment duration multiplied by the number of segments. There are at least three returned segments.
  The unit is second. Value ranges from `0` to `86,400`.

* `encryption` - (Optional, List) Specifies the encrypted information.
  The [encryption](#LiveChannel_Encryption) structure is documented below.

* `delay_segment` - (Optional, Int) Specifies the delayed playback time. The unit is second.

* `request_args` - (Optional, List) Specifies the play related configuration.
  The [request_args](#LiveChannel_RequestArgs) structure is documented below.

<a name="LiveChannel_StreamSelection"></a>
The `stream_selection` block supports:

* `key` - (Optional, String) Specifies the key used for bitrate filtering in streaming URLs.

* `max_bandwidth` - (Optional, Int) Specifies the maximum code rate. The unit is bps. Value ranges from `0` to `104,857,600`.

* `min_bandwidth` - (Optional, Int) Specifies the minimum code rate. The unit is bps. Value ranges from `0` to `104,857,600`.

<a name="LiveChannel_Encryption"></a>
The `encryption` block supports:

* `level` - (Optional, String) Specifies the level. Valid values are:
  + **content**: One channel corresponds to one key.
  + **profile**: One code rate corresponds to one key.

  Defaults to **content**.

* `resource_id` - (Required, String) Specifies the customer-generated DRM content ID.

* `system_ids` - (Required, List) Specifies the system ID enumeration values. Valid values are **FairPlay** (HLS),
  **Widevine** (DASH), **PlayReady** (DASH), and **PlayReady** (MSS).

* `url` - (Required, String) Specifies the DRM address of the key.

* `speke_version` - (Required, String) Specifies the DRM spec version number. Currently, only supports **1.0**.

* `request_mode` - (Required, String) Specifies the request mode. Valid values are:
  + **direct_http**: HTTP(S) direct access to DRM.
  + **functiongraph_proxy**: FunctionGraph proxy access to DRM.

* `http_headers` - (Optional, List) Specifies the authentication information that needs to be added to the DRM request header.
  Supports up to `5` configurations. Only the **direct_http** request mode supports configuring this field.
  The [http_headers](#LiveChannel_HttpHeader) structure is documented below.

* `urn` - (Optional, String) Specifies the URN of the function graph. The **functiongraph_proxy** request mode requires
  the function graph's urn to be provided.

<a name="LiveChannel_HttpHeader"></a>
The `http_headers` block supports:

* `key` - (Required, String) Specifies the key field name in the request header.

* `value` - (Required, String) Specifies the value corresponding to the key in the request header.

<a name="LiveChannel_RequestArgs"></a>
The `request_args` block supports:

* `record` - (Optional, List) Specifies the recording and playback related configuration.
  The [record](#LiveChannel_RequestArgsRecord) structure is documented below.

* `timeshift` - (Optional, List) Specifies the time-shift playback configuration.
  The [timeshift](#LiveChannel_RequestArgsTimeShift) structure is documented below.

* `live` - (Optional, List) Specifies the live broadcast configuration.
  The [live](#LiveChannel_RequestArgsLive) structure is documented below.

<a name="LiveChannel_RequestArgsRecord"></a>
The `record` block supports:

* `start_time` - (Optional, String) Specifies the start time.

* `end_time` - (Optional, String) Specifies the end time.

* `format` - (Optional, String) Specifies the format.

* `unit` - (Optional, String) Specifies the unit.

<a name="LiveChannel_RequestArgsTimeShift"></a>
The `timeshift` block supports:

* `back_time` - (Optional, String) Specifies the time shift duration field name.

* `unit` - (Optional, String) Specifies the unit.

<a name="LiveChannel_RequestArgsLive"></a>
The `live` block supports:

* `delay` - (Optional, String) Specifies the delay field.

* `unit` - (Optional, String) Specifies the unit.

<a name="LiveChannel_EncoderSettingsExpand"></a>
The `encoder_settings_expand` block supports:

* `audio_descriptions` - (Optional, List) Specifies the description of the audio output configuration.
  The [audio_descriptions](#LiveChannel_AudioDescriptions) structure is documented below.

<a name="LiveChannel_AudioDescriptions"></a>
The `audio_descriptions` block supports:

* `name` - (Required, String) Specifies the name of the audio output configuration. Only uppercase and lowercase letters,
  numbers, hyphens (-), and underscores (_) are supported.
  Different audio output configuration names for the same channel are not allowed to be duplicated.

* `audio_selector_name` - (Required, String) Specifies the audio selector name.

* `language_code_control` - (Optional, String) Specifies the language code control configuration.
  The settings here will not change the actual language of the audio, but only the language in which the audio is
  displayed externally. Valid values are:
  + **FOLLOW_INPUT**: If the output audio corresponding to the selected audio selector has a language, it will be
    consistent with it, otherwise it will be backed up by the language code and stream name configured here.
    The current option is recommended and is the default value.
  + **USE_CONFIGURED**: Users can customize the language and stream name of the output audio based on actual conditions.

* `language_code` - (Optional, String) Specifies the language code. The value could be `2` or `3` lowercase letters.

* `stream_name` - (Optional, String) Specifies the stream name.

<a name="LiveChannel_EncoderSettings"></a>
The `encoder_settings` block supports:

* `template_id` - (Optional, String) Specifies the transcoding template ID.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID (channel ID).

## Import

The live channel can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_live_channel.test <id>
```
