package live

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Live POST /v1/{project_id}/ott/channels
// @API Live GET /v1/{project_id}/ott/channels
// @API Live DELETE /v1/{project_id}/ott/channels
// @API Live PUT /v1/{project_id}/ott/channels/endpoints
// @API Live PUT /v1/{project_id}/ott/channels/input
// @API Live PUT /v1/{project_id}/ott/channels/record-settings
// @API Live PUT /v1/{project_id}/ott/channels/general
// @API Live PUT /v1/{project_id}/ott/channels/state
// @API Live PUT /v1/{project_id}/ott/channels/encorder-settings
func ResourceChannel() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceChannelCreate,
		UpdateContext: resourceChannelUpdate,
		ReadContext:   resourceChannelRead,
		DeleteContext: resourceChannelDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"domain_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the channel streaming domain name.`,
			},
			"app_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the group name or application name.`,
			},
			"state": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the channel status.`,
			},
			"input": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Elem:        channelInputSchema(),
				Required:    true,
				Description: `Specifies the channel input information.`,
			},
			"record_settings": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Elem:        channelRecordSettingsSchema(),
				Required:    true,
				Description: `Specifies the configuration for replaying a recording.`,
			},
			"endpoints": {
				Type:        schema.TypeList,
				Elem:        channelEndpointsSchema(),
				Required:    true,
				Description: `Specifies the channel outflow information.`,
			},
			"encoder_settings_expand": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Elem:        channelEncoderSettingsExpandSchema(),
				Optional:    true,
				Description: `Specifies the audio output configuration.`,
			},
			"encoder_settings": {
				Type:        schema.TypeList,
				Elem:        channelEncoderSettingsSchema(),
				Optional:    true,
				Description: `Specifies the transcoding template configuration.`,
			},
			// This field can be edited to be empty, so no `Computed` attribute is added.
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the channel name.`,
			},
			// This is a custom field. Users can pass in a custom resource ID through this field.
			"channel_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: `Specifies the unique channel ID.`,
			},
		},
	}
}

func channelInputSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"input_protocol": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the channel input protocol.`,
			},
			"sources": {
				Type:        schema.TypeList,
				Elem:        channelInputSourcesSchema(),
				Optional:    true,
				Description: `Specifies the channel main source stream information.`,
			},
			"secondary_sources": {
				Type:        schema.TypeList,
				Elem:        channelInputSecondarySourcesSchema(),
				Optional:    true,
				Description: `Specifies the prepared stream array.`,
			},
			// This field has a default value, so add a Computed attribute.
			"failover_conditions": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem:     channelInputFailoverConditionsSchema(),
				Optional: true,
				Computed: true,
				Description: `Specifies the configuration of switching between primary and backup audio and video
stream URLs.`,
			},
			"max_bandwidth_limit": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				Description: `Specifies the maximum bandwidth that needs to be configured when the inbound protocol is
**HLS_PULL**.`,
			},
			"ip_port_mode": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the IP port mode.`,
			},
			"ip_whitelist": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the IP whitelist when protocol is **SRT_PUSH**.`,
			},
			"scte35_source": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the advertisement scte35 signal source.`,
			},
			"ad_triggers": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
				Description: `Specifies the ad trigger configuration.`,
			},
			"audio_selectors": {
				Type:        schema.TypeList,
				Elem:        channelInputAudioSelectorsSchema(),
				Optional:    true,
				Description: `Specifies the audio selector configuration.`,
			},
		},
	}
	return &sc
}

func channelInputSourcesSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			// In some cases, this field will be overwritten by the URL generated by the cloud, which will cause the
			// field change to fail. We describe this scenario in the document and do not make multi-function design.
			"url": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the channel source stream URL, used for external streaming.`,
			},
			"bitrate": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the bitrate.`,
			},
			"width": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the resolution corresponds to the width value.`,
			},
			"height": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the resolution corresponds to the high value.`,
			},
			"enable_snapshot": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Specifies whether to use this stream to take screenshots.`,
			},
			"bitrate_for3u8": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Specifies whether to use bitrate to fix the bitrate.`,
			},
			"passphrase": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the encrypted information when the protocol is **SRT_PUSH**.`,
			},
			"backup_urls": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
				Description: `Specifies the list of backup stream addresses.`,
			},
			"stream_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the stream ID of the stream pull address when the channel type is **SRT_PULL**.`,
			},
			"latency": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the streaming delay when the channel type is **SRT_PULL**.`,
			},
		},
	}
	return &sc
}

func channelInputSecondarySourcesSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the channel source stream URL.`,
			},
			"bitrate": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the bitrate.`,
			},
			"width": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the resolution corresponds to the width value.`,
			},
			"height": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the resolution corresponds to the high value.`,
			},
			"bitrate_for3u8": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Specifies whether to use bitrate to fix the bitrate.`,
			},
			"passphrase": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the encrypted information when the protocol is **SRT_PUSH**.`,
			},
			"backup_urls": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
				Description: `Specifies the list of backup stream addresses.`,
			},
			"stream_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the stream ID of the stream pull address when the channel type is **SRT_PULL**.`,
			},
			"latency": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the streaming delay when the channel type is **SRT_PULL**.`,
			},
		},
	}
	return &sc
}

func channelInputFailoverConditionsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"input_loss_threshold_msec": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the duration threshold of inflow stop.`,
			},
			"input_preference": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the input preference type.`,
			},
		},
	}
	return &sc
}

func channelInputAudioSelectorsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the name of the audio selector.`,
			},
			"selector_settings": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Elem:        channelAudioSelectorsSelectorSettingsSchema(),
				Optional:    true,
				Description: `Specifies the audio selector configuration.`,
			},
		},
	}
	return &sc
}

func channelAudioSelectorsSelectorSettingsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"audio_language_selection": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Elem:        channelSelectorSettingsAudioLanguageSelectionSchema(),
				Optional:    true,
				Description: `Specifies the language selector configuration.`,
			},
			"audio_pid_selection": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Elem:        channelSelectorSettingsAudioPidSelectionSchema(),
				Optional:    true,
				Description: `Specifies the PID selector configuration.`,
			},
			"audio_hls_selection": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Elem:        channelSelectorSettingsAudioHlsSelectionSchema(),
				Optional:    true,
				Description: `Specifies the HLS selector configuration.`,
			},
		},
	}
	return &sc
}

func channelSelectorSettingsAudioLanguageSelectionSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"language_code": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the language abbreviation.`,
			},
			"language_selection_policy": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the language output strategy.`,
			},
		},
	}
	return &sc
}

func channelSelectorSettingsAudioPidSelectionSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"pid": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Specifies the value of PID.`,
			},
		},
	}
	return &sc
}

func channelSelectorSettingsAudioHlsSelectionSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the HLS audio selector name.`,
			},
			"group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the HLS audio selector gid.`,
			},
		},
	}
	return &sc
}

func channelRecordSettingsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"rollingbuffer_duration": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Specifies the maximum playback recording time.`,
			},
		},
	}
	return &sc
}

func channelEndpointsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"hls_package": {
				Type:        schema.TypeList,
				Elem:        channelEndpointsHlsPackageSchema(),
				Optional:    true,
				Description: `Specifies the HLS packaging information.`,
			},
			"dash_package": {
				Type:        schema.TypeList,
				Elem:        channelEndpointsDashPackageSchema(),
				Optional:    true,
				Description: `Specifies the DASH packaging information.`,
			},
			"mss_package": {
				Type:        schema.TypeList,
				Elem:        channelEndpointsMssPackageSchema(),
				Optional:    true,
				Description: `Specifies the MSS packaging information.`,
			},
		},
	}
	return &sc
}

func channelEndpointsHlsPackageSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the customer-defined streaming address.`,
			},
			"segment_duration_seconds": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Specifies the duration of the channel output segment.`,
			},
			"stream_selection": {
				Type:        schema.TypeList,
				Elem:        channelEndpointsStreamSelectionSchema(),
				Optional:    true,
				Description: `Specifies the stream selection.`,
			},
			"hls_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the HLS version number.`,
			},
			"playlist_window_seconds": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the window length of the channel live broadcast return shard.`,
			},
			"encryption": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Elem:        channelEndpointsEncryptionSchema(),
				Optional:    true,
				Description: `Specifies the encrypted information.`,
			},
			"request_args": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Elem:        channelEndpointsRequestArgsSchema(),
				Optional:    true,
				Description: `Specifies the play related configuration.`,
			},
			"ad_marker": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: `Specifies the advertising marker.`,
			},
		},
	}
	return &sc
}

func channelEndpointsStreamSelectionSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the key used for bitrate filtering in streaming URLs.`,
			},
			"max_bandwidth": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the maximum code rate.`,
			},
			"min_bandwidth": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the minimum code rate.`,
			},
		},
	}
	return &sc
}

func channelEndpointsEncryptionSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"resource_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the customer-generated DRM content ID.`,
			},
			"system_ids": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				Description: `Specifies the system ID enumeration values.`,
			},
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the DRM address of the key.`,
			},
			"speke_version": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the DRM spec version number.`,
			},
			"request_mode": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the request mode.`,
			},
			// Currently this field is reserved and does not support configuration.
			"key_rotation_interval_seconds": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `schema: Internal; Specifies the key rotation interval seconds.`,
			},
			// Currently this field is reserved and does not support configuration.
			"encryption_method": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `schema: Internal; Specifies the encryption method.`,
			},
			"level": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the level.`,
			},
			"http_headers": {
				Type:        schema.TypeList,
				Elem:        channelEncryptionHttpHeaderSchema(),
				Optional:    true,
				Description: `Specifies the authentication information that needs to be added to the DRM request header.`,
			},
			// This field can be edited to be empty, so no computed is needed.
			"urn": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the URN of the function graph.`,
			},
		},
	}
	return &sc
}

func channelEncryptionHttpHeaderSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the key field name in the request header.`,
			},
			"value": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the value corresponding to the key in the request header.`,
			},
		},
	}
	return &sc
}

func channelEndpointsRequestArgsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"record": {
				Type:        schema.TypeList,
				Elem:        channelRequestArgsRecordSchema(),
				Optional:    true,
				Description: `Specifies the recording and playback related configuration.`,
			},
			"timeshift": {
				Type:        schema.TypeList,
				Elem:        channelRequestArgsTimeShiftSchema(),
				Optional:    true,
				Description: `Specifies the time-shift playback configuration.`,
			},
			"live": {
				Type:        schema.TypeList,
				Elem:        channelRequestArgsLiveSchema(),
				Optional:    true,
				Description: `Specifies the live broadcast configuration.`,
			},
		},
	}
	return &sc
}

func channelRequestArgsRecordSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"start_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the start time.`,
			},
			"end_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the end time.`,
			},
			"format": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the format.`,
			},
			"unit": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the unit.`,
			},
		},
	}
	return &sc
}

func channelRequestArgsTimeShiftSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"back_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the time shift duration field name.`,
			},
			"unit": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the unit.`,
			},
		},
	}
	return &sc
}

func channelRequestArgsLiveSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"delay": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the delay field.`,
			},
			"unit": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the unit.`,
			},
		},
	}
	return &sc
}

func channelEndpointsDashPackageSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the customer-defined streaming address.`,
			},
			"stream_selection": {
				Type:        schema.TypeList,
				Elem:        channelEndpointsStreamSelectionSchema(),
				Optional:    true,
				Description: `Specifies the stream selection.`,
			},
			"segment_duration_seconds": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Specifies the duration of the channel output segment.`,
			},
			"playlist_window_seconds": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the window length of the channel live broadcast return shard.`,
			},
			"encryption": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Elem:        channelEndpointsEncryptionSchema(),
				Optional:    true,
				Description: `Specifies the encrypted information.`,
			},
			"request_args": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Elem:        channelEndpointsRequestArgsSchema(),
				Optional:    true,
				Description: `Specifies the play related configuration.`,
			},
			"ad_marker": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the advertising marker.`,
			},
		},
	}
	return &sc
}

func channelEndpointsMssPackageSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the customer-defined streaming address.`,
			},
			"stream_selection": {
				Type:        schema.TypeList,
				Elem:        channelEndpointsStreamSelectionSchema(),
				Optional:    true,
				Description: `Specifies the stream selection.`,
			},
			"segment_duration_seconds": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Specifies the duration of the channel output segment.`,
			},
			"playlist_window_seconds": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the window length of the channel live broadcast return shard.`,
			},
			"encryption": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Elem:        channelEndpointsEncryptionSchema(),
				Optional:    true,
				Description: `Specifies the encrypted information.`,
			},
			"delay_segment": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the delayed playback time.`,
			},
			"request_args": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Elem:        channelEndpointsRequestArgsSchema(),
				Optional:    true,
				Description: `Specifies the play related configuration.`,
			},
		},
	}
	return &sc
}

func channelEncoderSettingsExpandSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"audio_descriptions": {
				Type:        schema.TypeList,
				Elem:        channelEncoderSettingsExpandAudioDescriptionsSchema(),
				Optional:    true,
				Description: `Specifies the description of the audio output configuration.`,
			},
		},
	}
	return &sc
}

func channelEncoderSettingsExpandAudioDescriptionsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the name of the audio output configuration.`,
			},
			"audio_selector_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the audio selector name.`,
			},
			"language_code_control": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the language code control configuration.`,
			},
			"language_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the language code.`,
			},
			"stream_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the stream name.`,
			},
		},
	}
	return &sc
}

func channelEncoderSettingsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"template_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the transcoding template ID.`,
			},
		},
	}
	return &sc
}

// First try to get the value from field `channel_id`. If the value is empty, generate a random ID as the resource ID.
func generateChannelResourceID(d *schema.ResourceData) (string, error) {
	if v, ok := d.GetOk("channel_id"); ok {
		return v.(string), nil
	}

	return uuid.GenerateUUID()
}

func resourceChannelCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/ott/channels"
		product = "live"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating Live client: %s", err)
	}

	resourceID, err := generateChannelResourceID(d)
	if err != nil {
		return diag.Errorf("error generating Live channel resource ID: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(buildCreateChannelBodyParams(d, resourceID)),
	}

	_, err = client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating Live channel: %s", err)
	}

	d.SetId(resourceID)
	return resourceChannelRead(ctx, d, meta)
}

func buildCreateChannelBodyParams(d *schema.ResourceData, resourceID string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"domain":                  d.Get("domain_name"),
		"app_name":                d.Get("app_name"),
		"state":                   d.Get("state"),
		"input":                   buildChannelInputRequestBody(d.Get("input")),
		"record_settings":         buildChannelRecordSettingsRequestBody(d.Get("record_settings")),
		"endpoints":               buildChannelEndpointsRequestBody(d.Get("endpoints")),
		"encoder_settings_expand": buildChannelEncoderSettingsExpandRequestBody(d.Get("encoder_settings_expand")),
		"encoder_settings":        buildChannelEncoderSettingsRequestBody(d.Get("encoder_settings")),
		"name":                    utils.ValueIgnoreEmpty(d.Get("name")),
		"id":                      resourceID,
	}
	return bodyParams
}

func buildChannelInputRequestBody(rawParams interface{}) map[string]interface{} {
	rawArray, ok := rawParams.([]interface{})
	if !ok {
		return nil
	}

	if len(rawArray) == 0 {
		return nil
	}

	raw := rawArray[0].(map[string]interface{})
	return map[string]interface{}{
		"input_protocol":      utils.ValueIgnoreEmpty(raw["input_protocol"]),
		"sources":             buildInputSourcesRequestBody(raw["sources"]),
		"secondary_sources":   buildInputSecondarySourcesRequestBody(raw["secondary_sources"]),
		"failover_conditions": buildInputFailoverConditionsRequestBody(raw["failover_conditions"]),
		"max_bandwidth_limit": raw["max_bandwidth_limit"],
		"ip_port_mode":        utils.ValueIgnoreEmpty(raw["ip_port_mode"]),
		"ip_whitelist":        utils.ValueIgnoreEmpty(raw["ip_whitelist"]),
		"scte35_source":       utils.ValueIgnoreEmpty(raw["scte35_source"]),
		"ad_triggers":         utils.ValueIgnoreEmpty(raw["ad_triggers"]),
		"audio_selectors":     buildInputAudioSelectorsRequestBody(raw["audio_selectors"]),
	}
}

func buildInputSourcesRequestBody(rawParams interface{}) []map[string]interface{} {
	rawArray, ok := rawParams.([]interface{})
	if !ok {
		return nil
	}

	if len(rawArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, len(rawArray))
	for i, v := range rawArray {
		raw := v.(map[string]interface{})
		rst[i] = map[string]interface{}{
			"url":             utils.ValueIgnoreEmpty(raw["url"]),
			"bitrate":         raw["bitrate"],
			"width":           raw["width"],
			"height":          raw["height"],
			"enable_snapshot": utils.ValueIgnoreEmpty(raw["enable_snapshot"]),
			"bitrate_for3u8":  utils.ValueIgnoreEmpty(raw["bitrate_for3u8"]),
			"passphrase":      utils.ValueIgnoreEmpty(raw["passphrase"]),
			"backup_urls":     utils.ValueIgnoreEmpty(raw["backup_urls"]),
			"stream_id":       utils.ValueIgnoreEmpty(raw["stream_id"]),
			"latency":         raw["latency"],
		}
	}
	return rst
}

func buildInputSecondarySourcesRequestBody(rawParams interface{}) []map[string]interface{} {
	rawArray, ok := rawParams.([]interface{})
	if !ok {
		return nil
	}

	if len(rawArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, len(rawArray))
	for i, v := range rawArray {
		raw := v.(map[string]interface{})
		rst[i] = map[string]interface{}{
			"url":            utils.ValueIgnoreEmpty(raw["url"]),
			"bitrate":        raw["bitrate"],
			"width":          raw["width"],
			"height":         raw["height"],
			"bitrate_for3u8": utils.ValueIgnoreEmpty(raw["bitrate_for3u8"]),
			"passphrase":     utils.ValueIgnoreEmpty(raw["passphrase"]),
			"backup_urls":    utils.ValueIgnoreEmpty(raw["backup_urls"]),
			"stream_id":      utils.ValueIgnoreEmpty(raw["stream_id"]),
			"latency":        raw["latency"],
		}
	}
	return rst
}

func buildInputFailoverConditionsRequestBody(rawParams interface{}) map[string]interface{} {
	rawArray, ok := rawParams.([]interface{})
	if !ok {
		return nil
	}

	if len(rawArray) == 0 {
		return nil
	}

	raw := rawArray[0].(map[string]interface{})
	return map[string]interface{}{
		"input_loss_threshold_msec": raw["input_loss_threshold_msec"],
		"input_preference":          utils.ValueIgnoreEmpty(raw["input_preference"]),
	}
}

func buildInputAudioSelectorsRequestBody(rawParams interface{}) []map[string]interface{} {
	rawArray, ok := rawParams.([]interface{})
	if !ok {
		return nil
	}

	if len(rawArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, len(rawArray))
	for i, v := range rawArray {
		raw := v.(map[string]interface{})
		rst[i] = map[string]interface{}{
			"name":              utils.ValueIgnoreEmpty(raw["name"]),
			"selector_settings": buildAudioSelectorsSelectorSettingsRequestBody(raw["selector_settings"]),
		}
	}
	return rst
}

func buildAudioSelectorsSelectorSettingsRequestBody(rawParams interface{}) map[string]interface{} {
	rawArray, ok := rawParams.([]interface{})
	if !ok {
		return nil
	}

	if len(rawArray) == 0 {
		return nil
	}

	raw := rawArray[0].(map[string]interface{})
	return map[string]interface{}{
		"audio_language_selection": buildAudioLanguageSelectionRequestBody(raw["audio_language_selection"]),
		"audio_pid_selection":      buildAudioPidSelectionRequestBody(raw["audio_pid_selection"]),
		"audio_hls_selection":      buildAudioHlsSelectionRequestBody(raw["audio_hls_selection"]),
	}
}

func buildAudioLanguageSelectionRequestBody(rawParams interface{}) map[string]interface{} {
	rawArray, ok := rawParams.([]interface{})
	if !ok {
		return nil
	}

	if len(rawArray) == 0 {
		return nil
	}

	raw := rawArray[0].(map[string]interface{})
	return map[string]interface{}{
		"language_code":             utils.ValueIgnoreEmpty(raw["language_code"]),
		"language_selection_policy": utils.ValueIgnoreEmpty(raw["language_selection_policy"]),
	}
}

func buildAudioPidSelectionRequestBody(rawParams interface{}) map[string]interface{} {
	rawArray, ok := rawParams.([]interface{})
	if !ok {
		return nil
	}

	if len(rawArray) == 0 {
		return nil
	}

	raw := rawArray[0].(map[string]interface{})
	return map[string]interface{}{
		"pid": raw["pid"],
	}
}

func buildAudioHlsSelectionRequestBody(rawParams interface{}) map[string]interface{} {
	rawArray, ok := rawParams.([]interface{})
	if !ok {
		return nil
	}

	if len(rawArray) == 0 {
		return nil
	}

	raw := rawArray[0].(map[string]interface{})
	return map[string]interface{}{
		"name":     utils.ValueIgnoreEmpty(raw["name"]),
		"group_id": utils.ValueIgnoreEmpty(raw["group_id"]),
	}
}

func buildChannelRecordSettingsRequestBody(rawParams interface{}) map[string]interface{} {
	rawArray, ok := rawParams.([]interface{})
	if !ok {
		return nil
	}

	if len(rawArray) == 0 {
		return nil
	}

	raw := rawArray[0].(map[string]interface{})
	return map[string]interface{}{
		"rollingbuffer_duration": raw["rollingbuffer_duration"],
	}
}

func buildChannelEndpointsRequestBody(rawParams interface{}) []map[string]interface{} {
	rawArray, ok := rawParams.([]interface{})
	if !ok {
		return nil
	}

	if len(rawArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, len(rawArray))
	for i, v := range rawArray {
		raw := v.(map[string]interface{})
		rst[i] = map[string]interface{}{
			"hls_package":  buildEndpointsHlsPackageRequestBody(raw["hls_package"]),
			"dash_package": buildEndpointsDashPackageRequestBody(raw["dash_package"]),
			"mss_package":  buildEndpointsMssPackageRequestBody(raw["mss_package"]),
		}
	}
	return rst
}

func buildEndpointsHlsPackageRequestBody(rawParams interface{}) []map[string]interface{} {
	rawArray, ok := rawParams.([]interface{})
	if !ok {
		return nil
	}

	if len(rawArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, len(rawArray))
	for i, v := range rawArray {
		raw := v.(map[string]interface{})
		rst[i] = map[string]interface{}{
			"url":                      utils.ValueIgnoreEmpty(raw["url"]),
			"stream_selection":         buildStreamSelectionRequestBody(raw["stream_selection"]),
			"hls_version":              utils.ValueIgnoreEmpty(raw["hls_version"]),
			"segment_duration_seconds": raw["segment_duration_seconds"],
			"playlist_window_seconds":  raw["playlist_window_seconds"],
			"encryption":               buildEncryptionRequestBody(raw["encryption"]),
			"request_args":             buildRequestArgsRequestBody(raw["request_args"]),
			"ad_marker":                utils.ValueIgnoreEmpty(raw["ad_marker"]),
		}
	}
	return rst
}

func buildStreamSelectionRequestBody(rawParams interface{}) []map[string]interface{} {
	rawArray, ok := rawParams.([]interface{})
	if !ok {
		return nil
	}

	if len(rawArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, len(rawArray))
	for i, v := range rawArray {
		raw := v.(map[string]interface{})
		rst[i] = map[string]interface{}{
			"key":           utils.ValueIgnoreEmpty(raw["key"]),
			"max_bandwidth": raw["max_bandwidth"],
			"min_bandwidth": raw["min_bandwidth"],
		}
	}
	return rst
}

func buildEncryptionRequestBody(rawParams interface{}) map[string]interface{} {
	rawArray, ok := rawParams.([]interface{})
	if !ok {
		return nil
	}

	if len(rawArray) == 0 {
		return nil
	}

	raw := rawArray[0].(map[string]interface{})
	return map[string]interface{}{
		"key_rotation_interval_seconds": raw["key_rotation_interval_seconds"],
		"encryption_method":             utils.ValueIgnoreEmpty(raw["encryption_method"]),
		"level":                         utils.ValueIgnoreEmpty(raw["level"]),
		"resource_id":                   utils.ValueIgnoreEmpty(raw["resource_id"]),
		"system_ids":                    utils.ValueIgnoreEmpty(raw["system_ids"]),
		"url":                           utils.ValueIgnoreEmpty(raw["url"]),
		"speke_version":                 utils.ValueIgnoreEmpty(raw["speke_version"]),
		"request_mode":                  utils.ValueIgnoreEmpty(raw["request_mode"]),
		"http_headers":                  buildEncryptionHttpHeaderRequestBody(raw["http_headers"]),
		"urn":                           utils.ValueIgnoreEmpty(raw["urn"]),
	}
}

func buildEncryptionHttpHeaderRequestBody(rawParams interface{}) []map[string]interface{} {
	rawArray, ok := rawParams.([]interface{})
	if !ok {
		return nil
	}

	if len(rawArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, len(rawArray))
	for i, v := range rawArray {
		raw := v.(map[string]interface{})
		rst[i] = map[string]interface{}{
			"key":   utils.ValueIgnoreEmpty(raw["key"]),
			"value": utils.ValueIgnoreEmpty(raw["value"]),
		}
	}
	return rst
}

func buildRequestArgsRequestBody(rawParams interface{}) map[string]interface{} {
	rawArray, ok := rawParams.([]interface{})
	if !ok {
		return nil
	}

	if len(rawArray) == 0 {
		return nil
	}

	raw := rawArray[0].(map[string]interface{})
	return map[string]interface{}{
		"record":    buildRequestArgsRecordRequestBody(raw["record"]),
		"timeshift": buildRequestArgsTimeShiftRequestBody(raw["timeshift"]),
		"live":      buildRequestArgsLiveRequestBody(raw["live"]),
	}
}

func buildRequestArgsRecordRequestBody(rawParams interface{}) []map[string]interface{} {
	rawArray, ok := rawParams.([]interface{})
	if !ok {
		return nil
	}

	if len(rawArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, len(rawArray))
	for i, v := range rawArray {
		raw := v.(map[string]interface{})
		rst[i] = map[string]interface{}{
			"start_time": utils.ValueIgnoreEmpty(raw["start_time"]),
			"end_time":   utils.ValueIgnoreEmpty(raw["end_time"]),
			"format":     utils.ValueIgnoreEmpty(raw["format"]),
			"unit":       utils.ValueIgnoreEmpty(raw["unit"]),
		}
	}
	return rst
}

func buildRequestArgsTimeShiftRequestBody(rawParams interface{}) []map[string]interface{} {
	rawArray, ok := rawParams.([]interface{})
	if !ok {
		return nil
	}

	if len(rawArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, len(rawArray))
	for i, v := range rawArray {
		raw := v.(map[string]interface{})
		rst[i] = map[string]interface{}{
			"back_time": utils.ValueIgnoreEmpty(raw["back_time"]),
			"unit":      utils.ValueIgnoreEmpty(raw["unit"]),
		}
	}
	return rst
}

func buildRequestArgsLiveRequestBody(rawParams interface{}) []map[string]interface{} {
	rawArray, ok := rawParams.([]interface{})
	if !ok {
		return nil
	}

	if len(rawArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, len(rawArray))
	for i, v := range rawArray {
		raw := v.(map[string]interface{})
		rst[i] = map[string]interface{}{
			"delay": utils.ValueIgnoreEmpty(raw["delay"]),
			"unit":  utils.ValueIgnoreEmpty(raw["unit"]),
		}
	}
	return rst
}

func buildEndpointsDashPackageRequestBody(rawParams interface{}) []map[string]interface{} {
	rawArray, ok := rawParams.([]interface{})
	if !ok {
		return nil
	}

	if len(rawArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, len(rawArray))
	for i, v := range rawArray {
		raw := v.(map[string]interface{})
		rst[i] = map[string]interface{}{
			"url":                      utils.ValueIgnoreEmpty(raw["url"]),
			"stream_selection":         buildStreamSelectionRequestBody(raw["stream_selection"]),
			"segment_duration_seconds": raw["segment_duration_seconds"],
			"playlist_window_seconds":  raw["playlist_window_seconds"],
			"encryption":               buildEncryptionRequestBody(raw["encryption"]),
			"request_args":             buildRequestArgsRequestBody(raw["request_args"]),
			"ad_marker":                utils.ValueIgnoreEmpty(raw["ad_marker"]),
		}
	}
	return rst
}

func buildEndpointsMssPackageRequestBody(rawParams interface{}) []map[string]interface{} {
	rawArray, ok := rawParams.([]interface{})
	if !ok {
		return nil
	}

	if len(rawArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, len(rawArray))
	for i, v := range rawArray {
		raw := v.(map[string]interface{})
		rst[i] = map[string]interface{}{
			"url":                      utils.ValueIgnoreEmpty(raw["url"]),
			"stream_selection":         buildStreamSelectionRequestBody(raw["stream_selection"]),
			"segment_duration_seconds": raw["segment_duration_seconds"],
			"playlist_window_seconds":  raw["playlist_window_seconds"],
			"encryption":               buildEncryptionRequestBody(raw["encryption"]),
			// This field does not support a value of `0`, so it needs to be ignored.
			"delay_segment": utils.ValueIgnoreEmpty(raw["delay_segment"]),
			"request_args":  buildRequestArgsRequestBody(raw["request_args"]),
		}
	}
	return rst
}

func buildChannelEncoderSettingsExpandRequestBody(rawParams interface{}) map[string]interface{} {
	rawArray, ok := rawParams.([]interface{})
	if !ok {
		return nil
	}

	if len(rawArray) == 0 {
		return nil
	}

	raw := rawArray[0].(map[string]interface{})
	return map[string]interface{}{
		"audio_descriptions": buildAudioDescriptionsRequestBody(raw["audio_descriptions"]),
	}
}

func buildAudioDescriptionsRequestBody(rawParams interface{}) []map[string]interface{} {
	rawArray, ok := rawParams.([]interface{})
	if !ok {
		return nil
	}

	if len(rawArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, len(rawArray))
	for i, v := range rawArray {
		raw := v.(map[string]interface{})
		rst[i] = map[string]interface{}{
			"name":                  utils.ValueIgnoreEmpty(raw["name"]),
			"audio_selector_name":   utils.ValueIgnoreEmpty(raw["audio_selector_name"]),
			"language_code_control": utils.ValueIgnoreEmpty(raw["language_code_control"]),
			"language_code":         utils.ValueIgnoreEmpty(raw["language_code"]),
			"stream_name":           utils.ValueIgnoreEmpty(raw["stream_name"]),
		}
	}
	return rst
}

func buildChannelEncoderSettingsRequestBody(rawParams interface{}) []map[string]interface{} {
	rawArray, ok := rawParams.([]interface{})
	if !ok {
		return nil
	}

	if len(rawArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, len(rawArray))
	for i, v := range rawArray {
		raw := v.(map[string]interface{})
		rst[i] = map[string]interface{}{
			"template_id": utils.ValueIgnoreEmpty(raw["template_id"]),
		}
	}
	return rst
}

func buildChannelQueryParams(d *schema.ResourceData) string {
	return fmt.Sprintf("?id=%s", d.Id())
}

func resourceChannelRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		mErr    *multierror.Error
		httpUrl = "v1/{project_id}/ott/channels"
		product = "live"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating Live client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildChannelQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving Live channel: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	// When the target resource does not exist, the `channels` field of the response is empty.
	channelBody := utils.PathSearch("channels|[0]", respBody, nil)
	if channelBody == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("domain_name", utils.PathSearch("domain", channelBody, nil)),
		d.Set("app_name", utils.PathSearch("app_name", channelBody, nil)),
		d.Set("channel_id", utils.PathSearch("id", channelBody, nil)),
		d.Set("name", utils.PathSearch("name", channelBody, nil)),
		d.Set("state", utils.PathSearch("state", channelBody, nil)),
		d.Set("input", flattenChannelInputResponseBody(channelBody)),
		d.Set("encoder_settings", flattenChannelEncoderSettingsResponseBody(channelBody)),
		d.Set("record_settings", flattenChannelRecordSettingsResponseBody(channelBody)),
		d.Set("endpoints", flattenChannelEndpointsResponseBody(channelBody)),
		d.Set("encoder_settings_expand", flattenChannelEncoderSettingsExpandResponseBody(channelBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenChannelInputResponseBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	respBody := utils.PathSearch("input", resp, nil)
	if respBody == nil {
		return nil
	}

	return []interface{}{map[string]interface{}{
		"input_protocol":      utils.PathSearch("input_protocol", respBody, nil),
		"sources":             flattenInputSourcesResponseBody(respBody),
		"secondary_sources":   flattenInputSecondarySourcesResponseBody(respBody),
		"failover_conditions": flattenInputFailoverConditionsResponseBody(respBody),
		"max_bandwidth_limit": utils.PathSearch("max_bandwidth_limit", respBody, nil),
		"ip_port_mode":        utils.PathSearch("ip_port_mode", respBody, nil),
		"ip_whitelist":        utils.PathSearch("ip_whitelist", respBody, nil),
		"scte35_source":       utils.PathSearch("scte35_source", respBody, nil),
		"ad_triggers":         utils.PathSearch("ad_triggers", respBody, nil),
		"audio_selectors":     flattenInputAudioSelectorsResponseBody(respBody),
	}}
}

func flattenInputSourcesResponseBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("sources", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"url":             utils.PathSearch("url", v, nil),
			"bitrate":         utils.PathSearch("bitrate", v, nil),
			"width":           utils.PathSearch("width", v, nil),
			"height":          utils.PathSearch("height", v, nil),
			"enable_snapshot": utils.PathSearch("enable_snapshot", v, nil),
			"bitrate_for3u8":  utils.PathSearch("bitrate_for3u8", v, nil),
			"passphrase":      utils.PathSearch("passphrase", v, nil),
			"backup_urls":     utils.PathSearch("backup_urls", v, nil),
			"stream_id":       utils.PathSearch("stream_id", v, nil),
			"latency":         utils.PathSearch("latency", v, nil),
		})
	}
	return rst
}

func flattenInputSecondarySourcesResponseBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("secondary_sources", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"url":            utils.PathSearch("url", v, nil),
			"bitrate":        utils.PathSearch("bitrate", v, nil),
			"width":          utils.PathSearch("width", v, nil),
			"height":         utils.PathSearch("height", v, nil),
			"bitrate_for3u8": utils.PathSearch("bitrate_for3u8", v, nil),
			"passphrase":     utils.PathSearch("passphrase", v, nil),
			"backup_urls":    utils.PathSearch("backup_urls", v, nil),
			"stream_id":      utils.PathSearch("stream_id", v, nil),
			"latency":        utils.PathSearch("latency", v, nil),
		})
	}
	return rst
}

func flattenInputFailoverConditionsResponseBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	respBody := utils.PathSearch("failover_conditions", resp, nil)
	if respBody == nil {
		return nil
	}

	return []interface{}{map[string]interface{}{
		"input_loss_threshold_msec": utils.PathSearch("input_loss_threshold_msec", respBody, nil),
		"input_preference":          utils.PathSearch("input_preference", respBody, nil),
	}}
}

func flattenInputAudioSelectorsResponseBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("audio_selectors", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"name":              utils.PathSearch("name", v, nil),
			"selector_settings": flattenSelectorSettingsResponseBody(v),
		})
	}
	return rst
}

func flattenSelectorSettingsResponseBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	respBody := utils.PathSearch("selector_settings", resp, nil)
	if respBody == nil {
		return nil
	}

	return []interface{}{map[string]interface{}{
		"audio_language_selection": flattenAudioLanguageSelectionResponseBody(respBody),
		"audio_pid_selection":      flattenAudioPidSelectionResponseBody(respBody),
		"audio_hls_selection":      flattenAudioHlsSelectionResponseBody(respBody),
	}}
}

func flattenAudioLanguageSelectionResponseBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	respBody := utils.PathSearch("audio_language_selection", resp, nil)
	if respBody == nil {
		return nil
	}

	return []interface{}{map[string]interface{}{
		"language_code":             utils.PathSearch("language_code", respBody, nil),
		"language_selection_policy": utils.PathSearch("language_selection_policy", respBody, nil),
	}}
}

func flattenAudioPidSelectionResponseBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	respBody := utils.PathSearch("audio_pid_selection", resp, nil)
	if respBody == nil {
		return nil
	}

	return []interface{}{map[string]interface{}{
		"pid": utils.PathSearch("pid", respBody, nil),
	}}
}

func flattenAudioHlsSelectionResponseBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	respBody := utils.PathSearch("audio_hls_selection", resp, nil)
	if respBody == nil {
		return nil
	}

	return []interface{}{map[string]interface{}{
		"name":     utils.PathSearch("name", respBody, nil),
		"group_id": utils.PathSearch("group_id", respBody, nil),
	}}
}

func flattenChannelEncoderSettingsResponseBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("encoder_settings", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"template_id": utils.PathSearch("template_id", v, nil),
		})
	}
	return rst
}

func flattenChannelRecordSettingsResponseBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	respBody := utils.PathSearch("record_settings", resp, nil)
	if respBody == nil {
		return nil
	}

	return []interface{}{map[string]interface{}{
		"rollingbuffer_duration": utils.PathSearch("rollingbuffer_duration", respBody, nil),
	}}
}

func flattenChannelEndpointsResponseBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("endpoints", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"hls_package":  flattenHlsPackageResponseBody(v),
			"dash_package": flattenDashPackageResponseBody(v),
			"mss_package":  flattenMssPackageResponseBody(v),
		})
	}
	return rst
}

func flattenHlsPackageResponseBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("hls_package", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"url":                      utils.PathSearch("url", v, nil),
			"stream_selection":         flattenStreamSelectionResponseBody(v),
			"hls_version":              utils.PathSearch("hls_version", v, nil),
			"segment_duration_seconds": utils.PathSearch("segment_duration_seconds", v, nil),
			"playlist_window_seconds":  utils.PathSearch("playlist_window_seconds", v, nil),
			"encryption":               flattenEncryptionResponseBody(v),
			"request_args":             flattenRequestArgsResponseBody(v),
			"ad_marker":                utils.PathSearch("ad_marker", v, nil),
		})
	}
	return rst
}

func flattenStreamSelectionResponseBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("stream_selection", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"key":           utils.PathSearch("key", v, nil),
			"max_bandwidth": utils.PathSearch("max_bandwidth", v, nil),
			"min_bandwidth": utils.PathSearch("min_bandwidth", v, nil),
		})
	}
	return rst
}

func flattenEncryptionResponseBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	respBody := utils.PathSearch("encryption", resp, nil)
	if respBody == nil {
		return nil
	}

	return []interface{}{map[string]interface{}{
		"key_rotation_interval_seconds": utils.PathSearch("key_rotation_interval_seconds", respBody, nil),
		"encryption_method":             utils.PathSearch("encryption_method", respBody, nil),
		"level":                         utils.PathSearch("level", respBody, nil),
		"resource_id":                   utils.PathSearch("resource_id", respBody, nil),
		"system_ids":                    utils.PathSearch("system_ids", respBody, nil),
		"url":                           utils.PathSearch("url", respBody, nil),
		"speke_version":                 utils.PathSearch("speke_version", respBody, nil),
		"request_mode":                  utils.PathSearch("request_mode", respBody, nil),
		"http_headers":                  flattenEncryptionHttpHeadersResponseBody(respBody),
		"urn":                           utils.PathSearch("urn", respBody, nil),
	}}
}

func flattenEncryptionHttpHeadersResponseBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("http_headers", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"key":   utils.PathSearch("key", v, nil),
			"value": utils.PathSearch("value", v, nil),
		})
	}
	return rst
}

func flattenRequestArgsResponseBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	respBody := utils.PathSearch("request_args", resp, nil)
	if respBody == nil {
		return nil
	}

	return []interface{}{map[string]interface{}{
		"record":    flattenRequestArgsRecordResponseBody(respBody),
		"timeshift": flattenRequestArgsTimeShiftResponseBody(respBody),
		"live":      flattenRequestArgsLiveResponseBody(respBody),
	}}
}

func flattenRequestArgsRecordResponseBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("record", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"start_time": utils.PathSearch("start_time", v, nil),
			"end_time":   utils.PathSearch("end_time", v, nil),
			"format":     utils.PathSearch("format", v, nil),
			"unit":       utils.PathSearch("unit", v, nil),
		})
	}
	return rst
}

func flattenRequestArgsTimeShiftResponseBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("timeshift", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"back_time": utils.PathSearch("back_time", v, nil),
			"unit":      utils.PathSearch("unit", v, nil),
		})
	}
	return rst
}

func flattenRequestArgsLiveResponseBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("live", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"delay": utils.PathSearch("delay", v, nil),
			"unit":  utils.PathSearch("unit", v, nil),
		})
	}
	return rst
}

func flattenDashPackageResponseBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("dash_package", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"url":                      utils.PathSearch("url", v, nil),
			"stream_selection":         flattenStreamSelectionResponseBody(v),
			"segment_duration_seconds": utils.PathSearch("segment_duration_seconds", v, nil),
			"playlist_window_seconds":  utils.PathSearch("playlist_window_seconds", v, nil),
			"encryption":               flattenEncryptionResponseBody(v),
			"request_args":             flattenRequestArgsResponseBody(v),
			"ad_marker":                utils.PathSearch("ad_marker", v, nil),
		})
	}
	return rst
}

func flattenMssPackageResponseBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("mss_package", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"url":                      utils.PathSearch("url", v, nil),
			"stream_selection":         flattenStreamSelectionResponseBody(v),
			"segment_duration_seconds": utils.PathSearch("segment_duration_seconds", v, nil),
			"playlist_window_seconds":  utils.PathSearch("playlist_window_seconds", v, nil),
			"encryption":               flattenEncryptionResponseBody(v),
			"delay_segment":            utils.PathSearch("delay_segment", v, nil),
			"request_args":             flattenRequestArgsResponseBody(v),
		})
	}
	return rst
}

func flattenChannelEncoderSettingsExpandResponseBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	respBody := utils.PathSearch("encoder_settings_expand", resp, nil)
	if respBody == nil {
		return nil
	}

	return []interface{}{map[string]interface{}{
		"audio_descriptions": flattenAudioDescriptionsResponseBody(respBody),
	}}
}

func flattenAudioDescriptionsResponseBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("audio_descriptions", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"name":                  utils.PathSearch("name", v, nil),
			"audio_selector_name":   utils.PathSearch("audio_selector_name", v, nil),
			"language_code_control": utils.PathSearch("language_code_control", v, nil),
			"language_code":         utils.PathSearch("language_code", v, nil),
			"stream_name":           utils.PathSearch("stream_name", v, nil),
		})
	}
	return rst
}

func resourceChannelUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("live", region)
	if err != nil {
		return diag.Errorf("error creating Live client: %s", err)
	}

	if d.HasChanges("encoder_settings", "encoder_settings_expand") {
		if err := updateChannelEncorderSettings(client, d); err != nil {
			return diag.Errorf("error updating Live channel encorder settings configuration: %s", err)
		}
	}

	if d.HasChange("endpoints") {
		if err := updateChannelEndpoints(client, d); err != nil {
			return diag.Errorf("error updating Live channel endpoints configuration: %s", err)
		}
	}

	if d.HasChange("name") {
		if err := updateChannelGeneral(client, d); err != nil {
			return diag.Errorf("error updating Live channel general configuration: %s", err)
		}
	}

	if d.HasChange("input") {
		if err := updateChannelInput(client, d); err != nil {
			return diag.Errorf("error updating Live channel input configuration: %s", err)
		}
	}

	if d.HasChange("record_settings") {
		if err := updateChannelRecordSettings(client, d); err != nil {
			return diag.Errorf("error updating Live channel record settings configuration: %s", err)
		}
	}

	if d.HasChange("state") {
		state := d.Get("state").(string)
		if err := updateChannelState(client, d, state); err != nil {
			return diag.Errorf("error updating Live channel state: %s", err)
		}
	}
	return resourceChannelRead(ctx, d, meta)
}

func updateChannelEncorderSettings(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	requestPath := client.Endpoint + "v1/{project_id}/ott/channels/encorder-settings"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(buildUpdateChannelEncorderSettingsBodyParams(d)),
	}

	_, err := client.Request("PUT", requestPath, &requestOpt)
	return err
}

func buildUpdateChannelEncorderSettingsBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"domain":                  d.Get("domain_name"),
		"app_name":                d.Get("app_name"),
		"id":                      d.Id(),
		"encoder_settings":        buildChannelEncoderSettingsRequestBody(d.Get("encoder_settings")),
		"encoder_settings_expand": buildChannelEncoderSettingsExpandRequestBody(d.Get("encoder_settings_expand")),
	}
	return bodyParams
}

func updateChannelEndpoints(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	requestPath := client.Endpoint + "v1/{project_id}/ott/channels/endpoints"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(buildUpdateChannelEndpointsBodyParams(d)),
	}

	_, err := client.Request("PUT", requestPath, &requestOpt)
	return err
}

func buildUpdateChannelEndpointsBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"domain":    d.Get("domain_name"),
		"app_name":  d.Get("app_name"),
		"id":        d.Id(),
		"endpoints": buildChannelEndpointsRequestBody(d.Get("endpoints")),
	}
	return bodyParams
}

func updateChannelGeneral(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	requestPath := client.Endpoint + "v1/{project_id}/ott/channels/general"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(buildUpdateChannelGeneralBodyParams(d)),
	}

	_, err := client.Request("PUT", requestPath, &requestOpt)
	return err
}

func buildUpdateChannelGeneralBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"domain":   d.Get("domain_name"),
		"app_name": d.Get("app_name"),
		"id":       d.Id(),
		"name":     utils.ValueIgnoreEmpty(d.Get("name")),
	}
	return bodyParams
}

func updateChannelInput(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	requestPath := client.Endpoint + "v1/{project_id}/ott/channels/input"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(buildUpdateChannelInputBodyParams(d)),
	}

	_, err := client.Request("PUT", requestPath, &requestOpt)
	return err
}

func buildUpdateChannelInputBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"domain":   d.Get("domain_name"),
		"app_name": d.Get("app_name"),
		"id":       d.Id(),
		"input":    buildChannelInputRequestBody(d.Get("input")),
	}
	return bodyParams
}

func updateChannelRecordSettings(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	requestPath := client.Endpoint + "v1/{project_id}/ott/channels/record-settings"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(buildUpdateChannelRecordSettingsBodyParams(d)),
	}

	_, err := client.Request("PUT", requestPath, &requestOpt)
	return err
}

func buildUpdateChannelRecordSettingsBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"domain":          d.Get("domain_name"),
		"app_name":        d.Get("app_name"),
		"id":              d.Id(),
		"record_settings": buildChannelRecordSettingsRequestBody(d.Get("record_settings")),
	}
	return bodyParams
}

func updateChannelState(client *golangsdk.ServiceClient, d *schema.ResourceData, state string) error {
	requestPath := client.Endpoint + "v1/{project_id}/ott/channels/state"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(buildUpdateChannelStateBodyParams(d, state)),
	}

	_, err := client.Request("PUT", requestPath, &requestOpt)
	return err
}

func buildUpdateChannelStateBodyParams(d *schema.ResourceData, state string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"domain":   d.Get("domain_name"),
		"app_name": d.Get("app_name"),
		"id":       d.Id(),
		"state":    state,
	}
	return bodyParams
}

func buildDeleteChannelQueryParams(d *schema.ResourceData) string {
	return fmt.Sprintf("?domain=%[1]s&app_name=%[2]s&id=%[3]s",
		d.Get("domain_name").(string),
		d.Get("app_name").(string),
		d.Id())
}

func resourceChannelDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/ott/channels"
		product = "live"
		state   = d.Get("state").(string)
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating Live client: %s", err)
	}

	if state == "ON" {
		// Close the channel before deleting.
		if err := updateChannelState(client, d, "OFF"); err != nil {
			return diag.Errorf("error closing Live channel state in delete operation: %s", err)
		}
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildDeleteChannelQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	// If any query parameter of the delete API does not exist, the API will respond with a 404 error.
	_, err = client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting Live channel")
	}

	return nil
}
