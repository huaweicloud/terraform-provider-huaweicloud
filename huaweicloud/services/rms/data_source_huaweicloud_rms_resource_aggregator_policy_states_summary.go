package rms

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CONFIG POST /v1/resource-manager/domains/{domain_id}/aggregators/aggregate-data/policy-states/compliance-summary
func DataSourceAggregatorPolicyStatesSummary() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAggregatorPolicyStatesSummaryRead,

		Schema: map[string]*schema.Schema{
			"aggregator_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"account_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"group_by_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"results": {
				Type:     schema.TypeList,
				Elem:     aggregatorPolicyStatesSummaryResult(),
				Computed: true,
			},
		},
	}
}

func aggregatorPolicyStatesSummaryResult() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"resource_details": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     aggregatorPolicyStatesSummaryResultDetails(),
			},
			"assignment_details": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     aggregatorPolicyStatesSummaryResultDetails(),
			},
			"group_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"group_account_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func aggregatorPolicyStatesSummaryResultDetails() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"compliant_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"non_compliant_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
	return &sc
}

func dataSourceAggregatorPolicyStatesSummaryRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v1/resource-manager/domains/{domain_id}/aggregators/aggregate-data/policy-states/compliance-summary"
		product = "rms"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating Config client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{domain_id}", cfg.DomainID)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getOpt.JSONBody = utils.RemoveNil(buildAggregatorPolicyStatesSummaryQueryParams(d))

	getResp, err := client.Request("POST", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving Config resource aggregator policy states summary, %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		mErr,
		d.Set("results", flattenAggregatorPolicyStatesSummary(getRespBody)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildAggregatorPolicyStatesSummaryQueryParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"aggregator_id": d.Get("aggregator_id"),
		"account_id":    utils.ValueIgnoreEmpty(d.Get("account_id")),
		"group_by_key":  utils.ValueIgnoreEmpty(d.Get("group_by_key")),
	}
	return bodyParams
}

func flattenAggregatorPolicyStatesSummary(resp interface{}) []interface{} {
	curJson := utils.PathSearch("results", resp, nil)
	if curJson == nil {
		return nil
	}
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"group_name":         utils.PathSearch("group_name", v, nil),
			"group_account_name": utils.PathSearch("group_account_name", v, nil),
			"resource_details":   flattenAggregatorPolicyStatesSummaryDetails(utils.PathSearch("resource_details", v, nil)),
			"assignment_details": flattenAggregatorPolicyStatesSummaryDetails(utils.PathSearch("assignment_details", v, nil)),
		})
	}
	return rst
}

func flattenAggregatorPolicyStatesSummaryDetails(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	rst := []interface{}{
		map[string]interface{}{
			"compliant_count":     utils.PathSearch("compliant_count", resp, nil),
			"non_compliant_count": utils.PathSearch("non_compliant_count", resp, nil),
		},
	}
	return rst
}
