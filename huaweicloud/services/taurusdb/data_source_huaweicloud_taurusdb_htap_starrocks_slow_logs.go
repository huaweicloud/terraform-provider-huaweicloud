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

// @API TaurusDB POST /v3/{project_id}/instances/{instance_id}/starrocks/slow-logs
func DataSourceTaurusDBHtapStarrocksSlowLogs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTaurusDBHtapStarrocksSlowLogsRead,

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
			"operate_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"database": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"slow_log_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"count": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"lock_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rows_sent": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rows_examined": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"database": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"users": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"query_sample": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"start_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"client_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"slow_log_date": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceTaurusDBHtapStarrocksSlowLogsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/instances/{instance_id}/starrocks/slow-logs"
		limit   = 5000
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
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	for {
		getOpt.JSONBody = utils.RemoveNil(buildGetTaurusDBHtapStarrocksSlowLogsParams(d, lineNum, limit))
		getResp, err := client.Request("POST", getPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving TaurusDB HTAP StarRocks slow logs: %s", err)
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}
		slowLogs, nextLineNum := flattenTaurusDBHtapStarrocksSlowLogs(getRespBody)
		res = append(res, slowLogs...)
		if len(slowLogs) < limit {
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
		d.Set("slow_log_list", res),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetTaurusDBHtapStarrocksSlowLogsParams(d *schema.ResourceData, lineNum string, limit int) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"node_id":      d.Get("node_id").(string),
		"start_time":   d.Get("start_time").(string),
		"end_time":     d.Get("end_time").(string),
		"line_num":     utils.ValueIgnoreEmpty(lineNum),
		"limit":        limit,
		"operate_type": utils.ValueIgnoreEmpty(d.Get("operate_type").(string)),
		"database":     utils.ValueIgnoreEmpty(d.Get("database").(string)),
	}
	return bodyParams
}

func flattenTaurusDBHtapStarrocksSlowLogs(resp interface{}) ([]map[string]interface{}, string) {
	slowLogsJson := utils.PathSearch("slow_log_list", resp, make([]interface{}, 0))
	slowLogsArray := slowLogsJson.([]interface{})
	if len(slowLogsArray) == 0 {
		return nil, ""
	}

	result := make([]map[string]interface{}, 0, len(slowLogsArray))
	var lineNum string
	for _, slowLog := range slowLogsArray {
		result = append(result, map[string]interface{}{
			"node_id":       utils.PathSearch("node_id", slowLog, nil),
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
			"slow_log_date": utils.PathSearch("slow_log_date", slowLog, nil),
		})
		lineNum = utils.PathSearch("line_num", slowLog, "").(string)
	}
	return result, lineNum
}
