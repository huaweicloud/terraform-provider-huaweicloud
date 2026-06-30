package secmaster

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

// @API SecMaster POST /v1/{project_id}/workspaces/{workspace_id}/siem/moniter/metric/statistics
func DataSourceMoniterMetricStats() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMoniterMetricStatsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"workspace_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"dataspace_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"pipe_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"start_timestamp": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"end_timestamp": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"results": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     moniterMetricStatsSchema(),
			},
		},
	}
}

func moniterMetricStatsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"average_msg_bytes": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"subscribe_msgs": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func buildMoniterMetricStatsBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"dataspace_id":    d.Get("dataspace_id").(string),
		"pipe_id":         d.Get("pipe_id").(string),
		"start_timestamp": d.Get("start_timestamp").(int),
		"end_timestamp":   d.Get("end_timestamp").(int),
	}
}

func dataSourceMoniterMetricStatsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		product     = "secmaster"
		workspaceId = d.Get("workspace_id").(string)
		httpUrl     = "v1/{project_id}/workspaces/{workspace_id}/siem/moniter/metric/statistics"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", workspaceId)

	reqOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: buildMoniterMetricStatsBodyParams(d),
	}

	resp, err := client.Request("POST", requestPath, &reqOpt)
	if err != nil {
		return diag.Errorf("error retrieving SecMaster moniter metric stats: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("results", flattenMoniterMetricStats(utils.PathSearch("result", respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenMoniterMetricStats(result []interface{}) []interface{} {
	if len(result) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(result))
	for _, item := range result {
		rst = append(rst, map[string]interface{}{
			"average_msg_bytes": utils.PathSearch("average_msg_bytes", item, nil),
			"subscribe_msgs":    utils.PathSearch("subscribe_msgs", item, nil),
		})
	}

	return rst
}
