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

// @API RGC GET /v1/best-practice/detection-overview
func DataSourceBestPracticeOverview() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBestPracticeOverviewRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"total_score": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"detect_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"organization_account": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     bestPracticeOverviewItemSchema(),
			},
			"identity_permission": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     bestPracticeOverviewItemSchema(),
			},
			"network_planning": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     bestPracticeOverviewItemSchema(),
			},
			"compliance_audit": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     bestPracticeOverviewItemSchema(),
			},
			"financial_governance": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     bestPracticeOverviewItemSchema(),
			},
			"data_boundary": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     bestPracticeOverviewItemSchema(),
			},
			"om_monitor": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     bestPracticeOverviewItemSchema(),
			},
			"security_management": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     bestPracticeOverviewItemSchema(),
			},
		},
	}
}

func bestPracticeOverviewItemSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"score": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"detection_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"high_risk_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"medium_risk_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"low_risk_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"risk_item_description": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}

	return &sc
}

func dataSourceBestPracticeOverviewRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var getBestPracticeOverviewProduct = "rgc"
	getBestPracticeOverviewClient, err := cfg.NewServiceClient(getBestPracticeOverviewProduct, region)
	if err != nil {
		return diag.Errorf("error creating RGC client: %s", err)
	}

	getBestPracticeOverviewRespBody, err := getBestPracticeOverview(getBestPracticeOverviewClient)

	if err != nil {
		return diag.Errorf("error retrieving RGC best-practice overview: %s", err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	mErr = multierror.Append(mErr,
		d.Set("total_score", utils.PathSearch("total_score", getBestPracticeOverviewRespBody, nil)),
		d.Set("detect_time", utils.PathSearch("detect_time", getBestPracticeOverviewRespBody, nil)),
		d.Set("organization_account", parseBestPracticeOverviewItem(getBestPracticeOverviewRespBody, "organization_account")),
		d.Set("identity_permission", parseBestPracticeOverviewItem(getBestPracticeOverviewRespBody, "identity_permission")),
		d.Set("network_planning", parseBestPracticeOverviewItem(getBestPracticeOverviewRespBody, "network_planning")),
		d.Set("compliance_audit", parseBestPracticeOverviewItem(getBestPracticeOverviewRespBody, "compliance_audit")),
		d.Set("financial_governance", parseBestPracticeOverviewItem(getBestPracticeOverviewRespBody, "financial_governance")),
		d.Set("data_boundary", parseBestPracticeOverviewItem(getBestPracticeOverviewRespBody, "data_boundary")),
		d.Set("om_monitor", parseBestPracticeOverviewItem(getBestPracticeOverviewRespBody, "om_monitor")),
		d.Set("security_management", parseBestPracticeOverviewItem(getBestPracticeOverviewRespBody, "security_management")),
		d.Set("region", region),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func parseBestPracticeOverviewItem(respBody interface{}, key string) []interface{} {
	overviewItemList := make([]interface{}, 0)

	item := utils.PathSearch(key, respBody, nil)
	if item != nil {
		overviewItemList = append(overviewItemList, item)
	}
	return overviewItemList
}

func getBestPracticeOverview(client *golangsdk.ServiceClient) (interface{}, error) {
	var (
		getBestPracticeOverviewHttpUrl = "v1/best-practice/detection-overview"
	)
	getBestPracticeOverviewHttpPath := client.Endpoint + getBestPracticeOverviewHttpUrl

	getBestPracticeOverviewHttpOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getBestPracticeOverviewHttpResp, err := client.Request("GET", getBestPracticeOverviewHttpPath, &getBestPracticeOverviewHttpOpt)
	if err != nil {
		return nil, err
	}
	getBestPracticeOverviewRespBody, err := utils.FlattenResponse(getBestPracticeOverviewHttpResp)
	if err != nil {
		return nil, err
	}
	return getBestPracticeOverviewRespBody, nil
}
