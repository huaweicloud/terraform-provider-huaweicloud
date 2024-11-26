package iotda

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iotda/v5/model"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IoTDA POST /v5/iot/{project_id}/routing-rule/backlog-policy
// @API IoTDA GET /v5/iot/{project_id}/routing-rule/backlog-policy/{policy_id}
// @API IoTDA PUT /v5/iot/{project_id}/routing-rule/backlog-policy/{policy_id}
// @API IoTDA DELETE /v5/iot/{project_id}/routing-rule/backlog-policy/{policy_id}
func ResourceDataBacklogPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDataBacklogPolicyCreate,
		ReadContext:   resourceDataBacklogPolicyRead,
		UpdateContext: resourceDataBacklogPolicyUpdate,
		DeleteContext: resourceDataBacklogPolicyDelete,

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
			"backlog_size": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"backlog_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func convertStringValueToInt32(value string) *int32 {
	if value == "0" {
		return utils.Int32(0)
	}

	parsedValue := utils.StringToInt(&value)
	if parsedValue != nil {
		//nolint:gosec
		return utils.Int32IgnoreEmpty(int32(*parsedValue))
	}

	return nil
}

func buildDataBacklogPolicyCreateParams(d *schema.ResourceData) *model.CreateRoutingBacklogPolicyRequest {
	createOptsBody := model.AddBacklogPolicy{
		PolicyName:  utils.StringIgnoreEmpty(d.Get("name").(string)),
		Description: utils.StringIgnoreEmpty(d.Get("description").(string)),
		BacklogSize: convertStringValueToInt32(d.Get("backlog_size").(string)),
		BacklogTime: convertStringValueToInt32(d.Get("backlog_time").(string)),
	}

	return &model.CreateRoutingBacklogPolicyRequest{
		Body: &createOptsBody,
	}
}

func resourceDataBacklogPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		isDerived = WithDerivedAuth(cfg, region)
	)
	client, err := cfg.HcIoTdaV5Client(region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA v5 client: %s", err)
	}

	resp, err := client.CreateRoutingBacklogPolicy(buildDataBacklogPolicyCreateParams(d))
	if err != nil {
		return diag.Errorf("error creating IoTDA data backlog policy: %s", err)
	}

	if resp == nil || resp.PolicyId == nil {
		return diag.Errorf("error creating IoTDA data backlog policy: ID is not found in API response")
	}

	d.SetId(*resp.PolicyId)

	return resourceDataBacklogPolicyRead(ctx, d, meta)
}

func buildDataBacklogPolicyQueryParams(d *schema.ResourceData) *model.ShowRoutingBacklogPolicyRequest {
	return &model.ShowRoutingBacklogPolicyRequest{
		PolicyId: d.Id(),
	}
}

func resourceDataBacklogPolicyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		isDerived = WithDerivedAuth(cfg, region)
	)
	client, err := cfg.HcIoTdaV5Client(region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA v5 client: %s", err)
	}

	resp, err := client.ShowRoutingBacklogPolicy(buildDataBacklogPolicyQueryParams(d))
	// When the resource does not exist, query API will return `404` error code.
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving IoTDA data backlog policy")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", resp.PolicyName),
		d.Set("description", resp.Description),
		d.Set("backlog_size", parseBacklogValueToString(resp.BacklogSize)),
		d.Set("backlog_time", parseBacklogValueToString(resp.BacklogTime)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func parseBacklogValueToString(value *int32) string {
	if value != nil {
		return fmt.Sprintf("%v", *value)
	}

	return ""
}

func buildDataBacklogPolicyUpdateParams(d *schema.ResourceData) *model.UpdateRoutingBacklogPolicyRequest {
	updateOptsBody := model.UpdateBacklogPolicy{
		PolicyName:  utils.StringIgnoreEmpty(d.Get("name").(string)),
		Description: utils.StringIgnoreEmpty(d.Get("description").(string)),
		BacklogSize: convertStringValueToInt32(d.Get("backlog_size").(string)),
		BacklogTime: convertStringValueToInt32(d.Get("backlog_time").(string)),
	}

	return &model.UpdateRoutingBacklogPolicyRequest{
		PolicyId: d.Id(),
		Body:     &updateOptsBody,
	}
}

func resourceDataBacklogPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		isDerived = WithDerivedAuth(cfg, region)
	)
	client, err := cfg.HcIoTdaV5Client(region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA v5 client: %s", err)
	}

	_, err = client.UpdateRoutingBacklogPolicy(buildDataBacklogPolicyUpdateParams(d))
	if err != nil {
		return diag.Errorf("error updating IoTDA data backlog policy: %s", err)
	}

	return resourceDataBacklogPolicyRead(ctx, d, meta)
}

func buildDataBacklogPolicyDeleteParams(d *schema.ResourceData) *model.DeleteRoutingBacklogPolicyRequest {
	return &model.DeleteRoutingBacklogPolicyRequest{
		PolicyId: d.Id(),
	}
}

func resourceDataBacklogPolicyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		isDerived = WithDerivedAuth(cfg, region)
	)
	client, err := cfg.HcIoTdaV5Client(region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA v5 client: %s", err)
	}

	_, err = client.DeleteRoutingBacklogPolicy(buildDataBacklogPolicyDeleteParams(d))
	// When the resource does not exist, delete API will return `404` error code.
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting IoTDA data backlog policy")
	}

	return nil
}
