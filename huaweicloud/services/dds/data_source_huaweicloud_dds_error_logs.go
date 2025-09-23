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

// @API DDS POST /v3.1/{project_id}/instances/{instance_id}/error-logs
func DataSourceDDSErrorLogs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDDSErrorLogsRead,

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
			"severity": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the log level.`,
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
				Description: `Specifies the full-text log search based on multiple keywords, indicating that all keywords are matched.`,
			},
			"error_logs": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the list of the error logs.`,
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
						"log_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the time of the error log in the **yyyy-mm-ddThh:mm:ssZ** format.`,
						},
						"severity": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the error log level.`,
						},
						"raw_message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the error description.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceDDSErrorLogsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dds", region)
	if err != nil {
		return diag.Errorf("error creating DDS client: %s", err)
	}

	httpUrl := "v3.1/{project_id}/instances/{instance_id}/error-logs"
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
		getOpt.JSONBody = utils.RemoveNil(buildGetErrorLogsParams(d, lineNum, limit))
		getResp, err := client.Request("POST", getPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving DDS error logs: %s", err)
		}
		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}

		errorLogs, nextLineNum := flattenErrorLogs(getRespBody)
		rst = append(rst, errorLogs...)

		if len(errorLogs) < limit {
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
		d.Set("error_logs", rst),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetErrorLogsParams(d *schema.ResourceData, lineNum string, limit int) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"start_time": d.Get("start_time"),
		"end_time":   d.Get("end_time"),
		"line_num":   utils.ValueIgnoreEmpty(lineNum),
		"limit":      limit,
		"severity":   utils.ValueIgnoreEmpty(d.Get("severity")),
		"node_id":    utils.ValueIgnoreEmpty(d.Get("node_id")),
		"keywords":   utils.ValueIgnoreEmpty(d.Get("keywords")),
	}
	return bodyParams
}

func flattenErrorLogs(resp interface{}) ([]map[string]interface{}, string) {
	errorLogsArray := utils.PathSearch("error_logs", resp, make([]interface{}, 0)).([]interface{})
	if len(errorLogsArray) == 0 {
		return nil, ""
	}

	result := make([]map[string]interface{}, 0, len(errorLogsArray))
	var lineNum string
	for _, errorLog := range errorLogsArray {
		result = append(result, map[string]interface{}{
			"node_id":     utils.PathSearch("node_id", errorLog, nil),
			"node_name":   utils.PathSearch("node_name", errorLog, nil),
			"log_time":    utils.PathSearch("log_time", errorLog, nil),
			"severity":    utils.PathSearch("severity", errorLog, nil),
			"raw_message": utils.PathSearch("raw_message", errorLog, nil),
		})

		lineNum = utils.PathSearch("line_num", errorLog, "").(string)
	}
	return result, lineNum
}
