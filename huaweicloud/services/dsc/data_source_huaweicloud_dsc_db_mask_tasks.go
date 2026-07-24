package dsc

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DSC GET /v1/{project_id}/sdg/server/mask/dbs/templates/{template_id}/tasks
func DataSourceDscDbMaskTasks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDscDbMaskTasksRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"template_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tasks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dbMaskTaskSchema(),
			},
		},
	}
}

func dbMaskTaskSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"db_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"end_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"execute_line": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"progress": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"run_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"start_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"task_template_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildDbMaskTasksQueryParams(limit, offset int) string {
	return fmt.Sprintf("?limit=%d&offset=%d", limit, offset)
}

func dataSourceDscDbMaskTasksRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/sdg/server/mask/dbs/templates/{template_id}/tasks"
		product = "dsc"
		limit   = 1000
		offset  = 0
		result  = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{template_id}", d.Get("template_id").(string))

	for {
		currentPath := requestPath + buildDbMaskTasksQueryParams(limit, offset)
		requestOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		requestResp, err := client.Request("GET", currentPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving DSC db mask tasks: %s", err)
		}

		requestRespBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return diag.FromErr(err)
		}

		tasksList := utils.PathSearch("tasks", requestRespBody, make([]interface{}, 0)).([]interface{})
		result = append(result, tasksList...)
		if len(tasksList) < limit {
			break
		}
		offset += len(tasksList)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(dataSourceId.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("tasks", flattenDbMaskTasks(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDbMaskTasks(items []interface{}) []interface{} {
	if len(items) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(items))
	for _, v := range items {
		rst = append(rst, map[string]interface{}{
			"db_type":          utils.PathSearch("db_type", v, nil),
			"end_time":         utils.PathSearch("end_time", v, nil),
			"execute_line":     utils.PathSearch("execute_line", v, nil),
			"id":               utils.PathSearch("id", v, nil),
			"progress":         utils.PathSearch("progress", v, nil),
			"run_status":       utils.PathSearch("run_status", v, nil),
			"start_time":       utils.PathSearch("start_time", v, nil),
			"task_template_id": utils.PathSearch("task_template_id", v, nil),
			"type":             utils.PathSearch("type", v, nil),
		})
	}

	return rst
}
