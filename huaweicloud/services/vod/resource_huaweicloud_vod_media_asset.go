package vod

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	v1 "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/vod/v1"
	vod "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/vod/v1/model"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func ResourceMediaAsset() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMediaAssetCreate,
		ReadContext:   resourceMediaAssetRead,
		UpdateContext: resourceMediaAssetUpdate,
		DeleteContext: resourceMediaAssetDelete,
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
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 128),
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
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 1024),
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
							Type:         schema.TypeInt,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.IntBetween(1, 12),
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

func buildReviewOpts(reviewTemplateId string) *vod.Review {
	if reviewTemplateId == "" {
		return nil
	}

	return &vod.Review{
		TemplateId: reviewTemplateId,
	}
}

func buildThumbnailOpts(rawThumbnail []interface{}) *vod.Thumbnail {
	if len(rawThumbnail) != 1 {
		return nil
	}

	thumbnail := rawThumbnail[0].(map[string]interface{})

	var thumbnailType vod.ThumbnailType
	if thumbnail["type"].(string) == "time" {
		thumbnailType = vod.GetThumbnailTypeEnum().TIME
	}
	if thumbnail["type"].(string) == "dots" {
		thumbnailType = vod.GetThumbnailTypeEnum().DOTS
	}

	thumbnailOpts := vod.Thumbnail{
		Type:          thumbnailType,
		Time:          utils.Int32IgnoreEmpty(int32(thumbnail["time"].(int))),
		Dots:          utils.ExpandToInt32ListPointer(thumbnail["dots"].([]interface{})),
		CoverPosition: utils.Int32IgnoreEmpty(int32(thumbnail["cover_position"].(int))),
		Format:        utils.Int32IgnoreEmpty(int32(thumbnail["format"].(int))),
		AspectRatio:   utils.Int32IgnoreEmpty(int32(thumbnail["aspect_ratio"].(int))),
		MaxLength:     utils.Int32IgnoreEmpty(int32(thumbnail["max_length"].(int))),
	}

	return &thumbnailOpts
}

func createMediaAssetByUrl(client *v1.VodClient, d *schema.ResourceData) (string, error) {
	var videoType vod.UploadMetaDataByUrlVideoType
	if err := videoType.UnmarshalJSON([]byte(d.Get("media_type").(string))); err != nil {
		return "", fmt.Errorf("error parsing the argument media_type: %s", err)
	}

	createOpts := vod.UploadMetaDataByUrl{
		VideoType:         videoType,
		Title:             d.Get("name").(string),
		Url:               d.Get("url").(string),
		Description:       utils.StringIgnoreEmpty(d.Get("description").(string)),
		CategoryId:        utils.Int32IgnoreEmpty(int32(d.Get("category_id").(int))),
		Tags:              utils.StringIgnoreEmpty(d.Get("labels").(string)),
		TemplateGroupName: utils.StringIgnoreEmpty(d.Get("template_group_name").(string)),
		WorkflowName:      utils.StringIgnoreEmpty(d.Get("workflow_name").(string)),
		Review:            buildReviewOpts(d.Get("review_template_id").(string)),
		Thumbnail:         buildThumbnailOpts(d.Get("thumbnail").([]interface{})),
	}

	if d.Get("publish").(bool) {
		createOpts.AutoPublish = utils.Int32(int32(1))
	} else {
		createOpts.AutoPublish = utils.Int32(int32(0))
	}
	if d.Get("auto_encrypt").(bool) {
		createOpts.AutoEncrypt = utils.Int32(int32(1))
	}
	if d.Get("auto_preload").(bool) {
		createOpts.AutoPreheat = utils.Int32(int32(1))
	}

	uploadList := vod.UploadMetaDataByUrlReq{
		UploadMetadatas: []vod.UploadMetaDataByUrl{createOpts},
	}
	createReq := vod.UploadMetaDataByUrlRequest{
		Body: &uploadList,
	}
	resp, err := client.UploadMetaDataByUrl(&createReq)
	if err != nil {
		return "", fmt.Errorf("error creating VOD media asset: %s", err)
	}

	if resp.UploadAssets == nil {
		return "", fmt.Errorf("unable to find the asset after uploading")
	}
	assets := *resp.UploadAssets
	if len(assets) == 0 || assets[0].AssetId == nil {
		return "", fmt.Errorf("unable to find the asset after uploading")
	}
	return *assets[0].AssetId, nil
}

func createMediaAssetFromObs(client *v1.VodClient, d *schema.ResourceData, region string) (string, error) {
	bucketAuthOpts := vod.UpdateBucketAuthorizedReq{
		Bucket:    d.Get("input_bucket").(string),
		Operation: "1",
	}

	bucketAuthReq := vod.UpdateBucketAuthorizedRequest{
		Body: &bucketAuthOpts,
	}
	_, err := client.UpdateBucketAuthorized(&bucketAuthReq)
	if err != nil {
		return "", fmt.Errorf("error authorizing the OBS bucket to VOD: %s", err)
	}

	var videoType vod.PublishAssetFromObsReqVideoType
	if err = videoType.UnmarshalJSON([]byte(d.Get("media_type").(string))); err != nil {
		return "", fmt.Errorf("error parsing the argument media_type: %s", err)
	}

	createOpts := vod.PublishAssetFromObsReq{
		VideoType:         videoType,
		Title:             d.Get("name").(string),
		Description:       utils.StringIgnoreEmpty(d.Get("description").(string)),
		CategoryId:        utils.Int32IgnoreEmpty(int32(d.Get("category_id").(int))),
		Tags:              utils.StringIgnoreEmpty(d.Get("labels").(string)),
		TemplateGroupName: utils.StringIgnoreEmpty(d.Get("template_group_name").(string)),
		WorkflowName:      utils.StringIgnoreEmpty(d.Get("workflow_name").(string)),
		Review:            buildReviewOpts(d.Get("review_template_id").(string)),
		Thumbnail:         buildThumbnailOpts(d.Get("thumbnail").([]interface{})),
		Input: &vod.FileAddr{
			Bucket:   d.Get("input_bucket").(string),
			Object:   d.Get("input_path").(string),
			Location: region,
		},
		OutputBucket: utils.StringIgnoreEmpty(d.Get("output_bucket").(string)),
		OutputPath:   utils.StringIgnoreEmpty(d.Get("output_path").(string)),
		StorageMode:  utils.Int32IgnoreEmpty(int32(d.Get("storage_mode").(int))),
	}

	if d.Get("publish").(bool) {
		createOpts.AutoPublish = utils.Int32(int32(1))
	} else {
		createOpts.AutoPublish = utils.Int32(int32(0))
	}
	if d.Get("auto_encrypt").(bool) {
		createOpts.AutoEncrypt = utils.Int32(int32(1))
	}
	if d.Get("auto_preload").(bool) {
		createOpts.AutoPreheat = utils.Int32(int32(1))
	}

	createReq := vod.PublishAssetFromObsRequest{
		Body: &createOpts,
	}

	resp, err := client.PublishAssetFromObs(&createReq)
	if err != nil {
		return "", fmt.Errorf("error creating VOD media asset: %s", err)
	}

	if resp.AssetId == nil {
		return "", fmt.Errorf("unable to find the asset after uploading")
	}

	return *resp.AssetId, nil
}

func resourceMediaAssetCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.HcVodV1Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating VOD client: %s", err)
	}

	var AssetID string
	if _, ok := d.GetOk("url"); ok {
		AssetID, err = createMediaAssetByUrl(client, d)
		if err != nil {
			diag.FromErr(err)
		}
	} else {
		AssetID, err = createMediaAssetFromObs(client, d, config.GetRegion(d))
		if err != nil {
			diag.FromErr(err)
		}
	}

	d.SetId(AssetID)
	return resourceMediaAssetRead(ctx, d, meta)
}

func resourceMediaAssetRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.HcVodV1Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating VOD client: %s", err)
	}

	resp, err := client.ShowAssetDetail(&vod.ShowAssetDetailRequest{AssetId: d.Id()})
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving VOD media asset")
	}

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("name", resp.BaseInfo.Title),
		d.Set("media_name", resp.BaseInfo.VideoName),
		d.Set("description", resp.BaseInfo.Description),
		d.Set("category_id", resp.BaseInfo.CategoryId),
		d.Set("category_name", resp.BaseInfo.CategoryName),
		d.Set("media_type", resp.BaseInfo.VideoType),
		d.Set("labels", resp.BaseInfo.Tags),
		d.Set("media_url", resp.BaseInfo.VideoUrl),
	)

	if sourcePath := resp.BaseInfo.SourcePath; sourcePath != nil {
		mErr = multierror.Append(mErr,
			d.Set("input_bucket", sourcePath.Bucket),
			d.Set("input_path", sourcePath.Object),
		)
	}

	if outputPath := resp.BaseInfo.OutputPath; outputPath != nil {
		mErr = multierror.Append(mErr,
			d.Set("output_bucket", outputPath.Bucket),
			d.Set("output_path", outputPath.Object),
		)
	}

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting VOD media asset fields: %s", err)
	}

	return nil
}

func resourceMediaAssetUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.HcVodV1Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating VOD client: %s", err)
	}

	if d.HasChanges("name", "description", "category_id", "labels") {
		updateOpts := vod.UpdateAssetMetaReq{
			AssetId:     d.Id(),
			Title:       utils.String(d.Get("name").(string)),
			Description: utils.String(d.Get("description").(string)),
			CategoryId:  utils.Int32(int32(d.Get("category_id").(int))),
			Tags:        utils.String(d.Get("labels").(string)),
		}

		updateReq := vod.UpdateAssetMetaRequest{
			Body: &updateOpts,
		}

		_, err = client.UpdateAssetMeta(&updateReq)
		if err != nil {
			return diag.Errorf("error updating VOD media asset: %s", err)
		}
	}

	if d.HasChange("publish") {
		if d.Get("publish").(bool) {
			_, err = client.PublishAssets(&vod.PublishAssetsRequest{Body: &vod.PublishAssetReq{AssetId: []string{d.Id()}}})
			if err != nil {
				return diag.Errorf("error publishing VOD media asset: %s", err)
			}
		} else {
			_, err = client.UnpublishAssets(&vod.UnpublishAssetsRequest{Body: &vod.PublishAssetReq{AssetId: []string{d.Id()}}})
			if err != nil {
				return diag.Errorf("error unpublishing VOD media asset: %s", err)
			}
		}
	}

	return resourceMediaAssetRead(ctx, d, meta)
}

func resourceMediaAssetDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.HcVodV1Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating VOD client: %s", err)
	}

	_, err = client.DeleteAssets(&vod.DeleteAssetsRequest{AssetId: []string{d.Id()}})
	if err != nil {
		return diag.Errorf("error deleting VOD media asset: %s", err)
	}

	return nil
}
