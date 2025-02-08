package mpc

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	mpc "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/mpc/v1/model"

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

func buildCommonOpts(d *schema.ResourceData) *mpc.Common {
	commonOpts := mpc.Common{
		Pvc:          utils.Bool(d.Get("low_bitrate_hd").(bool)),
		HlsInterval:  utils.Int32(int32(d.Get("hls_segment_duration").(int))),
		DashInterval: utils.Int32(int32(d.Get("dash_segment_duration").(int))),
		PackType:     int32(d.Get("output_format").(int)),
	}

	return &commonOpts
}

func buildAudioOpts(rawAudio []interface{}) *mpc.Audio {
	if len(rawAudio) != 1 {
		return nil
	}

	audio := rawAudio[0].(map[string]interface{})
	audioOpts := mpc.Audio{
		Codec:        utils.Int32(int32(audio["codec"].(int))),
		SampleRate:   utils.Int32(int32(audio["sample_rate"].(int))),
		Channels:     int32(audio["channels"].(int)),
		Bitrate:      utils.Int32(int32(audio["bitrate"].(int))),
		OutputPolicy: buildAudioOutputPolicyOpts(audio["output_policy"].(string)),
	}

	return &audioOpts
}

func buildAudioOutputPolicyOpts(outputPolicy string) *mpc.AudioOutputPolicy {
	if outputPolicy == "" {
		return nil
	}

	var outputPolicyOpts mpc.AudioOutputPolicy
	switch outputPolicy {
	case "discard":
		outputPolicyOpts = mpc.GetAudioOutputPolicyEnum().DISCARD
	case "transcode":
		outputPolicyOpts = mpc.GetAudioOutputPolicyEnum().TRANSCODE
	default:
		log.Printf("[WARN] output_policy invalid: %s", outputPolicy)
		return nil
	}

	return &outputPolicyOpts
}

func buildVideoOpts(rawVideo []interface{}) *mpc.Video {
	if len(rawVideo) != 1 {
		return nil
	}

	video := rawVideo[0].(map[string]interface{})
	videoOpts := mpc.Video{
		Codec:              utils.Int32(int32(video["codec"].(int))),
		Bitrate:            utils.Int32(int32(video["bitrate"].(int))),
		Profile:            utils.Int32(int32(video["profile"].(int))),
		Level:              utils.Int32(int32(video["level"].(int))),
		Preset:             utils.Int32(int32(video["quality"].(int))),
		MaxIframesInterval: utils.Int32(int32(video["max_iframes_interval"].(int))),
		BframesCount:       utils.Int32(int32(video["max_consecutive_bframes"].(int))),
		FrameRate:          utils.Int32(int32(video["fps"].(int))),
		Width:              utils.Int32(int32(video["width"].(int))),
		Height:             utils.Int32(int32(video["height"].(int))),
		BlackCut:           utils.Int32(int32(video["black_bar_removal"].(int))),
		OutputPolicy:       buildVideoOutputPolicyOpts(video["output_policy"].(string)),
	}

	return &videoOpts
}

func buildVideoOutputPolicyOpts(outputPolicy string) *mpc.VideoOutputPolicy {
	if outputPolicy == "" {
		return nil
	}

	var outputPolicyOpts mpc.VideoOutputPolicy
	switch outputPolicy {
	case "discard":
		outputPolicyOpts = mpc.GetVideoOutputPolicyEnum().DISCARD
	case "transcode":
		outputPolicyOpts = mpc.GetVideoOutputPolicyEnum().TRANSCODE
	default:
		log.Printf("[WARN] output_policy invalid: %s", outputPolicy)
		return nil
	}

	return &outputPolicyOpts
}

func resourceTranscodingTemplateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.HcMpcV1Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating MPC client : %s", err)
	}

	createOpts := mpc.TransTemplate{
		TemplateName: d.Get("name").(string),
		Video:        buildVideoOpts(d.Get("video").([]interface{})),
		Audio:        buildAudioOpts(d.Get("audio").([]interface{})),
		Common:       buildCommonOpts(d),
	}

	createReq := mpc.CreateTransTemplateRequest{
		Body: &createOpts,
	}

	resp, err := client.CreateTransTemplate(&createReq)
	if err != nil {
		return diag.Errorf("error creating MPC transcoding template: %s", err)
	}

	d.SetId(strconv.FormatInt(int64(*resp.TemplateId), 10))

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
	config := meta.(*config.Config)
	client, err := config.HcMpcV1Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating MPC client : %s", err)
	}

	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return diag.FromErr(err)
	}

	updateOpts := mpc.ModifyTransTemplateReq{
		TemplateId:   id,
		TemplateName: d.Get("name").(string),
		Video:        buildVideoOpts(d.Get("video").([]interface{})),
		Audio:        buildAudioOpts(d.Get("audio").([]interface{})),
		Common:       buildCommonOpts(d),
	}

	updateReq := mpc.UpdateTransTemplateRequest{
		Body: &updateOpts,
	}

	_, err = client.UpdateTransTemplate(&updateReq)
	if err != nil {
		return diag.Errorf("error updating MPC transcoding template: %s", err)
	}

	return resourceTranscodingTemplateRead(ctx, d, meta)
}

func resourceTranscodingTemplateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.HcMpcV1Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating MPC client : %s", err)
	}

	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		diag.FromErr(err)
	}

	_, err = client.DeleteTemplate(&mpc.DeleteTemplateRequest{TemplateId: id})
	if err != nil {
		return diag.Errorf("error deleting MPC transcoding template: %s", err)
	}

	return nil
}
