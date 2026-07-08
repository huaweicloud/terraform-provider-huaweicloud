package taurusdb

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var taurusDBProxyUpgradeNonUpdatableParams = []string{
	"instance_id", "proxy_id", "source_version", "target_version",
}

// @API TaurusDB POST /v3/{project_id}/instances/{instance_id}/proxy/{proxy_id}/upgrade-version
// @API TaurusDB GET /v3/{project_id}/jobs
func ResourceTaurusDBProxyUpgrade() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTaurusDBProxyUpgradeCreate,
		ReadContext:   resourceTaurusDBProxyUpgradeRead,
		UpdateContext: resourceTaurusDBProxyUpgradeUpdate,
		DeleteContext: resourceTaurusDBProxyUpgradeDelete,

		CustomizeDiff: config.FlexibleForceNew(taurusDBProxyUpgradeNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"proxy_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"source_version": {
				Type:     schema.TypeString,
				Required: true,
			},
			"target_version": {
				Type:     schema.TypeString,
				Required: true,
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

func resourceTaurusDBProxyUpgradeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceId = d.Get("instance_id").(string)
		proxyId    = d.Get("proxy_id").(string)
		httpUrl    = "v3/{project_id}/instances/{instance_id}/proxy/{proxy_id}/upgrade-version"
	)

	client, err := cfg.NewServiceClient("gaussdb", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)
	createPath = strings.ReplaceAll(createPath, "{proxy_id}", proxyId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(buildCreateInstanceProxyUpgradeBodyParams(d)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error upgrading TaurusDB(%s) Proxy(%s) version: %s", instanceId, proxyId, err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	// Check the update_result
	state := utils.PathSearch("update_result|[0].state", createRespBody, "").(string)
	errorMessage := utils.PathSearch("update_result|[0].error_message", createRespBody, "").(string)
	if state != "ACCEPT" || errorMessage != "" {
		return diag.Errorf("error upgrading TaurusDB(%s) Proxy(%s) version: %s", instanceId, proxyId, errorMessage)
	}

	// Get workflow_id for job tracking when task accepted
	workflowId := utils.PathSearch("update_result|[0].workflow_id", createRespBody, "").(string)
	if workflowId == "" {
		return diag.Errorf("error upgrading TaurusDB(%s) Proxy(%s) version: workflow_id is not found in the response",
			instanceId, proxyId)
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("error generating UUID: %s", err)
	}
	d.SetId(id.String())
	if err = checkGaussDBMySQLJobFinish(ctx, client, workflowId, d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.Errorf("error waiting for TaurusDB Proxy upgrade job to complete: %s", err)
	}
	return nil
}

func buildCreateInstanceProxyUpgradeBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"source_version": d.Get("source_version"),
		"target_version": d.Get("target_version"),
	}
	return bodyParams
}

func resourceTaurusDBProxyUpgradeRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceTaurusDBProxyUpgradeUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceTaurusDBProxyUpgradeDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting TaurusDB instance proxy upgrade resource is not supported. The resource is only removed from the state," +
		" the TaurusDB instance proxy remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
