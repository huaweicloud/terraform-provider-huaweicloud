package rds

import (
	"context"
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

// @API RDS PUT /v3/{project_id}/instances/recycle-policy
// @API RDS GET /v3/{project_id}/instances/recycle-policy
func ResourceRecyclingPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRecyclingPolicyCreateOrUpdate,
		UpdateContext: resourceRecyclingPolicyCreateOrUpdate,
		ReadContext:   resourceRecyclingPolicyRead,
		DeleteContext: resourceRecyclingPolicyDelete,
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
			"retention_period_in_days": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(1),
				Description:  `Specifies the period of retaining deleted DB instances.`,
			},
		},
	}
}

func resourceRecyclingPolicyCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instances/recycle-policy"
		product = "rds"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RDS Client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createOpt.JSONBody = buildCreateRdsRecyclingPolicyBodyParams(d)

	createResp, err := client.Request("PUT", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating RDS recycling policy: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	result := utils.PathSearch("result", createRespBody, nil)
	if result == nil {
		return diag.Errorf("error creating RDS recycling policy: result is not found in API response")
	}
	if result != "success" {
		return diag.Errorf("error creating RDS recycling policy: result is: %s", result)
	}

	d.SetId(client.ProjectID)

	return resourceRecyclingPolicyRead(ctx, d, meta)
}

func buildCreateRdsRecyclingPolicyBodyParams(d *schema.ResourceData) map[string]interface{} {
	param := make(map[string]interface{})
	if v, ok := d.GetOk("retention_period_in_days"); ok {
		param["retention_period_in_days"] = v
	}
	bodyParams := map[string]interface{}{
		"recycle_policy": param,
	}
	return bodyParams
}

func resourceRecyclingPolicyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v3/{project_id}/instances/recycle-policy"
		product = "rds"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RDS Client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getResp, err := client.Request("GET", getPath, &getOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving RDS recycling policy")
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}
	retentionPeriod := utils.PathSearch("retention_period_in_days", getRespBody, nil)
	if retentionPeriod == nil {
		return diag.Errorf("error retrieving RDS recycling policy, retention_period_in_days is not found")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("retention_period_in_days", retentionPeriod),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceRecyclingPolicyDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting recycling policy is not supported. The recycling policy is only removed from the state," +
		" but it remains in the cloud. And the recycling policy doesn't return to the original state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
