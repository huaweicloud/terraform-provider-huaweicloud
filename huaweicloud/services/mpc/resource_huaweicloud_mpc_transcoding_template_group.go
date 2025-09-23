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

// @API MPC POST /v1/{project_id}/template_group/transcodings
// @API MPC GET /v1/{project_id}/template_group/transcodings
// @API MPC PUT /v1/{project_id}/template_group/transcodings
// @API MPC DELETE /v1/{project_id}/template_group/transcodings
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
			// If the requestBody not contain `audio`, after creating the resource,
			// execute `terraform plan` will trigger update.
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
			// If the requestBody not contain `video_common`, after creating the resource,
			// execute `terraform plan` will trigger update.
			"video_common": {
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
			"videos": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"width": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  0,
						},
						"height": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"bitrate": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
			"template_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func buildCreateTemplateGroupBodyParams(d *schema.ResourceData) map[string]interface{} {
	templateGroupParams := map[string]interface{}{
		"name":         d.Get("name"),
		"videos":       buildVideosBodyParams(d.Get("videos").([]interface{})),
		"audio":        buildAudioBodyParams(d.Get("audio").([]interface{})),
		"video_common": buildVideoCommonBodyParams(d.Get("video_common").([]interface{})),
		"common":       buildCommonBodyParams(d),
	}

	return templateGroupParams
}

func buildVideosBodyParams(rawVideos []interface{}) []map[string]interface{} {
	if len(rawVideos) == 0 {
		return nil
	}

	videos := make([]map[string]interface{}, 0, len(rawVideos))
	for _, v := range rawVideos {
		video := v.(map[string]interface{})
		params := map[string]interface{}{
			"width":   video["width"],
			"height":  video["height"],
			"bitrate": video["bitrate"],
		}
		videos = append(videos, params)
	}

	return videos
}

func buildAudioBodyParams(rawAudio []interface{}) map[string]interface{} {
	if len(rawAudio) == 0 {
		return nil
	}

	audio := rawAudio[0].(map[string]interface{})
	audioParams := map[string]interface{}{
		"output_policy": buildOutputPolicy(audio["output_policy"].(string)),
		"codec":         audio["codec"],
		"sample_rate":   audio["sample_rate"],
		"bitrate":       audio["bitrate"],
		"channels":      audio["channels"],
	}

	return audioParams
}

func buildVideoCommonBodyParams(rawVideoCommon []interface{}) map[string]interface{} {
	if len(rawVideoCommon) == 0 {
		return nil
	}

	videoCommon := rawVideoCommon[0].(map[string]interface{})
	videoCommonParams := map[string]interface{}{
		"output_policy":        buildOutputPolicy(videoCommon["output_policy"].(string)),
		"codec":                videoCommon["codec"],
		"profile":              videoCommon["profile"],
		"level":                videoCommon["level"],
		"preset":               videoCommon["quality"],
		"max_iframes_interval": videoCommon["max_iframes_interval"],
		"bframes_count":        videoCommon["max_consecutive_bframes"],
		"frame_rate":           videoCommon["fps"],
		"black_cut":            videoCommon["black_bar_removal"],
	}

	return videoCommonParams
}

func buildOutputPolicy(outputPolicy string) interface{} {
	if outputPolicy == "discard" || outputPolicy == "transcode" {
		return outputPolicy
	}

	return nil
}

func buildCommonBodyParams(d *schema.ResourceData) map[string]interface{} {
	commonParams := map[string]interface{}{
		"PVC":           d.Get("low_bitrate_hd"),
		"hls_interval":  d.Get("hls_segment_duration"),
		"dash_interval": d.Get("dash_segment_duration"),
		"pack_type":     d.Get("output_format"),
	}

	return commonParams
}

func resourceTranscodingTemplateGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/template_group/transcodings"
	)

	client, err := cfg.NewServiceClient("mpc", region)
	if err != nil {
		return diag.Errorf("error creating MPC client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateTemplateGroupBodyParams(d)),
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating MPC transcoding template group: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	groupId := utils.PathSearch("template_group.group_id", respBody, "").(string)
	if groupId == "" {
		return diag.Errorf("error creating MPC transcoding template group: ID is not found in API response")
	}

	d.SetId(groupId)

	return resourceTranscodingTemplateGroupRead(ctx, d, meta)
}

func resourceTranscodingTemplateGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("mpc", region)
	if err != nil {
		return diag.Errorf("error creating MPC client: %s", err)
	}

	// Use the `resource_id` filtter the resource, when the resource does not exist,
	// query API will return `403`, the error code is `MPC.10231`.
	respBody, err := GetTranscodingTemplateGroup(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected403ErrInto404Err(err, "error_code", "MPC.10231"),
			"error retrieving MPC transcoding template group")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", respBody, nil)),
		d.Set("videos", flattenTemplateGroupVideos(utils.PathSearch("videos", respBody, nil))),
		d.Set("audio", flattenTemplateGroupAudio(utils.PathSearch("audio", respBody, nil))),
		d.Set("video_common", flattenTemplateGroupVideoCommon(utils.PathSearch("video_common", respBody, nil))),
		d.Set("template_ids", flattenTemplateIds(utils.PathSearch("template_ids", respBody, make([]interface{}, 0)).([]interface{}))),
	)

	rawCommon := utils.PathSearch("common", respBody, nil)
	if rawCommon != nil {
		mErr = multierror.Append(mErr,
			d.Set("low_bitrate_hd", utils.PathSearch("common.PVC", respBody, nil)),
			d.Set("hls_segment_duration", utils.PathSearch("common.hls_interval", respBody, nil)),
			d.Set("dash_segment_duration", utils.PathSearch("common.dash_interval", respBody, nil)),
			d.Set("output_format", utils.PathSearch("common.pack_type", respBody, nil)),
		)
	}

	return diag.FromErr(mErr.ErrorOrNil())
}

func GetTranscodingTemplateGroup(client *golangsdk.ServiceClient, groupId string) (interface{}, error) {
	httpUrl := "v1/{project_id}/template_group/transcodings"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = fmt.Sprintf("%s?group_id=%s", getPath, groupId)
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

	templateGgroup := utils.PathSearch("template_group_list|[0]", respBody, nil)
	if templateGgroup == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return templateGgroup, nil
}

func flattenTemplateGroupVideos(resp interface{}) []map[string]interface{} {
	if resp == nil {
		return nil
	}

	rawArray := resp.([]interface{})
	rst := make([]map[string]interface{}, len(rawArray))
	for i, v := range rawArray {
		params := map[string]interface{}{
			"width":   utils.PathSearch("width", v, nil),
			"height":  utils.PathSearch("height", v, nil),
			"bitrate": utils.PathSearch("bitrate", v, nil),
		}
		rst[i] = params
	}

	return rst
}

func flattenTemplateGroupAudio(resp interface{}) []map[string]interface{} {
	if resp == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"output_policy": flattenOutputPolicy(utils.PathSearch("output_policy", resp, "").(string)),
			"codec":         utils.PathSearch("codec", resp, nil),
			"bitrate":       utils.PathSearch("bitrate", resp, nil),
			"sample_rate":   utils.PathSearch("sample_rate", resp, nil),
			"channels":      utils.PathSearch("channels", resp, nil),
		},
	}
}

func flattenOutputPolicy(outputPolicy string) string {
	if outputPolicy == "discard" || outputPolicy == "transcode" {
		return outputPolicy
	}

	return ""
}

func flattenTemplateGroupVideoCommon(resp interface{}) []map[string]interface{} {
	if resp == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"output_policy":           flattenOutputPolicy(utils.PathSearch("output_policy", resp, "").(string)),
			"codec":                   utils.PathSearch("codec", resp, nil),
			"profile":                 utils.PathSearch("profile", resp, nil),
			"level":                   utils.PathSearch("level", resp, nil),
			"quality":                 utils.PathSearch("preset", resp, nil),
			"max_reference_frames":    4,
			"max_iframes_interval":    utils.PathSearch("max_iframes_interval", resp, nil),
			"max_consecutive_bframes": utils.PathSearch("bframes_count", resp, nil),
			"fps":                     utils.PathSearch("frame_rate", resp, nil),
			"black_bar_removal":       utils.PathSearch("black_cut", resp, nil),
		},
	}
}

func flattenTemplateIds(rawtTemplateIds []interface{}) []string {
	templateIds := make([]string, len(rawtTemplateIds))
	for i, num := range rawtTemplateIds {
		templateIds[i] = strconv.Itoa(int(num.(float64)))
	}

	return templateIds
}

func resourceTranscodingTemplateGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/template_group/transcodings"
	)

	client, err := cfg.NewServiceClient("mpc", region)
	if err != nil {
		return diag.Errorf("error creating MPC client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updateOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200, 201, 204,
		},
		JSONBody: utils.RemoveNil(buildUpdateTemplateGroupBodyParams(d)),
	}

	_, err = client.Request("PUT", updatePath, &updateOpts)
	if err != nil {
		return diag.Errorf("error updating MPC transcoding template group: %s", err)
	}

	return resourceTranscodingTemplateGroupRead(ctx, d, meta)
}

func buildUpdateTemplateGroupBodyParams(d *schema.ResourceData) map[string]interface{} {
	templateGroupParams := map[string]interface{}{
		"group_id":     d.Id(),
		"name":         d.Get("name"),
		"videos":       buildVideosBodyParams(d.Get("videos").([]interface{})),
		"audio":        buildAudioBodyParams(d.Get("audio").([]interface{})),
		"video_common": buildVideoCommonBodyParams(d.Get("video_common").([]interface{})),
		"common":       buildCommonBodyParams(d),
	}

	return templateGroupParams
}

func resourceTranscodingTemplateGroupDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/template_group/transcodings"
	)

	client, err := cfg.NewServiceClient("mpc", region)
	if err != nil {
		return diag.Errorf("error creating MPC client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = fmt.Sprintf("%s?group_id=%s", deletePath, d.Id())
	deleteOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpts)
	if err != nil {
		// When the resource does not exist, delete API will return `403`, the error code is `MPC.10231`.
		return common.CheckDeletedDiag(d, common.ConvertExpected403ErrInto404Err(err, "error_code", "MPC.10231"),
			"error deleting MPC transcoding template group")
	}

	return nil
}
