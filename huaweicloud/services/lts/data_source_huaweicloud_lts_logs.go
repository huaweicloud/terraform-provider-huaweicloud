package lts

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API LTS POST /v2/{project_id}/groups/{log_group_id}/streams/{log_stream_id}/content/query
func DataSourceLogs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLogsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the LTS logs are located.`,
			},
			"log_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the log group to which the logs belong.`,
			},
			"log_stream_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the log stream to which the logs belong.`,
			},
			"start_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The start time for querying log list, in RFC3339 format.`,
			},
			"end_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The end time for querying log list, in RFC3339 format.`,
			},
			"labels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The labels in key/value format to be queried.`,
			},
			"keywords": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The keywords for exact search.`,
			},
			"is_custom_time_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether to enable the custom time function for the log stream structured configuration.`,
			},
			"highlight": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether to highlight the keyword in the logs.`,
			},
			"is_desc": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether to sort the logs in descending order.`,
			},
			"is_iterative": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether to enable iterative query.`,
			},
			"logs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"content": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The content of the log.`,
						},
						"labels": {
							Type:        schema.TypeMap,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The labels associated with the log.`,
						},
						"line_num": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The line number of the log.`,
						},
					},
				},
				Description: `The list of logs that match the filter parameters.`,
			},
		},
	}
}

func dataSourceLogsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("lts", region)
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}

	logs, err := queryLogs(client, d)
	if err != nil {
		return diag.Errorf("error querying LTS logs: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("logs", flattenLogs(logs)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func queryLogs(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl    = "v2/{project_id}/groups/{log_group_id}/streams/{log_stream_id}/content/query"
		result     = make([]interface{}, 0)
		limit      = 500
		lineNum    = ""
		customTime = ""
	)

	queryPath := client.Endpoint + httpUrl
	queryPath = strings.ReplaceAll(queryPath, "{project_id}", client.ProjectID)
	queryPath = strings.ReplaceAll(queryPath, "{log_group_id}", d.Get("log_group_id").(string))
	queryPath = strings.ReplaceAll(queryPath, "{log_stream_id}", d.Get("log_stream_id").(string))
	queryOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
	}

	for {
		queryOpt.JSONBody = utils.RemoveNil(buildLogsRequestBody(d, limit, lineNum, customTime))
		resp, err := client.Request("POST", queryPath, &queryOpt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		logs := utils.PathSearch("logs", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, logs...)
		if len(logs) < limit {
			break
		}
		lineNum = utils.PathSearch("[-1].line_num", logs, "").(string)
		if d.Get("is_custom_time_enabled").(bool) {
			customTime = utils.PathSearch("[-1].labels.__time__", logs, "").(string)
		}
	}

	return result, nil
}

func buildLogsRequestBody(d *schema.ResourceData, limit int, lineNum, customTime string) map[string]interface{} {
	return map[string]interface{}{
		"start_time":   fmt.Sprintf("%d", utils.ConvertTimeStrToNanoTimestamp(d.Get("start_time").(string))),
		"end_time":     fmt.Sprintf("%d", utils.ConvertTimeStrToNanoTimestamp(d.Get("end_time").(string))),
		"labels":       utils.ValueIgnoreEmpty(d.Get("labels")),
		"keywords":     utils.ValueIgnoreEmpty(d.Get("keywords")),
		"highlight":    utils.GetNestedObjectFromRawConfig(d.GetRawConfig(), "highlight"),
		"is_desc":      d.Get("is_desc"),
		"is_iterative": d.Get("is_iterative"),
		"is_count":     true,
		"limit":        limit,
		"line_num":     utils.ValueIgnoreEmpty(lineNum),
		// search_type: The type of the pagination query.
		// + `init`: Initial query (default).
		// + `forwards`: Query the next page.
		// + `backwards`: Query the previous page.
		"search_type": utils.ValueIgnoreEmpty(d.Get("search_type")),
		// If the structured configuration of this log stream has enabled the custom time function,
		// `__time__` parameter is required for paging.
		"__time__": utils.ValueIgnoreEmpty(customTime),
	}
}

func flattenLogs(logs []interface{}) []interface{} {
	if len(logs) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(logs))
	for _, v := range logs {
		rst = append(rst, map[string]interface{}{
			"content":  utils.PathSearch("content", v, nil),
			"line_num": utils.PathSearch("line_num", v, nil),
			"labels":   utils.PathSearch("labels", v, nil),
		})
	}
	return rst
}
