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

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API LIVE GET /v1/{project_id}/ott/channels
func DataSourceLiveChannels() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLiveChannelsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"domain_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the channel streaming domain name.`,
			},
			"app_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the group name or application name.`,
			},
			"channel_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the channel ID.`,
			},
			"channels": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The channel information.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The channel streaming domain name.`,
						},
						"app_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The group name or application name.`,
						},
						"state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The channel status.`,
						},
						"input": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        dataChannelInputSchema(),
							Description: `The channel input information.`,
						},
						"record_settings": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        dataChannelRecordSettingsSchema(),
							Description: `The configuration for replaying a recording.`,
						},
						"endpoints": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        dataChannelEndpointsSchema(),
							Description: `The channel outflow information.`,
						},
						"encoder_settings_expand": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        dataChannelEncoderSettingsExpandSchema(),
							Description: `The audio output configuration.`,
						},
						"encoder_settings": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        dataChannelEncoderSettingsSchema(),
							Description: `The transcoding template configuration.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The channel name.`,
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The channel ID.`,
						},
					},
				},
			},
		},
	}
}

func dataChannelEncoderSettingsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"template_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The transcoding template ID.`,
			},
		},
	}
	return &sc
}

func dataChannelEncoderSettingsExpandSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"audio_descriptions": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataChannelEncoderSettingsExpandAudioDescriptionsSchema(),
				Description: `The description of the audio output configuration.`,
			},
		},
	}
	return &sc
}

func dataChannelEncoderSettingsExpandAudioDescriptionsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the audio output configuration.`,
			},
			"audio_selector_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The audio selector name.`,
			},
			"language_code_control": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The language code control configuration.`,
			},
			"language_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The language code.`,
			},
			"stream_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The stream name.`,
			},
		},
	}
	return &sc
}

func dataChannelEndpointsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"hls_package": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataChannelEndpointsHlsPackageSchema(),
				Description: `The HLS packaging information.`,
			},
			"dash_package": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataChannelEndpointsDashPackageSchema(),
				Description: `The DASH packaging information.`,
			},
			"mss_package": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataChannelEndpointsMssPackageSchema(),
				Description: `The MSS packaging information.`,
			},
		},
	}
	return &sc
}

func dataChannelEndpointsMssPackageSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The customer-defined streaming address.`,
			},
			"stream_selection": {
				Type:        schema.TypeList,
				Elem:        dataChannelEndpointsStreamSelectionSchema(),
				Computed:    true,
				Description: `The stream selection.`,
			},
			"segment_duration_seconds": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The duration of the channel output segment.`,
			},
			"playlist_window_seconds": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The window length of the channel live broadcast return shard.`,
			},
			"encryption": {
				Type:        schema.TypeList,
				Elem:        dataChannelEndpointsEncryptionSchema(),
				Computed:    true,
				Description: `The encrypted information.`,
			},
			"delay_segment": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The delayed playback time.`,
			},
			"request_args": {
				Type:        schema.TypeList,
				Elem:        dataChannelEndpointsRequestArgsSchema(),
				Computed:    true,
				Description: `The play related configuration.`,
			},
		},
	}
	return &sc
}

func dataChannelEndpointsDashPackageSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The customer-defined streaming address.`,
			},
			"stream_selection": {
				Type:        schema.TypeList,
				Elem:        dataChannelEndpointsStreamSelectionSchema(),
				Computed:    true,
				Description: `The stream selection.`,
			},
			"segment_duration_seconds": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The duration of the channel output segment.`,
			},
			"playlist_window_seconds": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The window length of the channel live broadcast return shard.`,
			},
			"encryption": {
				Type:        schema.TypeList,
				Elem:        dataChannelEndpointsEncryptionSchema(),
				Computed:    true,
				Description: `The encrypted information.`,
			},
			"request_args": {
				Type:        schema.TypeList,
				Elem:        dataChannelEndpointsRequestArgsSchema(),
				Computed:    true,
				Description: `The play related configuration.`,
			},
			"ad_marker": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The advertising marker.`,
			},
		},
	}
	return &sc
}

func dataChannelEndpointsHlsPackageSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The customer-defined streaming address.`,
			},
			"segment_duration_seconds": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The duration of the channel output segment.`,
			},
			"stream_selection": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataChannelEndpointsStreamSelectionSchema(),
				Description: `The stream selection.`,
			},
			"hls_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The HLS version number.`,
			},
			"playlist_window_seconds": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The window length of the channel live broadcast return shard.`,
			},
			"encryption": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataChannelEndpointsEncryptionSchema(),
				Description: `The encrypted information.`,
			},
			"request_args": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataChannelEndpointsRequestArgsSchema(),
				Description: `The play related configuration.`,
			},
			"ad_marker": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The advertising marker.`,
			},
		},
	}
	return &sc
}

func dataChannelEndpointsRequestArgsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"record": {
				Type:        schema.TypeList,
				Elem:        dataChannelRequestArgsRecordSchema(),
				Computed:    true,
				Description: `The recording and playback related configuration.`,
			},
			"timeshift": {
				Type:        schema.TypeList,
				Elem:        dataChannelRequestArgsTimeShiftSchema(),
				Computed:    true,
				Description: `The time-shift playback configuration.`,
			},
			"live": {
				Type:        schema.TypeList,
				Elem:        dataChannelRequestArgsLiveSchema(),
				Computed:    true,
				Description: `The live broadcast configuration.`,
			},
		},
	}
	return &sc
}

func dataChannelRequestArgsLiveSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"delay": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The delay field.`,
			},
			"unit": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The unit.`,
			},
		},
	}
	return &sc
}

func dataChannelRequestArgsTimeShiftSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"back_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The time shift duration field name.`,
			},
			"unit": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The unit.`,
			},
		},
	}
	return &sc
}

func dataChannelRequestArgsRecordSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"start_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The start time.`,
			},
			"end_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The end time.`,
			},
			"format": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The format.`,
			},
			"unit": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The unit.`,
			},
		},
	}
	return &sc
}

func dataChannelEndpointsEncryptionSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"resource_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The customer-generated DRM content ID.`,
			},
			"system_ids": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The system ID enumeration values.`,
			},
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The DRM address of the key.`,
			},
			"speke_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The DRM spec version number.`,
			},
			"request_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The request mode.`,
			},
			"key_rotation_interval_seconds": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The key rotation interval seconds.`,
			},
			"encryption_method": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The encryption method.`,
			},
			"level": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The level.`,
			},
			"http_headers": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataChannelEncryptionHttpHeaderSchema(),
				Description: `The authentication information that needs to be added to the DRM request header.`,
			},
			"urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The URN of the function graph.`,
			},
		},
	}
	return &sc
}

func dataChannelEncryptionHttpHeaderSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The key field name in the request header.`,
			},
			"value": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The value corresponding to the key in the request header.`,
			},
		},
	}
	return &sc
}

func dataChannelEndpointsStreamSelectionSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The key used for bitrate filtering in streaming URLs.`,
			},
			"max_bandwidth": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The maximum code rate.`,
			},
			"min_bandwidth": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The minimum code rate.`,
			},
		},
	}
	return &sc
}

func dataChannelRecordSettingsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"rollingbuffer_duration": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The maximum playback recording time.`,
			},
		},
	}
	return &sc
}

func dataChannelInputSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"input_protocol": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The channel input protocol.`,
			},
			"sources": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataChannelInputSourcesSchema(),
				Description: `The channel main source stream information.`,
			},
			"secondary_sources": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataChannelInputSecondarySourcesSchema(),
				Description: `The prepared stream array.`,
			},
			"failover_conditions": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataChannelInputFailoverConditionsSchema(),
				Description: `The configuration of switching between primary and backup audio and video stream URLs.`,
			},
			"max_bandwidth_limit": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The maximum bandwidth that needs to be configured when the inbound protocol is **HLS_PULL**.`,
			},
			"ip_port_mode": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `The IP port mode.`,
			},
			"ip_whitelist": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The IP whitelist when protocol is **SRT_PUSH**.`,
			},
			"scte35_source": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The advertisement scte35 signal source.`,
			},
			"ad_triggers": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The ad trigger configuration.`,
			},
			"audio_selectors": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataChannelInputAudioSelectorsSchema(),
				Description: `The audio selector configuration.`,
			},
		},
	}
	return &sc
}

func dataChannelInputSourcesSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The channel source stream URL, used for external streaming.`,
			},
			"bitrate": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The bitrate.`,
			},
			"width": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The resolution corresponds to the width value.`,
			},
			"height": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The resolution corresponds to the high value.`,
			},
			"enable_snapshot": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether to use this stream to take screenshots.`,
			},
			"bitrate_for3u8": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether to use bitrate to fix the bitrate.`,
			},
			"passphrase": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The encrypted information when the protocol is **SRT_PUSH**.`,
			},
			"backup_urls": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of backup stream addresses.`,
			},
			"stream_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The stream ID of the stream pull address when the channel type is **SRT_PULL**.`,
			},
			"latency": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The streaming delay when the channel type is **SRT_PULL**.`,
			},
		},
	}
	return &sc
}

func dataChannelInputSecondarySourcesSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The channel source stream URL.`,
			},
			"bitrate": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The bitrate.`,
			},
			"width": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The resolution corresponds to the width value.`,
			},
			"height": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The resolution corresponds to the high value.`,
			},
			"bitrate_for3u8": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether to use bitrate to fix the bitrate.`,
			},
			"passphrase": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The encrypted information when the protocol is **SRT_PUSH**.`,
			},
			"backup_urls": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of backup stream addresses.`,
			},
			"stream_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The stream ID of the stream pull address when the channel type is **SRT_PULL**.`,
			},
			"latency": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The streaming delay when the channel type is **SRT_PULL**.`,
			},
		},
	}
	return &sc
}

func dataChannelInputFailoverConditionsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"input_loss_threshold_msec": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The duration threshold of inflow stop.`,
			},
			"input_preference": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The input preference type.`,
			},
		},
	}
	return &sc
}

func dataChannelInputAudioSelectorsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the audio selector.`,
			},
			"selector_settings": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataChannelAudioSelectorsSelectorSettingsSchema(),
				Description: `The audio selector configuration.`,
			},
		},
	}
	return &sc
}

func dataChannelAudioSelectorsSelectorSettingsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"audio_language_selection": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataChannelSelectorSettingsAudioLanguageSelectionSchema(),
				Description: `The language selector configuration.`,
			},
			"audio_pid_selection": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataChannelSelectorSettingsAudioPidSelectionSchema(),
				Description: `The PID selector configuration.`,
			},
			"audio_hls_selection": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataChannelSelectorSettingsAudioHlsSelectionSchema(),
				Description: `The HLS selector configuration.`,
			},
		},
	}
	return &sc
}

func dataChannelSelectorSettingsAudioLanguageSelectionSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"language_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The language abbreviation.`,
			},
			"language_selection_policy": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The language output strategy.`,
			},
		},
	}
	return &sc
}

func dataChannelSelectorSettingsAudioPidSelectionSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"pid": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The value of PID.`,
			},
		},
	}
	return &sc
}

func dataChannelSelectorSettingsAudioHlsSelectionSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The HLS audio selector name.`,
			},
			"group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The HLS audio selector gid.`,
			},
		},
	}
	return &sc
}

func buildDatasourceChannelQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("domain_name"); ok {
		res = fmt.Sprintf("%s&domain=%v", res, v)
	}

	if v, ok := d.GetOk("app_name"); ok {
		res = fmt.Sprintf("%s&app_name=%v", res, v)
	}

	if v, ok := d.GetOk("channel_id"); ok {
		res = fmt.Sprintf("%s&id=%v", res, v)
	}
	return res
}

func dataSourceLiveChannelsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		httpUrl       = "v1/{project_id}/ott/channels?limit=100"
		product       = "live"
		totalChannels []interface{}
		offset        = 0
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating Live client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildDatasourceChannelQueryParams(d)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	for {
		requestPathWithOffset := fmt.Sprintf("%s&offset=%d", requestPath, offset)
		resp, err := client.Request("GET", requestPathWithOffset, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving Live channels: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		channels := utils.PathSearch("channels", respBody, make([]interface{}, 0)).([]interface{})
		if len(channels) == 0 {
			break
		}

		totalChannels = append(totalChannels, channels...)
		offset += len(channels)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("channels", flattenDataSourceChannels(totalChannels)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDataSourceChannels(totalChannels []interface{}) []interface{} {
	if len(totalChannels) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(totalChannels))
	for _, v := range totalChannels {
		result = append(result, map[string]interface{}{
			"domain_name":             utils.PathSearch("domain", v, nil),
			"app_name":                utils.PathSearch("app_name", v, nil),
			"id":                      utils.PathSearch("id", v, nil),
			"name":                    utils.PathSearch("name", v, nil),
			"state":                   utils.PathSearch("state", v, nil),
			"input":                   flattenDataChannelInputResponseBody(v),
			"encoder_settings":        flattenDataChannelEncoderSettingsResponseBody(v),
			"record_settings":         flattenDataChannelRecordSettingsResponseBody(v),
			"endpoints":               flattenDataChannelEndpointsResponseBody(v),
			"encoder_settings_expand": flattenDataChannelEncoderSettingsExpandResponseBody(v),
		})
	}
	return result
}

func flattenDataChannelEncoderSettingsExpandResponseBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	respBody := utils.PathSearch("encoder_settings_expand", resp, nil)
	if respBody == nil {
		return nil
	}

	return []interface{}{map[string]interface{}{
		"audio_descriptions": flattenDataAudioDescriptionsResponseBody(respBody),
	}}
}

func flattenDataAudioDescriptionsResponseBody(resp interface{}) []interface{} {
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

func flattenDataChannelEndpointsResponseBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("endpoints", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"hls_package":  flattenDataHlsPackageResponseBody(v),
			"dash_package": flattenDataDashPackageResponseBody(v),
			"mss_package":  flattenDataMssPackageResponseBody(v),
		})
	}
	return rst
}

func flattenDataMssPackageResponseBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("mss_package", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"url":                      utils.PathSearch("url", v, nil),
			"stream_selection":         flattenDataStreamSelectionResponseBody(v),
			"segment_duration_seconds": utils.PathSearch("segment_duration_seconds", v, nil),
			"playlist_window_seconds":  utils.PathSearch("playlist_window_seconds", v, nil),
			"encryption":               flattenDataEncryptionResponseBody(v),
			"delay_segment":            utils.PathSearch("delay_segment", v, nil),
			"request_args":             flattenDataRequestArgsResponseBody(v),
		})
	}
	return rst
}

func flattenDataDashPackageResponseBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("dash_package", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"url":                      utils.PathSearch("url", v, nil),
			"stream_selection":         flattenDataStreamSelectionResponseBody(v),
			"segment_duration_seconds": utils.PathSearch("segment_duration_seconds", v, nil),
			"playlist_window_seconds":  utils.PathSearch("playlist_window_seconds", v, nil),
			"encryption":               flattenDataEncryptionResponseBody(v),
			"request_args":             flattenDataRequestArgsResponseBody(v),
			"ad_marker":                utils.PathSearch("ad_marker", v, nil),
		})
	}
	return rst
}

func flattenDataHlsPackageResponseBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("hls_package", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"url":                      utils.PathSearch("url", v, nil),
			"stream_selection":         flattenDataStreamSelectionResponseBody(v),
			"hls_version":              utils.PathSearch("hls_version", v, nil),
			"segment_duration_seconds": utils.PathSearch("segment_duration_seconds", v, nil),
			"playlist_window_seconds":  utils.PathSearch("playlist_window_seconds", v, nil),
			"encryption":               flattenDataEncryptionResponseBody(v),
			"request_args":             flattenDataRequestArgsResponseBody(v),
			"ad_marker":                utils.PathSearch("ad_marker", v, nil),
		})
	}
	return rst
}

func flattenDataRequestArgsResponseBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	respBody := utils.PathSearch("request_args", resp, nil)
	if respBody == nil {
		return nil
	}

	return []interface{}{map[string]interface{}{
		"record":    flattenDataRequestArgsRecordResponseBody(respBody),
		"timeshift": flattenDataRequestArgsTimeShiftResponseBody(respBody),
		"live":      flattenDataRequestArgsLiveResponseBody(respBody),
	}}
}

func flattenDataRequestArgsLiveResponseBody(resp interface{}) []interface{} {
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

func flattenDataRequestArgsTimeShiftResponseBody(resp interface{}) []interface{} {
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

func flattenDataRequestArgsRecordResponseBody(resp interface{}) []interface{} {
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

func flattenDataStreamSelectionResponseBody(resp interface{}) []interface{} {
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

func flattenDataEncryptionResponseBody(resp interface{}) []interface{} {
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
		"http_headers":                  flattenDataEncryptionHttpHeadersResponseBody(respBody),
		"urn":                           utils.PathSearch("urn", respBody, nil),
	}}
}

func flattenDataEncryptionHttpHeadersResponseBody(resp interface{}) []interface{} {
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

func flattenDataChannelRecordSettingsResponseBody(resp interface{}) []interface{} {
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

func flattenDataChannelEncoderSettingsResponseBody(resp interface{}) []interface{} {
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

func flattenDataChannelInputResponseBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	respBody := utils.PathSearch("input", resp, nil)
	if respBody == nil {
		return nil
	}

	return []interface{}{map[string]interface{}{
		"input_protocol":      utils.PathSearch("input_protocol", respBody, nil),
		"sources":             flattenDataInputSourcesResponseBody(respBody),
		"secondary_sources":   flattenDataInputSecondarySourcesResponseBody(respBody),
		"failover_conditions": flattenDataInputFailoverConditionsResponseBody(respBody),
		"max_bandwidth_limit": utils.PathSearch("max_bandwidth_limit", respBody, nil),
		"ip_port_mode":        utils.PathSearch("ip_port_mode", respBody, nil),
		"ip_whitelist":        utils.PathSearch("ip_whitelist", respBody, nil),
		"scte35_source":       utils.PathSearch("scte35_source", respBody, nil),
		"ad_triggers":         utils.PathSearch("ad_triggers", respBody, nil),
		"audio_selectors":     flattenDataInputAudioSelectorsResponseBody(respBody),
	}}
}

func flattenDataInputSourcesResponseBody(resp interface{}) []interface{} {
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

func flattenDataInputSecondarySourcesResponseBody(resp interface{}) []interface{} {
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

func flattenDataInputFailoverConditionsResponseBody(resp interface{}) []interface{} {
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

func flattenDataInputAudioSelectorsResponseBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("audio_selectors", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"name":              utils.PathSearch("name", v, nil),
			"selector_settings": flattenDataSelectorSettingsResponseBody(v),
		})
	}
	return rst
}

func flattenDataSelectorSettingsResponseBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	respBody := utils.PathSearch("selector_settings", resp, nil)
	if respBody == nil {
		return nil
	}

	return []interface{}{map[string]interface{}{
		"audio_language_selection": flattenDataAudioLanguageSelectionResponseBody(respBody),
		"audio_pid_selection":      flattenDataAudioPidSelectionResponseBody(respBody),
		"audio_hls_selection":      flattenDataAudioHlsSelectionResponseBody(respBody),
	}}
}

func flattenDataAudioLanguageSelectionResponseBody(resp interface{}) []interface{} {
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

func flattenDataAudioPidSelectionResponseBody(resp interface{}) []interface{} {
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

func flattenDataAudioHlsSelectionResponseBody(resp interface{}) []interface{} {
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
