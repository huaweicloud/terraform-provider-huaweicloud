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

// @API IoTDA GET /v5/iot/{project_id}/routing-rule/backlog-policy
func DataSourceDataBacklogPolicies() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDataBacklogPoliciesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
						"backlog_size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"backlog_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceDataBacklogPoliciesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		allPolicies []model.BacklogPolicyInfo
		limit       = int32(50)
		offset      int32
	)

	for {
		listOpts := model.ListRoutingBacklogPolicyRequest{
			PolicyName: utils.StringIgnoreEmpty(d.Get("policy_name").(string)),
			Limit:      utils.Int32(limit),
			Offset:     &offset,
		}

		listResp, listErr := client.ListRoutingBacklogPolicy(&listOpts)
		if listErr != nil {
			return diag.Errorf("error querying IoTDA data backlog policies: %s", listErr)
		}

		if listResp == nil || listResp.BacklogPolicies == nil {
			break
		}

		if len(*listResp.BacklogPolicies) == 0 {
			break
		}

		allPolicies = append(allPolicies, *listResp.BacklogPolicies...)
		//nolint:gosec
		offset += int32(len(*listResp.BacklogPolicies))
	}

	uuID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(uuID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("policies", flattenDataBacklogPolicies(allPolicies)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDataBacklogPolicies(policies []model.BacklogPolicyInfo) []interface{} {
	if len(policies) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(policies))
	for _, v := range policies {
		rst = append(rst, map[string]interface{}{
			"id":           v.PolicyId,
			"name":         v.PolicyName,
			"description":  v.Description,
			"backlog_size": v.BacklogSize,
			"backlog_time": v.BacklogTime,
		})
	}

	return rst
}
