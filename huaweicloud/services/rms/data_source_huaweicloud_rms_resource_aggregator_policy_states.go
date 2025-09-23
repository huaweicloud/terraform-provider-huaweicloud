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

// @API CONFIG POST /v1/resource-manager/domains/{domain_id}/aggregators/aggregate-data/policy-states/compliance-details
func DataSourceAggregatorPolicyStates() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAggregatorPolicyStatesRead,

		Schema: map[string]*schema.Schema{
			"aggregator_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the resource aggregator ID.`,
			},
			"account_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"policy_assignment_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"compliance_state": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"states": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_provider": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"trigger_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"compliance_state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_assignment_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_assignment_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_definition_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"evaluation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAggregatorPolicyStatesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("rms", region)
	if err != nil {
		return diag.Errorf("error creating RMS v1 client: %s", err)
	}

	getAggregatorPolicyStatesHttpUrl := "v1/resource-manager/domains/{domain_id}/aggregators/aggregate-data/policy-states/compliance-details"
	getAggregatorPolicyStatesPath := client.Endpoint + getAggregatorPolicyStatesHttpUrl
	getAggregatorPolicyStatesPath = strings.ReplaceAll(getAggregatorPolicyStatesPath, "{domain_id}", cfg.DomainID)

	states, err := getAggregatorPolicyStates(client, d, getAggregatorPolicyStatesPath)
	if err != nil {
		return diag.FromErr(err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr := multierror.Append(
		nil,
		d.Set("states", states),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildAggregatorPolicyStatesQueryParams(marker string) string {
	res := "?limit=200"

	if marker != "" {
		res += fmt.Sprintf("&marker=%v", marker)
	}

	return res
}

func getAggregatorPolicyStates(client *golangsdk.ServiceClient, d *schema.ResourceData,
	getAggregatorPolicyStatesPath string) ([]interface{}, error) {
	var resources []interface{}
	var marker string
	getAggregatorPolicyStatesOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildAggregatorPolicyStatesBodyParams(d)),
	}
	for {
		requestPath := getAggregatorPolicyStatesPath + buildAggregatorPolicyStatesQueryParams(marker)
		resp, err := client.Request("POST", requestPath, &getAggregatorPolicyStatesOpt)
		if err != nil {
			return nil, fmt.Errorf("error retrieving RMS aggregato policy states: %s", err)
		}

		getTrackedAggregatorPolicyStatesRespBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		resourcesTemp := flattenAggregatorPolicyStates(
			utils.PathSearch("policy_states", getTrackedAggregatorPolicyStatesRespBody, nil))
		if err != nil {
			return nil, err
		}
		resources = append(resources, resourcesTemp...)
		marker = utils.PathSearch("page_info.next_marker", getTrackedAggregatorPolicyStatesRespBody, "").(string)
		if marker == "" {
			break
		}
	}

	return resources, nil
}

func buildAggregatorPolicyStatesBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"aggregator_id":          d.Get("aggregator_id"),
		"account_id":             utils.ValueIgnoreEmpty(d.Get("account_id")),
		"policy_assignment_name": utils.ValueIgnoreEmpty(d.Get("policy_assignment_name")),
		"compliance_state":       utils.ValueIgnoreEmpty(d.Get("compliance_state")),
		"resource_name":          utils.ValueIgnoreEmpty(d.Get("resource_name")),
		"resource_id":            utils.ValueIgnoreEmpty(d.Get("resource_id")),
	}
	return bodyParams
}

func flattenAggregatorPolicyStates(resourcesAggregatorPolicyStatesRaw interface{}) []interface{} {
	if resourcesAggregatorPolicyStatesRaw == nil {
		return nil
	}

	resourcesAggregatorPolicyStates := resourcesAggregatorPolicyStatesRaw.([]interface{})
	res := make([]interface{}, len(resourcesAggregatorPolicyStates))
	for i, v := range resourcesAggregatorPolicyStates {
		policyState := v.(map[string]interface{})
		res[i] = map[string]interface{}{
			"domain_id":              policyState["domain_id"],
			"region_id":              policyState["region_id"],
			"resource_id":            policyState["resource_id"],
			"resource_name":          policyState["resource_name"],
			"resource_provider":      policyState["resource_provider"],
			"resource_type":          policyState["resource_type"],
			"trigger_type":           policyState["trigger_type"],
			"compliance_state":       policyState["compliance_state"],
			"policy_assignment_id":   policyState["policy_assignment_id"],
			"policy_assignment_name": policyState["policy_assignment_name"],
			"policy_definition_id":   policyState["policy_definition_id"],
			"evaluation_time":        policyState["evaluation_time"],
		}
	}

	return res
}
