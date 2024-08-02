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

// @API CONFIG GET /v1/resource-manager/domains/{domain_id}/policy-states
// @API CONFIG GET /v1/resource-manager/domains/{domain_id}/policy-assignments/{policy_assignment_id}/policy-states
func DataSourcePolicyStates() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePolicyStatesRead,

		Schema: map[string]*schema.Schema{
			"policy_assignment_id": {
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

func dataSourcePolicyStatesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("rms", region)
	if err != nil {
		return diag.Errorf("error creating RMS v1 client: %s", err)
	}

	var policyStates []interface{}
	if v, ok := d.GetOk("policy_assignment_id"); ok {
		getPolicyStatesHttpUrl := "v1/resource-manager/domains/{domain_id}/policy-assignments/{policy_assignment_id}/policy-states"
		getPolicyStatesPath := client.Endpoint + getPolicyStatesHttpUrl
		getPolicyStatesPath = strings.ReplaceAll(getPolicyStatesPath, "{domain_id}", cfg.DomainID)
		getPolicyStatesPath = strings.ReplaceAll(getPolicyStatesPath, "{policy_assignment_id}", v.(string))

		policyStates, err = getPolicyStates(client, d, getPolicyStatesPath)
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		getPolicyStatesHttpUrl := "v1/resource-manager/domains/{domain_id}/policy-states"
		getPolicyStatesPath := client.Endpoint + getPolicyStatesHttpUrl
		getPolicyStatesPath = strings.ReplaceAll(getPolicyStatesPath, "{domain_id}", cfg.DomainID)

		policyStates, err = getPolicyStates(client, d, getPolicyStatesPath)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr := multierror.Append(
		nil,
		d.Set("states", policyStates),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildPolicyStatesQueryParams(d *schema.ResourceData, marker string) string {
	res := ""

	if v, ok := d.GetOk("compliance_state"); ok {
		res = fmt.Sprintf("%s&compliance_state=%v", res, v)
	}
	if v, ok := d.GetOk("resource_name"); ok {
		res = fmt.Sprintf("%s&resource_name=%v", res, v)
	}
	if v, ok := d.GetOk("resource_id"); ok {
		res = fmt.Sprintf("%s&resource_id=%v", res, v)
	}

	if marker != "" {
		res = fmt.Sprintf("%s&marker=%v", res, marker)
	}

	if res != "" {
		res = "?" + res[1:]
	}

	return res
}

func getPolicyStates(client *golangsdk.ServiceClient, d *schema.ResourceData, getPolicyStatesPath string) ([]interface{}, error) {
	var policyStates []interface{}
	var marker string
	getPolicyStatesOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	for {
		requestPath := getPolicyStatesPath + buildPolicyStatesQueryParams(d, marker)
		resp, err := client.Request("GET", requestPath, &getPolicyStatesOpt)
		if err != nil {
			return nil, fmt.Errorf("error retrieving RMS policy states: %s", err)
		}

		getPolicyStatesRespBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		policyStatesTemp, err := flattenPolicyStates(utils.PathSearch("value", getPolicyStatesRespBody, nil))
		if err != nil {
			return nil, err
		}
		policyStates = append(policyStates, policyStatesTemp...)
		marker = utils.PathSearch("page_info.next_marker", getPolicyStatesRespBody, "").(string)
		if marker == "" {
			break
		}
	}

	return policyStates, nil
}

func flattenPolicyStates(policyStatesRaw interface{}) ([]interface{}, error) {
	if policyStatesRaw == nil {
		return nil, nil
	}

	policyStates := policyStatesRaw.([]interface{})
	res := make([]interface{}, len(policyStates))
	for i, v := range policyStates {
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

	return res, nil
}
