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

// @API RGC GET /v1/landing-zone/identity-center
func DataSourceLandingZoneIdentityCenter() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLandingZoneIdentityCenterRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"identity_store_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_portal_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"permission_sets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"permission_set_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"permission_set_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceLandingZoneIdentityCenterRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var getLandingZoneIdentityCenterProduct = "rgc"
	getLandingZoneIdentityCenterClient, err := cfg.NewServiceClient(getLandingZoneIdentityCenterProduct, region)
	if err != nil {
		return diag.Errorf("error creating RGC client: %s", err)
	}

	getLandingZoneIdentityCenterRespBody, err := getLandingZoneIdentityCenter(getLandingZoneIdentityCenterClient)

	if err != nil {
		return diag.Errorf("error retrieving RGC landing zone identity center: %s", err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	mErr = multierror.Append(mErr,
		d.Set("identity_store_id", utils.PathSearch("identity_store_id", getLandingZoneIdentityCenterRespBody, nil)),
		d.Set("user_portal_url", utils.PathSearch("user_portal_url", getLandingZoneIdentityCenterRespBody, nil)),
		d.Set("permission_sets", utils.PathSearch("permission_sets", getLandingZoneIdentityCenterRespBody, nil)),
		d.Set("groups", utils.PathSearch("groups", getLandingZoneIdentityCenterRespBody, nil)),
		d.Set("region", region),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func getLandingZoneIdentityCenter(client *golangsdk.ServiceClient) (interface{}, error) {
	var (
		getLandingZoneIdentityCenterHttpUrl = "v1/landing-zone/identity-center"
	)
	getLandingZoneIdentityCenterHttpPath := client.Endpoint + getLandingZoneIdentityCenterHttpUrl

	getLandingZoneIdentityCenterHttpOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getLandingZoneIdentityCenterHttpResp, err := client.Request("GET", getLandingZoneIdentityCenterHttpPath, &getLandingZoneIdentityCenterHttpOpt)
	if err != nil {
		return nil, err
	}
	getLandingZoneIdentityCenterRespBody, err := utils.FlattenResponse(getLandingZoneIdentityCenterHttpResp)
	if err != nil {
		return nil, err
	}
	return getLandingZoneIdentityCenterRespBody, nil
}
