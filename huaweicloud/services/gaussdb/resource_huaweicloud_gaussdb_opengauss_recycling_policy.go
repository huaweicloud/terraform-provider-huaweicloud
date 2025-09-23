package gaussdb

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

// @API GaussDB PUT /v3/{project_id}/recycle-policy
// @API GaussDB GET /v3/{project_id}/recycle-policy
func ResourceOpenGaussRecyclingPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOpenGaussRecyclingPolicyCreateOrUpdate,
		UpdateContext: resourceOpenGaussRecyclingPolicyCreateOrUpdate,
		ReadContext:   resourceOpenGaussRecyclingPolicyRead,
		DeleteContext: resourceOpenGaussRecyclingPolicyDelete,
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

func resourceOpenGaussRecyclingPolicyCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/recycle-policy"
		product = "opengauss"
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
	createOpt.JSONBody = buildCreateOpenGaussRecyclingPolicyBodyParams(d)

	_, err = client.Request("PUT", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating GaussDB OpenGauss recycling policy: %s", err)
	}

	d.SetId(client.ProjectID)

	return resourceOpenGaussRecyclingPolicyRead(ctx, d, meta)
}

func buildCreateOpenGaussRecyclingPolicyBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"recycle_policy": map[string]interface{}{
			"retention_period_in_days": d.Get("retention_period_in_days"),
		},
	}
	return bodyParams
}

func resourceOpenGaussRecyclingPolicyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v3/{project_id}/recycle-policy"
		product = "opengauss"
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
		return common.CheckDeletedDiag(d, err, "error retrieving GaussDB OpenGauss recycling policy")
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}
	retentionPeriod := utils.PathSearch("retention_period_in_days", getRespBody, nil)
	if retentionPeriod == nil {
		return diag.Errorf("error retrieving GaussDB OpenGauss recycling policy, retention_period_in_days is " +
			"not found")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("retention_period_in_days", retentionPeriod),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceOpenGaussRecyclingPolicyDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting recycling policy is not supported. The recycling policy is only removed from the state," +
		" but it remains in the cloud. And the recycling policy doesn't return to the original state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
