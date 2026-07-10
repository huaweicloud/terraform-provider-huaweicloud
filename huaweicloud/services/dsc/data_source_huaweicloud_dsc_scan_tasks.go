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

// @API DSC GET /v1/{project_id}/sdg/scan/job/{job_id}/task
func DataSourceDscScanTasks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDscScanTasksRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"job_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The scan job ID.",
			},
			"tasks": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The scan task list.",
				Elem:        dscScanTaskSchema(),
			},
		},
	}
}

func dscScanTaskSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The task ID.",
			},
			"category": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The asset type.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The task status.",
			},
			"progress": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The task progress.",
			},
			"asset_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The asset name.",
			},
			"asset_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The asset ID.",
			},
			"start_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The task start time.",
			},
			"end_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The task end time.",
			},
			"scanned_object_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of scanned objects.",
			},
			"to_be_scanned_object_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of objects to be scanned.",
			},
			"scan_speed": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The scan speed.",
			},
			"skip_object_num": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of skipped objects.",
			},
			"last_scan_risk": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The last scan risk result.",
			},
			"security_level_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The security level name.",
			},
			"security_level_color": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The security level color.",
			},
		},
	}
}

func buildDscScanTasksQueryParams(limit, offset int) string {
	return fmt.Sprintf("?limit=%d&offset=%d", limit, offset)
}

func dataSourceDscScanTasksRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "dsc"
		httpUrl = "v1/{project_id}/sdg/scan/job/{job_id}/task"
		offset  = 0
		limit   = 1000
		result  = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{job_id}", d.Get("job_id").(string))

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"content-type": "application/json;charset=UTF-8",
		},
	}

	for {
		currentPath := requestPath + buildDscScanTasksQueryParams(limit, offset)

		resp, err := client.Request("GET", currentPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving DSC scan tasks: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		tasksResp := utils.PathSearch("tasks", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, tasksResp...)

		if len(tasksResp) < limit {
			break
		}

		offset += len(tasksResp)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("tasks", flattenDscScanTasks(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDscScanTasks(tasks []interface{}) []interface{} {
	if len(tasks) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(tasks))
	for _, v := range tasks {
		rst = append(rst, map[string]interface{}{
			"id":                       utils.PathSearch("id", v, nil),
			"category":                 utils.PathSearch("category", v, nil),
			"status":                   utils.PathSearch("status", v, nil),
			"progress":                 utils.PathSearch("progress", v, nil),
			"asset_name":               utils.PathSearch("asset_name", v, nil),
			"asset_id":                 utils.PathSearch("asset_id", v, nil),
			"start_time":               utils.PathSearch("start_time", v, nil),
			"end_time":                 utils.PathSearch("end_time", v, nil),
			"scanned_object_num":       utils.PathSearch("scanned_object_num", v, nil),
			"to_be_scanned_object_num": utils.PathSearch("to_be_scanned_object_num", v, nil),
			"scan_speed":               utils.PathSearch("scan_speed", v, nil),
			"skip_object_num":          utils.PathSearch("skip_object_num", v, nil),
			"last_scan_risk":           utils.PathSearch("last_scan_risk", v, nil),
			"security_level_name":      utils.PathSearch("security_level_name", v, nil),
			"security_level_color":     utils.PathSearch("security_level_color", v, nil),
		})
	}

	return rst
}
