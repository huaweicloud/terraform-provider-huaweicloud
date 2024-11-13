package iotda

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iotda/v5/model"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IoTDA POST /v5/iot/{project_id}/routing-rule/flowcontrol-policy
// @API IoTDA GET /v5/iot/{project_id}/routing-rule/flowcontrol-policy/{policy_id}
// @API IoTDA PUT /v5/iot/{project_id}/routing-rule/flowcontrol-policy/{policy_id}
// @API IoTDA DELETE /v5/iot/{project_id}/routing-rule/flowcontrol-policy/{policy_id}
func ResourceDataFlowControlPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDataFlowControlPolicyCreate,
		ReadContext:   resourceDataFlowControlPolicyRead,
		UpdateContext: resourceDataFlowControlPolicyUpdate,
		DeleteContext: resourceDataFlowControlPolicyDelete,

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
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"scope": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"scope_value": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"limit": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func buildDataFlowControlPolicyCreateParams(d *schema.ResourceData) *model.CreateRoutingFlowControlPolicyRequest {
	createOptsBody := model.AddFlowControlPolicy{
		PolicyName:  utils.StringIgnoreEmpty(d.Get("name").(string)),
		Description: utils.StringIgnoreEmpty(d.Get("description").(string)),
		Scope:       utils.StringIgnoreEmpty(d.Get("scope").(string)),
		ScopeValue:  utils.StringIgnoreEmpty(d.Get("scope_value").(string)),
		//nolint:gosec
		Limit: utils.Int32IgnoreEmpty(int32(d.Get("limit").(int))),
	}

	return &model.CreateRoutingFlowControlPolicyRequest{
		Body: &createOptsBody,
	}
}

func resourceDataFlowControlPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		isDerived = WithDerivedAuth(cfg, region)
	)
	client, err := cfg.HcIoTdaV5Client(region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA v5 client: %s", err)
	}

	resp, err := client.CreateRoutingFlowControlPolicy(buildDataFlowControlPolicyCreateParams(d))
	if err != nil {
		return diag.Errorf("error creating IoTDA data flow control policy: %s", err)
	}

	if resp == nil || resp.PolicyId == nil {
		return diag.Errorf("error creating IoTDA data flow control policy: ID is not found in API response")
	}

	d.SetId(*resp.PolicyId)

	return resourceDataFlowControlPolicyRead(ctx, d, meta)
}

func buildDataFlowControlPolicyQueryParams(d *schema.ResourceData) *model.ShowRoutingFlowControlPolicyRequest {
	return &model.ShowRoutingFlowControlPolicyRequest{
		PolicyId: d.Id(),
	}
}

func resourceDataFlowControlPolicyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		isDerived = WithDerivedAuth(cfg, region)
	)
	client, err := cfg.HcIoTdaV5Client(region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA v5 client: %s", err)
	}

	resp, err := client.ShowRoutingFlowControlPolicy(buildDataFlowControlPolicyQueryParams(d))
	// When the resource does not exist, query API will return `404` error code.
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving IoTDA data flow control policy")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", resp.PolicyName),
		d.Set("description", resp.Description),
		d.Set("scope", resp.Scope),
		d.Set("scope_value", resp.ScopeValue),
		d.Set("limit", resp.Limit),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildDataFlowControlPolicyUpdateParams(d *schema.ResourceData) *model.UpdateRoutingFlowControlPolicyRequest {
	updateOptsBody := model.UpdateFlowControlPolicy{
		PolicyName:  utils.StringIgnoreEmpty(d.Get("name").(string)),
		Description: utils.StringIgnoreEmpty(d.Get("description").(string)),
		//nolint:gosec
		Limit: utils.Int32IgnoreEmpty(int32(d.Get("limit").(int))),
	}

	return &model.UpdateRoutingFlowControlPolicyRequest{
		PolicyId: d.Id(),
		Body:     &updateOptsBody,
	}
}

func resourceDataFlowControlPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		isDerived = WithDerivedAuth(cfg, region)
	)
	client, err := cfg.HcIoTdaV5Client(region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA v5 client: %s", err)
	}

	_, err = client.UpdateRoutingFlowControlPolicy(buildDataFlowControlPolicyUpdateParams(d))
	if err != nil {
		return diag.Errorf("error updating IoTDA data flow control policy: %s", err)
	}

	return resourceDataFlowControlPolicyRead(ctx, d, meta)
}

func buildDataFlowControlPolicyDeleteParams(d *schema.ResourceData) *model.DeleteRoutingFlowControlPolicyRequest {
	return &model.DeleteRoutingFlowControlPolicyRequest{
		PolicyId: d.Id(),
	}
}

func resourceDataFlowControlPolicyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		isDerived = WithDerivedAuth(cfg, region)
	)
	client, err := cfg.HcIoTdaV5Client(region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA v5 client: %s", err)
	}

	_, err = client.DeleteRoutingFlowControlPolicy(buildDataFlowControlPolicyDeleteParams(d))
	// When the resource does not exist, delete API will return `404` error code.
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting IoTDA data flow control policy")
	}

	return nil
}
