package mpc

import (
	"context"
	"fmt"
	"strconv"
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

// @API MPC POST /v1/{project_id}/template/transcodings
// @API MPC GET /v1/{project_id}/template/transcodings
// @API MPC PUT /v1/{project_id}/template/transcodings
// @API MPC DELETE /v1/{project_id}/template/transcodings
func ResourceTranscodingTemplate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTranscodingTemplateCreate,
		ReadContext:   resourceTranscodingTemplateRead,
		UpdateContext: resourceTranscodingTemplateUpdate,
		DeleteContext: resourceTranscodingTemplateDelete,
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
			"output_format": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"hls_segment_duration": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  5,
			},
			"dash_segment_duration": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  5,
			},
			"low_bitrate_hd": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"audio": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"codec": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"sample_rate": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"channels": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntInSlice([]int{1, 2, 6}),
						},
						"bitrate": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"output_policy": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"discard", "transcode",
							}, false),
							Default: "transcode",
						},
					},
				},
			},
			"video": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"output_policy": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"discard", "transcode",
							}, false),
							Default: "transcode",
						},
						"codec": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  1,
						},
						"bitrate": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"profile": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  3,
						},
						"level": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  15,
						},
						"quality": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  1,
						},
						"max_iframes_interval": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  5,
						},
						"max_consecutive_bframes": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  4,
						},
						"fps": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"width": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"height": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"black_bar_removal": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"max_reference_frames": {
							Type:        schema.TypeInt,
							Optional:    true,
							Computed:    true,
							Description: "schema: Deprecated; the SDK does not support it",
						},
					},
				},
			},
		},
	}
}

func buildCreateTemplateBodyParams(d *schema.ResourceData) map[string]interface{} {
	templateGroupParams := map[string]interface{}{
		"template_name": d.Get("name"),
		"video":         buildTemplateVideoBodyParams(d.Get("video").([]interface{})),
		"audio":         buildTemplateAudioBodyParams(d.Get("audio").([]interface{})),
		"common":        buildTemplateCommonBodyParams(d),
	}

	return templateGroupParams
}

func buildTemplateVideoBodyParams(rawVideo []interface{}) map[string]interface{} {
	if len(rawVideo) == 0 {
		return nil
	}

	video := rawVideo[0].(map[string]interface{})
	videoParams := map[string]interface{}{
		"output_policy":        buildTemplateOutputPolicy(video["output_policy"].(string)),
		"codec":                video["codec"],
		"bitrate":              video["bitrate"],
		"profile":              video["profile"],
		"level":                video["level"],
		"preset":               video["quality"],
		"max_iframes_interval": video["max_iframes_interval"],
		"bframes_count":        video["max_consecutive_bframes"],
		"frame_rate":           video["fps"],
		"width":                video["width"],
		"height":               video["height"],
		"black_cut":            video["black_bar_removal"],
	}

	return videoParams
}

func buildTemplateOutputPolicy(outputPolicy string) interface{} {
	if outputPolicy == "discard" || outputPolicy == "transcode" {
		return outputPolicy
	}

	return nil
}

func buildTemplateAudioBodyParams(rawAudio []interface{}) map[string]interface{} {
	if len(rawAudio) == 0 {
		return nil
	}

	audio := rawAudio[0].(map[string]interface{})
	audioParams := map[string]interface{}{
		"output_policy": buildTemplateOutputPolicy(audio["output_policy"].(string)),
		"codec":         audio["codec"],
		"sample_rate":   audio["sample_rate"],
		"bitrate":       audio["bitrate"],
		"channels":      audio["channels"],
	}

	return audioParams
}

func buildTemplateCommonBodyParams(d *schema.ResourceData) map[string]interface{} {
	commonParams := map[string]interface{}{
		"PVC":           d.Get("low_bitrate_hd"),
		"hls_interval":  d.Get("hls_segment_duration"),
		"dash_interval": d.Get("dash_segment_duration"),
		"pack_type":     d.Get("output_format"),
	}

	return commonParams
}

func resourceTranscodingTemplateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/template/transcodings"
	)

	client, err := cfg.NewServiceClient("mpc", region)
	if err != nil {
		return diag.Errorf("error creating MPC client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateTemplateBodyParams(d)),
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating MPC transcoding template: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	templateId := utils.PathSearch("template_id", respBody, float64(0)).(float64)
	if templateId == 0 {
		return diag.Errorf("error creating MPC transcoding template group: ID is not found in API response")
	}

	d.SetId(strconv.FormatInt(int64(templateId), 10))

	return resourceTranscodingTemplateRead(ctx, d, meta)
}

func resourceTranscodingTemplateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("mpc", region)
	if err != nil {
		return diag.Errorf("error creating MPC client: %s", err)
	}

	respBody, err := GetTranscodingTemplate(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving MPC transcoding template")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", utils.PathSearch("template_name", respBody, nil)),
		d.Set("video", flattenTemplateVideo(utils.PathSearch("video", respBody, nil))),
		d.Set("audio", flattenTemplateAudio(utils.PathSearch("audio", respBody, nil))),
	)

	rawCommon := utils.PathSearch("common", respBody, nil)
	if rawCommon != nil {
		mErr = multierror.Append(mErr,
			d.Set("low_bitrate_hd", utils.PathSearch("PVC", rawCommon, nil)),
			d.Set("hls_segment_duration", utils.PathSearch("hls_interval", rawCommon, nil)),
			d.Set("dash_segment_duration", utils.PathSearch("dash_interval", rawCommon, nil)),
			d.Set("output_format", utils.PathSearch("pack_type", rawCommon, nil)),
		)
	}

	return diag.FromErr(mErr.ErrorOrNil())
}

func GetTranscodingTemplate(client *golangsdk.ServiceClient, templateId string) (interface{}, error) {
	httpUrl := "v1/{project_id}/template/transcodings"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = fmt.Sprintf("%s?template_id=%s", getPath, templateId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	template := utils.PathSearch("template_array[0].template", respBody, nil)
	if template == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return template, nil
}

func flattenTemplateVideo(video interface{}) []map[string]interface{} {
	if video == nil {
		return nil
	}

	videoResult := map[string]interface{}{
		"output_policy":           flattenTemplateOutputPolicy(utils.PathSearch("output_policy", video, "").(string)),
		"codec":                   utils.PathSearch("codec", video, nil),
		"bitrate":                 utils.PathSearch("bitrate", video, nil),
		"profile":                 utils.PathSearch("profile", video, nil),
		"level":                   utils.PathSearch("level", video, nil),
		"quality":                 utils.PathSearch("preset", video, nil),
		"max_reference_frames":    4,
		"max_iframes_interval":    utils.PathSearch("max_iframes_interval", video, nil),
		"max_consecutive_bframes": utils.PathSearch("bframes_count", video, nil),
		"fps":                     utils.PathSearch("frame_rate", video, nil),
		"width":                   utils.PathSearch("width", video, nil),
		"height":                  utils.PathSearch("height", video, nil),
		"black_bar_removal":       utils.PathSearch("black_cut", video, nil),
	}

	return []map[string]interface{}{videoResult}
}

func flattenTemplateOutputPolicy(outputPolicy string) string {
	if outputPolicy == "discard" || outputPolicy == "transcode" {
		return outputPolicy
	}

	return ""
}

func flattenTemplateAudio(audio interface{}) []map[string]interface{} {
	if audio == nil {
		return nil
	}

	audioResult := map[string]interface{}{
		"output_policy": flattenTemplateOutputPolicy(utils.PathSearch("output_policy", audio, "").(string)),
		"codec":         utils.PathSearch("codec", audio, nil),
		"bitrate":       utils.PathSearch("bitrate", audio, nil),
		"sample_rate":   utils.PathSearch("sample_rate", audio, nil),
		"channels":      utils.PathSearch("channels", audio, nil),
	}

	return []map[string]interface{}{audioResult}
}

func resourceTranscodingTemplateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/template/transcodings"
	)

	client, err := cfg.NewServiceClient("mpc", region)
	if err != nil {
		return diag.Errorf("error creating MPC client: %s", err)
	}

	templateId, err := strconv.ParseInt(d.Id(), 10, 64)
	if err != nil {
		return diag.FromErr(err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updateOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200, 201, 204,
		},
		JSONBody: utils.RemoveNil(buildUpdateTemplateBodyParams(d, int(templateId))),
	}

	_, err = client.Request("PUT", updatePath, &updateOpts)
	if err != nil {
		return diag.Errorf("error updating MPC transcoding template: %s", err)
	}

	return resourceTranscodingTemplateRead(ctx, d, meta)
}

func buildUpdateTemplateBodyParams(d *schema.ResourceData, templateId int) map[string]interface{} {
	templateGroupParams := map[string]interface{}{
		"template_id":   templateId,
		"template_name": d.Get("name"),
		"video":         buildTemplateVideoBodyParams(d.Get("video").([]interface{})),
		"audio":         buildTemplateAudioBodyParams(d.Get("audio").([]interface{})),
		"common":        buildTemplateCommonBodyParams(d),
	}

	return templateGroupParams
}

func resourceTranscodingTemplateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/template/transcodings"
	)

	client, err := cfg.NewServiceClient("mpc", region)
	if err != nil {
		return diag.Errorf("error creating MPC client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = fmt.Sprintf("%s?template_id=%s", deletePath, d.Id())
	deleteOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpts)
	if err != nil {
		// When the resource does not exist, delete API will return `403`, the error code is `MPC.10239`.
		return common.CheckDeletedDiag(d, common.ConvertExpected403ErrInto404Err(err, "error_code", "MPC.10239"),
			"error deleting MPC transcoding template")
	}

	return nil
}
