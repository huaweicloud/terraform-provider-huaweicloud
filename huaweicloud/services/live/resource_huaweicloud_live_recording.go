package live

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Live POST /v1/{project_id}/record/rules
// @API Live GET /v1/{project_id}/record/rules/{id}
// @API Live PUT /v1/{project_id}/record/rules/{id}
// @API Live DELETE /v1/{project_id}/record/rules/{id}
func ResourceRecording() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRecordingCreate,
		ReadContext:   resourceRecordingRead,
		UpdateContext: resourceRecordingUpdate,
		DeleteContext: resourceRecordingDelete,
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
				Type:     schema.TypeString,
				Required: true,
			},
			"app_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"stream_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"obs": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"region": {
							Type:     schema.TypeString,
							Required: true,
						},
						"bucket": {
							Type:     schema.TypeString,
							Required: true,
						},
						"object": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"hls": {
				Type:         schema.TypeList,
				Optional:     true,
				MaxItems:     1,
				AtLeastOneOf: []string{"hls", "flv", "mp4"},
				Elem:         formatHlsSchema(),
			},
			"flv": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     formatFlvAndMp4Schema(),
			},
			"mp4": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     formatFlvAndMp4Schema(),
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
		},
	}
}

func formatHlsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"recording_length": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"file_naming": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ts_file_naming": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"record_slice_duration": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"max_stream_pause_length": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func formatFlvAndMp4Schema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"recording_length": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"file_naming": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"max_stream_pause_length": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func buildCreateOrUpdateRecordBodyParams(d *schema.ResourceData) map[string]interface{} {
	recordFormat := make([]string, 0)
	if _, ok := d.GetOk("hls"); ok {
		recordFormat = append(recordFormat, "HLS")
	}
	if _, ok := d.GetOk("flv"); ok {
		recordFormat = append(recordFormat, "FLV")
	}
	if _, ok := d.GetOk("mp4"); ok {
		recordFormat = append(recordFormat, "MP4")
	}

	return map[string]interface{}{
		"publish_domain": d.Get("domain_name"),
		"app":            d.Get("app_name"),
		"stream":         d.Get("stream_name"),
		"record_type":    utils.ValueIgnoreEmpty(d.Get("type")),
		"default_record_config": map[string]interface{}{
			"record_format": recordFormat,
			"obs_addr":      buildObsAddrBodyParams(d.Get("obs").([]interface{})),
			"hls_config":    buildHlsConfigBodyParams(d.Get("hls").([]interface{})),
			"flv_config":    buildFlvConfigBodyParams(d.Get("flv").([]interface{})),
			"mp4_config":    buildMp4ConfigBodyParams(d.Get("mp4").([]interface{})),
		},
	}
}

func buildObsAddrBodyParams(obsAddr []interface{}) map[string]interface{} {
	if len(obsAddr) == 0 {
		return nil
	}

	rawObs := obsAddr[0].(map[string]interface{})
	obsParams := map[string]interface{}{
		"bucket":   rawObs["bucket"],
		"location": rawObs["region"],
		"object":   utils.ValueIgnoreEmpty(rawObs["object"]),
	}

	return obsParams
}

func buildHlsConfigBodyParams(hlsConfig []interface{}) map[string]interface{} {
	if len(hlsConfig) == 0 {
		return nil
	}

	rawHls := hlsConfig[0].(map[string]interface{})
	hlsParams := map[string]interface{}{
		"record_cycle":                      rawHls["recording_length"].(int) * 60,
		"record_prefix":                     utils.ValueIgnoreEmpty(rawHls["file_naming"]),
		"record_ts_prefix":                  utils.ValueIgnoreEmpty(rawHls["ts_file_naming"]),
		"record_slice_duration":             utils.ValueIgnoreEmpty(rawHls["record_slice_duration"]),
		"record_max_duration_to_merge_file": utils.ValueIgnoreEmpty(rawHls["max_stream_pause_length"]),
	}

	return hlsParams
}

func buildFlvConfigBodyParams(flvConfig []interface{}) map[string]interface{} {
	if len(flvConfig) == 0 {
		return nil
	}

	rawFlv := flvConfig[0].(map[string]interface{})
	flvParams := map[string]interface{}{
		"record_cycle":                      rawFlv["recording_length"].(int) * 60,
		"record_prefix":                     utils.ValueIgnoreEmpty(rawFlv["file_naming"]),
		"record_max_duration_to_merge_file": utils.ValueIgnoreEmpty(rawFlv["max_stream_pause_length"]),
	}

	return flvParams
}

func buildMp4ConfigBodyParams(mp4Config []interface{}) map[string]interface{} {
	if len(mp4Config) == 0 {
		return nil
	}

	rawMp4 := mp4Config[0].(map[string]interface{})
	mp4Params := map[string]interface{}{
		"record_cycle":                      rawMp4["recording_length"].(int) * 60,
		"record_prefix":                     utils.ValueIgnoreEmpty(rawMp4["file_naming"]),
		"record_max_duration_to_merge_file": utils.ValueIgnoreEmpty(rawMp4["max_stream_pause_length"]),
	}

	return mp4Params
}

func resourceRecordingCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/record/rules"
	)

	client, err := cfg.NewServiceClient("live", region)
	if err != nil {
		return diag.Errorf("error creating Live client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateOrUpdateRecordBodyParams(d)),
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating Live recording: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	ruleId := utils.PathSearch("id", respBody, "").(string)
	if ruleId == "" {
		return diag.Errorf("error creating Live recording: ID is not found in API response")
	}
	d.SetId(ruleId)

	return resourceRecordingRead(ctx, d, meta)
}

func resourceRecordingRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/record/rules/{id}"
	)

	client, err := cfg.NewServiceClient("live", region)
	if err != nil {
		return diag.Errorf("error creating Live client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{id}", d.Id())

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Live recording")
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("domain_name", utils.PathSearch("publish_domain", getRespBody, nil)),
		d.Set("app_name", utils.PathSearch("app", getRespBody, nil)),
		d.Set("stream_name", utils.PathSearch("stream", getRespBody, nil)),
		d.Set("type", utils.PathSearch("record_type", getRespBody, nil)),
		d.Set("obs", flattenObsConfig(utils.PathSearch("default_record_config.obs_addr", getRespBody, nil))),
		d.Set("hls", flattenHlsConfig(utils.PathSearch("default_record_config.hls_config", getRespBody, nil))),
		d.Set("flv", flattenFlvConfig(utils.PathSearch("default_record_config.flv_config", getRespBody, nil))),
		d.Set("mp4", flattenMp4Config(utils.PathSearch("default_record_config.mp4_config", getRespBody, nil))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenObsConfig(obsConfig interface{}) []map[string]interface{} {
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

func flattenHlsConfig(hlsConfig interface{}) []map[string]interface{} {
	recordPrefix := utils.PathSearch("record_prefix", hlsConfig, "").(string)

	if hlsConfig != nil && recordPrefix != "" {
		return []map[string]interface{}{
			{
				"recording_length":        int(utils.PathSearch("record_cycle", hlsConfig, float64(0)).(float64)) / 60,
				"file_naming":             utils.PathSearch("record_prefix", hlsConfig, nil),
				"ts_file_naming":          utils.PathSearch("record_ts_prefix", hlsConfig, nil),
				"record_slice_duration":   utils.PathSearch("record_slice_duration", hlsConfig, nil),
				"max_stream_pause_length": utils.PathSearch("record_max_duration_to_merge_file", hlsConfig, nil),
			},
		}
	}

	return nil
}

func flattenFlvConfig(flvConfig interface{}) []map[string]interface{} {
	recordPrefix := utils.PathSearch("record_prefix", flvConfig, "").(string)

	if flvConfig != nil && recordPrefix != "" {
		return []map[string]interface{}{
			{
				"recording_length":        int(utils.PathSearch("record_cycle", flvConfig, float64(0)).(float64)) / 60,
				"file_naming":             utils.PathSearch("record_prefix", flvConfig, nil),
				"max_stream_pause_length": utils.PathSearch("record_max_duration_to_merge_file", flvConfig, nil),
			},
		}
	}

	return nil
}

func flattenMp4Config(mp4Config interface{}) []map[string]interface{} {
	recordPrefix := utils.PathSearch("record_prefix", mp4Config, "").(string)

	if mp4Config != nil && recordPrefix != "" {
		return []map[string]interface{}{
			{
				"recording_length":        int(utils.PathSearch("record_cycle", mp4Config, float64(0)).(float64)) / 60,
				"file_naming":             utils.PathSearch("record_prefix", mp4Config, nil),
				"max_stream_pause_length": utils.PathSearch("record_max_duration_to_merge_file", mp4Config, nil),
			},
		}
	}

	return nil
}

func resourceRecordingUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/record/rules/{id}"
	)

	client, err := cfg.NewServiceClient("live", region)
	if err != nil {
		return diag.Errorf("error creating Live client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{id}", d.Id())
	updateOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateOrUpdateRecordBodyParams(d)),
	}

	_, err = client.Request("PUT", updatePath, &updateOpts)
	if err != nil {
		return diag.Errorf("error updating Live recording: %s", err)
	}

	return resourceRecordingRead(ctx, d, meta)
}

func resourceRecordingDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/record/rules/{id}"
	)

	client, err := cfg.NewServiceClient("live", region)
	if err != nil {
		return diag.Errorf("error creating Live client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{id}", d.Id())
	deleteOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpts)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting Live recording")
	}

	return nil
}
