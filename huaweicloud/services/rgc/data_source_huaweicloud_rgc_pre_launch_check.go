package rgc

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API RGC POST /v1/landing-zone/pre-launch-check
func DataSourcePreLaunchCheck() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePreLaunchCheckRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"pre_launch_check": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourcePreLaunchCheckRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var getPreLaunchCheckProduct = "rgc"
	getPreLaunchCheckClient, err := cfg.NewServiceClient(getPreLaunchCheckProduct, region)
	if err != nil {
		return diag.Errorf("error creating RGC client: %s", err)
	}

	err = getPreLaunchCheck(getPreLaunchCheckClient)

	if err != nil {
		return diag.FromErr(err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	mErr = multierror.Append(mErr,
		d.Set("region", region),
		d.Set("pre_launch_check", true),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func getPreLaunchCheck(client *golangsdk.ServiceClient) error {
	var (
		getPreLaunchCheckHttpUrl = "v1/landing-zone/pre-launch-check"
	)
	getPreLaunchCheckHttpPath := client.Endpoint + getPreLaunchCheckHttpUrl

	getPreLaunchCheckHttpOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes:          []int{200},
	}
	_, err := client.Request("POST", getPreLaunchCheckHttpPath, &getPreLaunchCheckHttpOpt)
	if err != nil {
		return err
	}
	return nil
}
