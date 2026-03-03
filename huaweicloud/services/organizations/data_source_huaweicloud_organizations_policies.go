package organizations

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

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
				Description: `The build type of the policy.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the policy.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The type of the policy.`,
			},
			"policies": {
				Type:        schema.TypeList,
				Elem:        policiesPolicySchema(),
				Computed:    true,
				Description: `The list of policies that match the filter parameters.`,
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
				Description: `The unique ID of the policy.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the policy.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the policy.`,
			},
			"urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The uniform resource name of the policy.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The description of the policy.`,
			},
			"build_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The build type of the policy.`,
			},
		},
	}
	return &sc
}

func listPolicies(client *golangsdk.ServiceClient) ([]interface{}, error) {
	var (
		httpUrl = "v1/organizations/policies?limit=200"
		marker  = ""
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	for {
		listPathWithMarker := listPath
		if marker != "" {
			listPathWithMarker = fmt.Sprintf("%s&marker=%s", listPathWithMarker, marker)
		}

		resp, err := client.Request("GET", listPathWithMarker, &requestOpts)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		policies := utils.PathSearch("policies", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, policies...)
		marker = utils.PathSearch("page_info.next_marker", respBody, "").(string)
		if marker == "" {
			break
		}
	}

	return result, nil
}

func resourcePoliciesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("organizations", region)
	if err != nil {
		return diag.Errorf("error creating Organizations client: %s", err)
	}

	policies, err := listPolicies(client)
	if err != nil {
		return diag.Errorf("error retrieving Organizations policies: %s", err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(dataSourceId)

	return diag.FromErr(d.Set("policies", flattenPolicies(d, policies)))
}

func flattenPolicies(d *schema.ResourceData, policies []interface{}) []interface{} {
	if len(policies) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(policies))
	for _, v := range policies {
		rawBuildType := getBuildType(utils.PathSearch("is_builtin", v, false).(bool))
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
