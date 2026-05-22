package geminidb

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

// @API GeminiDB PUT /v3/{project_id}/instances/recycle-policy
// @API GeminiDB GET /v3/{project_id}/instances/recycle-policy
func ResourceGeminiDBRecyclingPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGeminiDBRecyclingPolicyCreateOrUpdate,
		UpdateContext: resourceGeminiDBRecyclingPolicyCreateOrUpdate,
		ReadContext:   resourceGeminiDBRecyclingPolicyRead,
		DeleteContext: resourceGeminiDBRecyclingPolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"retention_period_in_days": {
				Type:     schema.TypeInt,
				Required: true,
			},
		},
	}
}

func resourceGeminiDBRecyclingPolicyCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instances/recycle-policy"
		product = "geminidb"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         buildCreateGeminiDBRecyclingPolicyBodyParams(d),
	}

	_, err = client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating GeminiDB recycling policy: %s", err)
	}
	d.SetId(client.ProjectID)

	return resourceGeminiDBRecyclingPolicyRead(ctx, d, meta)
}

func buildCreateGeminiDBRecyclingPolicyBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"recycle_policy": map[string]interface{}{
			"retention_period_in_days": d.Get("retention_period_in_days"),
		},
	}
	return bodyParams
}

func resourceGeminiDBRecyclingPolicyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/instances/recycle-policy"
		product = "geminidb"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving GeminiDB recycling policy")
	}

	respBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("retention_period_in_days", utils.PathSearch("recycle_policy.retention_period_in_days", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceGeminiDBRecyclingPolicyDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting recycling policy is not supported. The recycling policy is only removed from the state," +
		" but it remains in the cloud. And the recycling policy doesn't return to the original state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
