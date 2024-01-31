// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product Organizations
// ---------------------------------------------------------------

package organizations

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Organizations GET /v1/organizations/policies
func DataSourcePolicies() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourcePoliciesRead,
		Schema: map[string]*schema.Schema{
			"build_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the build type of the policy.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the name of the policy.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the type of the policy.`,
			},
			"policies": {
				Type:        schema.TypeList,
				Elem:        policiesPolicySchema(),
				Computed:    true,
				Description: `List of policies in an organization.`,
			},
		},
	}
}

func policiesPolicySchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the unique ID of the policy.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the name of the policy.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the type of the policy.`,
			},
			"urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the uniform resource name of the policy.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the description of the policy.`,
			},
			"build_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the build type of the policy.`,
			},
		},
	}
	return &sc
}

func resourcePoliciesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// listPolicies: Query the list of policies in an organization.
	var (
		listPoliciesHttpUrl = "v1/organizations/policies"
		listPoliciesProduct = "organizations"
	)
	listPoliciesClient, err := cfg.NewServiceClient(listPoliciesProduct, region)
	if err != nil {
		return diag.Errorf("error creating Organizations client: %s", err)
	}

	listPoliciesPath := listPoliciesClient.Endpoint + listPoliciesHttpUrl

	listPoliciesResp, err := pagination.ListAllItems(
		listPoliciesClient,
		"marker",
		listPoliciesPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Organizations policies")
	}

	listPoliciesRespJson, err := json.Marshal(listPoliciesResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var listPoliciesRespBody interface{}
	err = json.Unmarshal(listPoliciesRespJson, &listPoliciesRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		mErr,
		d.Set("policies", flattenPoliciesRespPolicy(d, listPoliciesRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenPoliciesRespPolicy(d *schema.ResourceData, resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("policies", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		isBuiltin := utils.PathSearch("is_builtin", v, false).(bool)
		rawBuildType := getBuildType(isBuiltin)
		if buildType, ok := d.GetOk("build_type"); ok && buildType != rawBuildType {
			continue
		}
		rawName := utils.PathSearch("name", v, "").(string)
		if name, ok := d.GetOk("name"); ok && name != rawName {
			continue
		}
		rawPolicyType := utils.PathSearch("type", v, "").(string)
		if policyType, ok := d.GetOk("type"); ok && policyType != rawPolicyType {
			continue
		}
		rst = append(rst, map[string]interface{}{
			"id":          utils.PathSearch("id", v, nil),
			"name":        rawName,
			"type":        rawPolicyType,
			"urn":         utils.PathSearch("urn", v, nil),
			"description": utils.PathSearch("description", v, nil),
			"build_type":  rawBuildType,
		})
	}
	return rst
}

func getBuildType(isBuiltin bool) string {
	if isBuiltin {
		return "system"
	}
	return "custom"
}
