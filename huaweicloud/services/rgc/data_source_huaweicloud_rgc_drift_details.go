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

// @API RGC GET /v1/governance/drift-details
func DataSourceDriftDetails() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDriftDetailsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"drift_details": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     driftDetailSchema(),
			},
		},
	}
}

func driftDetailSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"managed_account_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"drift_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"drift_target_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"drift_target_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"drift_message": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"parent_organizational_unit_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"solve_solution": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

	return &sc
}

func dataSourceDriftDetailsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var getDriftDetailsProduct = "rgc"
	getDriftDetailsClient, err := cfg.NewServiceClient(getDriftDetailsProduct, region)
	if err != nil {
		return diag.Errorf("error creating RGC client: %s", err)
	}

	getDriftDetailsRespBody, err := getDriftDetails(getDriftDetailsClient)

	if err != nil {
		return diag.Errorf("error retrieving RGC drift details: %s", err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	mErr = multierror.Append(mErr,
		d.Set("drift_details", utils.PathSearch("drift_details", getDriftDetailsRespBody, nil)),
		d.Set("region", region),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func getDriftDetails(client *golangsdk.ServiceClient) (interface{}, error) {
	var (
		getDriftDetailsHttpUrl = "v1/governance/drift-details"
	)
	getDriftDetailsHttpPath := client.Endpoint + getDriftDetailsHttpUrl

	getDriftDetailsHttpOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getDriftDetailsHttpResp, err := client.Request("GET", getDriftDetailsHttpPath, &getDriftDetailsHttpOpt)
	if err != nil {
		return nil, err
	}
	getDriftDetailsRespBody, err := utils.FlattenResponse(getDriftDetailsHttpResp)
	if err != nil {
		return nil, err
	}
	return getDriftDetailsRespBody, nil
}
