package vod

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	vod "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/vod/v1/model"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func ResourceWatermarkTemplate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWatermarkTemplateCreate,
		ReadContext:   resourceWatermarkTemplateRead,
		UpdateContext: resourceWatermarkTemplateUpdate,
		DeleteContext: resourceWatermarkTemplateDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
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
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 128),
			},
			"image_file": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"image_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"PNG", "JPG", "JPEG",
				}, false),
			},
			"image_process": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ORIGINAL", "TRANSPARENT", "GRAYED",
				}, false),
				Default: "TRANSPARENT",
			},
			"horizontal_offset": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "0",
			},
			"vertical_offset": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "0",
			},
			"position": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"TOPRIGHT", "TOPLEFT", "BOTTOMRIGHT", "BOTTOMLEFT",
				}, false),
				Default: "TOPRIGHT",
			},
			"width": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "0.01",
			},
			"height": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "0.01",
			},
			"timeline_start": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"timeline_duration": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"watermark_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"image_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildCreateImageProcessOpts(imageProcess string) *vod.CreateWatermarkTemplateReqImageProcess {
	var imageProcessResult vod.CreateWatermarkTemplateReqImageProcess
	switch imageProcess {
	case "ORIGINAL":
		imageProcessResult = vod.GetCreateWatermarkTemplateReqImageProcessEnum().ORIGINAL
	case "TRANSPARENT":
		imageProcessResult = vod.GetCreateWatermarkTemplateReqImageProcessEnum().TRANSPARENT
	case "GRAYED":
		imageProcessResult = vod.GetCreateWatermarkTemplateReqImageProcessEnum().GRAYED
	default:
		return nil
	}
	return &imageProcessResult
}

func buildCreatePositionOpts(position string) *vod.CreateWatermarkTemplateReqPosition {
	var positionResult vod.CreateWatermarkTemplateReqPosition
	switch position {
	case "TOPRIGHT":
		positionResult = vod.GetCreateWatermarkTemplateReqPositionEnum().TOPRIGHT
	case "TOPLEFT":
		positionResult = vod.GetCreateWatermarkTemplateReqPositionEnum().TOPLEFT
	case "BOTTOMLEFT":
		positionResult = vod.GetCreateWatermarkTemplateReqPositionEnum().BOTTOMLEFT
	case "BOTTOMRIGHT":
		positionResult = vod.GetCreateWatermarkTemplateReqPositionEnum().BOTTOMRIGHT
	default:
		return nil
	}
	return &positionResult
}

func uploadImage(uploadUrl, fileName string, timeout time.Duration) error {
	data, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer data.Close()
	req, err := http.NewRequest("PUT", uploadUrl, data)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "text/plain")

	httpClient := &http.Client{Timeout: timeout}
	res, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}

func resourceWatermarkTemplateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.HcVodV1Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating VOD client: %s", err)
	}

	createOpts := vod.CreateWatermarkTemplateReq{
		Name:             d.Get("name").(string),
		Type:             d.Get("image_type").(string),
		ImageProcess:     buildCreateImageProcessOpts(d.Get("image_process").(string)),
		Dx:               utils.String(d.Get("horizontal_offset").(string)),
		Dy:               utils.String(d.Get("vertical_offset").(string)),
		Position:         buildCreatePositionOpts(d.Get("position").(string)),
		Width:            utils.String(d.Get("width").(string)),
		Height:           utils.String(d.Get("height").(string)),
		TimelineStart:    utils.String(d.Get("timeline_start").(string)),
		TimelineDuration: utils.String(d.Get("timeline_duration").(string)),
	}
	log.Printf("[DEBUG] Create VOD watermark template Options: %#v", createOpts)

	createReq := vod.CreateWatermarkTemplateRequest{
		Body: &createOpts,
	}

	resp, err := client.CreateWatermarkTemplate(&createReq)
	if err != nil {
		return diag.Errorf("error creating VOD watermark template: %s", err)
	}

	d.SetId(*resp.Id)

	err = uploadImage(*resp.UploadUrl, d.Get("image_file").(string), d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error uploading watermark image: %s", err)
	}

	confirmImageUploadReq := vod.ConfirmImageUploadReq{
		Id:     *resp.Id,
		Status: vod.GetConfirmImageUploadReqStatusEnum().SUCCEED,
	}
	_, err = client.ConfirmImageUpload(&vod.ConfirmImageUploadRequest{Body: &confirmImageUploadReq})
	if err != nil {
		return diag.Errorf("error condfirming watermark image upload: %s", err)
	}

	return resourceWatermarkTemplateRead(ctx, d, meta)
}

func resourceWatermarkTemplateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.HcVodV1Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating VOD client: %s", err)
	}

	resp, err := client.ListWatermarkTemplate(&vod.ListWatermarkTemplateRequest{Id: &[]string{d.Id()}})
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving VOD watermark template")
	}

	if resp.Templates == nil || len(*resp.Templates) == 0 {
		d.SetId("")
		return diag.Diagnostics{
			diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Resource not found",
				Detail:   fmt.Sprintf("unable to retrieve VOD watermark template: %s", d.Id()),
			},
		}
	}

	templateList := *resp.Templates
	template := templateList[0]

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("name", template.Name),
		d.Set("image_type", template.Type),
		d.Set("image_process", template.ImageProcess),
		d.Set("horizontal_offset", template.Dx),
		d.Set("vertical_offset", template.Dy),
		d.Set("position", template.Position),
		d.Set("width", template.Width),
		d.Set("height", template.Height),
		d.Set("timeline_start", template.TimelineStart),
		d.Set("timeline_duration", template.TimelineDuration),
		d.Set("watermark_type", template.WatermarkType),
		d.Set("image_url", template.ImageUrl),
	)

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting VOD watermark template fields: %s", err)
	}

	return nil
}

func buildUpdateImageProcessOpts(imageProcess string) *vod.UpdateWatermarkTemplateReqImageProcess {
	var imageProcessResult vod.UpdateWatermarkTemplateReqImageProcess
	switch imageProcess {
	case "ORIGINAL":
		imageProcessResult = vod.GetUpdateWatermarkTemplateReqImageProcessEnum().ORIGINAL
	case "TRANSPARENT":
		imageProcessResult = vod.GetUpdateWatermarkTemplateReqImageProcessEnum().TRANSPARENT
	case "GRAYED":
		imageProcessResult = vod.GetUpdateWatermarkTemplateReqImageProcessEnum().GRAYED
	default:
		return nil
	}
	return &imageProcessResult
}

func buildUpdatePositionOpts(position string) *vod.UpdateWatermarkTemplateReqPosition {
	var positionResult vod.UpdateWatermarkTemplateReqPosition
	switch position {
	case "TOPRIGHT":
		positionResult = vod.GetUpdateWatermarkTemplateReqPositionEnum().TOPRIGHT
	case "TOPLEFT":
		positionResult = vod.GetUpdateWatermarkTemplateReqPositionEnum().TOPLEFT
	case "BOTTOMLEFT":
		positionResult = vod.GetUpdateWatermarkTemplateReqPositionEnum().BOTTOMLEFT
	case "BOTTOMRIGHT":
		positionResult = vod.GetUpdateWatermarkTemplateReqPositionEnum().BOTTOMRIGHT
	default:
		return nil
	}
	return &positionResult
}

func resourceWatermarkTemplateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.HcVodV1Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating VOD client: %s", err)
	}

	updateOpts := vod.UpdateWatermarkTemplateReq{
		Id:               d.Id(),
		Name:             utils.String(d.Get("name").(string)),
		ImageProcess:     buildUpdateImageProcessOpts(d.Get("image_process").(string)),
		Dx:               utils.String(d.Get("horizontal_offset").(string)),
		Dy:               utils.String(d.Get("vertical_offset").(string)),
		Position:         buildUpdatePositionOpts(d.Get("position").(string)),
		Width:            utils.String(d.Get("width").(string)),
		Height:           utils.String(d.Get("height").(string)),
		TimelineStart:    utils.String(d.Get("timeline_start").(string)),
		TimelineDuration: utils.String(d.Get("timeline_duration").(string)),
	}
	log.Printf("[DEBUG] Update VOD watermark template Options: %#v", updateOpts)

	updateReq := vod.UpdateWatermarkTemplateRequest{
		Body: &updateOpts,
	}

	_, err = client.UpdateWatermarkTemplate(&updateReq)
	if err != nil {
		return diag.Errorf("error updating VOD watermark template: %s", err)
	}

	return resourceWatermarkTemplateRead(ctx, d, meta)
}

func resourceWatermarkTemplateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.HcVodV1Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating VOD client: %s", err)
	}

	_, err = client.DeleteWatermarkTemplate(&vod.DeleteWatermarkTemplateRequest{Id: d.Id()})
	if err != nil {
		return diag.Errorf("error deleting VOD watermark template: %s", err)
	}

	return nil
}
