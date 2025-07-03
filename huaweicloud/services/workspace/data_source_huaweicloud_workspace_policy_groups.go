package workspace

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

// @API Workspace GET /v2/{project_id}/policy-groups/detail
func DataSourcePolicyGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePolicyGroupsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region in which to query the data source.`,
			},
			"policy_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the policy group.`,
			},
			"policy_group_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The name of the policy group.`,
			},
			"priority": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The priority of the policy group.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The description of the policy group.`,
			},
			"policy_groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        policyGroupSchema(),
				Description: `The list of policy groups that match the filter parameters.`,
			},
		},
	}
}

func policyGroupSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"policy_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the policy group.`,
			},
			"policy_group_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the policy group.`,
			},
			"priority": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The priority of the policy group.`,
			},
			"update_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The update time of the policy group, in RFC3339 format.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The description of the policy group.`,
			},
			"policies": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        policySchema(),
				Description: `The list of policy configurations.`,
			},
			"targets": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        targetSchema(),
				Description: `The list of target configurations.`,
			},
		},
	}
}

func policySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"peripherals": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The peripheral device policies, in JSON format.`,
			},
			"audio": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The audio policies, in JSON format.`,
			},
			"client": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The client policies, in JSON format.`,
			},
			"display": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The display policies, in JSON format.`,
			},
			"file_and_clipboard": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The file and clipboard policies, in JSON format.`,
			},
			"session": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The session policies, in JSON format.`,
			},
			"virtual_channel": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The virtual channel policies, in JSON format.`,
			},
			"watermark": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The watermark policies, in JSON format.`,
			},
			"keyboard_mouse": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The keyboard and mouse policies, in JSON format.`,
			},
			"seamless": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The general audio and video bypass policies, in JSON format.`,
			},
			"personalized_data_mgmt": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The personalized data management policies, in JSON format.`,
			},
			"custom": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The custom policies, in JSON format.`,
			},
			"record_audit": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The screen recording audit policies, in JSON format.`,
			},
		},
	}
}

func targetSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"target_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the target.`,
			},
			"target_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the target.`,
			},
			"target_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the target.`,
			},
		},
	}
}

func flattenPolicyGroups(policyGroups []interface{}) []interface{} {
	if len(policyGroups) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(policyGroups))
	for _, item := range policyGroups {
		result = append(result, map[string]interface{}{
			"policy_group_id":   utils.PathSearch("policy_group_id", item, nil),
			"policy_group_name": utils.PathSearch("policy_group_name", item, nil),
			"priority":          utils.PathSearch("priority", item, nil),
			"update_time":       utils.PathSearch("update_time", item, nil),
			"description":       utils.PathSearch("description", item, nil),
			"policies":          flattenPolicyGroupsPolicies(utils.PathSearch("policies", item, nil)),
			"targets":           flattenPolicyGroupsTargets(utils.PathSearch("targets", item, make([]interface{}, 0)).([]interface{})),
		})
	}
	return result
}

func flattenPolicyGroupsPolicies(policies interface{}) []interface{} {
	if policies == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"peripherals":            utils.JsonToString(utils.PathSearch("peripherals", policies, nil)),
			"audio":                  utils.JsonToString(utils.PathSearch("audio", policies, nil)),
			"client":                 utils.JsonToString(utils.PathSearch("client", policies, nil)),
			"display":                utils.JsonToString(utils.PathSearch("display", policies, nil)),
			"file_and_clipboard":     utils.JsonToString(utils.PathSearch("file_and_clipboard", policies, nil)),
			"session":                utils.JsonToString(utils.PathSearch("session", policies, nil)),
			"virtual_channel":        utils.JsonToString(utils.PathSearch("virtual_channel", policies, nil)),
			"watermark":              utils.JsonToString(utils.PathSearch("watermark", policies, nil)),
			"keyboard_mouse":         utils.JsonToString(utils.PathSearch("keyboard_mouse", policies, nil)),
			"seamless":               utils.JsonToString(utils.PathSearch("seamless", policies, nil)),
			"personalized_data_mgmt": utils.JsonToString(utils.PathSearch("personalized_data_mgmt", policies, nil)),
			"custom":                 utils.JsonToString(utils.PathSearch("custom", policies, nil)),
			"record_audit":           utils.JsonToString(utils.PathSearch("record_audit", policies, nil)),
		},
	}
}

func flattenPolicyGroupsTargets(targets []interface{}) []interface{} {
	if len(targets) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(targets))
	for _, item := range targets {
		result = append(result, map[string]interface{}{
			"target_id":   utils.PathSearch("target_id", item, nil),
			"target_type": utils.PathSearch("target_type", item, nil),
			"target_name": utils.PathSearch("target_name", item, nil),
		})
	}
	return result
}

func buildListPolicyGroupsParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("policy_group_id"); ok {
		res = fmt.Sprintf("%s&policy_group_id=%v", res, v)
	}

	if v, ok := d.GetOk("policy_group_name"); ok {
		res = fmt.Sprintf("%s&policy_group_name=%v", res, v)
	}

	if v, ok := d.GetOk("priority"); ok {
		res = fmt.Sprintf("%s&priority=%v", res, v)
	}

	if v, ok := d.GetOk("description"); ok {
		res = fmt.Sprintf("%s&description=%v", res, v)
	}
	return res
}

func listPolicyGroups(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/policy-groups/detail?limit={limit}"
		offset  = 0
		limit   = 100
		result  = make([]interface{}, 0)
	)

	listPathWithLimit := client.Endpoint + httpUrl
	listPathWithLimit = strings.ReplaceAll(listPathWithLimit, "{project_id}", client.ProjectID)
	listPathWithLimit = strings.ReplaceAll(listPathWithLimit, "{limit}", strconv.Itoa(limit))
	listPathWithLimit += buildListPolicyGroupsParams(d)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPathWithLimit, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		policyGroups := utils.PathSearch("policy_groups", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, policyGroups...)
		if len(policyGroups) < limit {
			break
		}
		offset += len(policyGroups)
	}
	return result, nil
}

func dataSourcePolicyGroupsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	policyGroups, err := listPolicyGroups(client, d)
	if err != nil {
		return diag.Errorf("error querying Workspace policy groups: %s", err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr := multierror.Append(
		nil,
		d.Set("region", region),
		d.Set("policy_groups", flattenPolicyGroups(policyGroups)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
