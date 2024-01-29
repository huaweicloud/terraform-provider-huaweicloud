package rms

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/rms/v1/policyassignments"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Config GET /v1/resource-manager/policy-definitions
func DataSourcePolicyDefinitions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePolicyDefinitionsRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the policy definitions used to query definition list.",
			},
			"policy_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The policy type used to query definition list.",
			},
			"policy_rule_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The policy rule type used to query definition list.",
			},
			"trigger_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The trigger type used to query definition list.",
			},
			"keywords": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The keyword list used to query definition list.",
			},
			"definitions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the policy definition.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the policy definition.",
						},
						"policy_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The policy type of the policy definition.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of the policy definition.",
						},
						"policy_rule_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The policy rule type of the policy definition.",
						},
						"policy_rule": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The policy rule of the policy definition.",
						},
						"trigger_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The trigger type of the policy definition.",
						},
						"keywords": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "The keyword list of the policy definition.",
						},
						"parameters": {
							Type:        schema.TypeMap,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "The parameter reference map of the policy definition.",
						},
					},
				},
				Description: "The policy definition list.",
			},
		},
	}
}

func filterPolicyDefinitionsByKeywords(definitions []policyassignments.PolicyDefinition,
	keywords []interface{}) []policyassignments.PolicyDefinition {
	if len(keywords) < 1 {
		return definitions
	}

	filter := utils.ExpandToStringList(keywords)
	result := make([]policyassignments.PolicyDefinition, 0, len(definitions))
	for _, v := range definitions {
		if utils.StrSliceContainsAnother(v.Keywords, filter) {
			result = append(result, v)
		}
	}
	return result
}

func flattenDefinitionParameters(parameters map[string]policyassignments.PolicyParameterDefinition) (
	map[string]interface{}, error) {
	if len(parameters) < 1 {
		return nil, nil
	}

	result := make(map[string]interface{})
	for k, v := range parameters {
		jsonBytes, err := json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("generate json string failed: %s", err)
		}
		result[k] = string(jsonBytes)
	}
	return result, nil
}

func filterPolicyDefinitions(definitions []policyassignments.PolicyDefinition,
	d *schema.ResourceData) ([]map[string]interface{}, []string, error) {
	filter := map[string]interface{}{
		"Name":           d.Get("name"),
		"PolicyType":     d.Get("policy_type"),
		"PolicyRuleType": d.Get("policy_rule_type"),
		"TriggerType":    d.Get("trigger_type"),
	}
	filtResult, err := utils.FilterSliceWithField(definitions, filter)
	if err != nil {
		return nil, nil, fmt.Errorf("filter component runtimes failed: %s", err)
	}
	log.Printf("[DEBUG] Filter %d policy definitions from server through options: %v", len(filtResult), filter)

	result := make([]map[string]interface{}, len(filtResult))
	ids := make([]string, len(filtResult))
	for i, val := range filtResult {
		definition := val.(policyassignments.PolicyDefinition)
		ids[i] = definition.ID
		dm := map[string]interface{}{
			"id":               definition.ID,
			"name":             definition.Name,
			"policy_type":      definition.PolicyType,
			"description":      definition.Description,
			"policy_rule_type": definition.PolicyRuleType,
			"policy_rule":      definition.PolicyRule,
			"trigger_type":     definition.TriggerType,
			"keywords":         definition.Keywords,
		}

		params, err := flattenDefinitionParameters(definition.Parameters)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to flatten definition parameters: %s", err)
		}
		dm["parameters"] = params

		jsonBytes, err := json.Marshal(definition.PolicyRule)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to generate json string: %s", err)
		}
		dm["policy_rule"] = string(jsonBytes)

		result[i] = dm
	}
	return result, ids, nil
}

func dataSourcePolicyDefinitionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.RmsV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating RMS v1 client: %s", err)
	}

	definitions, err := policyassignments.ListDefinitions(client)
	if err != nil {
		return diag.Errorf("error getting the policy definition list form server: %s", err)
	}

	filterResult := filterPolicyDefinitionsByKeywords(definitions, d.Get("keywords").([]interface{}))
	dm, ids, err := filterPolicyDefinitions(filterResult, d)
	if err != nil {
		return diag.Errorf("error query policy definitions: %s", err)
	}
	d.SetId(hashcode.Strings(ids))

	if err = d.Set("definitions", dm); err != nil {
		return diag.Errorf("error saving the information of the policy definitions to state: %s", err)
	}
	return nil
}
