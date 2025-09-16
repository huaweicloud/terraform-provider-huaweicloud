package waf

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

// @API WAF GET /v1/{project_id}/waf/policy/{policy_id}/ip-reputation
func DataSourceRulesThreatIntelligence() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRulesThreatIntelligenceRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"policy_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"items": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policyid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"status": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"action": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"category": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func buildThreatIntelligenceQueryParams(cfg *config.Config, d *schema.ResourceData) string {
	epsId := cfg.GetEnterpriseProjectID(d)
	rst := "?limit=200"

	if epsId != "" {
		rst += fmt.Sprintf("&enterprise_project_id=%v", epsId)
	}

	return rst
}

func dataSourceRulesThreatIntelligenceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		policyId = d.Get("policy_id").(string)
		httpUrl  = "v1/{project_id}/waf/policy/{policy_id}/ip-reputation"
		offset   = 0
		result   = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient("waf", region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{policy_id}", policyId)
	listPath += buildThreatIntelligenceQueryParams(cfg, d)

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", listPath, offset)
		resp, err := client.Request("GET", currentPath, &listOpt)
		if err != nil {
			return diag.Errorf("error retrieving WAF threat intelligence rules: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		itemsResp := utils.PathSearch("items", respBody, make([]interface{}, 0)).([]interface{})
		if len(itemsResp) == 0 {
			break
		}

		result = append(result, itemsResp...)
		offset += len(itemsResp)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("items", flattenThreatIntelligenceRules(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenThreatIntelligenceRules(rawRules []interface{}) []interface{} {
	if len(rawRules) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(rawRules))
	for _, v := range rawRules {
		rst = append(rst, map[string]interface{}{
			"name":        utils.PathSearch("name", v, nil),
			"id":          utils.PathSearch("id", v, nil),
			"policyid":    utils.PathSearch("policyid", v, nil),
			"type":        utils.PathSearch("type", v, nil),
			"description": utils.PathSearch("description", v, nil),
			"tags":        utils.ExpandToStringList(utils.PathSearch("tags", v, make([]interface{}, 0)).([]interface{})),
			"status":      utils.PathSearch("status", v, nil),
			"action":      flattenThreatIntelligenceAction(utils.PathSearch("action", v, nil)),
		})
	}

	return rst
}

func flattenThreatIntelligenceAction(rawAction interface{}) []interface{} {
	if rawAction == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"category": utils.PathSearch("category", rawAction, nil),
		},
	}
}
