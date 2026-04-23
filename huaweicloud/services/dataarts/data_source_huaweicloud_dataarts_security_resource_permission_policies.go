package dataarts

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DataArtsStudio GET /v1/{project_id}/security/permission-resource
func DataSourceSecurityResourcePermissionPolicies() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSecurityResourcePermissionPoliciesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the resource permission policies are located.`,
			},

			// Required parameters.
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the workspace to which the resource permission policies belong.`,
			},

			// Optional parameters.
			"policy_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the resource permission policy to be queried.`,
			},
			"resource_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the authorized resource to be queried.`,
			},
			"member_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the authorized member to be queried.`,
			},

			// Attributes.
			"policies": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of resource permission policies that matched filter parameters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"policy_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the resource permission policy.`,
						},
						"policy_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the resource permission policy.`,
						},
						"resources": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The resource list of the resource permission policy.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"resource_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The ID of the resource.`,
									},
									"resource_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The name of the resource.`,
									},
									"resource_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The type of the resource.`,
									},
								},
							},
						},
						"members": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The member list of the resource permission policy.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"member_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The ID of the member.`,
									},
									"member_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The name of the member.`,
									},
									"member_type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The type of the member.`,
									},
								},
							},
						},
						"create_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the resource permission policy, in RFC3339 format.`,
						},
						"create_user": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creator of the resource permission policy.`,
						},
						"update_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The latest update time of the resource permission policy, in RFC3339 format.`,
						},
					},
				},
			},
		},
	}
}

func buildSecurityResourcePermissionPoliciesQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("policy_name"); ok {
		res = fmt.Sprintf("%s&policy_name=%v", res, v)
	}

	if v, ok := d.GetOk("resource_name"); ok {
		res = fmt.Sprintf("%s&resource_name=%v", res, v)
	}

	if v, ok := d.GetOk("member_name"); ok {
		res = fmt.Sprintf("%s&member_name=%v", res, v)
	}

	return res
}

func listSecurityResourcePermissionPolicies(client *golangsdk.ServiceClient, workspaceId string, queryParams ...string) ([]interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/security/permission-resource?limit={limit}"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))
	if len(queryParams) > 0 && queryParams[0] != "" {
		listPath += queryParams[0]
	}

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildSecurityMoreHeaders(workspaceId),
	}

	for {
		listPathWithOffset := listPath + fmt.Sprintf("&offset=%d", offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		policies := utils.PathSearch("policies", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, policies...)
		if len(policies) < limit {
			break
		}

		offset += len(policies)
	}

	return result, nil
}

func flattenSecurityResourcePermissionPoliciesResources(resources []interface{}) []map[string]interface{} {
	if len(resources) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(resources))
	for _, item := range resources {
		result = append(result, map[string]interface{}{
			"resource_id":   utils.PathSearch("resource_id", item, nil),
			"resource_name": utils.PathSearch("resource_name", item, nil),
			"resource_type": utils.PathSearch("resource_type", item, nil),
		})
	}

	return result
}

func flattenSecurityResourcePermissionPoliciesMembers(members []interface{}) []map[string]interface{} {
	if len(members) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(members))
	for _, member := range members {
		result = append(result, map[string]interface{}{
			"member_id":   utils.PathSearch("member_id", member, nil),
			"member_name": utils.PathSearch("member_name", member, nil),
			"member_type": utils.PathSearch("member_type", member, nil),
		})
	}

	return result
}

func flattenSecurityResourcePermissionPolicies(policies []interface{}) []map[string]interface{} {
	if len(policies) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(policies))
	for _, policy := range policies {
		result = append(result, map[string]interface{}{
			"policy_id":   utils.PathSearch("policy_id", policy, nil),
			"policy_name": utils.PathSearch("policy_name", policy, nil),
			"resources": flattenSecurityResourcePermissionPoliciesResources(utils.PathSearch("resources",
				policy, make([]interface{}, 0)).([]interface{})),
			"members": flattenSecurityResourcePermissionPoliciesMembers(utils.PathSearch("members",
				policy, make([]interface{}, 0)).([]interface{})),
			"create_user": utils.PathSearch("create_user", policy, nil),
			"create_time": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("create_time",
				policy, float64(0)).(float64))/1000, false),
			"update_time": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("update_time",
				policy, float64(0)).(float64))/1000, false),
		})
	}

	return result
}

func dataSourceSecurityResourcePermissionPoliciesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
	)

	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	policies, err := listSecurityResourcePermissionPolicies(client, workspaceId, buildSecurityResourcePermissionPoliciesQueryParams(d))
	if err != nil {
		return diag.Errorf("error querying DataArts Security resource permission policies: %s", err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(dataSourceId)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("policies", flattenSecurityResourcePermissionPolicies(policies)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
