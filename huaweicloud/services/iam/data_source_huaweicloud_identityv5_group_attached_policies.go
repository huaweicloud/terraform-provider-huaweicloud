package iam

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IAM GET /v5/groups/{group_id}/attached-policies
func DataSourceV5GroupAttachedPolicies() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceV5GroupAttachedPoliciesRead,

		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the IAM user group.",
			},
			"attached_policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
						"attached_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The creation time of the policy.",
						},
						"urn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The Uniform Resource Name (URN) of the policy.",
						},
					},
				},
				Description: "The list of policies attached to the user group.",
			},
		},
	}
}

func buildV5GroupAttachedPoliciesQueryParams(marker string) string {
	if marker != "" {
		return fmt.Sprintf("&marker=%v", marker)
	}
	return ""
}

func listV5GroupAttachedPolicies(client *golangsdk.ServiceClient, groupId string) ([]interface{}, error) {
	var (
		httpUrl = "v5/groups/{group_id}/attached-policies"
		limit   = 100
		marker  = ""
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{group_id}", groupId)
	listPath = fmt.Sprintf("%s?limit=%v", listPath, limit)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		listPathWithMarker := listPath + buildV5GroupAttachedPoliciesQueryParams(marker)
		resp, err := client.Request("GET", listPathWithMarker, &opt)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		policies := utils.PathSearch("attached_policies", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, policies...)

		if len(policies) < limit {
			break
		}

		marker = utils.PathSearch("page_info.next_marker", respBody, "").(string)
		if marker == "" {
			break
		}
	}

	return result, nil
}

func dataSourceV5GroupAttachedPoliciesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("iam", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	groupId := d.Get("group_id").(string)
	policies, err := listV5GroupAttachedPolicies(client, groupId)
	if err != nil {
		return diag.Errorf("error retrieving group (%s) attached policies: %s", groupId, err)
	}

	randomId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(randomId)

	return diag.FromErr(d.Set("attached_policies", flattenV5GroupAttachedPolicies(policies)))
}

func flattenV5GroupAttachedPolicies(policies []interface{}) []interface{} {
	if len(policies) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(policies))
	for _, policy := range policies {
		result = append(result, map[string]interface{}{
			"policy_id":   utils.PathSearch("policy_id", policy, nil),
			"policy_name": utils.PathSearch("policy_name", policy, nil),
			"urn":         utils.PathSearch("urn", policy, nil),
			"attached_at": utils.PathSearch("attached_at", policy, nil),
		})
	}
	return result
}
