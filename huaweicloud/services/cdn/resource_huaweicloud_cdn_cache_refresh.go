package cdn

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

var cacheRefreshNonUpdatableParams = []string{"urls", "enterprise_project_id", "type", "mode", "zh_url_encode"}

// @API CDN POST /v1.0/cdn/content/refresh-tasks
// @API CDN GET /v1.0/cdn/historytasks/{history_tasks_id}/detail
func ResourceCacheRefresh() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCacheRefreshCreate,
		ReadContext:   resourceCacheRefreshRead,
		UpdateContext: resourceCacheRefreshUpdate,
		DeleteContext: resourceCacheRefreshDelete,

		CustomizeDiff: config.FlexibleForceNew(cacheRefreshNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			// Required parameters
			"urls": {
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The URLs that need to be refreshed.`,
			},

			// Optional parameters
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The enterprise project ID to which the resource belongs.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The refresh type.`,
			},
			"mode": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The directory refresh mode.`,
			},
			"zh_url_encode": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether to encode Chinese characters in URLs before cache refresh.`,
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

func buildCacheQueryParams(d *schema.ResourceData) string {
	queryParams := ""

	if v, ok := d.GetOk("enterprise_project_id"); ok {
		queryParams = fmt.Sprintf("%s?enterprise_project_id=%v", queryParams, v)
	}

	return queryParams
}

func buildCacheRefreshBodyParams(d *schema.ResourceData) interface{} {
	refreshTaskMap := map[string]interface{}{
		"urls":          utils.ExpandToStringList(d.Get("urls").(*schema.Set).List()),
		"type":          utils.ValueIgnoreEmpty(d.Get("type")),
		"mode":          utils.ValueIgnoreEmpty(d.Get("mode")),
		"zh_url_encode": utils.ValueIgnoreEmpty(d.Get("zh_url_encode")),
	}

	return map[string]interface{}{
		"refresh_task": utils.RemoveNil(refreshTaskMap),
	}
}

func resourceCacheRefreshCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		httpUrl = "v1.0/cdn/content/refresh-tasks"
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
		JSONBody:         buildCacheRefreshBodyParams(d),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating cache refresh: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	refreshTaskId := utils.PathSearch("refresh_task", createRespBody, "").(string)
	if refreshTaskId == "" {
		return diag.Errorf("error creating cache refresh: ID is not found in API response")
	}
	d.SetId(refreshTaskId)

	if err = waitingForCacheCreateCompleted(ctx, client, refreshTaskId, d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.Errorf("error waiting for cache refresh (%s) creation to completed: %s", refreshTaskId, err)
	}
	return resourceCacheRefreshRead(ctx, d, meta)
}

func waitingForCacheCreateCompleted(ctx context.Context, client *golangsdk.ServiceClient, id string,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			getRespBody, err := GetCacheDetailById(client, id)
			if err != nil {
				return nil, "ERROR", err
			}

			processing := utils.PathSearch("processing", getRespBody, float64(0)).(float64)
			if processing == 0 {
				return getRespBody, "COMPLETED", nil
			}

			return getRespBody, "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        5 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)

	return err
}

func GetCacheDetailById(client *golangsdk.ServiceClient, id string) (interface{}, error) {
	getPath := client.Endpoint + "v1.0/cdn/historytasks/{history_tasks_id}/detail"
	getPath = strings.ReplaceAll(getPath, "{history_tasks_id}", id)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving cache detail: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	// When the resource does not exist, the query API still returns a `200` status code,
	// and the error message is as follows:
	// {
	//   "error": {
	//     "error_code": "CDN.0108",
	//     "error_msg": "The URL domain name is not the acceleration domain name of the current tenant."
	//   }
	// }
	// Return a `404` status code for handling this scenario.
	errorCode := utils.PathSearch("error.error_code", getRespBody, "").(string)
	if errorCode == "CDN.0108" {
		return nil, golangsdk.ErrDefault404{}
	}

	return getRespBody, nil
}

func flattenCacheUrls(urlsResp []interface{}) []string {
	result := make([]string, 0, len(urlsResp))
	for _, v := range urlsResp {
		url := utils.PathSearch("url", v, "").(string)
		if url != "" {
			result = append(result, url)
		}
	}

	return result
}

func resourceCacheRefreshRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("cdn", "")
	if err != nil {
		return diag.Errorf("error creating CDN client: %s", err)
	}

	getRespBody, err := GetCacheDetailById(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "CDN cache refresh")
	}

	mErr := multierror.Append(
		d.Set("urls", flattenCacheUrls(utils.PathSearch("urls", getRespBody, make([]interface{}, 0)).([]interface{}))),
		d.Set("type", utils.PathSearch("file_type", getRespBody, nil)),
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

func resourceCacheRefreshUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceCacheRefreshDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}
