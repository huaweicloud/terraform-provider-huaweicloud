package dsc

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

// @API DSC GET /v2/{project_id}/sec-ops/events/overview
func DataSourceDscEventOverview() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDscEventOverviewRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"block_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"not_overdue_event": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dscEventLevelInfoSchema(),
			},
			"overdue_event": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dscEventLevelInfoSchema(),
			},
			"turn_off_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"turn_on_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dscEventLevelInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"fatal_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"high_risk_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"middle_risk_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"low_risk_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"notice_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceDscEventOverviewRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "dsc"
		httpUrl = "v2/{project_id}/sec-ops/events/overview"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpts)
	if err != nil {
		return diag.Errorf("error retrieving DSC event overview: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("block_num", utils.PathSearch("block_num", respBody, nil)),
		d.Set("not_overdue_event", flattenDscEventLevelInfo(utils.PathSearch("not_overdue_event", respBody, nil))),
		d.Set("overdue_event", flattenDscEventLevelInfo(utils.PathSearch("overdue_event", respBody, nil))),
		d.Set("turn_off_num", utils.PathSearch("turn_off_num", respBody, nil)),
		d.Set("turn_on_num", utils.PathSearch("turn_on_num", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDscEventLevelInfo(eventLevelInfo interface{}) []map[string]interface{} {
	if eventLevelInfo == nil {
		return nil
	}

	result := map[string]interface{}{
		"fatal_num":       utils.PathSearch("fatal_num", eventLevelInfo, nil),
		"high_risk_num":   utils.PathSearch("high_risk_num", eventLevelInfo, nil),
		"middle_risk_num": utils.PathSearch("middle_risk_num", eventLevelInfo, nil),
		"low_risk_num":    utils.PathSearch("low_risk_num", eventLevelInfo, nil),
		"notice_num":      utils.PathSearch("notice_num", eventLevelInfo, nil),
	}

	return []map[string]interface{}{result}
}
