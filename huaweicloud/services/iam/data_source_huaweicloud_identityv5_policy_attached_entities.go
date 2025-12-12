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

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// DataSourceIdentityV5PolicyAttachedEntities
// @API IAM GET /v5/policies/{policy_id}/attached-entities
func DataSourceIdentityV5PolicyAttachedEntities() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityV5PolicyAttachedEntitiesRead,

		Schema: map[string]*schema.Schema{
			"policy_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"entity_type": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"policy_agencies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"agency_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"attached_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"policy_groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"attached_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"policy_users": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"attached_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceIdentityV5PolicyAttachedEntitiesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("iam", region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	policyID := d.Get("policy_id").(string)
	var allUsers []interface{}
	var allGroups []interface{}
	var allAgencies []interface{}
	var marker string

	for {
		path := client.Endpoint + "v5/policies/{policy_id}/attached-entities" + buildListEntitiesForPolicyV5Params(d, marker)
		path = strings.ReplaceAll(path, "{policy_id}", policyID)
		reqOpt := &golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		r, err := client.Request("GET", path, reqOpt)
		if err != nil {
			return diag.Errorf("error retrieving policy attached entities: %s", err)
		}
		resp, err := utils.FlattenResponse(r)
		if err != nil {
			return diag.FromErr(err)
		}

		users := flattenListEntitiesForPolicyV5ResponseUsers(resp)
		allUsers = append(allUsers, users...)

		groups := flattenListEntitiesForPolicyV5ResponseGroups(resp)
		allGroups = append(allGroups, groups...)

		agencies := flattenListEntitiesForPolicyV5ResponseAgencies(resp)
		allAgencies = append(allAgencies, agencies...)

		marker = utils.PathSearch("page_info.next_marker", resp, "").(string)
		if marker == "" {
			break
		}
	}

	id, _ := uuid.GenerateUUID()
	d.SetId(id)
	mErr := multierror.Append(nil,
		d.Set("policy_users", allUsers),
		d.Set("policy_groups", allGroups),
		d.Set("policy_agencies", allAgencies),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListEntitiesForPolicyV5ResponseUsers(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	users := utils.PathSearch("policy_users", resp, make([]interface{}, 0)).([]interface{})
	result := make([]interface{}, len(users))
	for i, user := range users {
		result[i] = map[string]interface{}{
			"user_id":     utils.PathSearch("user_id", user, nil),
			"attached_at": utils.PathSearch("attached_at", user, nil),
		}
	}
	return result
}

func flattenListEntitiesForPolicyV5ResponseGroups(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	groups := utils.PathSearch("policy_groups", resp, make([]interface{}, 0)).([]interface{})
	result := make([]interface{}, len(groups))
	for i, group := range groups {
		result[i] = map[string]interface{}{
			"group_id":    utils.PathSearch("group_id", group, nil),
			"attached_at": utils.PathSearch("attached_at", group, nil),
		}
	}
	return result
}

func flattenListEntitiesForPolicyV5ResponseAgencies(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	agencies := utils.PathSearch("policy_agencies", resp, make([]interface{}, 0)).([]interface{})
	result := make([]interface{}, len(agencies))
	for i, agency := range agencies {
		result[i] = map[string]interface{}{
			"agency_id":   utils.PathSearch("agency_id", agency, nil),
			"attached_at": utils.PathSearch("attached_at", agency, nil),
		}
	}
	return result
}

func buildListEntitiesForPolicyV5Params(d *schema.ResourceData, marker string) string {
	res := "?limit=100"
	if v, ok := d.GetOk("entity_type"); ok {
		res = fmt.Sprintf("%s&entity_type=%v", res, v)
	}
	if marker != "" {
		res = fmt.Sprintf("%s&marker=%v", res, marker)
	}
	return res
}
