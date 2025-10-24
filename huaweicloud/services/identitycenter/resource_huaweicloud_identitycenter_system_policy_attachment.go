// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product IdentityCenter
// ---------------------------------------------------------------

package identitycenter

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IdentityCenter POST /v1/instances/{instance_id}/permission-sets/{permission_set_id}/attach-managed-role
// @API IdentityCenter POST /v1/instances/{instance_id}/permission-sets/{permission_set_id}/detach-managed-role
// @API IdentityCenter GET /v1/instances/{instance_id}/permission-sets/{permission_set_id}/managed-roles
func ResourceSystemPolicyAttachment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSystemPolicyAttachmentCreate,
		UpdateContext: resourceSystemPolicyAttachmentUpdate,
		ReadContext:   resourceSystemPolicyAttachmentRead,
		DeleteContext: resourceSystemPolicyAttachmentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSystemPolicyAttachmentImport,
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
				Elem:     AttachedSystemPoliciesSchema(),
			},
		},
	}
}

func AttachedSystemPoliciesSchema() *schema.Resource {
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

func resourceSystemPolicyAttachmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	attachClient, err := cfg.NewServiceClient("identitycenter", region)
	if err != nil {
		return diag.Errorf("error creating Identity Center client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	psID := d.Get("permission_set_id").(string)
	policyList := utils.ExpandToStringListBySet(d.Get("policy_ids").(*schema.Set))

	if err := attachSystemPolicies(attachClient, instanceID, psID, policyList); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(psID)

	if diagErr := provisionPermissionSet(attachClient, instanceID, psID); diagErr != nil {
		return diagErr
	}

	return resourceSystemPolicyAttachmentRead(ctx, d, meta)
}

func resourceSystemPolicyAttachmentRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		listSystemPoliciesHttpUrl = "v1/instances/{instance_id}/permission-sets/{permission_set_id}/managed-roles"
		listSystemPoliciesProduct = "identitycenter"
	)
	listSystemPoliciesClient, err := cfg.NewServiceClient(listSystemPoliciesProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center client: %s", err)
	}

	listSystemPoliciesPath := listSystemPoliciesClient.Endpoint + listSystemPoliciesHttpUrl
	listSystemPoliciesPath = strings.ReplaceAll(listSystemPoliciesPath, "{instance_id}", d.Get("instance_id").(string))
	listSystemPoliciesPath = strings.ReplaceAll(listSystemPoliciesPath, "{permission_set_id}", d.Get("permission_set_id").(string))

	listSystemPoliciesOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	listSystemPoliciesResp, err := listSystemPoliciesClient.Request("GET", listSystemPoliciesPath, &listSystemPoliciesOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving attached system policies")
	}

	listSystemPoliciesRespBody, err := utils.FlattenResponse(listSystemPoliciesResp)
	if err != nil {
		return diag.FromErr(err)
	}

	policyIDs, policies := flattenAttachedSystemPolicies(listSystemPoliciesRespBody)
	if len(policyIDs) == 0 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "no system policies attached")
	}

	mErr := multierror.Append(nil,
		d.Set("policy_ids", policyIDs),
		d.Set("attached_policies", policies),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAttachedSystemPolicies(resp interface{}) ([]string, []interface{}) {
	if resp == nil {
		return nil, nil
	}

	curJson := utils.PathSearch("attached_managed_roles", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	ids := make([]string, len(curArray))
	objects := make([]interface{}, len(curArray))
	for i, v := range curArray {
		policyID := utils.PathSearch("role_id", v, "")
		policyName := utils.PathSearch("role_name", v, nil)

		ids[i] = policyID.(string)
		objects[i] = map[string]interface{}{
			"id":   policyID,
			"name": policyName,
		}
	}
	return ids, objects
}

func resourceSystemPolicyAttachmentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	if err := detachSystemPolicies(updateClient, instanceID, psID, removeList); err != nil {
		return diag.Errorf("error updating system policy attachment: %s", err)
	}

	addList := utils.ExpandToStringListBySet(addSet)
	if err := attachSystemPolicies(updateClient, instanceID, psID, addList); err != nil {
		return diag.Errorf("error updating system policy attachment: %s", err)
	}

	if diagErr := provisionPermissionSet(updateClient, instanceID, psID); diagErr != nil {
		return diagErr
	}

	return resourceSystemPolicyAttachmentRead(ctx, d, meta)
}

func resourceSystemPolicyAttachmentDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	detachClient, err := cfg.NewServiceClient("identitycenter", region)
	if err != nil {
		return diag.Errorf("error creating Identity Center client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	psID := d.Get("permission_set_id").(string)
	policyList := utils.ExpandToStringListBySet(d.Get("policy_ids").(*schema.Set))

	if err := detachSystemPolicies(detachClient, instanceID, psID, policyList); err != nil {
		return diag.FromErr(err)
	}

	//nolint:revive
	if diagErr := provisionPermissionSet(detachClient, instanceID, psID); diagErr != nil {
		return diagErr
	}

	return nil
}

func resourceSystemPolicyAttachmentImport(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
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

func attachSystemPolicies(client *golangsdk.ServiceClient, instanceID, psID string, policies []string) error {
	return requestSystemPolicies(client, "attach", instanceID, psID, policies)
}

func detachSystemPolicies(client *golangsdk.ServiceClient, instanceID, psID string, policies []string) error {
	return requestSystemPolicies(client, "detach", instanceID, psID, policies)
}

func requestSystemPolicies(client *golangsdk.ServiceClient, action, instanceID, psID string, policies []string) error {
	requestURI := fmt.Sprintf("v1/instances/%s/permission-sets/%s/%s-managed-role", instanceID, psID, action)
	requestPath := client.Endpoint + requestURI

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for _, policy := range policies {
		requestOpt.JSONBody = map[string]interface{}{
			"managed_role_id": policy,
		}
		_, err := client.Request("POST", requestPath, &requestOpt)
		if err != nil {
			return fmt.Errorf("failed to %s the system policy %s: %s", action, policy, err)
		}
	}

	return nil
}

// provisionPermissionSet: Provision the Permission Set to apply the corresponding updates to all assigned accounts
func provisionPermissionSet(client *golangsdk.ServiceClient, instanceID, psID string) diag.Diagnostics {
	accountDIs, err := getAssignededAccounts(client, instanceID, psID)
	if err != nil {
		return diag.Errorf("failed to get not provisioned accounts: %s", err)
	}

	log.Printf("[DEBUG] the following accounts need to provision: %v", accountDIs)
	if len(accountDIs) == 0 {
		return nil
	}

	requestURI := fmt.Sprintf("v1/instances/%s/permission-sets/%s/provision", instanceID, psID)
	requestPath := client.Endpoint + requestURI

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	var diags diag.Diagnostics
	for _, account := range accountDIs {
		requestOpt.JSONBody = map[string]interface{}{
			"target_type": "ACCOUNT",
			"target_id":   account,
		}

		_, err = client.Request("POST", requestPath, &requestOpt)
		if err != nil {
			diagIcagent := diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "failed to provision permission set",
				Detail:   fmt.Sprintf("failed to provision account %s with permission set %s: %s", account, psID, err),
			}
			diags = append(diags, diagIcagent)
		}
	}

	return diags
}
