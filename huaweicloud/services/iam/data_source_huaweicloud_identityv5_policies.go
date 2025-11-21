package iam

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IAM GET /v5/policies
func DataSourceIdentityV5Policies() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityV5PoliciesRead,

		Schema: map[string]*schema.Schema{
			"policy_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The type of identity policy, can be \"custom\" or \"system\".",
			},
			"path_prefix": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The resource path prefix.",
			},
			"only_attached": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to list only policies that are attached to entities.",
			},
			"policies": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of IAM policies.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"attachment_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of entities attached to the policy.",
						},
						"default_version_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the default version of the policy.",
						},
						"path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The path of the policy.",
						},
						"policy_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the policy.",
						},
						"policy_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the policy.",
						},
						"policy_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of the policy.",
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The time when the default version of the policy was last updated.",
						},
						"urn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URN of the policy.",
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The time when the policy was created.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of the policy.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIdentityV5PoliciesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("iam", region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	var allPolicies []interface{}
	var marker string
	var path string
	for {
		path = client.Endpoint + "v5/policies" + buildListPoliciesV5Params(d, marker)
		reqOpt := &golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		r, err := client.Request("GET", path, reqOpt)
		if err != nil {
			return common.CheckDeletedDiag(d, err, "error retrieving policies")
		}
		resp, err := utils.FlattenResponse(r)
		if err != nil {
			return diag.FromErr(err)
		}
		policies := flattenListPoliciesV5Response(resp)
		allPolicies = append(allPolicies, policies...)

		marker = utils.PathSearch("page_info.next_marker", resp, "").(string)
		if marker == "" {
			break
		}
	}

	id, _ := uuid.GenerateUUID()
	d.SetId(id)
	mErr := multierror.Append(nil,
		d.Set("policies", allPolicies),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListPoliciesV5Response(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	policies := utils.PathSearch("policies", resp, make([]interface{}, 0)).([]interface{})
	result := make([]interface{}, len(policies))
	for i, policy := range policies {
		result[i] = map[string]interface{}{
			"policy_type":        utils.PathSearch("policy_type", policy, nil),
			"policy_name":        utils.PathSearch("policy_name", policy, nil),
			"policy_id":          utils.PathSearch("policy_id", policy, nil),
			"urn":                utils.PathSearch("urn", policy, nil),
			"path":               utils.PathSearch("path", policy, nil),
			"default_version_id": utils.PathSearch("default_version_id", policy, nil),
			"attachment_count":   utils.PathSearch("attachment_count", policy, nil),
			"description":        utils.PathSearch("description", policy, nil),
			"created_at":         utils.PathSearch("created_at", policy, nil),
		}
	}
	return result
}

func buildListPoliciesV5Params(d *schema.ResourceData, marker string) string {
	res := "?limit=100"
	if v, ok := d.GetOk("policy_type"); ok {
		res = fmt.Sprintf("%s&policy_type=%v", res, v)
	}
	if v, ok := d.GetOk("path_prefix"); ok {
		res = fmt.Sprintf("%s&path_prefix=%v", res, v)
	}
	if v, ok := d.GetOk("only_attached"); ok {
		res = fmt.Sprintf("%s&only_attached=%v", res, v)
	}
	if marker != "" {
		res = fmt.Sprintf("%s&marker=%v", res, marker)
	}
	return res
}
