package er

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

// @API ER GET /v3/{project_id}/enterprise-router/{er_id}/flow-logs
func DataSourceFlowLogs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFlowLogsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The region where the flow logs are located.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the ER instance to which the flow logs belong.",
			},
			"resource_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The type of the flow logs.",
			},
			"resource_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the attachment to which the flow logs belong.",
			},
			"flow_log_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the flow log.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the flow log.",
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The status of the flow logs.",
			},
			"enabled": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The switch status of the flow log.",
			},
			"log_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the log group to which the flow logs belong.",
			},
			"log_stream_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the log stream to which the flow logs belong.",
			},
			"flow_logs": {
				Type:        schema.TypeList,
				Elem:        flowLogSchema(),
				Computed:    true,
				Description: "The list of the flow logs.",
			},
		},
	}
}

func flowLogSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the flow log.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the flow log.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the flow log.",
			},
			"resource_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of the flow log.",
			},
			"resource_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the attachment to which the flow log belongs.",
			},
			"log_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the log group to which the flow log belongs.",
			},
			"log_stream_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the log stream to which the flow log belongs.",
			},
			"log_store_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The storage type of the flow log.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of the flow log.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The latest update time of the flow log.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the flow log.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "The switch of the flow log.",
			},
		},
	}
	return &sc
}

func dataSourceFlowLogsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// listFlowLog: Query the flow log list
	listFlowLogsClient, err := cfg.ErV3Client(region)
	if err != nil {
		return diag.Errorf("error creating ER v3 client: %s", err)
	}

	listFlowLogsHttpUrl := "enterprise-router/{er_id}/flow-logs"
	listFlowLogsPath := listFlowLogsClient.ResourceBaseURL() + listFlowLogsHttpUrl
	listFlowLogsPath = strings.ReplaceAll(listFlowLogsPath, "{er_id}", d.Get("instance_id").(string))

	listFlowLogsQueryParams := buildListFlowLogsQueryParams(d, cfg)
	listFlowLogsPath += listFlowLogsQueryParams

	listFlowLogsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	listFlowLogsResp, err := listFlowLogsClient.Request("GET", listFlowLogsPath, &listFlowLogsOpt)
	if err != nil {
		return diag.Errorf("error retrieving flow logs: %s", err)
	}

	listFlowLogsRespBody, err := utils.FlattenResponse(listFlowLogsResp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	var mErr *multierror.Error
	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("flow_logs", filterListFlowLogsResponseBody(flattenListTransitIpsResponseBody(listFlowLogsRespBody), d)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListTransitIpsResponseBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	curJson := utils.PathSearch("flow_logs", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":             utils.PathSearch("id", v, nil),
			"name":           utils.PathSearch("name", v, nil),
			"description":    utils.PathSearch("description", v, nil),
			"resource_type":  utils.PathSearch("resource_type", v, nil),
			"resource_id":    utils.PathSearch("resource_id", v, nil),
			"log_group_id":   utils.PathSearch("log_group_id", v, nil),
			"log_stream_id":  utils.PathSearch("log_stream_id", v, nil),
			"log_store_type": utils.PathSearch("log_store_type", v, nil),
			// The time results are not the time in RF3339 format without milliseconds.
			"created_at": utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("created_at",
				v, "").(string))/1000, false),
			"updated_at": utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("updated_at",
				v, "").(string))/1000, false),
			"status":  utils.PathSearch("state", v, nil),
			"enabled": utils.PathSearch("enabled", v, nil),
		})
	}
	return rst
}

func filterListFlowLogsResponseBody(all []interface{}, d *schema.ResourceData) []interface{} {
	rst := make([]interface{}, 0, len(all))

	for _, v := range all {
		if param, ok := d.GetOk("flow_log_id"); ok &&
			fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("id", v, nil)) {
			continue
		}
		if param, ok := d.GetOk("name"); ok &&
			fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("name", v, nil)) {
			continue
		}
		if param, ok := d.GetOk("status"); ok &&
			fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("status", v, nil)) {
			continue
		}
		if param, ok := d.GetOk("enabled"); ok &&
			fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("enabled", v, nil)) {
			continue
		}
		if param, ok := d.GetOk("log_group_id"); ok &&
			fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("log_group_id", v, nil)) {
			continue
		}
		if param, ok := d.GetOk("log_stream_id"); ok &&
			fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("log_stream_id", v, nil)) {
			continue
		}

		rst = append(rst, v)
	}
	return rst
}

func buildListFlowLogsQueryParams(d *schema.ResourceData, _ *config.Config) string {
	res := "&sort_key=name"

	if v, ok := d.GetOk("resource_type"); ok {
		res = fmt.Sprintf("%s&resource_type=%v", res, v)
	}
	if v, ok := d.GetOk("resource_id"); ok {
		res = fmt.Sprintf("%s&resource_id=%v", res, v)
	}
	if res != "" {
		res = "?" + res[1:]
	}
	return res
}
