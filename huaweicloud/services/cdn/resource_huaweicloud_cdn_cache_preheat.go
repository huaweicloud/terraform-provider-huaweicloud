// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product CDN
// ---------------------------------------------------------------

package cdn

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cdn/v2/model"

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

func resourceCachePreheatCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	hcCdnClient, err := cfg.HcCdnV2Client(region)
	if err != nil {
		return diag.Errorf("error creating CDN v2 client: %s", err)
	}

	request := &model.CreatePreheatingTasksRequest{
		EnterpriseProjectId: utils.StringIgnoreEmpty(cfg.GetEnterpriseProjectID(d)),
		Body: &model.PreheatingTaskRequest{
			PreheatingTask: &model.PreheatingTaskRequestBody{
				ZhUrlEncode: utils.Bool(d.Get("zh_url_encode").(bool)),
				Urls:        utils.ExpandToStringList(d.Get("urls").(*schema.Set).List()),
			},
		},
	}

	resp, err := hcCdnClient.CreatePreheatingTasks(request)
	if err != nil {
		return diag.Errorf("error creating CDN cache preheat: %s", err)
	}

	if resp == nil || resp.PreheatingTask == nil || len(*resp.PreheatingTask) == 0 {
		return diag.Errorf("error creating CDN cache preheat: ID is not found in API response")
	}
	d.SetId(*resp.PreheatingTask)
	return resourceCachePreheatRead(ctx, d, meta)
}

func resourceCachePreheatRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		return diag.Errorf("error retrieving CDN cache preheat: %s", err)
	}

	if resp == nil {
		return diag.Errorf("error retrieving CDN cache preheat: Task is not found in API response")
	}

	var mErr *multierror.Error
	mErr = multierror.Append(
		mErr,
		d.Set("urls", flattenUrls(resp.Urls)),
		d.Set("status", resp.Status),
		d.Set("created_at", flattenCreatedAt(resp.CreateTime)),
		d.Set("processing", resp.Processing),
		d.Set("succeed", resp.Succeed),
		d.Set("failed", resp.Failed),
		d.Set("total", resp.Total),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenUrls(urlObjects *[]model.UrlObject) []string {
	if urlObjects == nil || len(*urlObjects) == 0 {
		return nil
	}
	urls := make([]string, 0, len(*urlObjects))
	for _, obj := range *urlObjects {
		if obj.Url != nil {
			urls = append(urls, *obj.Url)
		}
	}
	return urls
}

func flattenCreatedAt(createTime *int64) string {
	if createTime == nil {
		return ""
	}
	return utils.FormatTimeStampRFC3339(*createTime/1000, false)
}

func resourceCachePreheatDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}
