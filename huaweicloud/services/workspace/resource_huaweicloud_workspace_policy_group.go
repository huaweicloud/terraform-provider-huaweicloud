package workspace

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/workspace/v2/policygroups"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Workspace DELETE /v2/{project_id}/policy-groups/{policy_group_id}
// @API Workspace GET /v2/{project_id}/policy-groups/{policy_group_id}
// @API Workspace PUT /v2/{project_id}/policy-groups/{policy_group_id}
// @API Workspace POST /v2/{project_id}/policy-groups
func ResourcePolicyGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicyGroupCreate,
		ReadContext:   resourcePolicyGroupRead,
		UpdateContext: resourcePolicyGroupUpdate,
		DeleteContext: resourcePolicyGroupDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region where the policy group is located.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the policy group.",
			},
			"priority": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
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
							Description: "The target type.",
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The target name.",
						},
						"id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The target ID.",
						},
					},
				},
				Description: "The list of target objects.",
			},
			"policy": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_control": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ip_access_control": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The IP access configuration.",
									},
								},
							},
							Description: "The configuration of the access policy control.",
						},
					},
				},
				Description: "The configuration of the access policy",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The update time of the policy group.",
			},
		},
	}
}

func buildPolicyGroupTargets(targets *schema.Set) []policygroups.Target {
	if targets.Len() < 1 {
		return nil
	}
	result := make([]policygroups.Target, targets.Len())
	for i, val := range targets.List() {
		target := val.(map[string]interface{})
		result[i] = policygroups.Target{
			TargetType: target["type"].(string),
			TargetId:   target["id"].(string),
			TargetName: target["name"].(string),
		}
	}
	return result
}

func buildPolicyGroupAccessControl(rules []interface{}) *policygroups.AccessControl {
	if len(rules) < 1 {
		return nil
	}
	rule := rules[0].(map[string]interface{})
	return &policygroups.AccessControl{
		IpAccessControl: rule["ip_access_control"].(string),
	}
}

func buildPolicyGroupPolicy(policies []interface{}) *policygroups.Policy {
	if len(policies) < 1 {
		return nil
	}

	policy := policies[0].(map[string]interface{})
	return &policygroups.Policy{
		AccessControl: *buildPolicyGroupAccessControl(policy["access_control"].([]interface{})),
	}
}

func resourcePolicyGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.WorkspaceV2Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Workspace v2 client: %s", err)
	}

	opts := policygroups.CreateOpts{
		PolicyGroupName: d.Get("name").(string),
		Priority:        d.Get("priority").(int),
		Description:     d.Get("description").(string),
		Targets:         buildPolicyGroupTargets(d.Get("targets").(*schema.Set)),
		Policies:        *buildPolicyGroupPolicy(d.Get("policy").([]interface{})),
	}
	groupId, err := policygroups.Create(client, opts)
	if err != nil {
		return diag.Errorf("error creating Workspace policy group: %s", err)
	}

	d.SetId(groupId)

	return resourcePolicyGroupRead(ctx, d, meta)
}

func flattenPolicyGroupTargets(targets []policygroups.Target) []interface{} {
	if len(targets) < 1 {
		return nil
	}
	result := make([]interface{}, len(targets))
	for i, target := range targets {
		result[i] = map[string]interface{}{
			"type": target.TargetType,
			"id":   target.TargetId,
			"name": target.TargetName,
		}
	}
	return result
}

func flattenPolicyGroupPolicy(policy policygroups.Policy) []interface{} {
	return []interface{}{
		map[string]interface{}{
			"access_control": []map[string]interface{}{
				{
					"ip_access_control": policy.AccessControl.IpAccessControl,
				},
			},
		},
	}
}

func resourcePolicyGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	client, err := conf.WorkspaceV2Client(region)
	if err != nil {
		return diag.Errorf("error creating Workspace v2 client: %s", err)
	}

	resp, err := policygroups.Get(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "Workspace policy group")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", resp.PolicyGroupName),
		d.Set("priority", resp.Priority),
		d.Set("description", resp.Description),
		d.Set("targets", flattenPolicyGroupTargets(resp.Targets)),
		d.Set("policy", flattenPolicyGroupPolicy(resp.Policy)),
		d.Set("updated_at", resp.UpdateTime),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourcePolicyGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.WorkspaceV2Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Workspace v2 client: %s", err)
	}

	var (
		groupId = d.Id()
		opts    = policygroups.UpdateOpts{
			PolicyGroupId:   groupId,
			PolicyGroupName: d.Get("name").(string),
			Priority:        d.Get("priority").(int),
			Description:     utils.String(d.Get("description").(string)),
			Targets:         buildPolicyGroupTargets(d.Get("targets").(*schema.Set)),
			Policies:        buildPolicyGroupPolicy(d.Get("policy").([]interface{})),
		}
	)
	_, err = policygroups.Update(client, opts)
	if err != nil {
		return diag.Errorf("error updating Workspace policy group (%s): %s", groupId, err)
	}
	return resourcePolicyGroupRead(ctx, d, meta)
}

func resourcePolicyGroupDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.WorkspaceV2Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Workspace v2 client: %s", err)
	}

	groupId := d.Id()
	err = policygroups.Delete(client, groupId)
	if err != nil {
		return diag.Errorf("error deleting Workspace policy group (%s): %s", groupId, err)
	}
	return nil
}
