package vod

import (
	"context"
	"fmt"
	"log"
	"regexp"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	vod "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/vod/v1/model"

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
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 128),
					validation.StringMatch(regexp.MustCompile(`^\w+$`),
						"The name can only consist of letters, digits and underscores (_)."),
				),
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
										ValidateFunc: validation.Any(
											validation.IntInSlice([]int{0}),
											validation.IntBetween(128, 3840),
										),
									},
									"height": {
										Type:     schema.TypeInt,
										Optional: true,
										ValidateFunc: validation.Any(
											validation.IntInSlice([]int{0}),
											validation.IntBetween(128, 2160),
										),
									},
									"bitrate": {
										Type:     schema.TypeInt,
										Optional: true,
										ValidateFunc: validation.Any(
											validation.IntInSlice([]int{0}),
											validation.IntBetween(700, 3000),
										),
									},
									"frame_rate": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntBetween(1, 75),
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
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntBetween(1, 6),
									},
									"channels": {
										Type:         schema.TypeInt,
										Required:     true,
										ValidateFunc: validation.IntBetween(1, 2),
									},
									"bitrate": {
										Type:     schema.TypeInt,
										Optional: true,
										ValidateFunc: validation.Any(
											validation.IntInSlice([]int{0}),
											validation.IntBetween(8, 1000),
										),
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

func buildCommonOpts(d *schema.ResourceData) *vod.Common {
	commonOpts := vod.Common{
		HlsInterval: utils.Int32(int32(d.Get("hls_segment_duration").(int))),
	}

	if d.Get("low_bitrate_hd").(bool) {
		commonOpts.Pvc = vod.GetCommonPvcEnum().E_1
	} else {
		commonOpts.Pvc = vod.GetCommonPvcEnum().E_0
	}

	switch d.Get("video_codec").(string) {
	case "H265":
		videoCodec := vod.GetCommonVideoCodecEnum().H265
		commonOpts.VideoCodec = &videoCodec
	case "H264":
		videoCodec := vod.GetCommonVideoCodecEnum().H264
		commonOpts.VideoCodec = &videoCodec
	default:
		commonOpts.VideoCodec = nil
	}

	switch d.Get("audio_codec").(string) {
	case "HEAAC1":
		audioCodec := vod.GetCommonAudioCodecEnum().HEAAC1
		commonOpts.AudioCodec = &audioCodec
	case "HEAAC2":
		audioCodec := vod.GetCommonAudioCodecEnum().HEAAC2
		commonOpts.AudioCodec = &audioCodec
	case "MP3":
		audioCodec := vod.GetCommonAudioCodecEnum().MP3
		commonOpts.AudioCodec = &audioCodec
	case "AAC":
		audioCodec := vod.GetCommonAudioCodecEnum().AAC
		commonOpts.AudioCodec = &audioCodec
	default:
		commonOpts.AudioCodec = nil
	}

	return &commonOpts
}

func buildAudioOpts(rawAudio []interface{}) *vod.AudioTemplateInfo {
	if len(rawAudio) != 1 {
		return nil
	}

	audio := rawAudio[0].(map[string]interface{})
	audioOpts := vod.AudioTemplateInfo{
		SampleRate: int32(audio["sample_rate"].(int)),
		Channels:   int32(audio["channels"].(int)),
		Bitrate:    utils.Int32(int32(audio["bitrate"].(int))),
	}

	return &audioOpts
}

func buildVideoOpts(rawVideo []interface{}) *vod.VideoTemplateInfo {
	if len(rawVideo) != 1 {
		return nil
	}

	video := rawVideo[0].(map[string]interface{})
	var quality vod.VideoTemplateInfoQuality
	switch video["quality"].(string) {
	case "4K":
		quality = vod.GetVideoTemplateInfoQualityEnum().E_4_K
	case "2K":
		quality = vod.GetVideoTemplateInfoQualityEnum().E_2_K
	case "FHD":
		quality = vod.GetVideoTemplateInfoQualityEnum().FULL_HD
	case "SD":
		quality = vod.GetVideoTemplateInfoQualityEnum().SD
	case "LD":
		quality = vod.GetVideoTemplateInfoQualityEnum().FLUENT
	default:
		quality = vod.GetVideoTemplateInfoQualityEnum().HD
	}

	videoOpts := vod.VideoTemplateInfo{
		Quality:   quality,
		Width:     utils.Int32(int32(video["width"].(int))),
		Height:    utils.Int32(int32(video["height"].(int))),
		Bitrate:   utils.Int32(int32(video["bitrate"].(int))),
		FrameRate: utils.Int32(int32(video["frame_rate"].(int))),
	}

	return &videoOpts
}

func buildQualityInfoListOpts(qualityInfo []interface{}) *[]vod.QualityInfo {
	if len(qualityInfo) == 0 {
		return nil
	}

	qualityInfoOpts := make([]vod.QualityInfo, len(qualityInfo))
	for i, v := range qualityInfo {
		info := v.(map[string]interface{})
		var outputFormat vod.QualityInfoFormat
		switch info["output_format"].(string) {
		case "MP4":
			outputFormat = vod.GetQualityInfoFormatEnum().MP4
		case "DASH":
			outputFormat = vod.GetQualityInfoFormatEnum().DASH
		case "DASH_HLS":
			outputFormat = vod.GetQualityInfoFormatEnum().DASH_HLS
		case "MP3":
			outputFormat = vod.GetQualityInfoFormatEnum().MP3
		case "ADTS":
			outputFormat = vod.GetQualityInfoFormatEnum().ADTS
		default:
			outputFormat = vod.GetQualityInfoFormatEnum().HLS
		}
		qualityInfoOpts[i] = vod.QualityInfo{
			Video:  buildVideoOpts(info["video"].([]interface{})),
			Audio:  buildAudioOpts(info["audio"].([]interface{})),
			Format: outputFormat,
		}
	}

	return &qualityInfoOpts
}

func resourceTranscodingTemplateGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.HcVodV1Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating VOD client : %s", err)
	}

	createOpts := vod.TransTemplateGroup{
		Name:                 d.Get("name").(string),
		Type:                 vod.GetTransTemplateGroupTypeEnum().CUSTOM_TEMPLATE_GROUP,
		Description:          utils.String(d.Get("description").(string)),
		WatermarkTemplateIds: utils.ExpandToStringListPointer(d.Get("watermark_template_ids").([]interface{})),
		Common:               buildCommonOpts(d),
		QualityInfoList:      buildQualityInfoListOpts(d.Get("quality_info").(([]interface{}))),
	}

	if d.Get("auto_encrypt").(bool) {
		createOpts.AutoEncrypt = utils.Int32(int32(1))
	}

	if d.Get("is_default").(bool) {
		status := vod.GetTransTemplateGroupStatusEnum().E_1
		createOpts.Status = &status
	}

	createReq := vod.CreateTemplateGroupRequest{
		Body: &createOpts,
	}
	log.Printf("[DEBUG] Create VOD transcoding template group Options: %#v", createOpts)

	resp, err := client.CreateTemplateGroup(&createReq)
	if err != nil {
		return diag.Errorf("error creating VOD transcoding template group: %s", err)
	}

	d.SetId(*resp.GroupId)

	return resourceTranscodingTemplateGroupRead(ctx, d, meta)
}

func resourceTranscodingTemplateGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.HcVodV1Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating VOD client : %s", err)
	}

	resp, err := client.ListTemplateGroup(&vod.ListTemplateGroupRequest{GroupId: utils.String(d.Id())})
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving VOD transcoding template group")
	}

	if resp.TemplateGroupList == nil || len(*resp.TemplateGroupList) == 0 {
		d.SetId("")
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Resource not found",
				Detail:   fmt.Sprintf("unable to retrieve VOD transcoding template group: %s", d.Id()),
			},
		}
	}
	templateGroupList := *resp.TemplateGroupList
	templateGroup := templateGroupList[0]

	var autoEncrypt bool
	if templateGroup.AutoEncrypt != nil && *templateGroup.AutoEncrypt == 1 {
		autoEncrypt = true
	}

	var isDefault bool
	if templateGroup.Status != nil && *templateGroup.Status == "1" {
		isDefault = true
	}

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("name", templateGroup.Name),
		d.Set("auto_encrypt", autoEncrypt),
		d.Set("is_default", isDefault),
		d.Set("description", templateGroup.Description),
		d.Set("type", templateGroup.Type),
		d.Set("watermark_template_ids", templateGroup.WatermarkTemplateIds),
		d.Set("quality_info", flattenQualityInfoList(templateGroup.QualityInfoList)),
		setCommonAttrs(d, templateGroup.Common),
	)

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting VOD transcoding template group fields: %s", err)
	}

	return nil
}

func resourceTranscodingTemplateGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.HcVodV1Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating VOD client : %s", err)
	}

	updateOpts := vod.ModifyTransTemplateGroup{
		GroupId:              d.Id(),
		Name:                 d.Get("name").(string),
		Description:          utils.String(d.Get("description").(string)),
		WatermarkTemplateIds: utils.ExpandToStringListPointer(d.Get("watermark_template_ids").([]interface{})),
		Common:               buildCommonOpts(d),
		QualityInfoList:      buildQualityInfoListOpts(d.Get("quality_info").(([]interface{}))),
	}

	if d.Get("auto_encrypt").(bool) {
		updateOpts.AutoEncrypt = utils.Int32(int32(1))
	}

	var status vod.ModifyTransTemplateGroupStatus
	if d.Get("is_default").(bool) {
		status = vod.GetModifyTransTemplateGroupStatusEnum().E_1
	} else {
		status = vod.GetModifyTransTemplateGroupStatusEnum().E_0
	}
	updateOpts.Status = &status
	log.Printf("[DEBUG] Update VOD transcoding template group Options: %#v", updateOpts)

	updateReq := vod.UpdateTemplateGroupRequest{
		Body: &updateOpts,
	}

	_, err = client.UpdateTemplateGroup(&updateReq)
	if err != nil {
		return diag.Errorf("error updating VOD transcoding template group: %s", err)
	}

	return resourceTranscodingTemplateGroupRead(ctx, d, meta)
}

func resourceTranscodingTemplateGroupDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.HcVodV1Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating VOD client : %s", err)
	}

	_, err = client.DeleteTemplateGroup(&vod.DeleteTemplateGroupRequest{GroupId: d.Id()})
	if err != nil {
		return diag.Errorf("error deleting VOD transcoding template group: %s", err)
	}

	return nil
}

func setCommonAttrs(d *schema.ResourceData, common *vod.Common) error {
	if common == nil {
		return nil
	}

	var lowBitrateHd bool
	if common.Pvc == vod.GetCommonPvcEnum().E_1 {
		lowBitrateHd = true
	}

	var videoCodec string
	if common.VideoCodec != nil {
		if *common.VideoCodec == vod.GetCommonVideoCodecEnum().H265 {
			videoCodec = "H265"
		}
		if *common.VideoCodec == vod.GetCommonVideoCodecEnum().H264 {
			videoCodec = "H264"
		}
	}

	var audioCodec string
	if common.AudioCodec != nil {
		switch *common.AudioCodec {
		case vod.GetCommonAudioCodecEnum().HEAAC1:
			audioCodec = "HEAAC1"
		case vod.GetCommonAudioCodecEnum().HEAAC2:
			audioCodec = "HEAAC2"
		case vod.GetCommonAudioCodecEnum().MP3:
			audioCodec = "MP3"
		case vod.GetCommonAudioCodecEnum().AAC:
			audioCodec = "AAC"
		default:
			audioCodec = ""
		}
	}

	mErr := multierror.Append(nil,
		d.Set("low_bitrate_hd", lowBitrateHd),
		d.Set("video_codec", videoCodec),
		d.Set("audio_codec", audioCodec),
		d.Set("hls_segment_duration", common.HlsInterval),
	)

	return mErr
}

func flattenAudio(audio *vod.AudioTemplateInfo) []map[string]interface{} {
	if audio == nil {
		return nil
	}

	audioResult := map[string]interface{}{
		"sample_rate": audio.SampleRate,
		"channels":    audio.Channels,
		"bitrate":     audio.Bitrate,
	}

	return []map[string]interface{}{audioResult}
}

func flattenVideo(video *vod.VideoTemplateInfo) []map[string]interface{} {
	if video == nil {
		return nil
	}

	var quality string
	switch video.Quality {
	case vod.GetVideoTemplateInfoQualityEnum().E_4_K:
		quality = "4K"
	case vod.GetVideoTemplateInfoQualityEnum().E_2_K:
		quality = "2K"
	case vod.GetVideoTemplateInfoQualityEnum().FULL_HD:
		quality = "FHD"
	case vod.GetVideoTemplateInfoQualityEnum().SD:
		quality = "SD"
	case vod.GetVideoTemplateInfoQualityEnum().FLUENT:
		quality = "LD"
	case vod.GetVideoTemplateInfoQualityEnum().HD:
		quality = "HD"
	default:
		quality = ""
	}

	videoResult := map[string]interface{}{
		"quality":    quality,
		"bitrate":    video.Bitrate,
		"frame_rate": video.FrameRate,
		"width":      video.Width,
		"height":     video.Height,
	}

	return []map[string]interface{}{videoResult}
}

func flattenQualityInfoList(qualityInfo *[]vod.QualityInfo) []map[string]interface{} {
	if qualityInfo == nil {
		return nil
	}

	QualityInfoResult := make([]map[string]interface{}, len(*qualityInfo))
	for i, info := range *qualityInfo {
		var outputFormat string
		switch info.Format {
		case vod.GetQualityInfoFormatEnum().MP4:
			outputFormat = "MP4"
		case vod.GetQualityInfoFormatEnum().DASH:
			outputFormat = "DASH"
		case vod.GetQualityInfoFormatEnum().DASH_HLS:
			outputFormat = "DASH_HLS"
		case vod.GetQualityInfoFormatEnum().MP3:
			outputFormat = "MP3"
		case vod.GetQualityInfoFormatEnum().ADTS:
			outputFormat = "ADTS"
		case vod.GetQualityInfoFormatEnum().HLS:
			outputFormat = "HLS"
		default:
			outputFormat = ""
		}
		QualityInfoResult[i] = map[string]interface{}{
			"video":         flattenVideo(info.Video),
			"audio":         flattenAudio(info.Audio),
			"output_format": outputFormat,
		}
	}
	return QualityInfoResult
}
