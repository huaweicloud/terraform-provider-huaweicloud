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

// @API GeminiDB POST /v3/{project_id}/cassandra/instances/{instance_id}/slow-logs
func DataSourceCassandraSlowLogs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCassandraSlowLogsRead,

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
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"keyspace_keywords": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"table_keywords": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"max_cost_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"min_cost_time": {
				Type:     schema.TypeInt,
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
							Type:     schema.TypeInt,
							Computed: true,
						},
						"keyspace": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"table": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"log_time": {
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

func buildCassandraSlowLogsBodyParams(d *schema.ResourceData, limit int, lineNum string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"start_time":        d.Get("start_time"),
		"end_time":          d.Get("end_time"),
		"limit":             limit,
		"line_num":          utils.ValueIgnoreEmpty(lineNum),
		"operate_type":      utils.ValueIgnoreEmpty(d.Get("operate_type")),
		"node_id":           utils.ValueIgnoreEmpty(d.Get("node_id")),
		"keywords":          utils.ValueIgnoreEmpty(utils.ExpandToStringList(d.Get("keywords").([]interface{}))),
		"keyspace_keywords": utils.ValueIgnoreEmpty(utils.ExpandToStringList(d.Get("keyspace_keywords").([]interface{}))),
		"table_keywords":    utils.ValueIgnoreEmpty(utils.ExpandToStringList(d.Get("table_keywords").([]interface{}))),
		"max_cost_time":     utils.ValueIgnoreEmpty(d.Get("max_cost_time")),
		"min_cost_time":     utils.ValueIgnoreEmpty(d.Get("min_cost_time")),
	}

	return bodyParams
}

func dataSourceCassandraSlowLogsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		httpUrl    = "v3/{project_id}/cassandra/instances/{instance_id}/slow-logs"
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
		listOpt.JSONBody = utils.RemoveNil(buildCassandraSlowLogsBodyParams(d, limit, lineNum))
		resp, err := client.Request("POST", listPath, &listOpt)
		if err != nil {
			return diag.Errorf("error retrieving GeminiDB Cassandra instance slow logs: %s", err)
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
		d.Set("slow_logs", flattenCassandraSlowLogs(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenCassandraSlowLogs(instances []interface{}) []interface{} {
	if len(instances) == 0 {
		return nil
	}

	rst := make([]interface{}, len(instances))
	for i, v := range instances {
		rst[i] = map[string]interface{}{
			"node_id":       utils.PathSearch("node_id", v, nil),
			"node_name":     utils.PathSearch("node_name", v, nil),
			"whole_message": utils.PathSearch("whole_message", v, nil),
			"operate_type":  utils.PathSearch("operate_type", v, nil),
			"cost_time":     utils.PathSearch("cost_time", v, nil),
			"keyspace":      utils.PathSearch("keyspace", v, nil),
			"table":         utils.PathSearch("table", v, nil),
			"log_time":      utils.PathSearch("log_time", v, nil),
			"line_num":      utils.PathSearch("line_num", v, nil),
		}
	}
	return rst
}
