package live

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/live/v1/model"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const (
	recordingTypeAuto   = "CONTINUOUS_RECORD"
	recordingTypeManual = "COMMAND_RECORD"
)

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
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"recording_length": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(15, 720), // same with console
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

						"max_stream_pause_length": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(-1, 300),
							Computed:     true,
						},
					},
				},
			},

			"flv": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     formatSchema(),
			},

			"mp4": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     formatSchema(),
			},

			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.StringInSlice([]string{recordingTypeAuto, recordingTypeManual}, false),
			},
		},
	}
}

func formatSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"recording_length": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(15, 180), // same with console
			},

			"file_naming": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"max_stream_pause_length": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(0, 300),
				Computed:     true,
			},
		},
	}
}

func resourceRecordingCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	region := c.GetRegion(d)
	client, err := c.HcLiveV1Client(region)
	if err != nil {
		return diag.Errorf("error creating Live v1 client: %s", err)
	}

	createBody, err := buildRecordingParams(d)
	if err != nil {
		return diag.FromErr(err)
	}
	createOpts := &model.CreateRecordRuleRequest{
		Body: createBody,
	}

	log.Printf("[DEBUG] Create Live recording params: %#v", createOpts)

	resp, err := client.CreateRecordRule(createOpts)
	if err != nil {
		return diag.Errorf("error creating Live recording: %s", err)
	}

	if resp.Id == nil {
		return diag.Errorf("error creating Live recording: ID not found in the API response")
	}

	d.SetId(*resp.Id)

	return resourceRecordingRead(ctx, d, meta)
}

func resourceRecordingRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	region := c.GetRegion(d)
	client, err := c.HcLiveV1Client(region)
	if err != nil {
		return diag.Errorf("error creating Live v1 client: %s", err)
	}

	response, err := client.ShowRecordRule(&model.ShowRecordRuleRequest{Id: d.Id()})
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Live recording")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("domain_name", response.PublishDomain),
		d.Set("app_name", response.App),
		d.Set("stream_name", response.Stream),
		d.Set("type", utils.MarshalValue(response.RecordType)),
		d.Set("obs", flattenObs(response.DefaultRecordConfig)),
		d.Set("hls", flattenHls(response.DefaultRecordConfig)),
		d.Set("flv", flattenFlv(response.DefaultRecordConfig)),
		d.Set("mp4", flattenMp4(response.DefaultRecordConfig)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceRecordingUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	region := c.GetRegion(d)
	client, err := c.HcLiveV1Client(region)
	if err != nil {
		return diag.Errorf("error creating Live v1 client: %s", err)
	}

	updateBody, err := buildRecordingParams(d)
	if err != nil {
		return diag.FromErr(err)
	}
	updateOpts := &model.UpdateRecordRuleRequest{
		Id:   d.Id(),
		Body: updateBody,
	}

	_, err = client.UpdateRecordRule(updateOpts)
	if err != nil {
		return diag.Errorf("error updating Live recording: %s", err)
	}

	return resourceRecordingRead(ctx, d, meta)
}

func resourceRecordingDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	region := c.GetRegion(d)
	client, err := c.HcLiveV1Client(region)
	if err != nil {
		return diag.Errorf("error creating Live v1 client: %s", err)
	}

	deleteOpts := &model.DeleteRecordRuleRequest{
		Id: d.Id(),
	}
	_, err = client.DeleteRecordRule(deleteOpts)
	if err != nil {
		return diag.Errorf("error deleting Live recording: %s", err)
	}

	return nil
}

func buildRecordingParams(d *schema.ResourceData) (*model.RecordRuleRequest, error) {
	var recordType model.RecordRuleRequestRecordType
	if v, ok := d.GetOk("type"); ok {
		if err := recordType.UnmarshalJSON([]byte(v.(string))); err != nil {
			return nil, fmt.Errorf("error parsing the argument %q: %s", "type", err)
		}
	}

	r := d.Get("obs.0.region").(string)
	var region model.RecordObsFileAddrLocation
	if err := region.UnmarshalJSON([]byte(r)); err != nil {
		return nil, fmt.Errorf("error parsing the argument %q: %s", "obs region", err)
	}
	obs := model.RecordObsFileAddr{
		Location: region,
		Bucket:   d.Get("obs.0.bucket").(string),
		Object:   d.Get("obs.0.object").(string),
	}

	recordFormat := []model.VideoFormatVar{}
	var hlsConfig *model.HlsRecordConfig
	if _, ok := d.GetOk("hls"); ok {
		recordFormat = append(recordFormat, model.GetVideoFormatVarEnum().HLS)
		hlsConfig = &model.HlsRecordConfig{
			RecordCycle:                  int32(d.Get("hls.0.recording_length").(int) * 60),
			RecordPrefix:                 utils.String(d.Get("hls.0.file_naming").(string)),
			RecordTsPrefix:               utils.String(d.Get("hls.0.ts_file_naming").(string)),
			RecordMaxDurationToMergeFile: utils.Int32(int32(d.Get("hls.0.max_stream_pause_length").(int))),
		}
	}
	var flvConfig *model.FlvRecordConfig
	if _, ok := d.GetOk("flv"); ok {
		recordFormat = append(recordFormat, model.GetVideoFormatVarEnum().FLV)
		flvConfig = &model.FlvRecordConfig{
			RecordCycle:                  int32(d.Get("flv.0.recording_length").(int) * 60),
			RecordPrefix:                 utils.String(d.Get("flv.0.file_naming").(string)),
			RecordMaxDurationToMergeFile: utils.Int32(int32(d.Get("flv.0.max_stream_pause_length").(int))),
		}
	}
	var mp4Config *model.Mp4RecordConfig
	if _, ok := d.GetOk("mp4"); ok {
		recordFormat = append(recordFormat, model.GetVideoFormatVarEnum().MP4)
		mp4Config = &model.Mp4RecordConfig{
			RecordCycle:                  int32(d.Get("mp4.0.recording_length").(int) * 60),
			RecordPrefix:                 utils.String(d.Get("mp4.0.file_naming").(string)),
			RecordMaxDurationToMergeFile: utils.Int32(int32(d.Get("mp4.0.max_stream_pause_length").(int))),
		}
	}

	req := model.RecordRuleRequest{
		PublishDomain: d.Get("domain_name").(string),
		App:           d.Get("app_name").(string),
		Stream:        d.Get("stream_name").(string),
		RecordType:    &recordType,
		DefaultRecordConfig: &model.DefaultRecordConfig{
			ObsAddr:      &obs,
			RecordFormat: recordFormat,
			HlsConfig:    hlsConfig,
			FlvConfig:    flvConfig,
			Mp4Config:    mp4Config,
		},
	}

	return &req, nil
}

func flattenObs(r *model.DefaultRecordConfig) []interface{} {
	if r != nil && r.ObsAddr != nil {
		m := map[string]interface{}{
			"region": utils.MarshalValue(r.ObsAddr.Location),
			"bucket": r.ObsAddr.Bucket,
			"object": r.ObsAddr.Object,
		}
		return []interface{}{m}
	}

	return make([]interface{}, 0)
}

func flattenHls(r *model.DefaultRecordConfig) []interface{} {
	if r != nil && r.HlsConfig != nil && r.HlsConfig.RecordPrefix != nil && len(*r.HlsConfig.RecordPrefix) > 0 {
		m := map[string]interface{}{
			"recording_length":        r.HlsConfig.RecordCycle / 60,
			"file_naming":             r.HlsConfig.RecordPrefix,
			"ts_file_naming":          r.HlsConfig.RecordTsPrefix,
			"max_stream_pause_length": r.HlsConfig.RecordMaxDurationToMergeFile,
		}
		return []interface{}{m}
	}

	return make([]interface{}, 0)
}

func flattenFlv(r *model.DefaultRecordConfig) []interface{} {
	if r != nil && r.FlvConfig != nil && r.FlvConfig.RecordPrefix != nil && len(*r.FlvConfig.RecordPrefix) > 0 {
		m := map[string]interface{}{
			"recording_length":        r.FlvConfig.RecordCycle / 60,
			"file_naming":             r.FlvConfig.RecordPrefix,
			"max_stream_pause_length": r.FlvConfig.RecordMaxDurationToMergeFile,
		}
		return []interface{}{m}
	}

	return make([]interface{}, 0)
}

func flattenMp4(r *model.DefaultRecordConfig) []interface{} {
	if r != nil && r.Mp4Config != nil && r.Mp4Config.RecordPrefix != nil && len(*r.Mp4Config.RecordPrefix) > 0 {
		m := map[string]interface{}{
			"recording_length":        r.Mp4Config.RecordCycle / 60,
			"file_naming":             r.Mp4Config.RecordPrefix,
			"max_stream_pause_length": r.Mp4Config.RecordMaxDurationToMergeFile,
		}
		return []interface{}{m}
	}

	return make([]interface{}, 0)
}
