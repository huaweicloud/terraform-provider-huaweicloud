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
	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/throttles"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API APIG POST /v2/{project_id}/apigw/instances/{instance_id}/throttle-bindings
// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/throttle-bindings/unbinded-apis
// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/throttle-bindings/binded-apis
// @API APIG DELETE /v2/{project_id}/apigw/instances/{instance_id}/throttle-bindings/{throttle_binding_id}
func ResourceThrottlingPolicyAssociate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceThrottlingPolicyAssociateCreate,
		ReadContext:   resourceThrottlingPolicyAssociateRead,
		UpdateContext: resourceThrottlingPolicyAssociateUpdate,
		DeleteContext: resourceThrottlingPolicyAssociateDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceThrottlingPolicyAssociateImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(3 * time.Minute),
			Update: schema.DefaultTimeout(3 * time.Minute),
			Delete: schema.DefaultTimeout(3 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region where the dedicated instance and the throttling policy are located.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the dedicated instance to which the APIs and the throttling policy belongs.",
			},
			"policy_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the throttling policy.",
			},
			"publish_ids": {
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The publish IDs corresponding to the APIs bound by the throttling policy.",
			},
		},
	}
}

func throttlingPolicyBindingRefreshFunc(client *golangsdk.ServiceClient, instanceId, policyId string,
	publishIds []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var (
			httpUrl  = "v2/{project_id}/apigw/instances/{instance_id}/throttle-bindings/unbinded-apis"
			queryUrl = "?throttle_id={throttle_id}"
			offset   = 0
			result   = make([]interface{}, 0)
		)

		listPath := client.Endpoint + httpUrl
		listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
		listPath = strings.ReplaceAll(listPath, "{instance_id}", instanceId)
		queryUrl = strings.ReplaceAll(queryUrl, "{throttle_id}", policyId)
		listPath += queryUrl

		opt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}

		for {
			listPathWithOffset := fmt.Sprintf("%s&limit=100&offset=%d", listPath, offset)
			requestResp, err := client.Request("GET", listPathWithOffset, &opt)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error retrieving unassociated throttling policies: %s", err)
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

func bindThrottlingPolicyToApis(ctx context.Context, client *golangsdk.ServiceClient, opts throttles.BindOpts,
	timeout time.Duration) error {
	var (
		instanceId = opts.InstanceId
		policyId   = opts.ThrottleId
		publishIds = opts.PublishIds
	)

	_, err := throttles.Bind(client, opts)
	if err != nil {
		return fmt.Errorf("error binding throttling policy to one or more APIs: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: throttlingPolicyBindingRefreshFunc(client, instanceId, policyId, publishIds),
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

func resourceThrottlingPolicyAssociateCreate(ctx context.Context, d *schema.ResourceData,
	meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	var (
		instanceId = d.Get("instance_id").(string)
		policyId   = d.Get("policy_id").(string)
		publishIds = d.Get("publish_ids").(*schema.Set)

		opt = throttles.BindOpts{
			InstanceId: instanceId,
			ThrottleId: policyId,
			PublishIds: utils.ExpandToStringListBySet(publishIds),
		}
	)
	err = bindThrottlingPolicyToApis(ctx, client, opt, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(fmt.Sprintf("%s/%s", instanceId, policyId))

	return resourceThrottlingPolicyAssociateRead(ctx, d, meta)
}

func buildListOpts(instanceId, policyId string) throttles.ListBindOpts {
	return throttles.ListBindOpts{
		InstanceId: instanceId,
		ThrottleId: policyId,
		Limit:      500,
	}
}

func flattenApiPublishIds(apiList []throttles.ApiForThrottle) []string {
	if len(apiList) < 1 {
		return nil
	}

	result := make([]string, len(apiList))
	for i, val := range apiList {
		result[i] = val.PublishId
	}
	return result
}

func resourceThrottlingPolicyAssociateRead(_ context.Context, d *schema.ResourceData,
	meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %v", err)
	}

	var (
		instanceId = d.Get("instance_id").(string)
		policyId   = d.Get("policy_id").(string)
		opt        = buildListOpts(instanceId, policyId)
	)

	resp, err := throttles.ListBind(client, opt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error getting API information from server")
	}
	if len(resp) < 1 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr := multierror.Append(nil, d.Set("publish_ids", flattenApiPublishIds(resp)))
	return diag.FromErr(mErr.ErrorOrNil())
}

func throttlingPolicyUnbindingRefreshFunc(client *golangsdk.ServiceClient, instanceId, policyId string,
	publishIds []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var (
			httpUrl  = "v2/{project_id}/apigw/instances/{instance_id}/throttle-bindings/binded-apis"
			queryUrl = "?throttle_id={throttle_id}"
			offset   = 0
			result   = make([]interface{}, 0)
		)

		listPath := client.Endpoint + httpUrl
		listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
		listPath = strings.ReplaceAll(listPath, "{instance_id}", instanceId)
		queryUrl = strings.ReplaceAll(queryUrl, "{throttle_id}", policyId)
		listPath += queryUrl

		opt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}

		for {
			listPathWithOffset := fmt.Sprintf("%s&limit=100&offset=%d", listPath, offset)
			requestResp, err := client.Request("GET", listPathWithOffset, &opt)
			if err != nil {
				return nil, "ERROR", fmt.Errorf("error retrieving associated throttling policies: %s", err)
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

func unbindPolicy(ctx context.Context, client *golangsdk.ServiceClient, opt throttles.ListBindOpts, unbindSet *schema.Set,
	timeout time.Duration) error {
	var (
		instanceId = opt.InstanceId
		policyId   = opt.ThrottleId
	)
	resp, err := throttles.ListBind(client, opt)
	if err != nil {
		// The instance or throttling policy not exist.
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			return err
		}
		return fmt.Errorf("error getting API information from server: %s", err)
	}

	if len(resp) < 1 {
		log.Printf("[DEBUG] All APIs has been disassociated from the throttling policy (%s) under dedicated instance (%s)",
			policyId, instanceId)
		return nil
	}

	publishIds := make([]string, 0, unbindSet.Len())
	for _, rm := range unbindSet.List() {
		for _, api := range resp {
			// If the publish ID is not found, it means the policy has been unbound from the API by other ways.
			if rm == api.PublishId {
				publishIds = append(publishIds, api.PublishId)
				err = throttles.Unbind(client, instanceId, api.ThrottleApplyId)
				if err != nil {
					if _, ok := err.(golangsdk.ErrDefault404); ok {
						log.Printf("[DEBUG] All APIs has been disassociated from the throttling policy (%s)", policyId)
						continue
					}
					return fmt.Errorf("error unbound policy from the API: %s", err)
				}
				break
			}
		}
	}

	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: throttlingPolicyUnbindingRefreshFunc(client, instanceId, policyId, publishIds),
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

func resourceThrottlingPolicyAssociateUpdate(ctx context.Context, d *schema.ResourceData,
	meta interface{}) diag.Diagnostics {
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
		opt := buildListOpts(instanceId, policyId)
		err = unbindPolicy(ctx, client, opt, rmSet, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if addSet.Len() > 0 {
		opt := throttles.BindOpts{
			InstanceId: instanceId,
			ThrottleId: policyId,
			PublishIds: utils.ExpandToStringListBySet(addSet),
		}
		err = bindThrottlingPolicyToApis(ctx, client, opt, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceThrottlingPolicyAssociateRead(ctx, d, meta)
}

func resourceThrottlingPolicyAssociateDelete(ctx context.Context, d *schema.ResourceData,
	meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	var (
		instanceId = d.Get("instance_id").(string)
		policyId   = d.Get("policy_id").(string)
		publishIds = d.Get("publish_ids").(*schema.Set)
		opt        = buildListOpts(instanceId, policyId)
	)
	if err = unbindPolicy(ctx, client, opt, publishIds, d.Timeout(schema.TimeoutDelete)); err != nil {
		return common.CheckDeletedDiag(d, err, "error unbinding APIs from throttling policy")
	}

	return nil
}

func resourceThrottlingPolicyAssociateImportState(_ context.Context, d *schema.ResourceData,
	_ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, must be <instance_id>/<policy_id>")
	}

	d.Set("instance_id", parts[0])
	d.Set("policy_id", parts[1])

	return []*schema.ResourceData{d}, nil
}
