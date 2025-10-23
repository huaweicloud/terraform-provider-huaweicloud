package cdn

import (
	"context"
	"fmt"
	"strconv"
	"strings"

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
		ReadContext: dataSourceCacheHistoryTasksRead,

		Schema: map[string]*schema.Schema{
			// Optional parameters
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the enterprise project to which the resource belongs.`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The task status.`,
			},
			"start_date": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The query start time.`,
			},
			"end_date": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The query end time.`,
			},
			"order_field": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The field used for sorting.`,
			},
			"order_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The sorting type.`,
			},
			"file_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The file type.`,
			},
			"task_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The task type.`,
			},

			// Attributes
			"tasks": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        cacheHistoryTasksSchema(),
				Description: `The list of history tasks that matched filter parameters.`,
			},
		},
	}
}

func cacheHistoryTasksSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The task ID.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The task result.`,
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
				Description: `The total number of URLs in the task.`,
			},
			"task_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The task type.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time, in RFC3339 format.`,
			},
			"file_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The file type.`,
			},
		},
	}
}

func buildHistoryTaskQueryParams(d *schema.ResourceData, cfg *config.Config) string {
	res := ""

	if epsId := cfg.GetEnterpriseProjectID(d); epsId != "" {
		res = fmt.Sprintf("%s&enterprise_project_id=%s", res, epsId)
	}
	if v, ok := d.GetOk("status"); ok {
		res = fmt.Sprintf("%s&status=%v", res, v)
	}
	if v, ok := d.GetOk("start_date"); ok {
		res = fmt.Sprintf("%s&start_date=%v", res, v)
	}
	if v, ok := d.GetOk("end_date"); ok {
		res = fmt.Sprintf("%s&end_date=%v", res, v)
	}
	if v, ok := d.GetOk("order_field"); ok {
		res = fmt.Sprintf("%s&order_field=%v", res, v)
	}
	if v, ok := d.GetOk("order_type"); ok {
		res = fmt.Sprintf("%s&order_type=%v", res, v)
	}
	if v, ok := d.GetOk("file_type"); ok {
		res = fmt.Sprintf("%s&file_type=%v", res, v)
	}
	if v, ok := d.GetOk("task_type"); ok {
		res = fmt.Sprintf("%s&task_type=%v", res, v)
	}

	return res
}

func dataSourceCacheHistoryTasksRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		httpUrl    = "v1.0/cdn/historytasks?page_size={page_size}"
		pageNumber = 1
		pageSize   = 10000
		result     = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient("cdn", "")
	if err != nil {
		return diag.Errorf("error creating CDN client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{page_size}", strconv.Itoa(pageSize))
	getPath += buildHistoryTaskQueryParams(d, cfg)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	for {
		currentPath := fmt.Sprintf("%s&page_number=%v", getPath, pageNumber)
		resp, err := client.Request("GET", currentPath, &getOpt)
		if err != nil {
			return diag.Errorf("error querying cache history tasks: %s", err)
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

	mErr := multierror.Append(
		d.Set("tasks", flattenCacheHistoryTasks(result)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenCacheHistoryTasks(tasks []interface{}) []interface{} {
	result := make([]interface{}, 0, len(tasks))
	for _, v := range tasks {
		result = append(result, map[string]interface{}{
			"id":         utils.PathSearch("id", v, nil),
			"status":     utils.PathSearch("status", v, nil),
			"processing": utils.PathSearch("processing", v, nil),
			"succeed":    utils.PathSearch("succeed", v, nil),
			"failed":     utils.PathSearch("failed", v, nil),
			"total":      utils.PathSearch("total", v, nil),
			"task_type":  utils.PathSearch("task_type", v, nil),
			"created_at": utils.FormatTimeStampRFC3339(
				int64(utils.PathSearch("create_time", v, float64(0)).(float64))/1000, false),
			"file_type": utils.PathSearch("file_type", v, nil),
		})
	}

	return result
}
