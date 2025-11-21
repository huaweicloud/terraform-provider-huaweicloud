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

// @API RGC GET /v1/best-practice/account-info
func DataSourceBestPracticeAccountInfo() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBestPracticeAccountInfoRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"account_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"effective_start_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"effective_expiration_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"current_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceBestPracticeAccountInfoRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var getBestPracticeAccountInfoProduct = "rgc"
	getBestPracticeAccountInfoClient, err := cfg.NewServiceClient(getBestPracticeAccountInfoProduct, region)
	if err != nil {
		return diag.Errorf("error creating RGC client: %s", err)
	}

	getBestPracticeAccountInfoRespBody, err := getBestPracticeAccountInfo(getBestPracticeAccountInfoClient)
	if err != nil {
		return diag.Errorf("error retrieving RGC best practice account information: %s", err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	mErr = multierror.Append(mErr,
		d.Set("account_type", utils.PathSearch("account_type", getBestPracticeAccountInfoRespBody, nil)),
		d.Set("effective_start_time",
			utils.PathSearch("effective_start_time", getBestPracticeAccountInfoRespBody, nil)),
		d.Set("effective_expiration_time",
			utils.PathSearch("effective_expiration_time", getBestPracticeAccountInfoRespBody, nil)),
		d.Set("current_time", utils.PathSearch("current_time", getBestPracticeAccountInfoRespBody, nil)),
		d.Set("region", region),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func getBestPracticeAccountInfo(client *golangsdk.ServiceClient) (interface{}, error) {
	var (
		getBestPracticeAccountInfoHttpUrl = "v1/best-practice/account-info"
	)
	getBestPracticeAccountInfoHttpPath := client.Endpoint + getBestPracticeAccountInfoHttpUrl

	getBestPracticeAccountInfoHttpOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getBestPracticeAccountInfoHttpResp, err := client.Request("GET", getBestPracticeAccountInfoHttpPath,
		&getBestPracticeAccountInfoHttpOpt)
	if err != nil {
		return nil, err
	}
	getBestPracticeAccountInfoRespBody, err := utils.FlattenResponse(getBestPracticeAccountInfoHttpResp)
	if err != nil {
		return nil, err
	}
	return getBestPracticeAccountInfoRespBody, nil
}
