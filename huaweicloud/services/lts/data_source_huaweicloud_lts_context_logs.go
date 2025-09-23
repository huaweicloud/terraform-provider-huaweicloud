package lts

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

// @API LTS POST /v2/{project_id}/groups/{log_group_id}/streams/{log_stream_id}/context
func DataSourceContextLogs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceContextLogsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the context logs are located.`,
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
			"line_num": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The sequence number of a log line.`,
			},
			"time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The time field of the custom time function, in millisecond timestamp.`,
			},
			"backwards_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The number of logs before the start log.`,
			},
			"forwards_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The number of logs after the start log.`,
			},
			"logs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"content": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The original log data.`,
						},
						"line_num": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The log line sequence number.`,
						},
						"labels": {
							Type:        schema.TypeMap,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The labels contained in this log entry.`,
						},
					},
				},
				Description: `The context log information.`,
			},
			"backwards_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of logs queried backward based on 'line_num'.`,
			},
			"forwards_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of logs queried forward based on 'line_num'.`,
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The total number of logs, including the starting log specified in the request parameters.`,
			},
		},
	}
}

func dataSourceContextLogsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		httpUrl     = "v2/{project_id}/groups/{log_group_id}/streams/{log_stream_id}/context"
		logStreamId = d.Get("log_stream_id").(string)
	)

	client, err := cfg.NewServiceClient("lts", region)
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{log_group_id}", d.Get("log_group_id").(string))
	getPath = strings.ReplaceAll(getPath, "{log_stream_id}", logStreamId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildQueryContextLogsBodyParams(d)),
	}

	resp, err := client.Request("POST", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error querying context logs under log stream (%s): %s", logStreamId, err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("logs", flattenContextLogs(utils.PathSearch("logs", respBody, make([]interface{}, 0)).([]interface{}))),
		d.Set("backwards_count", utils.PathSearch("backwards_count", respBody, nil)),
		d.Set("forwards_count", utils.PathSearch("forwards_count", respBody, nil)),
		d.Set("total_count", utils.PathSearch("total_count", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildQueryContextLogsBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"line_num":      utils.ValueIgnoreEmpty(d.Get("line_num").(string)),
		"__time__":      utils.ValueIgnoreEmpty(d.Get("time").(string)),
		"backwardsSize": utils.ValueIgnoreEmpty(d.Get("backwards_size").(int)),
		"forwardsSize":  utils.ValueIgnoreEmpty(d.Get("forwards_size").(int)),
	}
}

func flattenContextLogs(contextLogs []interface{}) []interface{} {
	if len(contextLogs) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(contextLogs))
	for _, v := range contextLogs {
		result = append(result, map[string]interface{}{
			"content":  utils.PathSearch("content", v, nil),
			"line_num": utils.PathSearch("line_num", v, nil),
			"labels":   utils.PathSearch("labels", v, nil),
		})
	}

	return result
}
