package apig

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/acls"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API APIG GET /v2/{project_id}/apigw/instances/{instanceId}/acl-bindings/binded-apis
// @API APIG POST /v2/{project_id}/apigw/instances/{instanceId}/acl-bindings
func ResourceAclPolicyAssociate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAclPolicyAssociateCreate,
		ReadContext:   resourceAclPolicyAssociateRead,
		UpdateContext: resourceAclPolicyAssociateUpdate,
		DeleteContext: resourceAclPolicyAssociateDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceAclPolicyAssociateImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region where the ACL policy and the APIs are located.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the dedicated instance to which the APIs and the ACL policy belong.",
			},
			"policy_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ACL Policy ID for APIs binding.",
			},
			"publish_ids": {
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The publish IDs corresponding to the APIs bound by the ACL policy.",
			},
		},
	}
}

func resourceAclPolicyAssociateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	var (
		instanceId = d.Get("instance_id").(string)
		policyId   = d.Get("policy_id").(string)
		publishIds = d.Get("publish_ids").(*schema.Set)

		opts = acls.BindOpts{
			InstanceId: instanceId,
			PolicyId:   policyId,
			PublishIds: utils.ExpandToStringListBySet(publishIds),
		}
	)
	_, err = acls.Bind(client, opts)
	if err != nil {
		return diag.Errorf("error binding policy to the API: %s", err)
	}
	d.SetId(fmt.Sprintf("%s/%s", instanceId, policyId))

	return resourceAclPolicyAssociateRead(ctx, d, meta)
}

func buildAclPolicyListOpts(instanceId, policyId string) acls.ListBindOpts {
	return acls.ListBindOpts{
		InstanceId: instanceId,
		PolicyId:   policyId,
		Limit:      500, // This limitation can be removed after the parameter 'offset' is fixed.
	}
}

func flattenApiPublishIdsForAclPolicy(apiList []acls.AclBindApiInfo) []string {
	if len(apiList) < 1 {
		return nil
	}

	result := make([]string, len(apiList))
	for i, val := range apiList {
		result[i] = val.PublishId
	}
	return result
}

func resourceAclPolicyAssociateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.ApigV2Client(region)
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	var (
		instanceId = d.Get("instance_id").(string)
		policyId   = d.Get("policy_id").(string)
		opts       = buildAclPolicyListOpts(instanceId, policyId)
	)

	resp, err := acls.ListBind(client, opts)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "ACL policy association")
	}
	if len(resp) < 1 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	return diag.FromErr(d.Set("publish_ids", flattenApiPublishIdsForAclPolicy(resp)))
}

func unbindAclPolicy(client *golangsdk.ServiceClient, instanceId, policyId string, unbindSet *schema.Set) error {
	opt := buildAclPolicyListOpts(instanceId, policyId)
	resp, err := acls.ListBind(client, opt)
	if err != nil {
		return fmt.Errorf("error getting binding APIs based on ACL policy (%s): %s", policyId, err)
	}

	unbindList := make([]string, 0, unbindSet.Len())
	for _, rm := range unbindSet.List() {
		for _, api := range resp {
			// If the publish ID is not found, it means the policy has been unbound from the API by other ways.
			if rm == api.PublishId {
				unbindList = append(unbindList, api.BindId)
			}
		}
	}
	opts := acls.BatchUnbindOpts{
		InstanceId:  instanceId,
		AclBindings: unbindList,
	}
	unbindResp, err := acls.BatchUnbind(client, opts, "delete")
	if err != nil {
		return err
	}
	if len(unbindResp.Failures) > 0 {
		return fmt.Errorf("an error occurred during unbind: %#v", unbindResp.Failures)
	}
	return nil
}

func resourceAclPolicyAssociateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	var (
		instanceId     = d.Get("instance_id").(string)
		policyId       = d.Get("policy_id").(string)
		oldRaw, newRaw = d.GetChange("publish_ids")

		addSet = newRaw.(*schema.Set).Difference(oldRaw.(*schema.Set))
		rmSet  = oldRaw.(*schema.Set).Difference(newRaw.(*schema.Set))
	)

	if rmSet.Len() > 0 {
		err = unbindAclPolicy(client, instanceId, policyId, rmSet)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if addSet.Len() > 0 {
		opt := acls.BindOpts{
			InstanceId: instanceId,
			PolicyId:   policyId,
			PublishIds: utils.ExpandToStringListBySet(addSet),
		}
		_, err = acls.Bind(client, opt)
		if err != nil {
			return diag.Errorf("error binding published APIs to the ACL policy (%s): %s", policyId, err)
		}
	}

	return resourceAclPolicyAssociateRead(ctx, d, meta)
}

func resourceAclPolicyAssociateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	var (
		instanceId = d.Get("instance_id").(string)
		policyId   = d.Get("policy_id").(string)
		publishIds = d.Get("publish_ids").(*schema.Set)
	)

	return diag.FromErr(unbindAclPolicy(client, instanceId, policyId, publishIds))
}

func resourceAclPolicyAssociateImportState(_ context.Context, d *schema.ResourceData,
	_ interface{}) ([]*schema.ResourceData, error) {
	importedId := d.Id()
	parts := strings.SplitN(importedId, "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<instance_id>/<policy_id>', but got '%s'",
			importedId)
	}

	mErr := multierror.Append(nil,
		d.Set("instance_id", parts[0]),
		d.Set("policy_id", parts[1]),
	)
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
