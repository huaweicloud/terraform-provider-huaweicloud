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

// @API RGC GET /v1/landing-zone/home-region
func DataSourceHomeRegion() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceHomeRegionRead,
		Schema: map[string]*schema.Schema{
			"home_region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the home region`,
			},
		},
	}
}

func dataSourceHomeRegionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var getHomeRegionProduct = "rgc"
	getHomeRegionClient, err := cfg.NewServiceClient(getHomeRegionProduct, region)
	if err != nil {
		return diag.Errorf("error creating RGC client: %s", err)
	}

	getHomeRegionRespBody, err := getHomeRegion(getHomeRegionClient)

	if err != nil {
		return diag.FromErr(err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	mErr = multierror.Append(mErr,
		d.Set("home_region", utils.PathSearch("home_region", getHomeRegionRespBody, nil)))

	return diag.FromErr(mErr.ErrorOrNil())
}

func getHomeRegion(client *golangsdk.ServiceClient) (interface{}, error) {
	var (
		getHomeRegionHttpUrl = "v1/landing-zone/home-region"
	)
	getHomeRegionHttpPath := client.Endpoint + getHomeRegionHttpUrl

	getHomeRegionHttpOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getHomeRegionHttpResp, err := client.Request("GET", getHomeRegionHttpPath, &getHomeRegionHttpOpt)
	if err != nil {
		return nil, err
	}
	getHomeRegionRespBody, err := utils.FlattenResponse(getHomeRegionHttpResp)
	if err != nil {
		return nil, err
	}
	return getHomeRegionRespBody, nil
}
