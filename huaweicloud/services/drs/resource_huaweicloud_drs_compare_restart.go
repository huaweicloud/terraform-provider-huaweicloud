package drs

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var drsCompareRestartNonUpdatableParams = []string{
	"job_id",
	"compare_job_ids",
}

// @API DRS POST /v5/{project_id}/jobs/{job_id}/compare/restart
func ResourceDrsCompareRestart() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDrsCompareRestartCreate,
		ReadContext:   resourceDrsCompareRestartRead,
		UpdateContext: resourceDrsCompareRestartUpdate,
		DeleteContext: resourceDrsCompareRestartDelete,

		CustomizeDiff: config.FlexibleForceNew(drsCompareRestartNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"job_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Specifies the ID of the DRS job to restart compare tasks.",
			},
			"compare_job_ids": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Specifies the list of compare job IDs to restart.",
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func resourceDrsCompareRestartCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		httpUrl = "v5/{project_id}/jobs/{job_id}/compare/restart"
		product = "drs"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{job_id}", d.Get("job_id").(string))

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreateDrsCompareRestartBodyParams(d))

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error restarting DRS compare jobs: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("request_id", createRespBody, "").(string)
	if id == "" {
		return diag.Errorf("error restarting DRS compare jobs: request_id is not found in API response")
	}
	d.SetId(id)

	return resourceDrsCompareRestartRead(ctx, d, meta)
}

func buildCreateDrsCompareRestartBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"compare_job_ids": utils.ValueIgnoreEmpty(d.Get("compare_job_ids").([]interface{})),
	}
	return bodyParams
}

func resourceDrsCompareRestartRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDrsCompareRestartUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDrsCompareRestartDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting DRS compare restart resource is not supported. The resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
