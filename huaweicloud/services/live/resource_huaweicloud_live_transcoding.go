package live

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/live/v1/model"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func ResourceTranscoding() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTranscodingCreate,
		UpdateContext: resourceTranscodingUpdate,
		DeleteContext: resourceTranscodingDelete,
		ReadContext:   resourceTranscodingRead,
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
				ForceNew: true,
			},

			"app_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"video_encoding": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"H264", "H265"}, false),
			},

			"templates": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 4,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringMatch(regexp.MustCompile("[a-zA-Z0-9-]{1,64}$"),
								"The name can contain a maximum of 64 characters, and only contains letters,"+
									" digits and hyphens (-)."),
						},

						"width": {
							Type:     schema.TypeInt,
							Required: true,
						},

						"height": {
							Type:     schema.TypeInt,
							Required: true,
						},

						"bitrate": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(40, 30000),
						},

						"frame_rate": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(0, 30),
							Computed:     true,
						},
					},
				},
			},

			"low_bitrate_hd": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceTranscodingCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	region := c.GetRegion(d)
	client, err := c.HcLiveV1Client(region)
	if err != nil {
		return diag.Errorf("error creating Live v1 client: %s", err)
	}

	transcodingParams, err := buildTranscodingParams(d)
	if err != nil {
		return diag.FromErr(err)
	}

	createOpts := &model.CreateTranscodingsTemplateRequest{
		Body: transcodingParams,
	}
	log.Printf("[DEBUG] Create Live transcoding params: %#v", createOpts)

	_, err = client.CreateTranscodingsTemplate(createOpts)
	if err != nil {
		return diag.Errorf("error creating Live transcoding: %s", err)
	}

	d.SetId(fmt.Sprintf("%s/%s", createOpts.Body.Domain, createOpts.Body.AppName))

	return resourceTranscodingRead(ctx, d, meta)
}

func resourceTranscodingRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	region := c.GetRegion(d)
	client, err := c.HcLiveV1Client(region)
	if err != nil {
		return diag.Errorf("error creating Live v1 client: %s", err)
	}

	domain, appName, err := parseTranscodingId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	response, err := client.ShowTranscodingsTemplate(&model.ShowTranscodingsTemplateRequest{
		Domain:  domain,
		AppName: &appName,
	})
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Live transcoding")
	}

	if response.Templates == nil || len(*response.Templates) != 1 {
		return diag.Errorf("error retrieving Live transcoding")
	}
	r := *response.Templates
	detail := r[0]

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("domain_name", response.Domain),
		d.Set("app_name", detail.AppName),
		setTemplatesToState(d, detail.QualityInfo),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceTranscodingUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	region := c.GetRegion(d)
	client, err := c.HcLiveV1Client(region)
	if err != nil {
		return diag.Errorf("error creating Live v1 client: %s", err)
	}

	transcodingParams, err := buildTranscodingParams(d)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = client.UpdateTranscodingsTemplate(&model.UpdateTranscodingsTemplateRequest{
		Body: transcodingParams,
	})

	if err != nil {
		return diag.Errorf("error updating Live transcoding: %s", err)
	}

	return resourceTranscodingRead(ctx, d, meta)
}

func resourceTranscodingDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	region := c.GetRegion(d)
	client, err := c.HcLiveV1Client(region)
	if err != nil {
		return diag.Errorf("error creating Live v1 client: %s", err)
	}

	deleteOpts := &model.DeleteTranscodingsTemplateRequest{
		Domain:  d.Get("domain_name").(string),
		AppName: d.Get("app_name").(string),
	}
	_, err = client.DeleteTranscodingsTemplate(deleteOpts)
	if err != nil {
		return diag.Errorf("error deleting Live transcoding: %s", err)
	}

	return nil
}

func buildTranscodingParams(d *schema.ResourceData) (*model.StreamTranscodingTemplate, error) {
	var codec model.QualityInfoCodec
	if err := codec.UnmarshalJSON([]byte(d.Get("video_encoding").(string))); err != nil {
		return nil, fmt.Errorf("error parsing the argument %q: %s", "video_encoding", err)
	}

	var hdlb model.QualityInfoHdlb
	if v, ok := d.GetOk("low_bitrate_hd"); ok {
		if v.(bool) {
			hdlb = model.GetQualityInfoHdlbEnum().ON
		} else {
			hdlb = model.GetQualityInfoHdlbEnum().OFF
		}
	}

	templates := d.Get("templates").([]interface{})
	qualityInfo := make([]model.QualityInfo, len(templates))
	for i, v := range templates {
		template := v.(map[string]interface{})

		width := template["width"].(int)
		height := template["height"].(int)

		errFmt := "expected %s to be in the range (%d - %d) and must be a multiple of %d when " +
			"video_encoding is %s, got %d"

		if codec == model.GetQualityInfoCodecEnum().H264 {
			if width < 32 || width > 3840 || width%2 != 0 {
				return nil, fmt.Errorf(errFmt, "width", 32, 3840, 2, "H264", width)
			}
			if height < 32 || height > 2160 || height%2 != 0 {
				return nil, fmt.Errorf(errFmt, "height", 32, 2160, 2, "H264", height)
			}
		} else if codec == model.GetQualityInfoCodecEnum().H265 {
			if width < 320 || width > 3840 || width%4 != 0 {
				return nil, fmt.Errorf(errFmt, "width", 320, 3840, 4, "H265", width)
			}
			if height < 240 || height > 2160 || height%4 != 0 {
				return nil, fmt.Errorf(errFmt, "height", 240, 2160, 4, "H265", height)
			}
		}

		qualityInfo[i] = model.QualityInfo{
			TemplateName:   utils.String(template["name"].(string)),
			Quality:        template["name"].(string),
			Width:          int32(width),
			Height:         int32(height),
			Bitrate:        int32(template["bitrate"].(int)),
			VideoFrameRate: utils.Int32(int32(template["frame_rate"].(int))),
			Codec:          &codec,
			Hdlb:           &hdlb,
		}
	}

	req := model.StreamTranscodingTemplate{
		Domain:      d.Get("domain_name").(string),
		AppName:     d.Get("app_name").(string),
		QualityInfo: qualityInfo,
	}
	return &req, nil
}

func setTemplatesToState(d *schema.ResourceData, qualityInfo *[]model.QualityInfo) error {
	if qualityInfo != nil || len(*qualityInfo) > 0 {
		qualitys := *qualityInfo
		rst := make([]map[string]interface{}, len(qualitys))
		for i, v := range qualitys {
			rst[i] = map[string]interface{}{
				"name":       v.TemplateName,
				"width":      v.Width,
				"height":     v.Height,
				"bitrate":    v.Bitrate,
				"frame_rate": v.VideoFrameRate,
			}
		}

		var hdlb bool
		if utils.MarshalValue(qualitys[0].Hdlb) == "on" {
			hdlb = true
		}

		mErr := multierror.Append(
			d.Set("templates", rst),
			// this two attribute is same in one template group
			d.Set("video_encoding", utils.MarshalValue(qualitys[0].Codec)),
			d.Set("low_bitrate_hd", hdlb),
		)

		return mErr.ErrorOrNil()
	}
	return nil
}

func parseTranscodingId(id string) (string, string, error) {
	idArrays := strings.SplitN(id, "/", 2)
	if len(idArrays) != 2 {
		return "", "", fmt.Errorf("invalid format specified for import ID. Format must be <domain_name>/<app_name>")
	}
	return idArrays[0], idArrays[1], nil
}
