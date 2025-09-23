package dbss

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DBSS POST /v1/{project_id}/{instance_id}/dbss/audit/operate-log
func DataSourceOperationLogs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOperationLogsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the data source. If omitted, the provider-level region will be used.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the audit instance ID to which the user operation logs belong.`,
			},
			"user_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the name of the operation user.`,
			},
			"operate_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the name of the operation object.`,
			},
			"result": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the execution result of user operation.`,
			},
			"start_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the start time of the user operation.`,
			},
			"end_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the end time of the user operation.`,
			},
			"time_range": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the time segment.`,
			},
			"logs": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of the user operation logs.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the user operation log.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the operation object.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the user operation.`,
						},
						"result": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The execution result of user operation.`,
						},
						"action": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the user operation.`,
						},
						"function": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The function type of the operation record.`,
						},
						"user": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the operation user.`,
						},
						"time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The time of the operation record is generated, in UTC format.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceOperationLogsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		logsHttpUrl = "v1/{project_id}/{instance_id}/dbss/audit/operate-log"
		logsProduct = "dbss"
		mErr        *multierror.Error
	)

	client, err := cfg.NewServiceClient(logsProduct, region)
	if err != nil {
		return diag.Errorf("error creating DBSS client: %s", err)
	}

	logsPath := client.Endpoint + logsHttpUrl
	logsPath = strings.ReplaceAll(logsPath, "{project_id}", client.ProjectID)
	logsPath = strings.ReplaceAll(logsPath, "{instance_id}", d.Get("instance_id").(string))
	logsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	bodyParams := utils.RemoveNil(buildOperationLogsBodyParams(d))
	logs := make([]interface{}, 0)
	page := 1

	for {
		bodyParams["page"] = page
		logsOpt.JSONBody = bodyParams
		logsResp, err := client.Request("POST", logsPath, &logsOpt)
		if err != nil {
			return diag.FromErr(err)
		}
		logsRespBody, err := utils.FlattenResponse(logsResp)
		if err != nil {
			return diag.FromErr(err)
		}
		operateLog := utils.PathSearch("operate_log", logsRespBody, make([]interface{}, 0)).([]interface{})
		logs = append(logs, operateLog...)

		totalCount := utils.PathSearch("total_num", logsRespBody, float64(0))
		if len(logs) >= int(totalCount.(float64)) {
			break
		}
		page++
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	mErr = multierror.Append(
		nil,
		d.Set("region", region),
		d.Set("logs", flattenOperationLogsResponse(logs)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildOperationLogsBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParam := map[string]interface{}{
		"time": map[string]interface{}{
			"start_time": utils.ValueIgnoreEmpty(d.Get("start_time")),
			"end_time":   utils.ValueIgnoreEmpty(d.Get("end_time")),
			"time_range": utils.ValueIgnoreEmpty(d.Get("time_range")),
		},
		"user_name":    utils.ValueIgnoreEmpty(d.Get("user_name")),
		"operate_name": utils.ValueIgnoreEmpty(d.Get("operate_name")),
		"result":       utils.ValueIgnoreEmpty(d.Get("result")),
		"size":         100,
	}

	return bodyParam
}

func flattenOperationLogsResponse(rawParams []interface{}) []interface{} {
	if len(rawParams) == 0 {
		return nil
	}

	rst := make([]interface{}, len(rawParams))
	for i, v := range rawParams {
		rst[i] = map[string]interface{}{
			"id":          utils.PathSearch("id", v, nil),
			"name":        utils.PathSearch("name", v, nil),
			"description": utils.PathSearch("description", v, nil),
			"result":      utils.PathSearch("result", v, nil),
			"action":      utils.PathSearch("action", v, nil),
			"function":    utils.PathSearch("function", v, nil),
			"user":        utils.PathSearch("user", v, nil),
			"time":        utils.PathSearch("time", v, nil),
		}
	}
	return rst
}
