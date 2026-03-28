package rds

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

// @API RDS GET /v3/{project_id}/instances/{instance_id}/ops/intelligent-kill-session/statistic
func DataSourceIntelligentSessionKillStatistic() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIntelligentSessionKillStatisticRead,

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
			"node_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"statistics": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     intelligentSessionKillHistoryStatisticsSchema(),
			},
		},
	}
}

func intelligentSessionKillHistoryStatisticsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"keyword": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"raw_sql_text": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"total_time": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"avg_time": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"max_time": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"strategy": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"advice_concurrency": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func dataSourceIntelligentSessionKillStatisticRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/ops/intelligent-kill-session/statistic"
		product = "rds"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", d.Get("instance_id").(string))
	getPath += buildGetIntelligentSessionKillStatisticQueryParams(d)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving RDS intelligent session kill statistic: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("statistics", flattenIntelligentSessionKillStatistic(getRespBody)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetIntelligentSessionKillStatisticQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("node_id"); ok {
		res = fmt.Sprintf("%s&node_id=%s", res, v)
	}
	if res != "" {
		res = "?" + res[1:]
	}
	return res
}

func flattenIntelligentSessionKillStatistic(resp interface{}) []map[string]interface{} {
	curJson := utils.PathSearch("statistics", resp, nil)
	if curJson == nil {
		return nil
	}

	curArray := curJson.([]interface{})
	rst := make([]map[string]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"keyword":            utils.PathSearch("keyword", v, nil),
			"raw_sql_text":       utils.PathSearch("raw_sql_text", v, nil),
			"ids":                utils.PathSearch("ids", v, nil),
			"count":              utils.PathSearch("count", v, nil),
			"total_time":         utils.PathSearch("total_time", v, nil),
			"avg_time":           utils.PathSearch("avg_time", v, nil),
			"max_time":           utils.PathSearch("max_time", v, nil),
			"strategy":           utils.PathSearch("strategy", v, nil),
			"advice_concurrency": utils.PathSearch("advice_concurrency", v, nil),
			"type":               utils.PathSearch("type", v, nil),
		})
	}
	return rst
}
