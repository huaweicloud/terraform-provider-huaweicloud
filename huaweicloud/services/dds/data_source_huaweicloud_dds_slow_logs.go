package dds

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

// @API DDS POST /v3.1/{project_id}/instances/{instance_id}/slow-logs
func DataSourceDDSSlowLogs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDDSSlowLogsRead,

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
			"operate_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the statement type.`,
			},
			"node_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the node ID.`,
			},
			"keywords": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the full-text log search based on multiple keywords.`,
			},
			"database_keywords": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the fuzzy search for logs based on multiple database keywords.`,
			},
			"collection_keywords": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Specifies the fuzzy search for logs based on multiple database table name keywords.`,
			},
			"max_cost_time": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `Specifies the logs can be searched based on the maximum execution duration. Unit is ms.`,
			},
			"min_cost_time": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `Specifies the logs can be searched based on the minimum execution duration. Unit is ms.`,
			},
			"slow_logs": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the list of the slow logs.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"node_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the node ID.`,
						},
						"node_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the node name.`,
						},
						"whole_message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the statement.`,
						},
						"operate_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the statement type.`,
						},
						"cost_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the execution time. Unit is ms.`,
						},
						"lock_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the lock wait time. Unit is ms.`,
						},
						"docs_returned": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the number of returned documents.`,
						},
						"docs_scanned": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the number of scanned documents.`,
						},
						"database": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the name of the database which the log belongs to.`,
						},
						"collection": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the name of the database table which the log belongs to.`,
						},
						"log_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the time of the slow log in the **yyyy-mm-ddThh:mm:ssZ** format.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceDDSSlowLogsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dds", region)
	if err != nil {
		return diag.Errorf("error creating DDS client: %s", err)
	}

	httpUrl := "v3.1/{project_id}/instances/{instance_id}/slow-logs"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", d.Get("instance_id").(string))
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	lineNum := ""
	limit := 100
	rst := make([]map[string]interface{}, 0)
	for {
		getOpt.JSONBody = utils.RemoveNil(buildGetSlowLogsParams(d, lineNum, limit))
		getResp, err := client.Request("POST", getPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving DDS slow logs: %s", err)
		}
		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.Errorf("error flattening response: %s", err)
		}

		slowLogs, nextLineNum := flattenSlowLogs(getRespBody)
		rst = append(rst, slowLogs...)

		if len(slowLogs) < limit {
			break
		}

		lineNum = nextLineNum
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("slow_logs", rst),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetSlowLogsParams(d *schema.ResourceData, lineNum string, limit int) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"start_time":          d.Get("start_time"),
		"end_time":            d.Get("end_time"),
		"line_num":            utils.ValueIgnoreEmpty(lineNum),
		"limit":               limit,
		"operate_type":        utils.ValueIgnoreEmpty(d.Get("operate_type")),
		"node_id":             utils.ValueIgnoreEmpty(d.Get("node_id")),
		"keywords":            utils.ValueIgnoreEmpty(d.Get("keywords")),
		"database_keywords":   utils.ValueIgnoreEmpty(d.Get("database_keywords")),
		"collection_keywords": utils.ValueIgnoreEmpty(d.Get("collection_keywords")),
		"max_cost_time":       utils.ValueIgnoreEmpty(d.Get("max_cost_time")),
		"min_cost_time":       utils.ValueIgnoreEmpty(d.Get("min_cost_time")),
	}

	return bodyParams
}

func flattenSlowLogs(resp interface{}) ([]map[string]interface{}, string) {
	slowLogs := utils.PathSearch("slow_logs", resp, make([]interface{}, 0)).([]interface{})
	if len(slowLogs) == 0 {
		return nil, ""
	}

	result := make([]map[string]interface{}, 0, len(slowLogs))
	var lineNum string
	for _, slowLog := range slowLogs {
		result = append(result, map[string]interface{}{
			"node_id":       utils.PathSearch("node_id", slowLog, nil),
			"node_name":     utils.PathSearch("node_name", slowLog, nil),
			"whole_message": utils.PathSearch("whole_message", slowLog, nil),
			"operate_type":  utils.PathSearch("operate_type", slowLog, nil),
			"cost_time":     utils.PathSearch("cost_time", slowLog, nil),
			"lock_time":     utils.PathSearch("lock_time", slowLog, nil),
			"docs_returned": utils.PathSearch("docs_returned", slowLog, nil),
			"docs_scanned":  utils.PathSearch("docs_scanned", slowLog, nil),
			"database":      utils.PathSearch("database", slowLog, nil),
			"collection":    utils.PathSearch("collection", slowLog, nil),
			"log_time":      utils.PathSearch("log_time", slowLog, nil),
		})

		lineNum = utils.PathSearch("line_num", slowLog, "").(string)
	}

	return result, lineNum
}
