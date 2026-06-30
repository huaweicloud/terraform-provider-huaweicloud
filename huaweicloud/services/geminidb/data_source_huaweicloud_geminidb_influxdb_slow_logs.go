package geminidb

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

// @API GeminiDB POST /v3/{project_id}/influxdb/instances/{instance_id}/slow-logs
func DataSourceInfluxdbSlowLogs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceInfluxdbSlowLogsRead,

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
			"node_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"keywords": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"max_cost_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"min_cost_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"database": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"retention_policy": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"slow_logs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"node_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"whole_message": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"operate_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cost_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"log_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"database": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"retention_policy": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"line_num": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildInfluxdbSlowLogsBodyParams(d *schema.ResourceData, limit int, lineNum string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"start_time":       d.Get("start_time"),
		"end_time":         d.Get("end_time"),
		"limit":            limit,
		"line_num":         utils.ValueIgnoreEmpty(lineNum),
		"operate_type":     utils.ValueIgnoreEmpty(d.Get("operate_type")),
		"node_id":          utils.ValueIgnoreEmpty(d.Get("node_id")),
		"keywords":         utils.ValueIgnoreEmpty(d.Get("keywords")),
		"max_cost_time":    utils.ValueIgnoreEmpty(d.Get("max_cost_time")),
		"min_cost_time":    utils.ValueIgnoreEmpty(d.Get("min_cost_time")),
		"database":         utils.ValueIgnoreEmpty(d.Get("database")),
		"retention_policy": utils.ValueIgnoreEmpty(d.Get("retention_policy")),
	}

	return bodyParams
}

func dataSourceInfluxdbSlowLogsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		httpUrl    = "v3/{project_id}/influxdb/instances/{instance_id}/slow-logs"
		instanceId = d.Get("instance_id").(string)
		lineNum    = ""
		limit      = 100
		result     = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient("geminidb", region)
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", instanceId)
	listOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		KeepResponseBody: true,
	}

	for {
		listOpt.JSONBody = utils.RemoveNil(buildInfluxdbSlowLogsBodyParams(d, limit, lineNum))
		resp, err := client.Request("POST", listPath, &listOpt)
		if err != nil {
			return diag.Errorf("error retrieving GeminiDB InfluxDB instance slow logs: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		logs := utils.PathSearch("slow_logs", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, logs...)
		if len(logs) < limit {
			break
		}

		lineNum = utils.PathSearch("[-1].line_num", logs, "").(string)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("slow_logs", flattenInfluxdbSlowLogs(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenInfluxdbSlowLogs(logs []interface{}) []interface{} {
	if len(logs) == 0 {
		return nil
	}

	rst := make([]interface{}, len(logs))
	for i, v := range logs {
		rst[i] = map[string]interface{}{
			"node_id":          utils.PathSearch("node_id", v, nil),
			"node_name":        utils.PathSearch("node_name", v, nil),
			"whole_message":    utils.PathSearch("whole_message", v, nil),
			"operate_type":     utils.PathSearch("operate_type", v, nil),
			"cost_time":        utils.PathSearch("cost_time", v, nil),
			"log_time":         utils.PathSearch("log_time", v, nil),
			"database":         utils.PathSearch("database", v, nil),
			"retention_policy": utils.PathSearch("retention_policy", v, nil),
			"line_num":         utils.PathSearch("line_num", v, nil),
		}
	}
	return rst
}
