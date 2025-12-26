package iam

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IAM GET /v5/users/{user_id}/attached-policies
func DataSourceIdentityV5UserAttachedPolicies() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityV5UserAttachedPoliciesRead,

		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"attached_policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"attached_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"urn": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceIdentityV5UserAttachedPoliciesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("iam", region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	var allPolicies []interface{}
	var marker string
	var path string

	userID := d.Get("user_id").(string)

	for {
		path = client.Endpoint + "v5/users/" + userID + "/attached-policies" + buildListAttachedUserPoliciesV5Params(marker)
		reqOpt := &golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		r, err := client.Request("GET", path, reqOpt)
		if err != nil {
			return diag.Errorf("error retrieving attached policies: %s", err)
		}
		resp, err := utils.FlattenResponse(r)
		if err != nil {
			return diag.FromErr(err)
		}

		policies := flattenListAttachedUserPoliciesV5Response(resp)
		allPolicies = append(allPolicies, policies...)

		marker = utils.PathSearch("page_info.next_marker", resp, "").(string)
		if marker == "" {
			break
		}
	}

	id, _ := uuid.GenerateUUID()
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("attached_policies", allPolicies),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildListAttachedUserPoliciesV5Params(marker string) string {
	res := "?limit=100"
	if marker != "" {
		res = fmt.Sprintf("%s&marker=%v", res, marker)
	}
	return res
}

func flattenListAttachedUserPoliciesV5Response(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	policies := utils.PathSearch("attached_policies", resp, make([]interface{}, 0)).([]interface{})
	result := make([]interface{}, len(policies))
	for i, policy := range policies {
		result[i] = map[string]interface{}{
			"attached_at": utils.PathSearch("attached_at", policy, nil),
			"policy_id":   utils.PathSearch("policy_id", policy, nil),
			"policy_name": utils.PathSearch("policy_name", policy, nil),
			"urn":         utils.PathSearch("urn", policy, nil),
		}
	}
	return result
}
