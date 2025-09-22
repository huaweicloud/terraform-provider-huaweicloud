package ecs

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ECS PUT /v1/{project_id}/recycle-bin
// @API ECS PUT /v1/{project_id}/recycle-bin/policy
// @API ECS GET /v1/{project_id}/recycle-bin
func ResourceComputeRecyclePolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceComputeRecyclePolicyCreate,
		ReadContext:   resourceComputeRecyclePolicyRead,
		UpdateContext: resourceComputeRecyclePolicyUpdate,
		DeleteContext: resourceComputeRecyclePolicyDelete,
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
			"retention_hour": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"recycle_threshold_day": {
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

func resourceComputeRecyclePolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		product = "ecs"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ECS client: %s", err)
	}

	err = updateRecycleSwitch(client, "on")
	if err != nil {
		return diag.Errorf("error opening ECS recycle switch: %s", err)
	}

	d.SetId(client.ProjectID)

	err = updateRecyclePolicy(client, d)
	if err != nil {
		return diag.Errorf("error updating ECS recycle policy: %s", err)
	}

	return resourceComputeRecyclePolicyRead(ctx, d, meta)
}

func resourceComputeRecyclePolicyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v1/{project_id}/recycle-bin"
		product = "ecs"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ECS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving ECS recycle")
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	switchValue := utils.PathSearch("switch", getRespBody, "").(string)
	if switchValue == "off" {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving ECS recycle")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("retention_hour", utils.PathSearch("policy.retention_hour", getRespBody, nil)),
		d.Set("recycle_threshold_day", utils.PathSearch("policy.recycle_threshold_day", getRespBody, nil)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceComputeRecyclePolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		product = "ecs"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ECS client: %s", err)
	}

	if d.HasChanges("retention_hour", "recycle_threshold_day") {
		err = updateRecyclePolicy(client, d)
		if err != nil {
			return diag.Errorf("error updating ECS recycle policy: %s", err)
		}
	}
	return resourceComputeRecyclePolicyRead(ctx, d, meta)
}

func resourceComputeRecyclePolicyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		product = "ecs"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ECS client: %s", err)
	}

	err = updateRecycleSwitch(client, "off")
	if err != nil {
		return diag.Errorf("error closing ECS recycle switch: %s", err)
	}

	return nil
}

func updateRecycleSwitch(client *golangsdk.ServiceClient, switchValue string) error {
	httpUrl := "v1/{project_id}/recycle-bin"

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes:          []int{204},
	}
	updateOpt.JSONBody = buildRecycleSwitchBodyParams(switchValue)

	_, err := client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return err
	}
	return nil
}

func buildRecycleSwitchBodyParams(switchValue string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"switch": switchValue,
	}
	return map[string]interface{}{
		"recycle_bin": bodyParams,
	}
}

func updateRecyclePolicy(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	httpUrl := "v1/{project_id}/recycle-bin/policy"

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes:          []int{204},
	}
	updateOpt.JSONBody = buildRecyclePolicyBodyParams(d)

	_, err := client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return err
	}
	return nil
}

func buildRecyclePolicyBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"retention_hour":        d.Get("retention_hour"),
		"recycle_threshold_day": d.Get("recycle_threshold_day"),
	}
	return map[string]interface{}{
		"policy": bodyParams,
	}
}
