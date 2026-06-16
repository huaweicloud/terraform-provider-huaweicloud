package dcs

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

var dcsInstanceMinorVersionUpgradeNonUpdatableParams = []string{"instance_id", "engine_minor_version", "proxy_minor_version"}

// @API DCS POST /v2/{project_id}/instances/{instance_id}/minor-version/upgrade
func ResourceInstanceMinorVersionUpgrade() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceInstanceMinorVersionUpgradeCreate,
		ReadContext:   resourceInstanceMinorVersionUpgradeRead,
		UpdateContext: resourceInstanceMinorVersionUpgradeUpdate,
		DeleteContext: resourceInstanceMinorVersionUpgradeDelete,

		CustomizeDiff: config.FlexibleForceNew(dcsInstanceMinorVersionUpgradeNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"engine_minor_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"proxy_minor_version": {
				Type:     schema.TypeString,
				Optional: true,
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

func resourceInstanceMinorVersionUpgradeCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v2/{project_id}/instances/{instance_id}/minor-version/upgrade"
		product = "dcs"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DCS client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreateInstanceMinorVersionUpgradeBodyParams(d))

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating DCS instance minor version upgrade: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	taskID := utils.PathSearch("task_id", createRespBody, "").(string)
	if taskID == "" {
		return diag.Errorf("task_id not found in upgrade response")
	}

	generateUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID.String())

	return nil
}

func resourceInstanceMinorVersionUpgradeRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceInstanceMinorVersionUpgradeUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceInstanceMinorVersionUpgradeDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting DCS instance minor version upgrade resource is not supported. The resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

func buildCreateInstanceMinorVersionUpgradeBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"engine_minor_version": utils.ValueIgnoreEmpty(d.Get("engine_minor_version")),
		"proxy_minor_version":  utils.ValueIgnoreEmpty(d.Get("proxy_minor_version")),
	}
}
