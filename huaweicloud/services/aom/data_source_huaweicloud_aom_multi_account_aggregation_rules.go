package aom

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

// @API AOM GET /v1/{project_id}/aom/aggr-config
func DataSourceMultiAccountAggregationRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMultiAccountAggregationRulesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"accounts": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"urn": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"join_method": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"joined_at": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"services": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"service": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"metrics": {
										Type:     schema.TypeList,
										Computed: true,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"send_to_source_account": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildHeadersForDataSource(cfg *config.Config, d *schema.ResourceData) map[string]string {
	moreHeaders := map[string]string{
		"Content-Type":          "application/json",
		"Enterprise-Project-Id": cfg.GetEnterpriseProjectID(d, "all_granted_eps"),
	}
	return moreHeaders
}

func dataSourceMultiAccountAggregationRulesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("aom", region)
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	results, err := listMultiAccountAggregationRules(cfg, client, d)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("error generating UUID")
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("rules", flattenRules(results.([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func listMultiAccountAggregationRules(cfg *config.Config, client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	listHttpUrl := "v1/{project_id}/aom/aggr-config"
	listPath := client.Endpoint + listHttpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildHeadersForDataSource(cfg, d),
	}

	listResp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving multi account aggregation rule: %s", err)
	}
	listRespBody, err := utils.FlattenResponse(listResp)
	if err != nil {
		return nil, fmt.Errorf("error flattening multi account aggregation rule: %s", err)
	}

	return listRespBody, nil
}

func flattenRules(rules []interface{}) []interface{} {
	if len(rules) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(rules))
	for _, rule := range rules {
		result = append(result, map[string]interface{}{
			"instance_id":            utils.PathSearch("dest_prometheus_id", rule, nil),
			"accounts":               flattenMultiAccountAggregationRuleResponseBodyAccounts(rule),
			"services":               flattenMultiAccountAggregationRuleResponseBodyMetrics(rule),
			"send_to_source_account": utils.PathSearch("send_to_source_account", rule, nil),
		})
	}
	return result
}
