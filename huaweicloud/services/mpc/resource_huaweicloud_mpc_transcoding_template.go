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
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 2),
							Default:      1,
						},
						"bitrate": {
							Type:     schema.TypeInt,
							Optional: true,
							ValidateFunc: validation.Any(
								validation.IntInSlice([]int{0}),
								validation.IntBetween(40, 30000),
							),
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
						"width": {
							Type:     schema.TypeInt,
							Optional: true,
							ValidateFunc: validation.Any(
								validation.IntInSlice([]int{0}),
								validation.IntBetween(32, 4096),
							),
						},
						"height": {
							Type:     schema.TypeInt,
							Optional: true,
							ValidateFunc: validation.Any(
								validation.IntInSlice([]int{0}),
								validation.IntBetween(32, 2880),
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
		},
	}
}

func buildCommonOpts(d *schema.ResourceData) *mpc.Common {
	commonOpts := mpc.Common{
		Pvc:          d.Get("low_bitrate_hd").(bool),
		HlsInterval:  int32(d.Get("hls_segment_duration").(int)),
		DashInterval: int32(d.Get("dash_segment_duration").(int)),
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
		Codec:        int32(audio["codec"].(int)),
		SampleRate:   int32(audio["sample_rate"].(int)),
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
		RefFramesCount:     utils.Int32(int32(video["max_reference_frames"].(int))),
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
	config := meta.(*config.Config)
	client, err := config.HcMpcV1Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating MPC client : %s", err)
	}

	id, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.ListTemplate(&mpc.ListTemplateRequest{TemplateId: &[]int32{int32(id)}})
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving MPC transcoding template")
	}

	templateList := *resp.TemplateArray
	// templateList is never empty but still check it in case the api change
	if len(templateList) == 0 || templateList[0].Template == nil {
		log.Printf("unable to retrieve MPC transcoding template: %s", d.Id())
		d.SetId("")
		return nil
	}
	template := templateList[0].Template

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("name", template.TemplateName),
		d.Set("audio", flattenAudio(template.Audio)),
		d.Set("video", flattenVideo(template.Video)),
		setCommonAttrs(d, template.Common),
	)

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting MPC transcoding template fields: %s", err)
	}

	return nil
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

func setCommonAttrs(d *schema.ResourceData, common *mpc.Common) error {
	if common == nil {
		return nil
	}
	mErr := multierror.Append(nil,
		d.Set("low_bitrate_hd", common.Pvc),
		d.Set("hls_segment_duration", common.HlsInterval),
		d.Set("dash_segment_duration", common.DashInterval),
		d.Set("output_format", common.PackType),
	)

	return mErr
}

func flattenAudio(audio *mpc.Audio) []map[string]interface{} {
	if audio == nil {
		return nil
	}
	audioResult := map[string]interface{}{
		"codec":       audio.Codec,
		"sample_rate": audio.SampleRate,
		"channels":    audio.Channels,
		"bitrate":     audio.Bitrate,
	}

	var outputPolicy string
	switch *audio.OutputPolicy {
	case mpc.GetAudioOutputPolicyEnum().DISCARD:
		outputPolicy = "discard"
	case mpc.GetAudioOutputPolicyEnum().TRANSCODE:
		outputPolicy = "transcode"
	default:
		outputPolicy = ""
	}
	audioResult["output_policy"] = outputPolicy

	return []map[string]interface{}{audioResult}
}

func flattenVideo(video *mpc.Video) []map[string]interface{} {
	if video == nil {
		return nil
	}
	videoResult := map[string]interface{}{
		"output_policy":           video.OutputPolicy,
		"codec":                   video.Codec,
		"bitrate":                 video.Bitrate,
		"profile":                 video.Profile,
		"level":                   video.Level,
		"quality":                 video.Preset,
		"max_reference_frames":    video.RefFramesCount,
		"max_iframes_interval":    video.MaxIframesInterval,
		"max_consecutive_bframes": video.BframesCount,
		"fps":                     video.FrameRate,
		"width":                   video.Width,
		"height":                  video.Height,
		"black_bar_removal":       video.BlackCut,
	}

	var outputPolicy string
	switch *video.OutputPolicy {
	case mpc.GetVideoOutputPolicyEnum().DISCARD:
		outputPolicy = "discard"
	case mpc.GetVideoOutputPolicyEnum().TRANSCODE:
		outputPolicy = "transcode"
	default:
		outputPolicy = ""
	}
	videoResult["output_policy"] = outputPolicy

	return []map[string]interface{}{videoResult}
}
