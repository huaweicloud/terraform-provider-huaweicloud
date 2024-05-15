// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product CDN
// ---------------------------------------------------------------

package cdn

import (
	"context"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cdn/v2/model"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CDN POST /v1.0/cdn/content/refresh-tasks
// @API CDN GET /v1.0/cdn/historytasks/{history_tasks_id}/detail
func ResourceCacheRefresh() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCacheRefreshCreate,
		ReadContext:   resourceCacheRefreshRead,
		DeleteContext: resourceCacheRefreshDelete,

		Schema: map[string]*schema.Schema{
			"urls": {
				Type:        schema.TypeSet,
				Required:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the URLs that need to be refreshed.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the enterprise project ID.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the refresh type.`,
			},
			"mode": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies the directory refresh mode.`,
			},
			"zh_url_encode": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: `Specifies whether to encode Chinese characters in URLs before cache refresh.`,
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

func resourceCacheRefreshCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	hcCdnClient, err := cfg.HcCdnV2Client(region)
	if err != nil {
		return diag.Errorf("error creating CDN v2 client: %s", err)
	}

	request := &model.CreateRefreshTasksRequest{
		EnterpriseProjectId: utils.StringIgnoreEmpty(cfg.GetEnterpriseProjectID(d)),
		Body: &model.RefreshTaskRequest{
			RefreshTask: &model.RefreshTaskRequestBody{
				Type:        buildRefreshTaskRequestBodyTypeOpts(d.Get("type").(string)),
				Mode:        buildRefreshTaskRequestBodyModeOpts(d.Get("mode").(string)),
				ZhUrlEncode: utils.Bool(d.Get("zh_url_encode").(bool)),
				Urls:        utils.ExpandToStringList(d.Get("urls").(*schema.Set).List()),
			},
		},
	}

	resp, err := hcCdnClient.CreateRefreshTasks(request)
	if err != nil {
		return diag.Errorf("error creating CDN cache refresh: %s", err)
	}

	if resp == nil || resp.RefreshTask == nil || len(*resp.RefreshTask) == 0 {
		return diag.Errorf("error creating CDN cache refresh: ID is not found in API response")
	}
	d.SetId(*resp.RefreshTask)
	return resourceCacheRefreshRead(ctx, d, meta)
}

func buildRefreshTaskRequestBodyTypeOpts(refreshTask string) *model.RefreshTaskRequestBodyType {
	if refreshTask == "" {
		return nil
	}

	refreshTaskToReq := new(model.RefreshTaskRequestBodyType)
	if err := refreshTaskToReq.UnmarshalJSON([]byte(refreshTask)); err != nil {
		log.Printf("[WARN] failed to parse task %s: %s", refreshTask, err)
		return nil
	}
	return refreshTaskToReq
}

func buildRefreshTaskRequestBodyModeOpts(mode string) *model.RefreshTaskRequestBodyMode {
	if mode == "" {
		return nil
	}

	modeToReq := new(model.RefreshTaskRequestBodyMode)
	if err := modeToReq.UnmarshalJSON([]byte(mode)); err != nil {
		log.Printf("[WARN] failed to parse mode %s: %s", mode, err)
		return nil
	}
	return modeToReq
}

func resourceCacheRefreshRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	hcCdnClient, err := cfg.HcCdnV2Client(region)
	if err != nil {
		return diag.Errorf("error creating CDN v2 client: %s", err)
	}

	request := &model.ShowHistoryTaskDetailsRequest{
		EnterpriseProjectId: utils.StringIgnoreEmpty(cfg.GetEnterpriseProjectID(d)),
		HistoryTasksId:      d.Id(),
	}

	resp, err := hcCdnClient.ShowHistoryTaskDetails(request)
	if err != nil {
		return diag.Errorf("error retrieving CDN cache refresh: %s", err)
	}

	if resp == nil {
		return diag.Errorf("error retrieving CDN cache refresh: Task is not found in API response")
	}

	var mErr *multierror.Error
	mErr = multierror.Append(
		mErr,
		d.Set("urls", flattenUrls(resp.Urls)),
		d.Set("type", resp.FileType),
		d.Set("status", resp.Status),
		d.Set("created_at", flattenCreatedAt(resp.CreateTime)),
		d.Set("processing", resp.Processing),
		d.Set("succeed", resp.Succeed),
		d.Set("failed", resp.Failed),
		d.Set("total", resp.Total),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceCacheRefreshDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}
