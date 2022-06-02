package mpc

import (
	"context"
	"log"
	"strconv"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	mpc "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/mpc/v1/model"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

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
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 6),
			},
			"hls_segment_duration": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(2, 10),
				Default:      5,
			},
			"dash_segment_duration": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(2, 10),
				Default:      5,
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
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 4),
						},
						"sample_rate": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 6),
						},
						"channels": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntInSlice([]int{1, 2, 6}),
						},
						"bitrate": {
							Type:     schema.TypeInt,
							Optional: true,
							ValidateFunc: validation.Any(
								validation.IntInSlice([]int{0}),
								validation.IntBetween(8, 1000),
							),
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
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 2),
							Default:      1,
						},
						"profile": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 4),
							Default:      3,
						},
						"level": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 15),
							Default:      15,
						},
						"quality": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 3),
							Default:      1,
						},
						"max_reference_frames": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 8),
							Default:      4,
						},
						"max_iframes_interval": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(2, 10),
							Default:      5,
						},
						"max_consecutive_bframes": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(0, 7),
							Default:      4,
						},
						"fps": {
							Type:     schema.TypeInt,
							Optional: true,
							ValidateFunc: validation.Any(
								validation.IntInSlice([]int{0}),
								validation.IntBetween(5, 30),
							),
						},
						"black_bar_removal": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(0, 2),
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
							ValidateFunc: validation.Any(
								validation.IntInSlice([]int{0}),
								validation.IntBetween(32, 4096),
							),
							Default: 0,
						},
						"height": {
							Type:     schema.TypeInt,
							Optional: true,
							ValidateFunc: validation.Any(
								validation.IntInSlice([]int{0}),
								validation.IntBetween(32, 2880),
							),
						},
						"bitrate": {
							Type:     schema.TypeInt,
							Optional: true,
							ValidateFunc: validation.Any(
								validation.IntInSlice([]int{0}),
								validation.IntBetween(40, 30000),
							),
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

func buildVideosOpts(rawVideos []interface{}) *[]mpc.VideoObj {
	if len(rawVideos) == 0 {
		return nil
	}

	videos := make([]mpc.VideoObj, len(rawVideos))
	for i, rawVideo := range rawVideos {
		video := rawVideo.(map[string]interface{})
		videos[i] = mpc.VideoObj{
			Width:   int32(video["width"].(int)),
			Height:  int32(video["height"].(int)),
			Bitrate: int32(video["bitrate"].(int)),
		}
	}

	return &videos
}

func buildVideoCommonOpts(rawVideoCommon []interface{}) *mpc.VideoCommon {
	if len(rawVideoCommon) != 1 {
		return nil
	}

	videoCommon := rawVideoCommon[0].(map[string]interface{})
	videoCommonOpts := mpc.VideoCommon{
		Codec:              utils.Int32(int32(videoCommon["codec"].(int))),
		Profile:            utils.Int32(int32(videoCommon["profile"].(int))),
		Level:              utils.Int32(int32(videoCommon["level"].(int))),
		Preset:             utils.Int32(int32(videoCommon["quality"].(int))),
		RefFramesCount:     utils.Int32(int32(videoCommon["max_reference_frames"].(int))),
		MaxIframesInterval: utils.Int32(int32(videoCommon["max_iframes_interval"].(int))),
		BframesCount:       utils.Int32(int32(videoCommon["max_consecutive_bframes"].(int))),
		FrameRate:          utils.Int32(int32(videoCommon["fps"].(int))),
		BlackCut:           utils.Int32(int32(videoCommon["black_bar_removal"].(int))),
		OutputPolicy:       buildVideoCommonOutputPolicyOpts(videoCommon["output_policy"].(string)),
	}

	return &videoCommonOpts
}

func buildVideoCommonOutputPolicyOpts(outputPolicy string) *mpc.VideoCommonOutputPolicy {
	if outputPolicy == "" {
		return nil
	}

	var outputPolicyOpts mpc.VideoCommonOutputPolicy
	switch outputPolicy {
	case "discard":
		outputPolicyOpts = mpc.GetVideoCommonOutputPolicyEnum().DISCARD
	case "transcode":
		outputPolicyOpts = mpc.GetVideoCommonOutputPolicyEnum().TRANSCODE
	default:
		log.Printf("[WARN] output_policy invalid: %s", outputPolicy)
		return nil
	}

	return &outputPolicyOpts
}

func resourceTranscodingTemplateGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.HcMpcV1Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating MPC client : %s", err)
	}

	createOpts := mpc.TransTemplateGroup{
		Name:        utils.String(d.Get("name").(string)),
		Videos:      buildVideosOpts(d.Get("videos").([]interface{})),
		VideoCommon: buildVideoCommonOpts(d.Get("video_common").([]interface{})),
		Audio:       buildAudioOpts(d.Get("audio").([]interface{})),
		Common:      buildCommonOpts(d),
	}

	createReq := mpc.CreateTemplateGroupRequest{
		Body: &createOpts,
	}

	resp, err := client.CreateTemplateGroup(&createReq)
	if err != nil {
		return diag.Errorf("error creating MPC transcoding template group: %s", err)
	}

	d.SetId(*resp.TemplateGroup.GroupId)

	return resourceTranscodingTemplateGroupRead(ctx, d, meta)
}

func resourceTranscodingTemplateGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.HcMpcV1Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating MPC client : %s", err)
	}

	resp, err := client.ListTemplateGroup(&mpc.ListTemplateGroupRequest{GroupId: &[]string{d.Id()}})
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving MPC transcoding template group")
	}

	templateGroupList := *resp.TemplateGroupList

	if len(templateGroupList) == 0 {
		log.Printf("unable to retrieve MPC transcoding template group: %s", d.Id())
		d.SetId("")
		return nil
	}
	templateGroup := templateGroupList[0]

	templateIds := make([]string, len(*templateGroup.TemplateIds))
	for i, v := range *templateGroup.TemplateIds {
		templateIds[i] = strconv.FormatInt(int64(v), 10)
	}

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("name", templateGroup.Name),
		d.Set("audio", flattenAudio(templateGroup.Audio)),
		d.Set("video_common", flattenVideoCommon(templateGroup.VideoCommon)),
		d.Set("videos", flattenVideos(templateGroup.Videos)),
		d.Set("template_ids", templateIds),
		setCommonAttrs(d, templateGroup.Common),
	)

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting MPC transcoding template group fields: %s", err)
	}

	return nil
}

func resourceTranscodingTemplateGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.HcMpcV1Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating MPC client : %s", err)
	}

	updateOpts := mpc.ModifyTransTemplateGroup{
		GroupId:     utils.String(d.Id()),
		Name:        utils.String(d.Get("name").(string)),
		Videos:      buildVideosOpts(d.Get("videos").([]interface{})),
		VideoCommon: buildVideoCommonOpts(d.Get("video_common").([]interface{})),
		Audio:       buildAudioOpts(d.Get("audio").([]interface{})),
		Common:      buildCommonOpts(d),
	}

	updateReq := mpc.UpdateTemplateGroupRequest{
		Body: &updateOpts,
	}

	_, err = client.UpdateTemplateGroup(&updateReq)
	if err != nil {
		return diag.Errorf("error updating MPC transcoding template group: %s", err)
	}

	return resourceTranscodingTemplateGroupRead(ctx, d, meta)
}

func resourceTranscodingTemplateGroupDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.HcMpcV1Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating MPC client : %s", err)
	}

	_, err = client.DeleteTemplateGroup(&mpc.DeleteTemplateGroupRequest{GroupId: d.Id()})
	if err != nil {
		return diag.Errorf("error deleting MPC transcoding template group: %s", err)
	}

	return nil
}

func flattenVideoCommon(video *mpc.VideoCommon) []map[string]interface{} {
	if video == nil {
		return nil
	}
	videoResult := map[string]interface{}{
		"output_policy":           video.OutputPolicy,
		"codec":                   video.Codec,
		"profile":                 video.Profile,
		"level":                   video.Level,
		"quality":                 video.Preset,
		"max_reference_frames":    video.RefFramesCount,
		"max_iframes_interval":    video.MaxIframesInterval,
		"max_consecutive_bframes": video.BframesCount,
		"fps":                     video.FrameRate,
		"black_bar_removal":       video.BlackCut,
	}

	var outputPolicy string
	switch *video.OutputPolicy {
	case mpc.GetVideoCommonOutputPolicyEnum().DISCARD:
		outputPolicy = "discard"
	case mpc.GetVideoCommonOutputPolicyEnum().TRANSCODE:
		outputPolicy = "transcode"
	default:
		outputPolicy = ""
	}
	videoResult["output_policy"] = outputPolicy

	return []map[string]interface{}{videoResult}
}

func flattenVideos(videos *[]mpc.VideoAndTemplate) []map[string]interface{} {
	if videos == nil {
		return nil
	}

	videosResult := make([]map[string]interface{}, len(*videos))
	for i, v := range *videos {
		video := map[string]interface{}{
			"width":   v.Width,
			"height":  v.Height,
			"bitrate": v.Bitrate,
		}
		videosResult[i] = video
	}

	return videosResult
}
