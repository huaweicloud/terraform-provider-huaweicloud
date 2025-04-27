package vod

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API VOD POST /v1.0/{project_id}/template/watermark
// @API VOD POST /v1.0/{project_id}/watermark/status/uploaded
// @API VOD GET /v1.0/{project_id}/template/watermark
// @API VOD PUT /v1.0/{project_id}/template/watermark
// @API VOD DELETE /v1.0/{project_id}/template/watermark
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
				Type:     schema.TypeString,
				Required: true,
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

func buildWatermarkTemplateBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name":              d.Get("name"),
		"type":              d.Get("image_type"),
		"image_process":     utils.ValueIgnoreEmpty(d.Get("image_process")),
		"dx":                d.Get("horizontal_offset"),
		"dy":                d.Get("vertical_offset"),
		"position":          utils.ValueIgnoreEmpty(d.Get("position")),
		"width":             d.Get("width"),
		"height":            d.Get("height"),
		"timeline_start":    d.Get("timeline_start"),
		"timeline_duration": d.Get("timeline_duration"),
	}
}

func createWatermarkTemplate(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	requestPath := client.Endpoint + "v1.0/{project_id}/template/watermark"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildWatermarkTemplateBodyParams(d)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func confirmImageUpload(client *golangsdk.ServiceClient, id string) error {
	requestPath := client.Endpoint + "v1.0/{project_id}/watermark/status/uploaded"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"id":     id,
			"status": "SUCCEED",
		},
	}

	_, err := client.Request("POST", requestPath, &requestOpt)
	return err
}

func resourceWatermarkTemplateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "vod"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating VOD client: %s", err)
	}

	respBody, err := createWatermarkTemplate(client, d)
	if err != nil {
		return diag.Errorf("error creating VOD watermark template: %s", err)
	}

	id := utils.PathSearch("id", respBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating VOD watermark template: ID is not found in API response")
	}
	d.SetId(id)

	uploadUrl := utils.PathSearch("upload_url", respBody, "").(string)
	if uploadUrl == "" {
		return diag.Errorf("error creating VOD watermark template: upload_url is not found in API response")
	}

	err = uploadImage(uploadUrl, d.Get("image_file").(string), d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error uploading watermark image: %s", err)
	}

	if err := confirmImageUpload(client, d.Id()); err != nil {
		return diag.Errorf("error confirming watermark image upload: %s", err)
	}

	return resourceWatermarkTemplateRead(ctx, d, meta)
}

func resourceWatermarkTemplateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1.0/{project_id}/template/watermark"
		product = "vod"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating VOD client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath += fmt.Sprintf("?id=%s", d.Id())
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving VOD watermark template")
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	template := utils.PathSearch("templates|[0]", respBody, nil)
	if template == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", template, nil)),
		d.Set("image_type", utils.PathSearch("type", template, nil)),
		d.Set("image_process", utils.PathSearch("image_process", template, nil)),
		d.Set("horizontal_offset", utils.PathSearch("dx", template, nil)),
		d.Set("vertical_offset", utils.PathSearch("dy", template, nil)),
		d.Set("position", utils.PathSearch("position", template, nil)),
		d.Set("width", utils.PathSearch("width", template, nil)),
		d.Set("height", utils.PathSearch("height", template, nil)),
		d.Set("timeline_start", utils.PathSearch("timeline_start", template, nil)),
		d.Set("timeline_duration", utils.PathSearch("timeline_duration", template, nil)),
		d.Set("watermark_type", utils.PathSearch("watermark_type", template, nil)),
		d.Set("image_url", utils.PathSearch("image_url", template, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildUpdateWatermarkTemplateBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"id":                d.Id(),
		"name":              d.Get("name"),
		"image_process":     utils.ValueIgnoreEmpty(d.Get("image_process")),
		"dx":                d.Get("horizontal_offset"),
		"dy":                d.Get("vertical_offset"),
		"position":          utils.ValueIgnoreEmpty(d.Get("position")),
		"width":             d.Get("width"),
		"height":            d.Get("height"),
		"timeline_start":    d.Get("timeline_start"),
		"timeline_duration": d.Get("timeline_duration"),
	}
}

func resourceWatermarkTemplateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1.0/{project_id}/template/watermark"
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
		JSONBody:         utils.RemoveNil(buildUpdateWatermarkTemplateBodyParams(d)),
	}

	_, err = client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error updating VOD watermark template: %s", err)
	}

	return resourceWatermarkTemplateRead(ctx, d, meta)
}

func resourceWatermarkTemplateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1.0/{project_id}/template/watermark"
		product = "vod"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating VOD client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath += fmt.Sprintf("?id=%s", d.Id())
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error deleting VOD watermark template: %s", err)
	}

	return nil
}
