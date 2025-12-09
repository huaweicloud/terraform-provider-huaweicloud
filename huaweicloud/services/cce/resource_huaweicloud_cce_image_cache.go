package cce

import (
	"context"
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

// @API CCE POST /v5/imagecaches
// @API CCE GET /v5/imagecaches/{image_cache_id}
// @API CCE DELETE /v5/imagecaches/{image_cache_id}

var imageCacheNonUpdatableParams = []string{
	"name",
	"images",
	"building_config", "building_config.*.cluster", "building_config.*.image_pull_secrets",
	"image_cache_size",
	"retention_days",
}

func ResourceImageCache() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceImageCacheCreate,
		ReadContext:   resourceImageCacheRead,
		UpdateContext: resourceImageCacheUpdate,
		DeleteContext: resourceImageCacheDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(imageCacheNonUpdatableParams),

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
			"images": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"building_config": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster": {
							Type:     schema.TypeString,
							Required: true,
						},
						"image_pull_secrets": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"image_cache_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"retention_days": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildImageCacheBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":             d.Get("name"),
		"images":           d.Get("images"),
		"image_cache_size": d.Get("image_cache_size"),
		"retention_days":   d.Get("retention_days"),
		"building_config":  buildImageCacheBuildingConfigParams(d),
	}

	return bodyParams
}

func buildImageCacheBuildingConfigParams(d *schema.ResourceData) map[string]interface{} {
	buildingConfig := d.Get("building_config")
	bodyParams := map[string]interface{}{
		"cluster":            utils.PathSearch("[0].cluster", buildingConfig, nil),
		"image_pull_secrets": utils.PathSearch("[0].image_pull_secrets", buildingConfig, nil),
	}

	return bodyParams
}

func resourceImageCacheCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		createImageCacheHttpUrl = "v5/imagecaches"
		createImageCacheProduct = "cce"
	)
	createImageCacheClient, err := cfg.NewServiceClient(createImageCacheProduct, region)
	if err != nil {
		return diag.Errorf("error creating CCE Client: %s", err)
	}

	createImageCachePath := createImageCacheClient.Endpoint + createImageCacheHttpUrl

	createImageCacheOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	createImageCacheOpt.JSONBody = utils.RemoveNil(buildImageCacheBodyParams(d))
	createImageCacheResp, err := createImageCacheClient.Request("POST", createImageCachePath, &createImageCacheOpt)
	if err != nil {
		return diag.Errorf("error creating CCE image cache: %s", err)
	}

	createImageCacheRespBody, err := utils.FlattenResponse(createImageCacheResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("image_cache.id", createImageCacheRespBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating CCE image cache: ID is not found in API response")
	}
	d.SetId(id)

	return resourceImageCacheRead(ctx, d, meta)
}

func resourceImageCacheRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		getImageCacheHttpUrl = "v5/imagecaches/{image_cache_id}"
		getImageCacheProduct = "cce"
	)
	getImageCacheClient, err := cfg.NewServiceClient(getImageCacheProduct, region)
	if err != nil {
		return diag.Errorf("error creating CCE Client: %s", err)
	}

	getImageCachePath := getImageCacheClient.Endpoint + getImageCacheHttpUrl
	getImageCachePath = strings.ReplaceAll(getImageCachePath, "{image_cache_id}", d.Id())

	getImageCacheOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getImageCacheResp, err := getImageCacheClient.Request("GET", getImageCachePath, &getImageCacheOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CCE image cache")
	}

	getImageCacheRespBody, err := utils.FlattenResponse(getImageCacheResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(nil,
		d.Set("region", cfg.GetRegion(d)),
		d.Set("name", utils.PathSearch("image_cache.name", getImageCacheRespBody, nil)),
		d.Set("images", utils.PathSearch("image_cache.images", getImageCacheRespBody, nil)),
		d.Set("image_cache_size", utils.PathSearch("image_cache.image_cache_size", getImageCacheRespBody, nil)),
		d.Set("retention_days", utils.PathSearch("image_cache.retention_days", getImageCacheRespBody, nil)),
		d.Set("building_config", flattenImageCacheBuildingConfig(getImageCacheRespBody)),
		d.Set("created_at", utils.PathSearch("image_cache.created_at", getImageCacheRespBody, nil)),
		d.Set("status", utils.PathSearch("image_cache.status", getImageCacheRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenImageCacheBuildingConfig(respBody interface{}) []map[string]interface{} {
	cluster := utils.PathSearch("image_cache.building_config.cluster", respBody, "").(string)
	imagePullSecrets := utils.PathSearch("image_cache.building_config.image_pull_secrets", respBody, make([]interface{}, 0))
	return []map[string]interface{}{
		{
			"cluster":            cluster,
			"image_pull_secrets": imagePullSecrets,
		},
	}
}

func resourceImageCacheUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceImageCacheDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		deleteImageCacheHttpUrl = "v5/imagecaches/{image_cache_id}"
		deleteImageCacheProduct = "cce"
	)
	deleteImageCacheClient, err := cfg.NewServiceClient(deleteImageCacheProduct, region)
	if err != nil {
		return diag.Errorf("error creating CCE Client: %s", err)
	}

	deleteImageCachePath := deleteImageCacheClient.Endpoint + deleteImageCacheHttpUrl
	deleteImageCachePath = strings.ReplaceAll(deleteImageCachePath, "{image_cache_id}", d.Id())

	deleteImageCacheOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = deleteImageCacheClient.Request("DELETE", deleteImageCachePath, &deleteImageCacheOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting CCE image cache")
	}

	err = imageCacheWaitingForDeleted(ctx, d, meta, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.Errorf("error waiting for deleting CCE image cache (%s) to complete: %s", d.Id(), err)
	}

	return nil
}

func imageCacheWaitingForDeleted(ctx context.Context, d *schema.ResourceData, meta interface{}, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			cfg := meta.(*config.Config)
			region := cfg.GetRegion(d)
			var (
				imageCacheWaitingHttpUrl = "v5/imagecaches/{image_cache_id}"
				imageCacheWaitingProduct = "cce"
			)
			imageCacheWaitingClient, err := cfg.NewServiceClient(imageCacheWaitingProduct, region)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error creating CCE client: %s", err)
			}

			imageCacheWaitingPath := imageCacheWaitingClient.Endpoint + imageCacheWaitingHttpUrl
			imageCacheWaitingPath = strings.ReplaceAll(imageCacheWaitingPath, "{image_cache_id}", d.Id())

			imageCacheWaitingOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
			}
			imageCacheWaitingResp, err := imageCacheWaitingClient.Request("GET", imageCacheWaitingPath, &imageCacheWaitingOpt)
			if err != nil {
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					return imageCacheWaitingResp, "COMPLETED", nil
				}
				return nil, "ERROR", err
			}

			imageCacheWaitingRespBody, err := utils.FlattenResponse(imageCacheWaitingResp)
			if err != nil {
				return nil, "ERROR", err
			}

			return imageCacheWaitingRespBody, "PENDING", nil
		},
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}
