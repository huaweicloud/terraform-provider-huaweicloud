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

// @API DSC GET /v2/{project_id}/sec-ops/situation-dashboard/handling-statistics
func DataSourceDscAlarmHandlingTrend() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDscAlarmHandlingTrendRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"time_axis": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"total_alarm_variation": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"total_event_variation": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"untreated_alarm_variation": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"untreated_event_variation": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
		},
	}
}

func dataSourceDscAlarmHandlingTrendRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "dsc"
		httpUrl = "v2/{project_id}/sec-ops/situation-dashboard/handling-statistics"
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
		return diag.Errorf("error retrieving DSC alarm handling trend: %s", err)
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
		d.Set("time_axis", utils.PathSearch("time_axis", respBody, nil)),
		d.Set("total_alarm_variation", utils.PathSearch("total_alarm_variation", respBody, nil)),
		d.Set("total_event_variation", utils.PathSearch("total_event_variation", respBody, nil)),
		d.Set("untreated_alarm_variation", utils.PathSearch("untreated_alarm_variation", respBody, nil)),
		d.Set("untreated_event_variation", utils.PathSearch("untreated_event_variation", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
