package vod

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API VOD POST /v1.0/{project_id}/asset/template_group/transcodings
// @API VOD GET /v1.0/{project_id}/asset/template_group/transcodings
// @API VOD PUT /v1.0/{project_id}/asset/template_group/transcodings
// @API VOD DELETE /v1.0/{project_id}/asset/template_group/transcodings
func ResourceTranscodingTemplateGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTranscodingTemplateGroupCreate,
		ReadContext:   resourceTranscodingTemplateGroupRead,
		UpdateContext: resourceTranscodingTemplateGroupUpdate,
		DeleteContext: resourceTranscodingTemplateGroupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		//request and response parameters
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"auto_encrypt": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"is_default": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"watermark_template_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"low_bitrate_hd": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"audio_codec": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"AAC", "HEAAC1",
				}, false),
			},
			"video_codec": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"H264", "H265",
				}, false),
			},
			"hls_segment_duration": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntInSlice([]int{2, 3, 5, 10}),
			},
			"quality_info": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"output_format": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"MP4", "DASH", "DASH_HLS", "MP3", "ADTS", "HLS",
							}, false),
						},
						"video": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"quality": {
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: validation.StringInSlice([]string{
											"4K", "2K", "FHD", "SD", "LD", "HD",
										}, false),
									},
									"width": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"height": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"bitrate": {
										Type:     schema.TypeInt,
										Optional: true,
									},
									"frame_rate": {
										Type:     schema.TypeInt,
										Optional: true,
									},
								},
							},
						},
						"audio": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"sample_rate": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"channels": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"bitrate": {
										Type:     schema.TypeInt,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildCommonPvcParams(d *schema.ResourceData) string {
	if d.Get("low_bitrate_hd").(bool) {
		return "1"
	}

	return "0"
}

func buildTemplateGroupCommonBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"hls_interval": d.Get("hls_segment_duration"),
		"pvc":          buildCommonPvcParams(d),
		"video_codec":  utils.ValueIgnoreEmpty(d.Get("video_codec")),
		"audio_codec":  utils.ValueIgnoreEmpty(d.Get("audio_codec")),
	}
}

func buildQualityInfoListVideoQualityParams(quality string) string {
	if quality == "FHD" {
		return "FULL_HD"
	}

	if quality == "LD" {
		return "FLUENT"
	}

	return quality
}

func buildQualityInfoListVideoParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap := rawArray[0].(map[string]interface{})
	return map[string]interface{}{
		"quality":    buildQualityInfoListVideoQualityParams(rawMap["quality"].(string)),
		"width":      rawMap["width"],
		"height":     rawMap["height"],
		"bitrate":    rawMap["bitrate"],
		"frame_rate": rawMap["frame_rate"],
	}
}

func buildQualityInfoListAudioParams(rawArray []interface{}) map[string]interface{} {
	if len(rawArray) == 0 {
		return nil
	}

	rawMap := rawArray[0].(map[string]interface{})
	return map[string]interface{}{
		"sample_rate": rawMap["sample_rate"],
		"channels":    rawMap["channels"],
		"bitrate":     rawMap["bitrate"],
	}
}

func buildQualityInfoListBodyParams(d *schema.ResourceData) []map[string]interface{} {
	rawArray := d.Get("quality_info").([]interface{})
	if len(rawArray) == 0 {
		return nil
	}

	rst := make([]map[string]interface{}, 0, len(rawArray))
	for _, v := range rawArray {
		rawMap := v.(map[string]interface{})
		rst = append(rst, map[string]interface{}{
			"format": rawMap["output_format"],
			"video":  buildQualityInfoListVideoParams(rawMap["video"].([]interface{})),
			"audio":  buildQualityInfoListAudioParams(rawMap["audio"].([]interface{})),
		})
	}

	return rst
}

func buildAutoEncryptParam(d *schema.ResourceData) interface{} {
	if d.Get("auto_encrypt").(bool) {
		return 1
	}

	return nil
}

func buildStatusParam(d *schema.ResourceData) interface{} {
	if d.Get("is_default").(bool) {
		return "1"
	}

	return nil
}

func buildCreateTemplateGroupBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name":                   d.Get("name"),
		"type":                   "custom_template_group",
		"description":            d.Get("description"),
		"watermark_template_ids": d.Get("watermark_template_ids"),
		"common":                 buildTemplateGroupCommonBodyParams(d),
		"quality_info_list":      buildQualityInfoListBodyParams(d),
		"auto_encrypt":           buildAutoEncryptParam(d),
		"status":                 buildStatusParam(d),
	}
}

func resourceTranscodingTemplateGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1.0/{project_id}/asset/template_group/transcodings"
		product = "vod"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating VOD client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateTemplateGroupBodyParams(d)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating VOD transcoding template group: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	groupID := utils.PathSearch("group_id", respBody, "").(string)
	if groupID == "" {
		return diag.Errorf("error creating VOD transcoding template group: ID is not found in API response")
	}
	d.SetId(groupID)

	return resourceTranscodingTemplateGroupRead(ctx, d, meta)
}

func flattenVideoQualityAttribute(quality string) string {
	if quality == "FULL_HD" {
		return "FHD"
	}

	if quality == "FLUENT" {
		return "LD"
	}

	return quality
}

func flattenVideoAttribute(respBody interface{}) []interface{} {
	rawMap := utils.PathSearch("video", respBody, nil)
	if rawMap == nil {
		return nil
	}

	videoResult := map[string]interface{}{
		"quality":    flattenVideoQualityAttribute(utils.PathSearch("quality", rawMap, "").(string)),
		"bitrate":    utils.PathSearch("bitrate", rawMap, nil),
		"frame_rate": utils.PathSearch("frame_rate", rawMap, nil),
		"width":      utils.PathSearch("width", rawMap, nil),
		"height":     utils.PathSearch("height", rawMap, nil),
	}

	return []interface{}{videoResult}
}

func flattenAudioAttribute(respBody interface{}) []interface{} {
	rawMap := utils.PathSearch("audio", respBody, nil)
	if rawMap == nil {
		return nil
	}

	audioResult := map[string]interface{}{
		"sample_rate": utils.PathSearch("sample_rate", rawMap, nil),
		"channels":    utils.PathSearch("channels", rawMap, nil),
		"bitrate":     utils.PathSearch("bitrate", rawMap, nil),
	}

	return []interface{}{audioResult}
}

func flattenQualityInfoAttribute(respBody interface{}) []interface{} {
	rawArray := utils.PathSearch("quality_info_list", respBody, make([]interface{}, 0)).([]interface{})
	if len(rawArray) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(rawArray))
	for _, v := range rawArray {
		rst = append(rst, map[string]interface{}{
			"output_format": utils.PathSearch("format", v, nil),
			"video":         flattenVideoAttribute(v),
			"audio":         flattenAudioAttribute(v),
		})
	}

	return rst
}

func flattenVideoCodecAttribute(commonResp interface{}) string {
	rawString := utils.PathSearch("video_codec", commonResp, "").(string)
	if rawString == "H265" || rawString == "H264" {
		return rawString
	}

	return ""
}

func setCommonAttributes(d *schema.ResourceData, meta interface{}) error {
	if meta == nil {
		return nil
	}

	mErr := multierror.Append(
		d.Set("low_bitrate_hd", utils.PathSearch("pvc", meta, "").(string) == "1"),
		d.Set("video_codec", flattenVideoCodecAttribute(meta)),
		d.Set("audio_codec", utils.PathSearch("audio_codec", meta, nil)),
		d.Set("hls_segment_duration", utils.PathSearch("hls_interval", meta, nil)),
	)

	return mErr
}

func resourceTranscodingTemplateGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1.0/{project_id}/asset/template_group/transcodings"
		product = "vod"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating VOD client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath += fmt.Sprintf("?group_id=%s", d.Id())
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "error_code", "VOD.10053"),
			"error retrieving VOD transcoding template group")
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	templateGroup := utils.PathSearch("template_group_list|[0]", respBody, nil)
	if templateGroup == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", templateGroup, nil)),
		d.Set("auto_encrypt", int(utils.PathSearch("auto_encrypt", templateGroup, float64(0)).(float64)) == 1),
		d.Set("is_default", utils.PathSearch("status", templateGroup, "").(string) == "1"),
		d.Set("description", utils.PathSearch("description", templateGroup, nil)),
		d.Set("type", utils.PathSearch("type", templateGroup, nil)),
		d.Set("watermark_template_ids", utils.PathSearch("watermark_template_ids", templateGroup, nil)),
		d.Set("quality_info", flattenQualityInfoAttribute(templateGroup)),
		setCommonAttributes(d, utils.PathSearch("common", templateGroup, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildUpdateStatusParam(d *schema.ResourceData) string {
	if d.Get("is_default").(bool) {
		return "1"
	}

	return "0"
}

func buildUpdateTemplateGroupBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"group_id":               d.Id(),
		"name":                   d.Get("name"),
		"description":            d.Get("description"),
		"watermark_template_ids": d.Get("watermark_template_ids"),
		"common":                 buildTemplateGroupCommonBodyParams(d),
		"quality_info_list":      buildQualityInfoListBodyParams(d),
		"auto_encrypt":           buildAutoEncryptParam(d),
		"status":                 buildUpdateStatusParam(d),
	}
}

func resourceTranscodingTemplateGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1.0/{project_id}/asset/template_group/transcodings"
		product = "vod"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating VOD client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes:          []int{200, 201, 204},
		JSONBody:         utils.RemoveNil(buildUpdateTemplateGroupBodyParams(d)),
	}

	_, err = client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error updating VOD transcoding template group: %s", err)
	}

	return resourceTranscodingTemplateGroupRead(ctx, d, meta)
}

func resourceTranscodingTemplateGroupDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1.0/{project_id}/asset/template_group/transcodings"
		product = "vod"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating VOD client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath += fmt.Sprintf("?group_id=%s", d.Id())
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error deleting VOD transcoding template group: %s", err)
	}

	return nil
}
