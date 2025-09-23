// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product CDN
// ---------------------------------------------------------------

package cdn

import (
	"context"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CDN POST /v1.0/cdn/content/preheating-tasks
// @API CDN GET /v1.0/cdn/historytasks/{history_tasks_id}/detail
func ResourceCachePreheat() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCachePreheatCreate,
		ReadContext:   resourceCachePreheatRead,
		DeleteContext: resourceCachePreheatDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"urls": {
				Type:        schema.TypeSet,
				Required:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the URLs that need to be preheated.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the enterprise project ID.`,
			},
			"zh_url_encode": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies whether to encode Chinese characters in URLs before cache preheat.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The task execution result.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time.`,
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
		},
	}
}

func buildCachePreheatBodyParams(d *schema.ResourceData) interface{} {
	refreshTaskMap := map[string]interface{}{
		"urls": utils.ExpandToStringList(d.Get("urls").(*schema.Set).List()),
	}
	if v, ok := d.GetOk("zh_url_encode"); ok {
		refreshTaskMap["zh_url_encode"] = v
	}
	bodyParams := map[string]interface{}{
		"preheating_task": refreshTaskMap,
	}

	return bodyParams
}

func resourceCachePreheatCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		product       = "cdn"
		createHttpUrl = "v1.0/cdn/content/preheating-tasks"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CDN client: %s", err)
	}

	createPath := client.Endpoint + createHttpUrl
	createPath += buildCacheQueryParams(d)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         buildCachePreheatBodyParams(d),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating CDN cache preheat: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	preheatingTask := utils.PathSearch("preheating_task", createRespBody, "").(string)
	if preheatingTask == "" {
		return diag.Errorf("error creating CDN cache preheat: ID is not found in API response")
	}

	d.SetId(preheatingTask)

	if err = waitingForCacheCreateCompleted(ctx, client, preheatingTask, d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.Errorf("error waiting for CDN cache preheat (%s) creation to completed: %s", preheatingTask, err)
	}

	return resourceCachePreheatRead(ctx, d, meta)
}

func resourceCachePreheatRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "cdn"
		mErr    *multierror.Error
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CDN client: %s", err)
	}

	getRespBody, err := GetCacheDetailById(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "CDN cache preheat")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("urls", flattenCacheUrls(utils.PathSearch("urls", getRespBody, make([]interface{}, 0)).([]interface{}))),
		d.Set("status", utils.PathSearch("status", getRespBody, nil)),
		d.Set("created_at", flattenCreatedAt(getRespBody)),
		d.Set("processing", utils.PathSearch("processing", getRespBody, float64(0)).(float64)),
		d.Set("succeed", utils.PathSearch("succeed", getRespBody, float64(0)).(float64)),
		d.Set("failed", utils.PathSearch("failed", getRespBody, float64(0)).(float64)),
		d.Set("total", utils.PathSearch("total", getRespBody, float64(0)).(float64)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceCachePreheatDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}
