package gaussdb

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

var gaussDbDrConfigurationResetNonUpdatableParams = []string{
	"instance_id",
	"opposite_data_cidr",
}

// @API GaussDB POST /v3.5/{project_id}/instances/{instance_id}/reset-dr-config
func ResourceGaussDbDrConfigurationReset() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGaussDbDrConfigurationResetCreate,
		ReadContext:   resourceGaussDbDrConfigurationResetRead,
		UpdateContext: resourceGaussDbDrConfigurationResetUpdate,
		DeleteContext: resourceGaussDbDrConfigurationResetDelete,

		CustomizeDiff: config.FlexibleForceNew(gaussDbDrConfigurationResetNonUpdatableParams),

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
			"opposite_data_cidr": {
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

func resourceGaussDbDrConfigurationResetCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		httpUrl = "v3.5/{project_id}/instances/{instance_id}/reset-dr-config"
		product = "opengauss"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", d.Get("instance_id").(string))

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	createOpt.JSONBody = utils.RemoveNil(buildGaussDbDrConfigurationResetBodyParams(d))

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error resetting GaussDB DR configuration: %s", err)
	}

	d.SetId(d.Get("instance_id").(string))

	return nil
}

func buildGaussDbDrConfigurationResetBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"opposite_data_cidr": utils.ValueIgnoreEmpty(d.Get("opposite_data_cidr")),
	}
	return bodyParams
}

func resourceGaussDbDrConfigurationResetRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceGaussDbDrConfigurationResetUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceGaussDbDrConfigurationResetDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting GaussDB DR configuration reset resource is not supported. The restart resource is only " +
		"removed from the state, the GaussDB instance remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
