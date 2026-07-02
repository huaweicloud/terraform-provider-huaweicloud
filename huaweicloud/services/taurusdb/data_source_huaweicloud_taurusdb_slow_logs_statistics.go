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

// @API TaurusDB POST /v3/{project_id}/instances/{instance_id}/slow-logs/statistics
func DataSourceTaurusDBSlowLogsStatistics() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTaurusDBSlowLogsStatisticsRead,

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
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"database": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sort": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"order": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"slow_log_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"client_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"count": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"database": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"lock_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"node_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"query_sample": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rows_examined": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"rows_sent": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"users": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceTaurusDBSlowLogsStatisticsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/instances/{instance_id}/slow-logs/statistics"
		limit   = 100
		offset  = 0
		res     = make([]map[string]interface{}, 0)
	)
	client, err := cfg.NewServiceClient("gaussdb", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB Client: %s", err)
	}

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", d.Get("instance_id").(string))

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listOpt.JSONBody = utils.RemoveNil(buildGetSlowLogsStatisticsBodyParams(d, limit, offset))
		listResp, err := client.Request("POST", listPath, &listOpt)
		if err != nil {
			return diag.Errorf("error retrieving TaurusDB slow logs statistics: %s", err)
		}

		listRespBody, err := utils.FlattenResponse(listResp)
		if err != nil {
			return diag.FromErr(err)
		}

		slowLogs := flattenGetSlowLogsStatisticsBody(listRespBody)
		res = append(res, slowLogs...)

		if len(slowLogs) < limit {
			break
		}
		offset += len(slowLogs)
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("slow_log_list", res),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetSlowLogsStatisticsBodyParams(d *schema.ResourceData, limit, offset int) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"node_id":    d.Get("node_id").(string),
		"start_time": d.Get("start_time").(string),
		"end_time":   d.Get("end_time").(string),
		"limit":      limit,
		"offset":     offset,
		"type":       utils.ValueIgnoreEmpty(d.Get("type").(string)),
		"database":   utils.ValueIgnoreEmpty(d.Get("database").(string)),
		"sort":       utils.ValueIgnoreEmpty(d.Get("sort").(string)),
		"order":      utils.ValueIgnoreEmpty(d.Get("order").(string)),
	}
	return bodyParams
}

func flattenGetSlowLogsStatisticsBody(resp interface{}) []map[string]interface{} {
	curJson := utils.PathSearch("slow_log_list", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]map[string]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"client_ip":     utils.PathSearch("client_ip", v, nil),
			"count":         utils.PathSearch("count", v, nil),
			"database":      utils.PathSearch("database", v, nil),
			"lock_time":     utils.PathSearch("lock_time", v, nil),
			"node_id":       utils.PathSearch("node_id", v, nil),
			"query_sample":  utils.PathSearch("query_sample", v, nil),
			"rows_examined": utils.PathSearch("rows_examined", v, nil),
			"rows_sent":     utils.PathSearch("rows_sent", v, nil),
			"time":          utils.PathSearch("time", v, nil),
			"type":          utils.PathSearch("type", v, nil),
			"users":         utils.PathSearch("users", v, nil),
		})
	}
	return res
}
