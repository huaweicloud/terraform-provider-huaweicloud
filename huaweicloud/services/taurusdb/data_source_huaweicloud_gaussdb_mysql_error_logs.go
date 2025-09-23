package taurusdb

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

// @API GaussDBforMySQL POST /v3.1/{project_id}/instances/{instance_id}/error-logs
func DataSourceGaussDBMysqlErrorLogs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGaussDBMysqlErrorLogsRead,

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
				Description: `Specifies the ID of the GaussDB MySQL instance.`,
			},
			"node_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the GaussDB MySQL instance node.`,
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
			"level": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the log level.`,
			},
			"error_log_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the list of the error logs.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the ID of the GaussDB MySQL instance node.`,
						},
						"time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the execution time.`,
						},
						"level": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the error log level.`,
						},
						"content": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the error log content.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceGaussDBMysqlErrorLogsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v3.1/{project_id}/instances/{instance_id}/error-logs"
		product = "gaussdb"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB Client: %s", err)
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
		getOpt.JSONBody = utils.RemoveNil(buildGetErrorLogsParams(d, lineNum, limit))
		getResp, err := client.Request("POST", getPath, &getOpt)

		if err != nil {
			return diag.Errorf("error retrieving GaussDB MySQL error logs: %s", err)
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}
		errorLogs, nextLineNum := flattenGaussDBMysqlGetErrorLogs(getRespBody)
		res = append(res, errorLogs...)
		if len(errorLogs) < limit {
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
		d.Set("error_log_list", res),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetErrorLogsParams(d *schema.ResourceData, lineNum string, limit int) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"node_id":    d.Get("node_id").(string),
		"start_time": d.Get("start_time").(string),
		"end_time":   d.Get("end_time").(string),
		"line_num":   utils.ValueIgnoreEmpty(lineNum),
		"limit":      limit,
		"level":      utils.ValueIgnoreEmpty(d.Get("level").(string)),
	}
	return bodyParams
}

func flattenGaussDBMysqlGetErrorLogs(resp interface{}) ([]map[string]interface{}, string) {
	errorLogsJson := utils.PathSearch("error_log_list", resp, make([]interface{}, 0))
	errorLogsArray := errorLogsJson.([]interface{})
	if len(errorLogsArray) == 0 {
		return nil, ""
	}

	result := make([]map[string]interface{}, 0, len(errorLogsArray))
	var lineNum string
	for _, errorLog := range errorLogsArray {
		result = append(result, map[string]interface{}{
			"node_id": utils.PathSearch("node_id", errorLog, nil),
			"time":    utils.PathSearch("time", errorLog, nil),
			"level":   utils.PathSearch("level", errorLog, nil),
			"content": utils.PathSearch("content", errorLog, nil),
		})
		lineNum = utils.PathSearch("line_num", errorLog, "").(string)
	}
	return result, lineNum
}
