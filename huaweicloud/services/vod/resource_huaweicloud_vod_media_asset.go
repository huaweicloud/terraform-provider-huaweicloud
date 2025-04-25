package vod

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API VOD POST /v1.0/{project_id}/asset/upload_by_url
// @API VOD PUT /v1.0/{project_id}/asset/authority
// @API VOD POST /v1.0/{project_id}/asset/reproduction
// @API VOD GET /v1.0/{project_id}/asset/details
// @API VOD PUT /v1.0/{project_id}/asset/info
// @API VOD POST /v1.0/{project_id}/asset/status/publish
// @API VOD POST /v1.0/{project_id}/asset/status/unpublish
// @API VOD DELETE /v1.0/{project_id}/asset
func ResourceMediaAsset() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMediaAssetCreate,
		ReadContext:   resourceMediaAssetRead,
		UpdateContext: resourceMediaAssetUpdate,
		DeleteContext: resourceMediaAssetDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(60 * time.Second),
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
			"media_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"MP4", "TS", "MOV", "MXF", "MPG", "FLV", "WMV", "AVI", "M4V", "F4V", "MPEG", "3GP", "ASF",
					"MKV", "HLS", "M3U8", "MP3", "OGG", "WAV", "WMA", "APE", "FLAC", "AAC", "AC3", "MMF", "AMR",
					"M4A", "M4R", "WV", "MP2",
				}, false),
			},
			"url": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"input_bucket"},
				ConflictsWith: []string{
					"input_path",
					"output_bucket",
					"output_path",
					"storage_mode",
				},
			},
			"input_bucket": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"input_path"}},
			"input_path": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"output_bucket": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"output_path": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"storage_mode": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"category_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"labels": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"template_group_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ConflictsWith: []string{
					"workflow_name",
				},
			},
			"workflow_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"publish": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"auto_encrypt": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"auto_preload": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"review_template_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"thumbnail": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								"time", "dots",
							}, false),
						},
						"time": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"dots": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem:     &schema.Schema{Type: schema.TypeInt},
						},
						"cover_position": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"format": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"aspect_ratio": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"max_length": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"media_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"media_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"category_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildVodMediaAssetByUrlReviewBodyParams(d *schema.ResourceData) map[string]interface{} {
	templateID := d.Get("review_template_id").(string)
	if templateID == "" {
		return nil
	}

	return map[string]interface{}{
		"template_id": templateID,
	}
}

func buildVodMediaAssetByUrlThumbnailBodyParams(d *schema.ResourceData) map[string]interface{} {
	rawArray := d.Get("thumbnail").([]interface{})
	if len(rawArray) == 0 {
		return nil
	}
	rawMap := rawArray[0].(map[string]interface{})
	return map[string]interface{}{
		"type":           rawMap["type"],
		"time":           utils.ValueIgnoreEmpty(rawMap["time"]),
		"dots":           rawMap["dots"],
		"cover_position": utils.ValueIgnoreEmpty(rawMap["cover_position"]),
		"format":         utils.ValueIgnoreEmpty(rawMap["format"]),
		"aspect_ratio":   utils.ValueIgnoreEmpty(rawMap["aspect_ratio"]),
		"max_length":     utils.ValueIgnoreEmpty(rawMap["max_length"]),
	}
}

func buildVodMediaAssetByUrlAutoPublishParam(d *schema.ResourceData) int {
	if d.Get("publish").(bool) {
		return 1
	}

	return 0
}

func buildVodMediaAssetByUrlAutoEncryptParam(d *schema.ResourceData) interface{} {
	if d.Get("auto_encrypt").(bool) {
		return 1
	}

	return nil
}

func buildVodMediaAssetByUrlAutoPreheatParam(d *schema.ResourceData) interface{} {
	if d.Get("auto_preload").(bool) {
		return 1
	}

	return nil
}

func buildVodMediaAssetByUrlBodyParams(d *schema.ResourceData) map[string]interface{} {
	metadataObject := map[string]interface{}{
		"video_type":          d.Get("media_type"),
		"title":               d.Get("name"),
		"url":                 d.Get("url"),
		"description":         utils.ValueIgnoreEmpty(d.Get("description")),
		"category_id":         utils.ValueIgnoreEmpty(d.Get("category_id")),
		"tags":                utils.ValueIgnoreEmpty(d.Get("labels")),
		"template_group_name": utils.ValueIgnoreEmpty(d.Get("template_group_name")),
		"workflow_name":       utils.ValueIgnoreEmpty(d.Get("workflow_name")),
		"review":              buildVodMediaAssetByUrlReviewBodyParams(d),
		"thumbnail":           buildVodMediaAssetByUrlThumbnailBodyParams(d),
		"auto_publish":        buildVodMediaAssetByUrlAutoPublishParam(d),
		"auto_encrypt":        buildVodMediaAssetByUrlAutoEncryptParam(d),
		"auto_preheat":        buildVodMediaAssetByUrlAutoPreheatParam(d),
	}

	return map[string]interface{}{
		"upload_metadatas": []map[string]interface{}{metadataObject},
	}
}

func createVodMediaAssetByUrl(client *golangsdk.ServiceClient, d *schema.ResourceData) (string, error) {
	requestPath := client.Endpoint + "v1.0/{project_id}/asset/upload_by_url"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildVodMediaAssetByUrlBodyParams(d)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return "", err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return "", err
	}

	assetID := utils.PathSearch("upload_assets|[0].asset_id", respBody, "").(string)
	if assetID == "" {
		return "", errors.New("asset_id is not found in API response")
	}

	return assetID, nil
}

func buildAuthorizeAssetBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"bucket":    d.Get("input_bucket"),
		"operation": "1",
	}
}

func authorizeAsset(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	requestPath := client.Endpoint + "v1.0/{project_id}/asset/authority"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildAuthorizeAssetBodyParams(d),
	}

	_, err := client.Request("PUT", requestPath, &requestOpt)
	return err
}

func buildVodMediaAssetByObsReviewBodyParams(d *schema.ResourceData) map[string]interface{} {
	templateID := d.Get("review_template_id").(string)
	if templateID == "" {
		return nil
	}

	return map[string]interface{}{
		"template_id": templateID,
	}
}

func buildVodMediaAssetByObsThumbnailBodyParams(d *schema.ResourceData) map[string]interface{} {
	rawArray := d.Get("thumbnail").([]interface{})
	if len(rawArray) == 0 {
		return nil
	}

	rawMap := rawArray[0].(map[string]interface{})
	return map[string]interface{}{
		"type":           rawMap["type"],
		"time":           utils.ValueIgnoreEmpty(rawMap["time"]),
		"dots":           rawMap["dots"],
		"cover_position": utils.ValueIgnoreEmpty(rawMap["cover_position"]),
		"format":         utils.ValueIgnoreEmpty(rawMap["format"]),
		"aspect_ratio":   utils.ValueIgnoreEmpty(rawMap["aspect_ratio"]),
		"max_length":     utils.ValueIgnoreEmpty(rawMap["max_length"]),
	}
}

func buildVodMediaAssetByObsInputBodyParams(d *schema.ResourceData, region string) map[string]interface{} {
	return map[string]interface{}{
		"bucket":   d.Get("input_bucket"),
		"object":   d.Get("input_path"),
		"location": region,
	}
}

func buildVodMediaAssetByObsAutoPublishBodyParams(d *schema.ResourceData) int {
	if d.Get("publish").(bool) {
		return 1
	}

	return 0
}

func buildVodMediaAssetByObsAutoEncryptBodyParams(d *schema.ResourceData) interface{} {
	if d.Get("auto_encrypt").(bool) {
		return 1
	}

	return nil
}

func buildVodMediaAssetByObsAutoPreheatBodyParams(d *schema.ResourceData) interface{} {
	if d.Get("auto_preload").(bool) {
		return 1
	}

	return nil
}

func buildVodMediaAssetByObsBodyParams(d *schema.ResourceData, region string) map[string]interface{} {
	return map[string]interface{}{
		"video_type":          d.Get("media_type"),
		"title":               d.Get("name"),
		"description":         utils.ValueIgnoreEmpty(d.Get("description")),
		"category_id":         utils.ValueIgnoreEmpty(d.Get("category_id")),
		"tags":                utils.ValueIgnoreEmpty(d.Get("labels")),
		"template_group_name": utils.ValueIgnoreEmpty(d.Get("template_group_name")),
		"workflow_name":       utils.ValueIgnoreEmpty(d.Get("workflow_name")),
		"review":              utils.ValueIgnoreEmpty(buildVodMediaAssetByObsReviewBodyParams(d)),
		"thumbnail":           buildVodMediaAssetByObsThumbnailBodyParams(d),
		"input":               buildVodMediaAssetByObsInputBodyParams(d, region),
		"output_bucket":       utils.ValueIgnoreEmpty(d.Get("output_bucket")),
		"output_path":         utils.ValueIgnoreEmpty(d.Get("output_path")),
		"storage_mode":        utils.ValueIgnoreEmpty(d.Get("storage_mode")),
		"auto_publish":        buildVodMediaAssetByObsAutoPublishBodyParams(d),
		"auto_encrypt":        buildVodMediaAssetByObsAutoEncryptBodyParams(d),
		"auto_preheat":        buildVodMediaAssetByObsAutoPreheatBodyParams(d),
	}
}

func createVodMediaAssetByObs(client *golangsdk.ServiceClient, d *schema.ResourceData, region string) (string, error) {
	if err := authorizeAsset(client, d); err != nil {
		return "", fmt.Errorf("error authorizing the OBS bucket to VOD: %s", err)
	}

	requestPath := client.Endpoint + "v1.0/{project_id}/asset/reproduction"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildVodMediaAssetByObsBodyParams(d, region)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return "", err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return "", err
	}

	assetID := utils.PathSearch("asset_id", respBody, "").(string)
	if assetID == "" {
		return "", errors.New("asset_id is not found in API response")
	}

	return assetID, nil
}

func resourceMediaAssetCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "vod"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating VOD client: %s", err)
	}

	var assetID string
	if _, ok := d.GetOk("url"); ok {
		assetID, err = createVodMediaAssetByUrl(client, d)
		if err != nil {
			return diag.Errorf("error creating VOD media asset by URL: %s", err)
		}
	} else {
		assetID, err = createVodMediaAssetByObs(client, d, region)
		if err != nil {
			return diag.Errorf("error creating VOD media asset by OBS: %s", err)
		}
	}

	d.SetId(assetID)
	return resourceMediaAssetRead(ctx, d, meta)
}

func ReadMediaAssetDetail(client *golangsdk.ServiceClient, assetID string) (interface{}, error) {
	requestPath := client.Endpoint + "v1.0/{project_id}/asset/details"
	requestPath += fmt.Sprintf("?asset_id=%s", assetID)
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func resourceMediaAssetRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "vod"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating VOD client: %s", err)
	}

	respBody, err := ReadMediaAssetDetail(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving VOD media asset")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", utils.PathSearch("base_info.title", respBody, nil)),
		d.Set("media_name", utils.PathSearch("base_info.video_name", respBody, nil)),
		d.Set("description", utils.PathSearch("base_info.description", respBody, nil)),
		d.Set("category_id", utils.PathSearch("base_info.category_id", respBody, nil)),
		d.Set("category_name", utils.PathSearch("base_info.category_name", respBody, nil)),
		d.Set("media_type", utils.PathSearch("base_info.video_type", respBody, nil)),
		d.Set("labels", utils.PathSearch("base_info.tags", respBody, nil)),
		d.Set("media_url", utils.PathSearch("base_info.video_url", respBody, nil)),
	)

	sourcePath := utils.PathSearch("base_info.source_path", respBody, nil)
	if sourcePath != nil {
		mErr = multierror.Append(mErr,
			d.Set("input_bucket", utils.PathSearch("bucket", sourcePath, nil)),
			d.Set("input_path", utils.PathSearch("object", sourcePath, nil)),
		)
	}

	outputPath := utils.PathSearch("base_info.output_path", respBody, nil)
	if outputPath != nil {
		mErr = multierror.Append(mErr,
			d.Set("output_bucket", utils.PathSearch("bucket", outputPath, nil)),
			d.Set("output_path", utils.PathSearch("object", outputPath, nil)),
		)
	}

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildUpdateAssetMetaBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"asset_id":    d.Id(),
		"title":       d.Get("name"),
		"description": d.Get("description"),
		"category_id": d.Get("category_id"),
		"tags":        d.Get("labels"),
	}
}

func updateAssetMeta(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	requestPath := client.Endpoint + "v1.0/{project_id}/asset/info"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes:          []int{200, 201, 204},
		JSONBody:         buildUpdateAssetMetaBodyParams(d),
	}

	_, err := client.Request("PUT", requestPath, &requestOpt)
	return err
}

func publishAssets(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	requestPath := client.Endpoint + "v1.0/{project_id}/asset/status/publish"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"asset_id": []string{d.Id()},
		},
	}

	_, err := client.Request("POST", requestPath, &requestOpt)
	return err
}

func unPublishAssets(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	requestPath := client.Endpoint + "v1.0/{project_id}/asset/status/unpublish"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"asset_id": []string{d.Id()},
		},
	}

	_, err := client.Request("POST", requestPath, &requestOpt)
	return err
}

func resourceMediaAssetUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "vod"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating VOD client: %s", err)
	}

	if d.HasChanges("name", "description", "category_id", "labels") {
		if err := updateAssetMeta(client, d); err != nil {
			return diag.Errorf("error updating VOD media asset: %s", err)
		}
	}

	if d.HasChange("publish") {
		if d.Get("publish").(bool) {
			if err := publishAssets(client, d); err != nil {
				return diag.Errorf("error publishing VOD media asset: %s", err)
			}
		} else {
			if err := unPublishAssets(client, d); err != nil {
				return diag.Errorf("error un-publishing VOD media asset: %s", err)
			}
		}
	}

	return resourceMediaAssetRead(ctx, d, meta)
}

func waitingForMediaAssetDelete(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			respBody, err := ReadMediaAssetDetail(client, d.Id())
			if err != nil {
				var errDefault404 golangsdk.ErrDefault404
				if errors.As(err, &errDefault404) {
					return "deleted", "COMPLETED", nil
				}
				return nil, "ERROR", err
			}

			return respBody, "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceMediaAssetDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "vod"
		httpUrl = "v1.0/{project_id}/asset"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating VOD client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath += fmt.Sprintf("?asset_id=%s", d.Id())
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error deleting VOD media asset: %s", err)
	}

	if err := waitingForMediaAssetDelete(ctx, client, d, d.Timeout(schema.TimeoutDelete)); err != nil {
		return diag.Errorf("error waiting for VOD media asset (%s) deleted: %s", d.Id(), err)
	}

	return nil
}
