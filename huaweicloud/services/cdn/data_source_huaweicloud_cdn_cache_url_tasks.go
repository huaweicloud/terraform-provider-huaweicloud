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

// @API CDN GET /v1.0/cdn/contentgateway/url-tasks
func DataSourceCacheUrlTasks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCacheUrlTasksRead,

		Schema: map[string]*schema.Schema{
			// Optional parameters
			"start_time": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The start timestamp.`,
			},
			"end_time": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The end timestamp.`,
			},
			"url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The refresh or preheat URL.`,
			},
			"task_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The task type.`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The URL status.`,
			},
			"file_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The file type.`,
			},

			// Attributes
			"tasks": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        cacheUrlTasksSchema(),
				Description: `The list of URL task information.`,
			},
		},
	}
}

func cacheUrlTasksSchema() *schema.Resource {
	return &schema.Resource{
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
				Description: `The modification time, in RFC3339 format.`,
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

func buildCacheUrlTasksQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("start_time"); ok {
		res = fmt.Sprintf("%s&start_time=%v", res, v)
	}
	if v, ok := d.GetOk("end_time"); ok {
		res = fmt.Sprintf("%s&end_time=%v", res, v)
	}
	if v, ok := d.GetOk("url"); ok {
		res = fmt.Sprintf("%s&url=%v", res, v)
	}
	if v, ok := d.GetOk("task_type"); ok {
		res = fmt.Sprintf("%s&task_type=%v", res, v)
	}
	if v, ok := d.GetOk("status"); ok {
		res = fmt.Sprintf("%s&status=%v", res, v)
	}
	if v, ok := d.GetOk("file_type"); ok {
		res = fmt.Sprintf("%s&file_type=%v", res, v)
	}

	return res
}

func dataSourceCacheUrlTasksRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		httpUrl = "v1.0/cdn/contentgateway/url-tasks?limit={limit}"
		offset  = 0
		limit   = 100
		result  = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient("cdn", "")
	if err != nil {
		return diag.Errorf("error creating CDN client: %s", err)
	}

	listPathWithLimit := client.Endpoint + httpUrl
	listPathWithLimit = strings.ReplaceAll(listPathWithLimit, "{limit}", strconv.Itoa(limit))
	listPathWithLimit += buildCacheUrlTasksQueryParams(d)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%v", listPathWithLimit, offset)

		resp, err := client.Request("GET", listPathWithOffset, &getOpt)
		if err != nil {
			return diag.Errorf("error querying cache url tasks: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		tasks := utils.PathSearch("result", respBody, make([]interface{}, 0)).([]interface{})
		if len(tasks) == 0 {
			break
		}
		result = append(result, tasks...)
		offset += len(tasks)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(generateUUID)

	mErr := multierror.Append(
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
