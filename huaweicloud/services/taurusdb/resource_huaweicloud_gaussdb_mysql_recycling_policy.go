package taurusdb

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

// @API GaussDBforMySQL PUT /v3/{project_id}/instances/recycle-policy
// @API GaussDBforMySQL GET /v3/{project_id}/instances/recycle-policy
func ResourceGaussDBRecyclingPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGaussDBRecyclingPolicyCreateOrUpdate,
		UpdateContext: resourceGaussDBRecyclingPolicyCreateOrUpdate,
		ReadContext:   resourceGaussDBMysqlRecyclingPolicyRead,
		DeleteContext: resourceGaussDBMysqlRecyclingPolicyDelete,
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
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceGaussDBRecyclingPolicyCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instances/recycle-policy"
		product = "gaussdb"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createOpt.JSONBody = buildCreateGaussDBRecyclingPolicyBodyParams(d)

	_, err = client.Request("PUT", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating GaussDB recycling policy: %s", err)
	}

	d.SetId(client.ProjectID)

	return resourceGaussDBMysqlRecyclingPolicyRead(ctx, d, meta)
}

func buildCreateGaussDBRecyclingPolicyBodyParams(d *schema.ResourceData) map[string]interface{} {
	param := make(map[string]interface{})
	if v, ok := d.GetOk("retention_period_in_days"); ok {
		param["retention_period_in_days"] = v
	}
	bodyParams := map[string]interface{}{
		"recycle_policy": param,
	}
	return bodyParams
}

func resourceGaussDBMysqlRecyclingPolicyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v3/{project_id}/instances/recycle-policy"
		product = "gaussdb"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
	}
	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving GaussDB recycling policy")
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}
	retentionPeriod := utils.PathSearch("recycle_policy.retention_period_in_days", getRespBody, nil)
	if retentionPeriod == nil {
		return diag.Errorf("error retrieving GaussDB recycling policy, retention_period_in_days is not found")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("retention_period_in_days", retentionPeriod),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceGaussDBMysqlRecyclingPolicyDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting recycling policy is not supported. The recycling policy is only removed from the state," +
		" but it remains in the cloud. And the recycling policy doesn't return to the original state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
