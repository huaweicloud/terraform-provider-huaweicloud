package waf

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

// @API WAF GET /v1/{project_id}/waf/overviews/bandwidth/timeline
func DataSourceWafOverviewsBandwidthTimeline() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceWafOverviewsBandwidthTimelineRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"from": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"to": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"group_by": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"display_option": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"hosts": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instances": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"bandwidths": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"timeline": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"time": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"num": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func buildOverviewsBandwidthQueryParams(cfg *config.Config, d *schema.ResourceData) string {
	epsId := cfg.GetEnterpriseProjectID(d)
	rst := fmt.Sprintf("?from=%v&to=%v", d.Get("from"), d.Get("to"))

	if groupBy, ok := d.GetOk("group_by"); ok {
		rst += fmt.Sprintf("&group_by=%v", groupBy)
	}

	if displayOption, ok := d.GetOk("display_option"); ok {
		rst += fmt.Sprintf("&display_option=%v", displayOption)
	}

	if hostId, ok := d.GetOk("hosts"); ok {
		rst += fmt.Sprintf("&hosts=%v", hostId)
	}

	if instanceId, ok := d.GetOk("instances"); ok {
		rst += fmt.Sprintf("&instances=%v", instanceId)
	}

	if epsId != "" {
		rst += fmt.Sprintf("&enterprise_project_id=%v", epsId)
	}

	return rst
}

func dataSourceWafOverviewsBandwidthTimelineRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/waf/overviews/bandwidth/timeline"
	)

	client, err := cfg.NewServiceClient("waf", region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath += buildOverviewsBandwidthQueryParams(cfg, d)

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
	}

	listResp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return diag.Errorf("error retrieving bandwidth usage statistics: %s", err)
	}

	listRespBody, err := utils.FlattenResponse(listResp)
	if err != nil {
		return diag.FromErr(err)
	}

	listArray, ok := listRespBody.([]interface{})
	if !ok {
		return diag.Errorf("convert inteface array failed")
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("bandwidths", flattenOverviewsBandwidths(listArray)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenOverviewsBandwidths(rawBandwidths []interface{}) []interface{} {
	if len(rawBandwidths) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(rawBandwidths))
	for _, v := range rawBandwidths {
		rst = append(rst, map[string]interface{}{
			"key":      utils.PathSearch("key", v, nil),
			"timeline": flattenOverviewsBandwidthsTimeline(utils.PathSearch("timeline", v, make([]interface{}, 0)).([]interface{})),
		})
	}

	return rst
}

func flattenOverviewsBandwidthsTimeline(rawTimeline []interface{}) []interface{} {
	if len(rawTimeline) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(rawTimeline))
	for _, v := range rawTimeline {
		rst = append(rst, map[string]interface{}{
			"time": utils.PathSearch("time", v, nil),
			"num":  utils.PathSearch("num", v, nil),
		})
	}

	return rst
}
