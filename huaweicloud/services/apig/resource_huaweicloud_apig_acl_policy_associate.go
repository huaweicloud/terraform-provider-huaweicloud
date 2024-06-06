package apig

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/acls"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API APIG POST /v2/{project_id}/apigw/instances/{instance_id}/acl-bindings
// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/acl-bindings/unbinded-apis
// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/acl-bindings/binded-apis
// @API APIG PUT /v2/{project_id}/apigw/instances/{instance_id}/acl-bindings
func ResourceAclPolicyAssociate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAclPolicyAssociateCreate,
		ReadContext:   resourceAclPolicyAssociateRead,
		UpdateContext: resourceAclPolicyAssociateUpdate,
		DeleteContext: resourceAclPolicyAssociateDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
			Update: schema.DefaultTimeout(3 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},

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

func aclPolicyBindingRefreshFunc(client *golangsdk.ServiceClient, instanceId, policyId string,
	publishIds []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var (
			httpUrl  = "v2/{project_id}/apigw/instances/{instance_id}/acl-bindings/unbinded-apis"
			queryUrl = "?acl_id={acl_id}"
			offset   = 0
			result   = make([]interface{}, 0)
		)

		listPath := client.Endpoint + httpUrl
		listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
		listPath = strings.ReplaceAll(listPath, "{instance_id}", instanceId)
		queryUrl = strings.ReplaceAll(queryUrl, "{acl_id}", policyId)
		listPath += queryUrl

		opt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}

		for {
			listPathWithOffset := fmt.Sprintf("%s&limit=100&offset=%d", listPath, offset)
			requestResp, err := client.Request("GET", listPathWithOffset, &opt)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error retrieving unassociated ACL policies: %s", err)
			}
			respBody, err := utils.FlattenResponse(requestResp)
			if err != nil {
				return nil, "ERROR", err
			}
			unbindPublishIds := utils.PathSearch("apis[*].publish_id", respBody, make([]interface{}, 0)).([]interface{})
			if len(unbindPublishIds) < 1 {
				break
			}
			result = append(result, unbindPublishIds...)
			offset += len(unbindPublishIds)
		}

		if utils.IsSliceContainsAnyAnotherSliceElement(utils.ExpandToStringList(result), publishIds, false, true) {
			return result, "PENDING", nil
		}
		return result, "COMPLETED", nil
	}
}

func bindAclPolicyToApis(ctx context.Context, client *golangsdk.ServiceClient, opts acls.BindOpts,
	timeout time.Duration) error {
	var (
		instanceId = opts.InstanceId
		policyId   = opts.PolicyId
		publishIds = opts.PublishIds
	)

	_, err := acls.Bind(client, opts)
	if err != nil {
		return fmt.Errorf("error binding ACL policy to one or more APIs: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: aclPolicyBindingRefreshFunc(client, instanceId, policyId, publishIds),
		Timeout: timeout,
		// In most cases, the bind operation will be completed immediately, but in a few cases, it needs to wait
		// for a short period of time, and the polling is performed by incrementing the time here.
		MinTimeout: 2 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for the binding completed: %s", err)
	}
	return nil
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

		opt = acls.BindOpts{
			InstanceId: instanceId,
			PolicyId:   policyId,
			PublishIds: utils.ExpandToStringListBySet(publishIds),
		}
	)
	err = bindAclPolicyToApis(ctx, client, opt, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
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

func aclPolicyUnbindingRefreshFunc(client *golangsdk.ServiceClient, instanceId, policyId string,
	publishIds []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var (
			httpUrl  = "v2/{project_id}/apigw/instances/{instance_id}/acl-bindings/binded-apis"
			queryUrl = "?acl_id={acl_id}"
			offset   = 0
			result   = make([]interface{}, 0)
		)

		listPath := client.Endpoint + httpUrl
		listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
		listPath = strings.ReplaceAll(listPath, "{instance_id}", instanceId)
		queryUrl = strings.ReplaceAll(queryUrl, "{acl_id}", policyId)
		listPath += queryUrl

		opt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}

		for {
			listPathWithOffset := fmt.Sprintf("%s&limit=100&offset=%d", listPath, offset)
			requestResp, err := client.Request("GET", listPathWithOffset, &opt)
			if err != nil {
				// The API returns a 404 error, which means that the instance or ACL policy has been deleted.
				// In this case, there's no need to disassociate API, also this action has been completed.
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					return "instance_or_ACL_policy_not_exist", "COMPLETED", nil
				}
				return nil, "ERROR", fmt.Errorf("error retrieving associated ACL policies: %s", err)
			}
			respBody, err := utils.FlattenResponse(requestResp)
			if err != nil {
				return nil, "ERROR", err
			}
			unbindPublishIds := utils.PathSearch("apis[*].publish_id", respBody, make([]interface{}, 0)).([]interface{})
			if len(unbindPublishIds) < 1 {
				break
			}
			result = append(result, unbindPublishIds...)
			offset += len(unbindPublishIds)
		}

		if utils.IsSliceContainsAnyAnotherSliceElement(utils.ExpandToStringList(result), publishIds, false, true) {
			return result, "PENDING", nil
		}
		return result, "COMPLETED", nil
	}
}

func unbindAclPolicy(ctx context.Context, client *golangsdk.ServiceClient, opt acls.ListBindOpts, publishIds *schema.Set,
	timeout time.Duration) error {
	var (
		instanceId = opt.InstanceId
		policyId   = opt.PolicyId
		unbindOpt  = acls.BatchUnbindOpts{
			InstanceId: instanceId,
		}
	)
	resp, err := acls.ListBind(client, opt)
	if err != nil {
		// The instance or ACL policy not exist.
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			return err
		}
		return fmt.Errorf("error getting binding APIs based on ACL policy (%s): %s", policyId, err)
	}

	bindIds := make([]string, 0, publishIds.Len())
	for _, publishId := range publishIds.List() {
		for _, api := range resp {
			// If the publish ID is not found, it means the policy has been unbound from the API by other ways.
			if publishId == api.PublishId {
				bindIds = append(bindIds, api.BindId)
			}
		}
	}

	if len(bindIds) < 1 {
		log.Printf("[DEBUG] All APIs has been disassociated from the ACL policy (%s)", policyId)
		return nil
	}

	unbindOpt.AclBindings = bindIds
	unbindResp, err := acls.BatchUnbind(client, unbindOpt, "delete")
	if err != nil {
		return err
	}
	if len(unbindResp.Failures) > 0 {
		return fmt.Errorf("an error occurred during unbind: %#v", unbindResp.Failures)
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: aclPolicyUnbindingRefreshFunc(client, instanceId, policyId, utils.ExpandToStringList(publishIds.List())),
		Timeout: timeout,
		// In most cases, the bind operation will be completed immediately, but in a few cases, it needs to wait
		// for a short period of time, and the polling is performed by incrementing the time here.
		MinTimeout: 2 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for the binding completed: %s", err)
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
		opt := buildAclPolicyListOpts(instanceId, policyId)
		err = unbindAclPolicy(ctx, client, opt, rmSet, d.Timeout(schema.TimeoutUpdate))
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
		err = bindAclPolicyToApis(ctx, client, opt, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceAclPolicyAssociateRead(ctx, d, meta)
}

func resourceAclPolicyAssociateDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	var (
		instanceId = d.Get("instance_id").(string)
		policyId   = d.Get("policy_id").(string)
		publishIds = d.Get("publish_ids").(*schema.Set)
		opt        = buildAclPolicyListOpts(instanceId, policyId)
	)

	if err = unbindAclPolicy(ctx, client, opt, publishIds, d.Timeout(schema.TimeoutUpdate)); err != nil {
		return common.CheckDeletedDiag(d, err, "error unbinding APIs from ACL policy")
	}
	return nil
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
