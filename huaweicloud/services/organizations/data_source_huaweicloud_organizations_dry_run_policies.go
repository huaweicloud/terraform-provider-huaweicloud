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

// @API Organizations GET /v1/organizations/dry-run-policies
func DataSourceDryRunPolicies() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDryRunPoliciesRead,
		Schema: map[string]*schema.Schema{
			// Optional parameters.
			"policy_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The type of the dry run policies to be queried.`,
			},
			"attached_entity_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the entity associated with the dry run policy.`,
			},

			// Attributes.
			"policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The unique ID of the dry run policy.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the dry run policy.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the dry run policy.`,
						},
						"urn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The uniform resource name of the dry run policy.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the dry run policy.`,
						},
						"is_builtin": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether the dry run policy is a built-in policy.`,
						},
					},
				},
				Description: `The list of dry run policies that matched the filter parameters.`,
			},
		},
	}
}

func buildDryRunPoliciesBodyParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("policy_type"); ok {
		res = fmt.Sprintf("&policy_type=%v", v)
	}

	if v, ok := d.GetOk("attached_entity_id"); ok {
		res = fmt.Sprintf("%s&attached_entity_id=%v", res, v)
	}

	if res != "" {
		res = "?" + res[1:]
	}

	return res
}

func listDryRunPolicies(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v1/organizations/dry-run-policies"
		marker  = ""
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath += buildDryRunPoliciesBodyParams(d)
	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	for {
		listPathWithMarker := listPath
		if marker != "" {
			listPathWithMarker = fmt.Sprintf("%s&marker=%s", listPathWithMarker, marker)
		}

		resp, err := client.Request("GET", listPathWithMarker, &listOpt)
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

func dataSourceDryRunPoliciesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("organizations", region)
	if err != nil {
		return diag.Errorf("error creating Organizations client: %s", err)
	}

	policies, err := listDryRunPolicies(client, d)
	if err != nil {
		return diag.Errorf("error retrieving dry run policies: %s", err)
	}

	dataSourceID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(dataSourceID)

	return diag.FromErr(d.Set("policies", flattenDryRunPolicies(policies)))
}

func flattenDryRunPolicies(policies []interface{}) []interface{} {
	if len(policies) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(policies))
	for _, v := range policies {
		rst = append(rst, map[string]interface{}{
			"id":          utils.PathSearch("id", v, nil),
			"name":        utils.PathSearch("name", v, nil),
			"type":        utils.PathSearch("type", v, nil),
			"urn":         utils.PathSearch("urn", v, nil),
			"description": utils.PathSearch("description", v, nil),
			"is_builtin":  utils.PathSearch("is_builtin", v, nil),
		})
	}

	return rst
}
