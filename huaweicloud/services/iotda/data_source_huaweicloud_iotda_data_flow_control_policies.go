package iotda

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iotda/v5/model"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IoTDA GET /v5/iot/{project_id}/routing-rule/flowcontrol-policy
func DataSourceDataFlowControlPolicies() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDataFlowControlPoliciesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"scope": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"scope_value": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"policy_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"policies": {
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
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scope": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scope_value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"limit": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceDataFlowControlPoliciesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		isDerived = WithDerivedAuth(cfg, region)
	)

	client, err := cfg.HcIoTdaV5Client(region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA v5 client: %s", err)
	}

	var (
		allPolicies []model.FlowControlPolicyInfo
		limit       = int32(50)
		offset      int32
	)

	for {
		listOpts := model.ListRoutingFlowControlPolicyRequest{
			Scope:      utils.StringIgnoreEmpty(d.Get("scope").(string)),
			ScopeValue: utils.StringIgnoreEmpty(d.Get("scope_value").(string)),
			PolicyName: utils.StringIgnoreEmpty(d.Get("policy_name").(string)),
			Limit:      utils.Int32(limit),
			Offset:     &offset,
		}

		listResp, listErr := client.ListRoutingFlowControlPolicy(&listOpts)
		if listErr != nil {
			return diag.Errorf("error querying IoTDA data flow control policies: %s", listErr)
		}

		if listResp == nil || listResp.FlowcontrolPolicies == nil {
			break
		}

		if len(*listResp.FlowcontrolPolicies) == 0 {
			break
		}

		allPolicies = append(allPolicies, *listResp.FlowcontrolPolicies...)
		//nolint:gosec
		offset += int32(len(*listResp.FlowcontrolPolicies))
	}

	uuID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(uuID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("policies", flattenDataFlowControlPolicies(allPolicies)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDataFlowControlPolicies(policies []model.FlowControlPolicyInfo) []interface{} {
	if len(policies) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(policies))
	for _, v := range policies {
		rst = append(rst, map[string]interface{}{
			"id":          v.PolicyId,
			"name":        v.PolicyName,
			"description": v.Description,
			"scope":       v.Scope,
			"scope_value": v.ScopeValue,
			"limit":       v.Limit,
		})
	}

	return rst
}
