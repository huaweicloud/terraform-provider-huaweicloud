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

var gaussDbInstancePasswordResetNonUpdatableParams = []string{
	"instance_id",
	"password",
}

// @API GaussDB POST /v3/{project_id}/instances/{instance_id}/password
func ResourceGaussDbInstancePasswordReset() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGaussDbInstancePasswordResetCreate,
		ReadContext:   resourceGaussDbInstancePasswordResetRead,
		UpdateContext: resourceGaussDbInstancePasswordResetUpdate,
		DeleteContext: resourceGaussDbInstancePasswordResetDelete,

		CustomizeDiff: config.FlexibleForceNew(gaussDbInstancePasswordResetNonUpdatableParams),

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
			"password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
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

func resourceGaussDbInstancePasswordResetCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/password"
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
	createOpt.JSONBody = utils.RemoveNil(buildGaussDbInstancePasswordResetBodyParams(d))

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error resetting GaussDB instance password: %s", err)
	}

	d.SetId(d.Get("instance_id").(string))

	return nil
}

func buildGaussDbInstancePasswordResetBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"password": d.Get("password"),
	}
	return bodyParams
}

func resourceGaussDbInstancePasswordResetRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceGaussDbInstancePasswordResetUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceGaussDbInstancePasswordResetDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting GaussDB instance password reset resource is not supported. The resource is only removed from" +
		" the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
