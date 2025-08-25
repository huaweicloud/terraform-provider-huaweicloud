package rms

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CONFIG GET /v1/resource-manager/domains/{domain_id}/policy-states/summary
func DataSourcePolicyStatesSummary() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePolicyStatesSummaryRead,

		Schema: map[string]*schema.Schema{
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"results": {
				Type:     schema.TypeList,
				Elem:     policyStatesSummaryResult(),
				Computed: true,
			},
		},
	}
}

func policyStatesSummaryResult() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"resource_details": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     policyStatesSummaryResultDetails(),
			},
			"assignment_details": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     policyStatesSummaryResultDetails(),
			},
		},
	}
	return &sc
}

func policyStatesSummaryResultDetails() *schema.Resource {
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

func dataSourcePolicyStatesSummaryRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v1/resource-manager/domains/{domain_id}/policy-states/summary"
		product = "rms"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating Config client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{domain_id}", cfg.DomainID)

	getPath += buildPolicyStatesSummaryQueryParams(d)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving Config policy states summary, %s", err)
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
		d.Set("results", flattenPolicyStatesSummary(getRespBody)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildPolicyStatesSummaryQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("tags"); ok {
		for _, tag := range v.([]interface{}) {
			res = fmt.Sprintf("%s&tags=%v", res, tag)
		}
	}
	if res != "" {
		res = "?" + res[1:]
	}
	return res
}

func flattenPolicyStatesSummary(resp interface{}) []interface{} {
	curJson := utils.PathSearch("results", resp, nil)
	if curJson == nil {
		return nil
	}

	rst := []interface{}{
		map[string]interface{}{
			"resource_details":   flattenPolicyStatesSummaryDetails(utils.PathSearch("resource_details", curJson, nil)),
			"assignment_details": flattenPolicyStatesSummaryDetails(utils.PathSearch("assignment_details", curJson, nil)),
		},
	}
	return rst
}

func flattenPolicyStatesSummaryDetails(resp interface{}) []interface{} {
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
