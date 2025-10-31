package rgc

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API RGC GET /v1/landing-zone/available-updates
func DataSourceLandingZoneAvailableUpdates() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLandingZoneAvailableUpdatesRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"baseline_update_available": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"control_update_available": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"landing_zone_update_available": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"service_landing_zone_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_landing_zone_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceLandingZoneAvailableUpdatesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var getLandingZoneAvailableUpdatesProduct = "rgc"
	getLandingZoneAvailableUpdatesClient, err := cfg.NewServiceClient(getLandingZoneAvailableUpdatesProduct, region)
	if err != nil {
		return diag.Errorf("Error creating RGC client: %s", err)
	}

	getLandingZoneAvailableUpdatesRespBody, err := getLandingZoneAvailableUpdates(getLandingZoneAvailableUpdatesClient)

	if err != nil {
		return diag.Errorf("error retrieving RGC landing zone available updates: %s", err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	mErr = multierror.Append(mErr,
		d.Set("baseline_update_available", utils.PathSearch("baseline_update_available", getLandingZoneAvailableUpdatesRespBody, nil)),
		d.Set("control_update_available", utils.PathSearch("control_update_available", getLandingZoneAvailableUpdatesRespBody, nil)),
		d.Set("landing_zone_update_available", utils.PathSearch("landing_zone_update_available", getLandingZoneAvailableUpdatesRespBody, nil)),
		d.Set("service_landing_zone_version", utils.PathSearch("service_landing_zone_version", getLandingZoneAvailableUpdatesRespBody, nil)),
		d.Set("user_landing_zone_version", utils.PathSearch("user_landing_zone_version", getLandingZoneAvailableUpdatesRespBody, nil)),
		d.Set("region", region),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func getLandingZoneAvailableUpdates(client *golangsdk.ServiceClient) (interface{}, error) {
	var (
		getLandingZoneAvailableUpdatesHttpUrl = "v1/landing-zone/available-updates"
	)
	getLandingZoneAvailableUpdatesHttpPath := client.Endpoint + getLandingZoneAvailableUpdatesHttpUrl

	getLandingZoneAvailableUpdatesHttpOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getLandingZoneAvailableUpdatesHttpResp, err := client.Request("GET", getLandingZoneAvailableUpdatesHttpPath,
		&getLandingZoneAvailableUpdatesHttpOpt)
	if err != nil {
		return nil, err
	}
	getLandingZoneAvailableUpdatesRespBody, err := utils.FlattenResponse(getLandingZoneAvailableUpdatesHttpResp)
	if err != nil {
		return nil, err
	}
	return getLandingZoneAvailableUpdatesRespBody, nil
}
