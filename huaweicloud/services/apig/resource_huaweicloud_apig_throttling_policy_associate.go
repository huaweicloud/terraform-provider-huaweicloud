package apig

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/throttles"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/throttle-bindings/binded-apis
// @API APIG DELETE /v2/{project_id}/apigw/instances/{instance_id}/throttle-bindings/{throttle_binding_id}
// @API APIG POST /v2/{project_id}/apigw/instances/{instance_id}/throttle-bindings
func ResourceThrottlingPolicyAssociate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceThrottlingPolicyAssociateCreate,
		ReadContext:   resourceThrottlingPolicyAssociateRead,
		UpdateContext: resourceThrottlingPolicyAssociateUpdate,
		DeleteContext: resourceThrottlingPolicyAssociateDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceThrottlingPolicyAssociateImportState,
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
	_, err = throttles.Bind(client, opt)
	if err != nil {
		return diag.Errorf("error binding policy to the API: %s", err)
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

func unbindPolicy(client *golangsdk.ServiceClient, instanceId, policyId string, unbindSet *schema.Set) error {
	opt := buildListOpts(instanceId, policyId)
	resp, err := throttles.ListBind(client, opt)
	if err != nil {
		return fmt.Errorf("error getting API information from server: %s", err)
	}

	for _, rm := range unbindSet.List() {
		for _, api := range resp {
			// If the publish ID is not found, it means the policy has been unbound from the API by other ways.
			if rm == api.PublishId {
				err = throttles.Unbind(client, instanceId, api.ThrottleApplyId)
				if err != nil {
					return fmt.Errorf("error unbound policy from the API: %s", err)
				}
				break
			}
		}
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
		err = unbindPolicy(client, instanceId, policyId, rmSet)
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
		_, err = throttles.Bind(client, opt)
		if err != nil {
			return diag.Errorf("error binding policy to the API: %v", err)
		}
	}

	return resourceThrottlingPolicyAssociateRead(ctx, d, meta)
}

func resourceThrottlingPolicyAssociateDelete(_ context.Context, d *schema.ResourceData,
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
	)

	return diag.FromErr(unbindPolicy(client, instanceId, policyId, publishIds))
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
