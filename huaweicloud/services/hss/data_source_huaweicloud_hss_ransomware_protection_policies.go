package hss

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	hssv5model "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/hss/v5/model"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API HSS GET /v5/{project_id}/ransomware/protection/policy
func DataSourceRansomwareProtectionPolicies() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceRansomwareProtectionPoliciesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"policy_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"operating_system": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
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
						"protection_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"bait_protection_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"deploy_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"protection_directory": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"protection_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"exclude_directory": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"runtime_detection_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"count_associated_server": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"operating_system": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"process_whitelist": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"path": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"hash": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"default_policy": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func datasourceRansomwareProtectionPoliciesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		epsId       = cfg.GetEnterpriseProjectID(d, "all_granted_eps")
		limit       = int32(20)
		offset      int32
		allPolicies []hssv5model.ProtectionPolicyInfo
	)

	client, err := cfg.HcHssV5Client(region)
	if err != nil {
		return diag.Errorf("error creating HSS v5 client: %s", err)
	}

	for {
		request := hssv5model.ListProtectionPolicyRequest{
			Region:              region,
			EnterpriseProjectId: utils.String(epsId),
			Limit:               utils.Int32(limit),
			Offset:              utils.Int32(offset),
			ProtectPolicyId:     utils.StringIgnoreEmpty(d.Get("policy_id").(string)),
			PolicyName:          utils.StringIgnoreEmpty(d.Get("name").(string)),
			OperatingSystem:     utils.StringIgnoreEmpty(d.Get("operating_system").(string)),
		}

		listResp, listErr := client.ListProtectionPolicy(&request)
		if listErr != nil {
			return diag.Errorf("error querying HSS ransomware protection policies: %s", listErr)
		}

		if listResp == nil || listResp.DataList == nil {
			break
		}
		if len(*listResp.DataList) == 0 {
			break
		}

		allPolicies = append(allPolicies, *listResp.DataList...)
		offset += limit
	}

	uuId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(uuId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("policies", flattenPolicies(allPolicies)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenPolicies(policies []hssv5model.ProtectionPolicyInfo) []interface{} {
	if len(policies) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(policies))
	for _, v := range policies {
		rst = append(rst, map[string]interface{}{
			"id":                       v.PolicyId,
			"name":                     v.PolicyName,
			"protection_mode":          v.ProtectionMode,
			"bait_protection_status":   v.BaitProtectionStatus,
			"deploy_mode":              v.DeployMode,
			"protection_directory":     v.ProtectionDirectory,
			"protection_type":          v.ProtectionType,
			"exclude_directory":        v.ExcludeDirectory,
			"runtime_detection_status": v.RuntimeDetectionStatus,
			"count_associated_server":  v.CountAssociatedServer,
			"operating_system":         v.OperatingSystem,
			"process_whitelist":        flattenProcessWhitelist(v.ProcessWhitelist),
			"default_policy":           v.DefaultPolicy,
		})
	}

	return rst
}

func flattenProcessWhitelist(processWhitelist *[]hssv5model.TrustProcessInfo) []interface{} {
	if processWhitelist == nil {
		return nil
	}

	rst := make([]interface{}, 0, len(*processWhitelist))
	for _, v := range *processWhitelist {
		rst = append(rst, map[string]interface{}{
			"path": v.Path,
			"hash": v.Hash,
		})
	}

	return rst
}
