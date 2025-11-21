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

// @API RGC GET /v1/best-practice/detection-details
func DataSourceBestPracticeDetails() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBestPracticeDetailsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"details": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     BestPracticeDetailsSchema(),
			},
		},
	}
}

func BestPracticeDetailsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"check_item": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"check_item_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"risk_description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"result": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"scene": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"risk_level": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

	return &sc
}

func dataSourceBestPracticeDetailsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var getBestPracticeDetailsProduct = "rgc"
	getBestPracticeDetailsClient, err := cfg.NewServiceClient(getBestPracticeDetailsProduct, region)
	if err != nil {
		return diag.Errorf("error creating RGC client: %s", err)
	}

	getBestPracticeDetailsRespBody, err := getBestPracticeDetails(getBestPracticeDetailsClient)
	if err != nil {
		return diag.Errorf("error retrieving RGC best practice details: %s", err)
	}

	bestPracticeDetails := getBestPracticeDetailsRespBody.([]interface{})

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	mErr = multierror.Append(mErr,
		d.Set("details", bestPracticeDetails),
		d.Set("region", region),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func getBestPracticeDetails(client *golangsdk.ServiceClient) (interface{}, error) {
	var (
		getBestPracticeDetailsHttpUrl = "v1/best-practice/detection-details"
	)
	getBestPracticeDetailsHttpPath := client.Endpoint + getBestPracticeDetailsHttpUrl

	getBestPracticeDetailsHttpOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getBestPracticeDetailsHttpResp, err := client.Request("GET", getBestPracticeDetailsHttpPath, &getBestPracticeDetailsHttpOpt)
	if err != nil {
		return nil, err
	}
	getBestPracticeDetailsRespBody, err := utils.FlattenResponse(getBestPracticeDetailsHttpResp)
	if err != nil {
		return nil, err
	}
	return getBestPracticeDetailsRespBody, nil
}
