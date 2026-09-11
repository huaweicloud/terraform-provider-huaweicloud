package waf

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API WAF GET /v1/{project_id}/waf/policy/{rule_type}/rules
func DataSourceRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRulesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the rules.`,
			},
			"rule_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the type of the rules to be queried.`,
			},
			"policy_ids": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the ID list of the policies to be queried.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the enterprise project ID to which the rules belong.`,
			},
			"rules": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `All rules that matched the filter parameters.`,
			},
		},
	}
}

func buildListRulesQueryParams(policyIds, epsId string) string {
	res := "?pagesize=1000"
	if policyIds != "" {
		res = fmt.Sprintf("%s&policyids=%v", res, policyIds)
	}
	if epsId != "" {
		res = fmt.Sprintf("%s&enterprise_project_id=%v", res, epsId)
	}

	return res
}

func listRules(client *golangsdk.ServiceClient, ruleType, policyIds, epsId string) ([]interface{}, error) {
	var (
		allRules    []interface{}
		currentPage = 1
	)

	httpUrl := "v1/{project_id}/waf/policy/{rule_type}/rules"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{rule_type}", ruleType)
	listPath += buildListRulesQueryParams(policyIds, epsId)
	requestOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		KeepResponseBody: true,
	}

	for {
		requestPathWithPage := fmt.Sprintf("%s&page=%d", listPath, currentPage)
		resp, err := client.Request("GET", requestPathWithPage, &requestOpt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		rulesResp := utils.PathSearch("items", respBody, make([]interface{}, 0)).([]interface{})
		if len(rulesResp) == 0 {
			break
		}

		allRules = append(allRules, rulesResp...)
		currentPage++
	}

	return allRules, nil
}

func dataSourceRulesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		epsId     = cfg.GetEnterpriseProjectID(d)
		ruleType  = d.Get("rule_type").(string)
		policyIds = d.Get("policy_ids").(string)
	)

	client, err := cfg.NewServiceClient("waf", region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	allRules, err := listRules(client, ruleType, policyIds, epsId)
	if err != nil {
		return diag.Errorf("error retrieving WAF protection rules: %s", err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("rules", flattenRules(allRules)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenRules(rules []interface{}) []interface{} {
	rst := make([]interface{}, 0, len(rules))
	for _, v := range rules {
		rst = append(rst, utils.JsonToString(v))
	}

	return rst
}
