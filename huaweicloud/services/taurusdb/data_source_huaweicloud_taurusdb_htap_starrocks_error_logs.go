package taurusdb

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API TaurusDB POST /v3/{project_id}/instances/{instance_id}/starrocks/error-logs
func DataSourceTaurusDBHtapStarrocksErrorLogs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTaurusDBHtapStarrocksErrorLogsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"node_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"start_time": {
				Type:     schema.TypeString,
				Required: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Required: true,
			},
			"level": {
				Type:     schema.TypeString,
				Required: true,
			},
			"error_log_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"level": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"content": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceTaurusDBHtapStarrocksErrorLogsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/instances/{instance_id}/starrocks/error-logs"
		limit   = 100
		lineNum = ""
		res     = make([]map[string]interface{}, 0)
	)
	client, err := cfg.NewServiceClient("gaussdb", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB Client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", d.Get("instance_id").(string))

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		getOpt.JSONBody = utils.RemoveNil(buildGetTaurusDBHtapStarrocksErrorLogsParams(d, lineNum, limit))
		getResp, err := client.Request("POST", getPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving TaurusDB HTAP StarRocks error logs: %s", err)
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}
		errorLogs, nextLineNum := flattenTaurusDBHtapStarrocksErrorLogs(getRespBody)
		res = append(res, errorLogs...)
		if len(errorLogs) < limit {
			break
		}
		lineNum = nextLineNum
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("error_log_list", res),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetTaurusDBHtapStarrocksErrorLogsParams(d *schema.ResourceData, lineNum string, limit int) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"node_id":    d.Get("node_id").(string),
		"start_time": d.Get("start_time").(string),
		"end_time":   d.Get("end_time").(string),
		"level":      d.Get("level").(string),
		"line_num":   utils.ValueIgnoreEmpty(lineNum),
		"limit":      limit,
	}
	return bodyParams
}

func flattenTaurusDBHtapStarrocksErrorLogs(resp interface{}) ([]map[string]interface{}, string) {
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
