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

// @API CONFIG POST /v1/resource-manager/domains/{domain_id}/aggregators/aggregate-data/policy-assignments/compliance
func DataSourceAggregatorPolicyAssignments() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAggregatorPolicyAssignmentsRead,

		Schema: map[string]*schema.Schema{
			"aggregator_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the resource aggregator ID.`,
			},
			"filter": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"account_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"compliance_state": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"policy_assignment_name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"assignments": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"policy_assignment_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_assignment_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"account_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"account_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"compliance": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"compliance_state": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"resource_details": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
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
										},
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

func dataSourceAggregatorPolicyAssignmentsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("rms", region)
	if err != nil {
		return diag.Errorf("error creating RMS v1 client: %s", err)
	}

	getAggregatorPolicyAssignmentsHttpUrl := "v1/resource-manager/domains/{domain_id}/aggregators/aggregate-data/policy-assignments/compliance"
	getAggregatorPolicyAssignmentsPath := client.Endpoint + getAggregatorPolicyAssignmentsHttpUrl
	getAggregatorPolicyAssignmentsPath = strings.ReplaceAll(getAggregatorPolicyAssignmentsPath, "{domain_id}", cfg.DomainID)

	assignments, err := getAggregatorPolicyAssignments(client, d, getAggregatorPolicyAssignmentsPath)
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
		d.Set("assignments", assignments),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildAggregatorPolicyAssignmentsQueryParams(marker string) string {
	res := "?limit=200"

	if marker != "" {
		res += fmt.Sprintf("&marker=%v", marker)
	}

	return res
}

func getAggregatorPolicyAssignments(client *golangsdk.ServiceClient, d *schema.ResourceData,
	getAggregatorPolicyAssignmentsPath string) ([]interface{}, error) {
	var assignments []interface{}
	var marker string
	getAggregatorPolicyAssignmentsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildAggregatorPolicyAssignmentsBodyParams(d)),
	}
	for {
		requestPath := getAggregatorPolicyAssignmentsPath + buildAggregatorPolicyAssignmentsQueryParams(marker)
		resp, err := client.Request("POST", requestPath, &getAggregatorPolicyAssignmentsOpt)
		if err != nil {
			return nil, fmt.Errorf("error retrieving aggregator policy assignments: %s", err)
		}

		getAggregatorPolicyAssignmentsRespBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		assignmentsTemp := flattenAggregatorPolicyAssignments(
			utils.PathSearch("aggregate_policy_assignments", getAggregatorPolicyAssignmentsRespBody, nil))
		if err != nil {
			return nil, err
		}
		assignments = append(assignments, assignmentsTemp...)
		marker = utils.PathSearch("page_info.next_marker", getAggregatorPolicyAssignmentsRespBody, "").(string)
		if marker == "" {
			break
		}
	}

	return assignments, nil
}

func buildAggregatorPolicyAssignmentsBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"aggregator_id": d.Get("aggregator_id"),
		"filter":        buildAggregatorPolicyAssignmentsFilterBodyParams(d),
	}
	return bodyParams
}

func buildAggregatorPolicyAssignmentsFilterBodyParams(d *schema.ResourceData) map[string]interface{} {
	filterRaw := d.Get("filter").([]interface{})
	if len(filterRaw) == 0 {
		return nil
	}

	if filter, ok := filterRaw[0].(map[string]interface{}); ok {
		bodyParams := map[string]interface{}{
			"account_id":             utils.ValueIgnoreEmpty(filter["account_id"]),
			"compliance_state":       utils.ValueIgnoreEmpty(filter["compliance_state"]),
			"policy_assignment_name": utils.ValueIgnoreEmpty(filter["policy_assignment_name"]),
		}

		return bodyParams
	}
	return nil
}

func flattenAggregatorPolicyAssignments(aggregatorPolicyAssignmentsRaw interface{}) []interface{} {
	if aggregatorPolicyAssignmentsRaw == nil {
		return nil
	}

	aggregatorPolicyAssignments := aggregatorPolicyAssignmentsRaw.([]interface{})
	res := make([]interface{}, len(aggregatorPolicyAssignments))
	for i, v := range aggregatorPolicyAssignments {
		assignment := v.(map[string]interface{})
		res[i] = map[string]interface{}{
			"policy_assignment_id":   assignment["policy_assignment_id"],
			"policy_assignment_name": assignment["policy_assignment_name"],
			"account_id":             assignment["account_id"],
			"account_name":           assignment["account_name"],
			"compliance": flattenAggregatorPolicyAssignmentsCompliance(
				utils.PathSearch("compliance", assignment, nil)),
		}
	}

	return res
}

func flattenAggregatorPolicyAssignmentsCompliance(complianceRaw interface{}) []interface{} {
	if complianceRaw == nil {
		return nil
	}

	compliance := complianceRaw.(map[string]interface{})
	res := []interface{}{
		map[string]interface{}{
			"compliance_state": compliance["compliance_state"],
			"resource_details": flattenAggregatorPolicyAssignmentsComplianceResourceDetails(
				utils.PathSearch("resource_details", compliance, nil)),
		},
	}

	return res
}

func flattenAggregatorPolicyAssignmentsComplianceResourceDetails(resourceDetailsRaw interface{}) []interface{} {
	if resourceDetailsRaw == nil {
		return nil
	}

	details := resourceDetailsRaw.(map[string]interface{})
	res := []interface{}{
		map[string]interface{}{
			"compliant_count":     details["compliant_count"],
			"non_compliant_count": details["non_compliant_count"],
		},
	}

	return res
}
