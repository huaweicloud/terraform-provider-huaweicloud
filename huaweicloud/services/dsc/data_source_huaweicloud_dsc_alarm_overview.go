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

// @API DSC GET /v2/{project_id}/sec-ops/alarms/overview
func DataSourceDscAlarmOverview() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDscAlarmOverviewRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"alarm_source_info": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dscAlarmOverviewSourceInfoSchema(),
			},
			"total_alarm": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dscAlarmOverviewLevelInfoSchema(),
			},
			"turn_off_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"turn_on_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"untreated_alarm": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dscAlarmOverviewLevelInfoSchema(),
			},
		},
	}
}

func dscAlarmOverviewSourceInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"api_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cbh_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"database_encrypt_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"database_op_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"dbss_num": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dscAlarmOverviewLevelInfoSchema() *schema.Resource {
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

func dataSourceDscAlarmOverviewRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "dsc"
		httpUrl = "v2/{project_id}/sec-ops/alarms/overview"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving DSC alarm overview: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("alarm_source_info", flattenDscAlarmOverviewSourceInfo(
			utils.PathSearch("alarm_source_info", respBody, nil))),
		d.Set("total_alarm", flattenDscAlarmOverviewLevelInfo(
			utils.PathSearch("total_alarm", respBody, nil))),
		d.Set("turn_off_num", utils.PathSearch("turn_off_num", respBody, nil)),
		d.Set("turn_on_num", utils.PathSearch("turn_on_num", respBody, nil)),
		d.Set("untreated_alarm", flattenDscAlarmOverviewLevelInfo(
			utils.PathSearch("untreated_alarm", respBody, nil))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDscAlarmOverviewSourceInfo(sourceInfo interface{}) []interface{} {
	if sourceInfo == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"api_num":              utils.PathSearch("api_num", sourceInfo, nil),
			"cbh_num":              utils.PathSearch("cbh_num", sourceInfo, nil),
			"database_encrypt_num": utils.PathSearch("database_encrypt_num", sourceInfo, nil),
			"database_op_num":      utils.PathSearch("database_op_num", sourceInfo, nil),
			"dbss_num":             utils.PathSearch("dbss_num", sourceInfo, nil),
		},
	}
}

func flattenDscAlarmOverviewLevelInfo(levelInfo interface{}) []interface{} {
	if levelInfo == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"fatal_num":       utils.PathSearch("fatal_num", levelInfo, nil),
			"high_risk_num":   utils.PathSearch("high_risk_num", levelInfo, nil),
			"middle_risk_num": utils.PathSearch("middle_risk_num", levelInfo, nil),
			"low_risk_num":    utils.PathSearch("low_risk_num", levelInfo, nil),
			"notice_num":      utils.PathSearch("notice_num", levelInfo, nil),
		},
	}
}
