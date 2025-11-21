package workspace

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Workspace POST /v1/{project_id}/policy-groups
// @API Workspace GET /v1/{project_id}/policy-groups
// @API Workspace PATCH /v1/{project_id}/policy-groups/{policy_group_id}
// @API Workspace DELETE /v1/{project_id}/policy-groups/{policy_group_id}
func ResourceAppPolicyGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppPolicyGroupCreate,
		ReadContext:   resourceAppPolicyGroupRead,
		UpdateContext: resourceAppPolicyGroupUpdate,
		DeleteContext: resourceAppPolicyGroupDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the policy group.",
			},
			"priority": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The priority of the policy group.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the policy group.",
			},
			"targets": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The object type.",
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The object name.",
						},
						"id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The object ID.",
						},
					},
				},
				Description: "The list of target objects.",
			},
			"policies": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsJSON,
				Description:  `The policies of the policy group, in JSON format.`,
			},
		},
	}
}

func resourceAppPolicyGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	httpUrl := "v1/{project_id}/policy-groups"
	client, err := cfg.NewServiceClient("appstream", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildAppPolicyGroupParams(d),
	}
	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating Workspace APP policy group: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	policyGroupId := utils.PathSearch("id", respBody, "").(string)
	if policyGroupId == "" {
		return diag.Errorf("unable to find the policy group ID from the API response")
	}
	d.SetId(policyGroupId)

	return resourceAppPolicyGroupRead(ctx, d, meta)
}

func buildAppPolicyGroupParams(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"policy_group_name": d.Get("name"),
		"description":       d.Get("description"),
		"targets":           buildTargets(d.Get("targets").(*schema.Set)),
		"policies":          utils.StringToJson(d.Get("policies").(string)),
	}

	if priority, ok := d.GetOk("priority"); ok {
		params["priority"] = priority
	}

	return map[string]interface{}{
		"policy_group": params,
	}
}

// The `targets` parameter can be set an empty list.
func buildTargets(targets *schema.Set) []map[string]interface{} {
	result := make([]map[string]interface{}, targets.Len())
	for i, v := range targets.List() {
		result[i] = map[string]interface{}{
			"target_id":   utils.PathSearch("id", v, nil),
			"target_name": utils.PathSearch("name", v, nil),
			"target_type": utils.PathSearch("type", v, nil),
		}
	}
	return result
}

func resourceAppPolicyGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("appstream", region)
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	policyGroup, err := GetAppGroupPolicy(client, d.Get("name").(string), d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Workspace APP policy group")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", policyGroup, nil)),
		d.Set("description", utils.PathSearch("description", policyGroup, nil)),
		d.Set("targets", flattenAppPolicyGroupTargets(utils.PathSearch("targets", policyGroup, make([]interface{}, 0)).([]interface{}))),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

// GetAppGroupPolicy is a method used to query policy group detail.
func GetAppGroupPolicy(client *golangsdk.ServiceClient, policyGroupName, policyGroupId string) (interface{}, error) {
	// limit: The default value is `10`.
	httpUrl := "v1/{project_id}/policy-groups?limit=100"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	// The "policy_group_name" parameter is fuzzy search.
	if policyGroupName != "" {
		getPath += fmt.Sprintf("&policy_group_name=%v", policyGroupName)
	}

	getOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	offset := 0
	for {
		getPathWithOffset := fmt.Sprintf("%s&offset=%d", getPath, offset)
		requestResp, err := client.Request("GET", getPathWithOffset, &getOpts)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		policyGroups := utils.PathSearch("items", respBody, make([]interface{}, 0)).([]interface{})
		if len(policyGroups) == 0 {
			break
		}

		policyGroup := utils.PathSearch(fmt.Sprintf("[?id=='%s']|[0]", policyGroupId), policyGroups, nil)
		if policyGroup != nil {
			return policyGroup, nil
		}
		offset += len(policyGroups)
	}

	return nil, golangsdk.ErrDefault404{}
}

func flattenAppPolicyGroupTargets(targets []interface{}) []interface{} {
	if len(targets) < 1 {
		return nil
	}
	result := make([]interface{}, len(targets))
	for i, target := range targets {
		result[i] = map[string]interface{}{
			"id":   utils.PathSearch("target_id", target, nil),
			"name": utils.PathSearch("target_name", target, nil),
			"type": utils.PathSearch("target_type", target, nil),
		}
	}
	return result
}

func resourceAppPolicyGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("appstream", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	httpUrl := "v1/{project_id}/policy-groups/{policy_group_id}"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{policy_group_id}", d.Id())
	updateOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildAppPolicyGroupParams(d),
	}

	_, err = client.Request("PATCH", updatePath, &updateOpts)
	if err != nil {
		return diag.Errorf("error updating Workspace APP policy group (%s): %s", d.Get("name").(string), err)
	}
	return resourceAppPolicyGroupRead(ctx, d, meta)
}

func resourceAppPolicyGroupDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("appstream", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	httpUrl := "v1/{project_id}/policy-groups/{policy_group_id}"
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{policy_group_id}", d.Id())

	deleteOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = client.Request("DELETE", deletePath, &deleteOpts)
	if err != nil {
		return diag.Errorf("error deleting Workspace APP policy group (%s): %s", d.Get("name").(string), err)
	}
	return nil
}
