package iam

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IAM GET /v5/groups/{group_id}/attached-policies
func DataSourceIdentityV5GroupAttachedPolicies() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityV5GroupAttachedPoliciesRead,

		Schema: map[string]*schema.Schema{
			"group_id": {
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

func dataSourceIdentityV5GroupAttachedPoliciesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("iam", region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	groupId := d.Get("group_id").(string)

	var allPolicies []interface{}
	var marker string
	var path string

	for {
		path = client.Endpoint + "v5/groups/{group_id}/attached-policies" + buildListGroupAttachedPoliciesV5Params(marker)
		path = strings.ReplaceAll(path, "{group_id}", groupId)
		reqOpt := &golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		r, err := client.Request("GET", path, reqOpt)
		if err != nil {
			return common.CheckDeletedDiag(d, err, "error retrieving attached policies")
		}
		resp, err := utils.FlattenResponse(r)
		if err != nil {
			return diag.FromErr(err)
		}
		policies := flattenListGroupAttachedPoliciesV5Response(resp)
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

func flattenListGroupAttachedPoliciesV5Response(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	policies := utils.PathSearch("attached_policies", resp, make([]interface{}, 0)).([]interface{})
	result := make([]interface{}, len(policies))
	for i, policy := range policies {
		result[i] = map[string]interface{}{
			"policy_name": utils.PathSearch("policy_name", policy, nil),
			"policy_id":   utils.PathSearch("policy_id", policy, nil),
			"urn":         utils.PathSearch("urn", policy, nil),
			"attached_at": utils.PathSearch("attached_at", policy, nil),
		}
	}
	return result
}

func buildListGroupAttachedPoliciesV5Params(marker string) string {
	res := "?limit=100"
	if marker != "" {
		res = fmt.Sprintf("%s&marker=%v", res, marker)
	}
	return res
}
