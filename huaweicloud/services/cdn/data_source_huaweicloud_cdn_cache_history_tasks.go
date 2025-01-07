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

// @API CDN GET /v1.0/cdn/historytasks
func DataSourceCacheHistoryTasks() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceCacheHistoryTasksRead,
		Schema: map[string]*schema.Schema{
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the ID of the enterprise project to which the resource belongs.`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the task status.`,
			},
			"start_date": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `Specifies the query start time.`,
			},
			"end_date": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `Specifies the query end time.`,
			},
			"order_field": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the field used for sorting.`,
			},
			"order_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the sorting type.`,
			},
			"file_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the file type.`,
			},
			"task_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the task type.`,
			},
			"tasks": {
				Type:        schema.TypeList,
				Elem:        cacheHistoryTasksSchema(),
				Computed:    true,
				Description: `The history task list.`,
			},
		},
	}
}

func cacheHistoryTasksSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the task ID.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the task result.`,
			},
			"processing": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the number of URLs that are being processed.`,
			},
			"succeed": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the number of URLs processed.`,
			},
			"failed": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the number of URLs that failed to be processed.`,
			},
			"total": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Indicates the total number of URLs in the task.`,
			},
			"task_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the task type.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the time when the task was created.`,
			},
			"file_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the file type.`,
			},
		},
	}
	return &sc
}

func buildHistoryTaskQueryParams(d *schema.ResourceData, cfg *config.Config) string {
	queryParams := "?page_size=10000"
	if epsId := cfg.GetEnterpriseProjectID(d); epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%s", queryParams, epsId)
	}
	if v, ok := d.GetOk("status"); ok {
		queryParams = fmt.Sprintf("%s&status=%v", queryParams, v)
	}
	if v, ok := d.GetOk("start_date"); ok {
		queryParams = fmt.Sprintf("%s&start_date=%v", queryParams, v)
	}
	if v, ok := d.GetOk("end_date"); ok {
		queryParams = fmt.Sprintf("%s&end_date=%v", queryParams, v)
	}
	if v, ok := d.GetOk("order_field"); ok {
		queryParams = fmt.Sprintf("%s&order_field=%v", queryParams, v)
	}
	if v, ok := d.GetOk("order_type"); ok {
		queryParams = fmt.Sprintf("%s&order_type=%v", queryParams, v)
	}
	if v, ok := d.GetOk("file_type"); ok {
		queryParams = fmt.Sprintf("%s&file_type=%v", queryParams, v)
	}
	if v, ok := d.GetOk("task_type"); ok {
		queryParams = fmt.Sprintf("%s&task_type=%v", queryParams, v)
	}

	return queryParams
}

func resourceCacheHistoryTasksRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		product    = "cdn"
		httpUrl    = "v1.0/cdn/historytasks"
		pageNumber = 1
		result     = make([]interface{}, 0)
		mErr       *multierror.Error
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CDN client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath += buildHistoryTaskQueryParams(d, cfg)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	for {
		currentPath := fmt.Sprintf("%s&page_number=%v", getPath, pageNumber)
		resp, err := client.Request("GET", currentPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving CDN cache history tasks: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		tasksResp := utils.PathSearch("tasks", respBody, make([]interface{}, 0)).([]interface{})
		if len(tasksResp) == 0 {
			break
		}
		result = append(result, tasksResp...)
		pageNumber++
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr = multierror.Append(
		mErr,
		d.Set("tasks", flattenCacheHistoryTasks(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenCacheHistoryTasks(tasks []interface{}) []interface{} {
	rst := make([]interface{}, 0, len(tasks))
	for _, v := range tasks {
		createdAt := utils.PathSearch("create_time", v, float64(0)).(float64)
		rst = append(rst, map[string]interface{}{
			"id":         utils.PathSearch("id", v, nil),
			"status":     utils.PathSearch("status", v, nil),
			"processing": utils.PathSearch("processing", v, nil),
			"succeed":    utils.PathSearch("succeed", v, nil),
			"failed":     utils.PathSearch("failed", v, nil),
			"total":      utils.PathSearch("total", v, nil),
			"task_type":  utils.PathSearch("task_type", v, nil),
			"created_at": utils.FormatTimeStampRFC3339(int64(createdAt)/1000, false),
			"file_type":  utils.PathSearch("file_type", v, nil),
		})
	}

	return rst
}
