package geminidb

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GeminiDB GET /v3/{project_id}/instances/{instance_id}/slowlog
func DataSourceGeminiDBSlowLogs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGeminiDBSlowLogsRead,

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
			"start_date": {
				Type:     schema.TypeString,
				Required: true,
			},
			"end_date": {
				Type:     schema.TypeString,
				Required: true,
			},
			"node_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"slow_log_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"database": {
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
					},
				},
			},
		},
	}
}

func dataSourceGeminiDBSlowLogsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("geminidb", region)
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)

	httpUrl := "v3/{project_id}/instances/{instance_id}/slowlog"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceID)

	// Build query parameters
	queryParams := buildSlowLogsQueryParams(d)
	getPath += queryParams

	resp, err := pagination.ListAllItems(
		client,
		"offset",
		getPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return diag.Errorf("error retrieving scheduled jobs: %s", err)
	}

	respJson, err := json.Marshal(resp)
	if err != nil {
		return diag.FromErr(err)
	}
	var respBody interface{}
	err = json.Unmarshal(respJson, &respBody)
	if err != nil {
		return diag.FromErr(err)
	}

	slowLogList := utils.PathSearch("slow_log_list", respBody, []interface{}{}).([]interface{})

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("slow_log_list", flattenSlowLogList(slowLogList)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildSlowLogsQueryParams(d *schema.ResourceData) string {
	queryParams := ""

	queryParams = fmt.Sprintf("%s&start_date=%v", queryParams, d.Get("start_date"))
	queryParams = fmt.Sprintf("%s&end_date=%v", queryParams, d.Get("end_date"))
	if v, ok := d.GetOk("node_id"); ok {
		queryParams = fmt.Sprintf("%s&node_id=%v", queryParams, v)
	}
	if v, ok := d.GetOk("type"); ok {
		queryParams = fmt.Sprintf("%s&type=%v", queryParams, v)
	}

	if queryParams != "" {
		queryParams = "?" + queryParams[1:]
	}

	return queryParams
}

func flattenSlowLogList(slowLogList []interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(slowLogList))

	for _, log := range slowLogList {
		logMap := log.(map[string]interface{})

		result = append(result, map[string]interface{}{
			"time":         utils.PathSearch("time", logMap, ""),
			"database":     utils.PathSearch("database", logMap, ""),
			"query_sample": utils.PathSearch("query_sample", logMap, ""),
			"type":         utils.PathSearch("type", logMap, ""),
			"start_time":   utils.PathSearch("start_time", logMap, ""),
		})
	}

	return result
}
