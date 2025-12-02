package workspace

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

// @API Workspace GET /v1/{project_id}/policy-groups
func DataSourceAppPolicyGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAppPolicyGroupsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the policy groups are located.`,
			},
			"policy_group_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the policy group.`,
			},
			"policy_group_type": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The type of the policy group.`,
			},
			"policy_groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the policy group.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the policy group.`,
						},
						"priority": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The priority of the policy group.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the policy group.`,
						},
						"targets": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The list of target objects.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The ID of the target object.`,
									},
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The name of the target object.`,
									},
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The type of the target object.`,
									},
								},
							},
						},
						"policies": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The policies of the policy group, in JSON format.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the policy group, in RFC3339 format.`,
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The latest update time of the policy group, in RFC3339 format.`,
						},
					},
				},
				Description: `The list of policy groups that match the filter parameters.`,
			},
		},
	}
}

func listAppPolicyGroups(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/policy-groups"
		offset  = 0
		limit   = 100
		result  = make([]interface{}, 0)
	)
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = fmt.Sprintf("%s?limit=%d", listPath, limit)
	listPath += buildAppPolicyGroupsQueryParams(d)
	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=UTF-8",
		},
	}

	for {
		listPathWithOffset := listPath + fmt.Sprintf("&offset=%d", offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &listOpt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		policyGroups := utils.PathSearch("items", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, policyGroups...)
		if len(policyGroups) < limit {
			break
		}
		offset += len(policyGroups)
	}

	return result, nil
}

func dataSourceAppPolicyGroupsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("appstream", region)
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	policyGroups, err := listAppPolicyGroups(client, d)
	if err != nil {
		return diag.Errorf("error retrieving policy groups: %s", err)
	}

	randomId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(randomId)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("policy_groups", flattenAppPolicyGroups(policyGroups)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildAppPolicyGroupsQueryParams(d *schema.ResourceData) string {
	var res string

	if v, ok := d.GetOk("policy_group_name"); ok {
		res = fmt.Sprintf("%s&policy_group_name=%v", res, v)
	}
	if v, ok := d.GetOk("policy_group_type"); ok {
		res = fmt.Sprintf("%s&policy_group_type=%v", res, v)
	}

	return res
}

func flattenAppPolicyGroups(policyGroups []interface{}) []interface{} {
	if len(policyGroups) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(policyGroups))
	for _, policyGroup := range policyGroups {
		result = append(result, map[string]interface{}{
			"id":          utils.PathSearch("id", policyGroup, nil),
			"name":        utils.PathSearch("name", policyGroup, nil),
			"priority":    utils.PathSearch("priority", policyGroup, 0),
			"description": utils.PathSearch("description", policyGroup, nil),
			"targets": flattenAppPolicyGroupsTargets(utils.PathSearch("targets",
				policyGroup, make([]interface{}, 0)).([]interface{})),
			"policies": utils.JsonToString(utils.PathSearch("policies", policyGroup, nil)),
			"created_at": utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(
				utils.PathSearch("create_time", policyGroup, "").(string))/1000, false),
			"updated_at": utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(
				utils.PathSearch("update_time", policyGroup, "").(string))/1000, false),
		})
	}

	return result
}

func flattenAppPolicyGroupsTargets(targets []interface{}) []interface{} {
	if len(targets) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(targets))
	for _, target := range targets {
		result = append(result, map[string]interface{}{
			"id":   utils.PathSearch("target_id", target, nil),
			"name": utils.PathSearch("target_name", target, nil),
			"type": utils.PathSearch("target_type", target, nil),
		})
	}

	return result
}
