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

// @API LIVE GET /v1/{project_id}/record/rules
func DataSourceLiveRecordings() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLiveRecordingsRead,

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
				Description: `Specifies the ingest domain name to which the recording rules belong.`,
			},
			"app_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the application name of the recording rule.`,
			},
			"stream_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the stream name of the recording rule.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the recording type of the recording rule.`,
			},
			"rules": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of the recording rules.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The recording rule ID.`,
						},
						"domain_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ingest domain name to which the recording rule belongs.`,
						},
						"app_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The application name of the recording rule.`,
						},
						"stream_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The stream name of the recording rule.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The recording type of the recording rule.`,
						},
						"default_record_config": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The default recording configuration rule.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"record_format": {
										Type:        schema.TypeList,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: `The recording format.`,
									},
									"obs": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: `The OBS bucket information for storing recordings.`,
										Elem:        rulesDefaultRecordConfigObsElem(),
									},
									"hls": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: `The HLS configuration rule.`,
										Elem:        rulesDefaultRecordConfigHlsElem(),
									},
									"flv": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: `The FLV configuration rule.`,
										Elem:        rulesDefaultRecordConfigFlvElem(),
									},
									"mp4": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: `The MP4 configuration rule.`,
										Elem:        rulesDefaultRecordConfigMp4Elem(),
									},
								},
							},
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the recording rule.`,
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The lasted update time of the recording rule.`,
						},
					},
				},
			},
		},
	}
}

func rulesDefaultRecordConfigObsElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The OBS bucket name.`,
			},
			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The region to which the OBS bucket belongs.`,
			},
			"object": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The OBS object storage path.`,
			},
		},
	}
}

func rulesDefaultRecordConfigHlsElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"recording_length": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The periodic recording duration, in seconds.`,
			},
			"file_naming": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The file path and file name prefix of the recorded M3U8 file.`,
			},
			"ts_file_naming": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The file name prefix of recorded TS file.`,
			},
			"record_slice_duration": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The TS slicing duration during HLS recording, in seconds.`,
			},
			"max_stream_pause_length": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The recording HLS file concatenation duration, in seconds.`,
			},
		},
	}
}

func rulesDefaultRecordConfigFlvElem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"recording_length": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The periodic recording duration, in seconds.`,
			},
			"file_naming": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The file path and file name prefix of the recorded FLV file.`,
			},
			"max_stream_pause_length": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The recording FLV file concatenation duration, in seconds.`,
			},
		},
	}
}

func rulesDefaultRecordConfigMp4Elem() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"recording_length": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The periodic recording duration, in seconds.`,
			},
			"file_naming": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The file path and file name prefix of the recorded MP4 file.`,
			},
			"max_stream_pause_length": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The recording MP4 file concatenation duration, in seconds.`,
			},
		},
	}
}

func dataSourceLiveRecordingsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("live", region)
	if err != nil {
		return diag.Errorf("error creating Live client: %s", err)
	}

	recordings, err := queryLiveRecordings(d, client)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("rules", flattenLiveRecordings(recordings)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func queryLiveRecordings(d *schema.ResourceData, client *golangsdk.ServiceClient) ([]interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/record/rules?limit=100"
		page    = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath += buildRecordingsQueryParams(d)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		// The offset indicates the page number.
		// The default value is 0, which represents the first page.
		listPathWithPage := fmt.Sprintf("%s&offset=%d", listPath, page)
		requestResp, err := client.Request("GET", listPathWithPage, &opt)
		if err != nil {
			return nil, fmt.Errorf("error retrieving recordings: %s", err)
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		recording := utils.PathSearch("record_config", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, recording...)
		if len(recording) == 0 {
			break
		}
		page++
	}
	return result, nil
}

func buildRecordingsQueryParams(d *schema.ResourceData) string {
	res := ""
	if domainName, ok := d.GetOk("domain_name"); ok {
		res = fmt.Sprintf("%s&publish_domain=%v", res, domainName)
	}
	if appName, ok := d.GetOk("app_name"); ok {
		res = fmt.Sprintf("%s&app=%v", res, appName)
	}
	if streamName, ok := d.GetOk("stream_name"); ok {
		res = fmt.Sprintf("%s&stream=%v", res, streamName)
	}
	if recordType, ok := d.GetOk("type"); ok {
		res = fmt.Sprintf("%s&record_type=%v", res, recordType)
	}
	return res
}

func flattenLiveRecordings(recordings []interface{}) []map[string]interface{} {
	if len(recordings) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, len(recordings))
	for i, v := range recordings {
		result[i] = map[string]interface{}{
			"id":          utils.PathSearch("id", v, nil),
			"domain_name": utils.PathSearch("publish_domain", v, nil),
			"app_name":    utils.PathSearch("app", v, nil),
			"stream_name": utils.PathSearch("stream", v, nil),
			"type":        utils.PathSearch("record_type", v, nil),
			"created_at":  utils.PathSearch("create_time", v, nil),
			"updated_at":  utils.PathSearch("update_time", v, nil),
			"default_record_config": flattenDefaultRecordConfig(
				utils.PathSearch("default_record_config", v, make(map[string]interface{})).(map[string]interface{})),
		}
	}
	return result
}

func flattenDefaultRecordConfig(recordCofig interface{}) []map[string]interface{} {
	if recordCofig == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"record_format": utils.PathSearch("record_format", recordCofig, nil),
			"obs": flattenObsCofig(
				utils.PathSearch("obs_addr", recordCofig, make(map[string]interface{})).(map[string]interface{})),
			"hls": flattenHlsCofig(
				utils.PathSearch("hls_config", recordCofig, make(map[string]interface{})).(map[string]interface{})),
			"flv": flattenFlvCofig(
				utils.PathSearch("flv_config", recordCofig, make(map[string]interface{})).(map[string]interface{})),
			"mp4": flattenMp4Cofig(
				utils.PathSearch("mp4_config", recordCofig, make(map[string]interface{})).(map[string]interface{})),
		},
	}
}

func flattenObsCofig(obsConfig interface{}) []map[string]interface{} {
	if obsConfig == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"bucket": utils.PathSearch("bucket", obsConfig, nil),
			"region": utils.PathSearch("location", obsConfig, nil),
			"object": utils.PathSearch("object", obsConfig, nil),
		},
	}
}

func flattenHlsCofig(hlsConfig interface{}) []map[string]interface{} {
	if hlsConfig == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"recording_length":        utils.PathSearch("record_cycle", hlsConfig, nil),
			"file_naming":             utils.PathSearch("record_prefix", hlsConfig, nil),
			"ts_file_naming":          utils.PathSearch("record_ts_prefix", hlsConfig, nil),
			"record_slice_duration":   utils.PathSearch("record_slice_duration", hlsConfig, nil),
			"max_stream_pause_length": utils.PathSearch("record_max_duration_to_merge_file", hlsConfig, nil),
		},
	}
}

func flattenFlvCofig(flvConfig interface{}) []map[string]interface{} {
	if flvConfig == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"recording_length":        utils.PathSearch("record_cycle", flvConfig, nil),
			"file_naming":             utils.PathSearch("record_prefix", flvConfig, nil),
			"max_stream_pause_length": utils.PathSearch("record_max_duration_to_merge_file", flvConfig, nil),
		},
	}
}

func flattenMp4Cofig(mp4Config interface{}) []map[string]interface{} {
	if mp4Config == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"recording_length":        utils.PathSearch("record_cycle", mp4Config, nil),
			"file_naming":             utils.PathSearch("record_prefix", mp4Config, nil),
			"max_stream_pause_length": utils.PathSearch("record_max_duration_to_merge_file", mp4Config, nil),
		},
	}
}
