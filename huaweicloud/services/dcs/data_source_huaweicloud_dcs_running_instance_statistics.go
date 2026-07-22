package dcs

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

// @API DCS GET /v2/{project_id}/instances/statistic
func DataSourceDcsRunningInstanceStatistics() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDcsRunningInstanceStatisticsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"statistics": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     runningInstanceStatisticsSchema(),
			},
		},
	}
}

func runningInstanceStatisticsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"input_kbps": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"output_kbps": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"keys": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"used_memory": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"max_memory": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cmd_get_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cmd_set_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"used_cpu": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceDcsRunningInstanceStatisticsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("dcs", region)
	if err != nil {
		return diag.Errorf("error creating DCS client: %s", err)
	}

	httpUrl := "v2/{project_id}/instances/statistic"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving DCS running instance statistics: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("statistics", flattenGetRunningInstanceStatisticsBody(getRespBody)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetRunningInstanceStatisticsBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("statistics", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"instance_id":   utils.PathSearch("instance_id", v, nil),
			"input_kbps":    utils.PathSearch("input_kbps", v, nil),
			"output_kbps":   utils.PathSearch("output_kbps", v, nil),
			"keys":          utils.PathSearch("keys", v, nil),
			"used_memory":   utils.PathSearch("used_memory", v, nil),
			"max_memory":    utils.PathSearch("max_memory", v, nil),
			"cmd_get_count": utils.PathSearch("cmd_get_count", v, nil),
			"cmd_set_count": utils.PathSearch("cmd_set_count", v, nil),
			"used_cpu":      utils.PathSearch("used_cpu", v, nil),
		})
	}
	return res
}
