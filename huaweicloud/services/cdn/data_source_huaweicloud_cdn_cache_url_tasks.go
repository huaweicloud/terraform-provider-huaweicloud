// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product CDN
// ---------------------------------------------------------------

package cdn

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CDN GET /v1.0/cdn/contentgateway/url-tasks
func DataSourceCacheUrlTasks() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceCacheUrlTasksRead,
		Schema: map[string]*schema.Schema{
			"start_time": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `Specifies the start timestamp, in milliseconds.`,
			},
			"end_time": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `Specifies the end timestamp, in milliseconds.`,
			},
			"url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the refresh or preheat URL.`,
			},
			"task_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the task type.`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the URL status.`,
			},
			"file_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the file type.`,
			},
			"tasks": {
				Type:        schema.TypeList,
				Elem:        cacheUrlTasksSchema(),
				Computed:    true,
				Description: `The list of URL task information.`,
			},
		},
	}
}

func cacheUrlTasksSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The URL ID.`,
			},
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The URL.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The URL status.`,
			},
			"task_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The task type.`,
			},
			"mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The directory refresh mode.`,
			},
			"task_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The task ID.`,
			},
			"modify_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The modification time.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time.`,
			},
			"file_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The file type.`,
			},
		},
	}
	return &sc
}

func buildCacheUrlTasksQueryParams(d *schema.ResourceData) string {
	queryParams := "?limit=100"
	if v, ok := d.GetOk("start_time"); ok {
		queryParams = fmt.Sprintf("%s&start_time=%v", queryParams, v)
	}
	if v, ok := d.GetOk("end_time"); ok {
		queryParams = fmt.Sprintf("%s&end_time=%v", queryParams, v)
	}
	if v, ok := d.GetOk("url"); ok {
		queryParams = fmt.Sprintf("%s&url=%v", queryParams, v)
	}
	if v, ok := d.GetOk("task_type"); ok {
		queryParams = fmt.Sprintf("%s&task_type=%v", queryParams, v)
	}
	if v, ok := d.GetOk("status"); ok {
		queryParams = fmt.Sprintf("%s&status=%v", queryParams, v)
	}
	if v, ok := d.GetOk("file_type"); ok {
		queryParams = fmt.Sprintf("%s&file_type=%v", queryParams, v)
	}

	return queryParams
}

func resourceCacheUrlTasksRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "cdn"
		httpUrl = "v1.0/cdn/contentgateway/url-tasks"
		offset  = 0
		result  = make([]interface{}, 0)
		mErr    *multierror.Error
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CDN client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath += buildCacheUrlTasksQueryParams(d)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", getPath, offset)
		getResp, err := client.Request("GET", currentPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving CDN cache url tasks: %s", err)
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}

		tasksResp := utils.PathSearch("result", getRespBody, make([]interface{}, 0)).([]interface{})
		if len(tasksResp) == 0 {
			break
		}
		result = append(result, tasksResp...)
		offset += len(tasksResp)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr = multierror.Append(
		mErr,
		d.Set("tasks", flattenCacheUrlTasks(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenCacheUrlTasks(tasksResp []interface{}) []interface{} {
	rst := make([]interface{}, 0, len(tasksResp))
	for _, v := range tasksResp {
		modifyTime := utils.PathSearch("modify_time", v, float64(0)).(float64)
		createdAt := utils.PathSearch("create_time", v, float64(0)).(float64)
		rst = append(rst, map[string]interface{}{
			"id":          utils.PathSearch("id", v, nil),
			"url":         utils.PathSearch("url", v, nil),
			"status":      utils.PathSearch("status", v, nil),
			"task_type":   utils.PathSearch("type", v, nil),
			"mode":        utils.PathSearch("mode", v, nil),
			"task_id":     utils.PathSearch("task_id", v, nil),
			"modify_time": utils.FormatTimeStampRFC3339(int64(modifyTime)/1000, false),
			"created_at":  utils.FormatTimeStampRFC3339(int64(createdAt)/1000, false),
			"file_type":   utils.PathSearch("file_type", v, nil),
		})
	}

	return rst
}
