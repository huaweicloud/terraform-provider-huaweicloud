package cdn

import (
	"context"
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

var cachePreheatNonUpdatableParams = []string{"urls", "enterprise_project_id", "zh_url_encode"}

// @API CDN POST /v1.0/cdn/content/preheating-tasks
// @API CDN GET /v1.0/cdn/historytasks/{history_tasks_id}/detail
func ResourceCachePreheat() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCachePreheatCreate,
		ReadContext:   resourceCachePreheatRead,
		UpdateContext: resourceCachePreheatUpdate,
		DeleteContext: resourceCachePreheatDelete,

		CustomizeDiff: config.FlexibleForceNew(cachePreheatNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			// Required parameters
			"urls": {
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The URLs that need to be preheated.`,
			},

			// Optional parameters
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The enterprise project ID to which the resource belongs.`,
			},
			"zh_url_encode": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether to encode Chinese characters in URLs before cache preheat.`,
			},

			// Attributes
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The task execution result.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time, in RFC3339 format.`,
			},
			"processing": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of URLs that are being processed.`,
			},
			"succeed": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of URLs processed.`,
			},
			"failed": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of URLs that failed to be processed.`,
			},
			"total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The total number of URLs in historical tasks.`,
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func buildCachePreheatBodyParams(d *schema.ResourceData) interface{} {
	refreshTaskMap := map[string]interface{}{
		"urls":          utils.ExpandToStringList(d.Get("urls").(*schema.Set).List()),
		"zh_url_encode": utils.ValueIgnoreEmpty(d.Get("zh_url_encode")),
	}

	return map[string]interface{}{
		"preheating_task": utils.RemoveNil(refreshTaskMap),
	}
}

func resourceCachePreheatCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		httpUrl = "v1.0/cdn/content/preheating-tasks"
	)

	client, err := cfg.NewServiceClient("cdn", "")
	if err != nil {
		return diag.Errorf("error creating CDN client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath += buildCacheQueryParams(d)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         buildCachePreheatBodyParams(d),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating cache preheat: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	preheatingTask := utils.PathSearch("preheating_task", createRespBody, "").(string)
	if preheatingTask == "" {
		return diag.Errorf("error creating cache preheat: ID is not found in API response")
	}
	d.SetId(preheatingTask)

	if err = waitingForCacheCreateCompleted(ctx, client, preheatingTask, d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.Errorf("error waiting for cache preheat (%s) creation to completed: %s", preheatingTask, err)
	}
	return resourceCachePreheatRead(ctx, d, meta)
}

func resourceCachePreheatRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("cdn", "")
	if err != nil {
		return diag.Errorf("error creating CDN client: %s", err)
	}

	getRespBody, err := GetCacheDetailById(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "CDN cache preheat")
	}

	mErr := multierror.Append(
		d.Set("urls", flattenCacheUrls(utils.PathSearch("urls", getRespBody, make([]interface{}, 0)).([]interface{}))),
		d.Set("status", utils.PathSearch("status", getRespBody, nil)),
		d.Set("created_at", utils.FormatTimeStampRFC3339(
			int64(utils.PathSearch("create_time", getRespBody, float64(0)).(float64))/1000, false)),
		d.Set("processing", utils.PathSearch("processing", getRespBody, float64(0)).(float64)),
		d.Set("succeed", utils.PathSearch("succeed", getRespBody, float64(0)).(float64)),
		d.Set("failed", utils.PathSearch("failed", getRespBody, float64(0)).(float64)),
		d.Set("total", utils.PathSearch("total", getRespBody, float64(0)).(float64)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceCachePreheatUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceCachePreheatDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}
