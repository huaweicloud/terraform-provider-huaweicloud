package rds

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

// @API RDS POST /v3/{project_id}/instances/{instance_id}/slow-logs
func DataSourceRdsSlowLogs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRdsSlowLogsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the instance.`,
			},
			"start_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the start time in the **yyyy-mm-ddThh:mm:ssZ** format.`,
			},
			"end_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the end time in the **yyyy-mm-ddThh:mm:ssZ** format.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the statement type.`,
			},
			"database": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the name of the database.`,
			},
			"users": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the name of the account.`,
			},
			"slow_logs": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the list of the slow logs.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"count": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the number of execution times.`,
						},
						"time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the execution time.`,
						},
						"lock_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the wait lock time.`,
						},
						"rows_sent": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the number of result lines.`,
						},
						"rows_examined": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the number of rows scanned.`,
						},
						"database": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the name of the database.`,
						},
						"users": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the name of the account.`,
						},
						"query_sample": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the execution syntax.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the statement type.`,
						},
						"start_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the start time in the **yyyy-mm-ddThh:mm:ssZ** format.`,
						},
						"client_ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the IP address of the client.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceRdsSlowLogsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/slow-logs"
		product = "rds"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RDS Client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", d.Get("instance_id").(string))

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	lineNum := ""
	limit := 100
	res := make([]map[string]interface{}, 0)
	for {
		getOpt.JSONBody = utils.RemoveNil(buildGetSlowLogsParams(d, lineNum, limit))
		getResp, err := client.Request("POST", getPath, &getOpt)

		if err != nil {
			return diag.Errorf("error retrieving RDS slow logs: %s", err)
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}
		slowLogs, nextLineNum := flattenRdsSlowLogs(getRespBody)
		filterSlowLogs := filterRdsSlowLogs(d, slowLogs)
		res = append(res, filterSlowLogs...)
		if len(slowLogs) < limit {
			break
		}
		lineNum = nextLineNum
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("slow_logs", res),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetSlowLogsParams(d *schema.ResourceData, lineNum string, limit int) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"start_time": d.Get("start_time").(string),
		"end_time":   d.Get("end_time").(string),
		"line_num":   utils.ValueIgnoreEmpty(lineNum),
		"limit":      limit,
		"type":       utils.ValueIgnoreEmpty(d.Get("type").(string)),
		"database":   utils.ValueIgnoreEmpty(d.Get("database").(string)),
	}
	return bodyParams
}

func filterRdsSlowLogs(d *schema.ResourceData, resp []map[string]interface{}) []map[string]interface{} {
	if resp == nil {
		return nil
	}
	result := make([]map[string]interface{}, 0)

	rawUsers, rawUsersOK := d.GetOk("users")

	for _, slowLog := range resp {
		users := utils.PathSearch("users", slowLog, nil)
		if rawUsersOK && rawUsers != users {
			continue
		}
		result = append(result, slowLog)
	}
	return result
}

func flattenRdsSlowLogs(resp interface{}) ([]map[string]interface{}, string) {
	slowLogsJson := utils.PathSearch("slow_log_list", resp, make([]interface{}, 0))
	slowLogsArray := slowLogsJson.([]interface{})
	if len(slowLogsArray) == 0 {
		return nil, ""
	}

	result := make([]map[string]interface{}, 0, len(slowLogsArray))
	var lineNum string
	for _, slowLog := range slowLogsArray {
		result = append(result, map[string]interface{}{
			"count":         utils.PathSearch("count", slowLog, nil),
			"time":          utils.PathSearch("time", slowLog, nil),
			"lock_time":     utils.PathSearch("lock_time", slowLog, nil),
			"rows_sent":     utils.PathSearch("rows_sent", slowLog, nil),
			"rows_examined": utils.PathSearch("rows_examined", slowLog, nil),
			"database":      utils.PathSearch("database", slowLog, nil),
			"users":         utils.PathSearch("users", slowLog, nil),
			"query_sample":  utils.PathSearch("query_sample", slowLog, nil),
			"type":          utils.PathSearch("type", slowLog, nil),
			"start_time":    utils.PathSearch("start_time", slowLog, nil),
			"client_ip":     utils.PathSearch("client_ip", slowLog, nil),
		})
		lineNum = utils.PathSearch("line_num", slowLog, "").(string)
	}
	return result, lineNum
}
