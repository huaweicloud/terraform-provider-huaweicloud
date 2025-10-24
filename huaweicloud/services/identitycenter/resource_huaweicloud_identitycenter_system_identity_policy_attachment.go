package identitycenter

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IdentityCenter POST /v1/instances/{instance_id}/permission-sets/{permission_set_id}/attach-managed-policy
// @API IdentityCenter POST /v1/instances/{instance_id}/permission-sets/{permission_set_id}/detach-managed-policy
// @API IdentityCenter GET /v1/instances/{instance_id}/permission-sets/{permission_set_id}/managed-policies
func ResourceSystemIdentityPolicyAttachment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSystemIdentityPolicyAttachmentCreate,
		UpdateContext: resourceSystemIdentityPolicyAttachmentUpdate,
		ReadContext:   resourceSystemIdentityPolicyAttachmentRead,
		DeleteContext: resourceSystemIdentityPolicyAttachmentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSystemIdentityPolicyAttachmentImport,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"permission_set_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"policy_ids": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"attached_policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     AttachedSystemIdentityPoliciesSchema(),
			},
		},
	}
}

func AttachedSystemIdentityPoliciesSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func resourceSystemIdentityPolicyAttachmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	attachClient, err := cfg.NewServiceClient("identitycenter", region)
	if err != nil {
		return diag.Errorf("error creating Identity Center client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	psID := d.Get("permission_set_id").(string)
	policyList := utils.ExpandToStringListBySet(d.Get("policy_ids").(*schema.Set))

	if err := attachedSystemIdentityPolicies(attachClient, instanceID, psID, policyList); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(psID)

	if diagErr := provisionPermissionSet(attachClient, instanceID, psID); diagErr != nil {
		return diagErr
	}

	return resourceSystemIdentityPolicyAttachmentRead(ctx, d, meta)
}

func resourceSystemIdentityPolicyAttachmentRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		listSystemIdentityPoliciesHttpUrl = "v1/instances/{instance_id}/permission-sets/{permission_set_id}/managed-policies"
		listSystemIdentityPoliciesProduct = "identitycenter"
	)
	client, err := cfg.NewServiceClient(listSystemIdentityPoliciesProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center client: %s", err)
	}

	listSystemIdentityPoliciesPath := client.Endpoint + listSystemIdentityPoliciesHttpUrl
	listSystemIdentityPoliciesPath = strings.ReplaceAll(listSystemIdentityPoliciesPath, "{instance_id}", d.Get("instance_id").(string))
	listSystemIdentityPoliciesPath = strings.ReplaceAll(listSystemIdentityPoliciesPath, "{permission_set_id}", d.Get("permission_set_id").(string))

	listSystemIdentityPoliciesOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	listSystemIdentityPoliciesResp, err := client.Request("GET", listSystemIdentityPoliciesPath, &listSystemIdentityPoliciesOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving attached system identity policies")
	}

	listSystemIdentityPoliciesRespBody, err := utils.FlattenResponse(listSystemIdentityPoliciesResp)
	if err != nil {
		return diag.FromErr(err)
	}

	policyIDs, policies := flattenAttachedSystemIdentityPolicies(listSystemIdentityPoliciesRespBody)
	if len(policyIDs) == 0 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "no system identity policies attached")
	}

	mErr := multierror.Append(nil,
		d.Set("policy_ids", policyIDs),
		d.Set("attached_policies", policies),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAttachedSystemIdentityPolicies(resp interface{}) ([]string, []interface{}) {
	if resp == nil {
		return nil, nil
	}

	curJson := utils.PathSearch("attached_managed_policies", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	ids := make([]string, len(curArray))
	objects := make([]interface{}, len(curArray))
	for i, v := range curArray {
		policyID := utils.PathSearch("policy_id", v, "")
		policyName := utils.PathSearch("policy_name", v, nil)

		ids[i] = policyID.(string)
		objects[i] = map[string]interface{}{
			"id":   policyID,
			"name": policyName,
		}
	}
	return ids, objects
}

func resourceSystemIdentityPolicyAttachmentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	updateClient, err := cfg.NewServiceClient("identitycenter", region)
	if err != nil {
		return diag.Errorf("error creating Identity Center client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	psID := d.Get("permission_set_id").(string)

	oldRaw, newRaw := d.GetChange("policy_ids")
	rmSet := oldRaw.(*schema.Set).Difference(newRaw.(*schema.Set))
	addSet := newRaw.(*schema.Set).Difference(oldRaw.(*schema.Set))

	removeList := utils.ExpandToStringListBySet(rmSet)
	if err := detachSystemIdentityPolicies(updateClient, instanceID, psID, removeList); err != nil {
		return diag.Errorf("error updating system identity policy attachment: %s", err)
	}

	addList := utils.ExpandToStringListBySet(addSet)
	if err := attachedSystemIdentityPolicies(updateClient, instanceID, psID, addList); err != nil {
		return diag.Errorf("error updating system identity policy attachment: %s", err)
	}

	if diagErr := provisionPermissionSet(updateClient, instanceID, psID); diagErr != nil {
		return diagErr
	}

	return resourceSystemIdentityPolicyAttachmentRead(ctx, d, meta)
}

func resourceSystemIdentityPolicyAttachmentDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	detachClient, err := cfg.NewServiceClient("identitycenter", region)
	if err != nil {
		return diag.Errorf("error creating Identity Center client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	psID := d.Get("permission_set_id").(string)
	policyList := utils.ExpandToStringListBySet(d.Get("policy_ids").(*schema.Set))

	if err := detachSystemIdentityPolicies(detachClient, instanceID, psID, policyList); err != nil {
		return diag.FromErr(err)
	}

	//nolint:revive
	if diagErr := provisionPermissionSet(detachClient, instanceID, psID); diagErr != nil {
		return diagErr
	}

	return nil
}

func resourceSystemIdentityPolicyAttachmentImport(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		err := fmt.Errorf("invalid format: the format must be <instance id>/<permission set id>")
		return nil, err
	}

	instanceID := parts[0]
	psID := parts[1]

	d.SetId(psID)
	d.Set("instance_id", instanceID)
	d.Set("permission_set_id", psID)

	return []*schema.ResourceData{d}, nil
}

func attachedSystemIdentityPolicies(client *golangsdk.ServiceClient, instanceID, psID string, policies []string) error {
	return requestSystemIdentityPolicies(client, "attach", instanceID, psID, policies)
}

func detachSystemIdentityPolicies(client *golangsdk.ServiceClient, instanceID, psID string, policies []string) error {
	return requestSystemIdentityPolicies(client, "detach", instanceID, psID, policies)
}

func requestSystemIdentityPolicies(client *golangsdk.ServiceClient, action, instanceID, psID string, policies []string) error {
	requestURI := fmt.Sprintf("v1/instances/%s/permission-sets/%s/%s-managed-policy", instanceID, psID, action)
	requestPath := client.Endpoint + requestURI

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for _, policy := range policies {
		requestOpt.JSONBody = map[string]interface{}{
			"managed_policy_id": policy,
		}
		_, err := client.Request("POST", requestPath, &requestOpt)
		if err != nil {
			return fmt.Errorf("failed to %s the system policy %s: %s", action, policy, err)
		}
	}

	return nil
}
