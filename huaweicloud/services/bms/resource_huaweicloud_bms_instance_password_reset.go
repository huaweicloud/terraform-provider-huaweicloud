package bms

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

var instancePasswordResetNonUpdatableParams = []string{"server_id", "new_password"}

// @API BMS PUT /v1/{project_id}/baremetalservers/{server_id}/os-reset-password
func ResourceInstancePasswordReset() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceInstancePasswordResetCreate,
		ReadContext:   resourceInstancePasswordResetRead,
		UpdateContext: resourceInstancePasswordResetUpdate,
		DeleteContext: resourceInstancePasswordResetDelete,

		CustomizeDiff: config.FlexibleForceNew(instancePasswordResetNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"server_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"new_password": {
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

func resourceInstancePasswordResetCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v1/{project_id}/baremetalservers/{server_id}/os-reset-password"
		product = "bms"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating BMS client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{server_id}", d.Get("server_id").(string))

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes:          []int{204},
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreateInstancePasswordResetBodyParams(d))
	_, err = client.Request("PUT", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating BMS instance password reset: %s", err)
	}

	d.SetId(d.Get("server_id").(string))

	return nil
}

func buildCreateInstancePasswordResetBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"new_password": d.Get("new_password").(string),
	}
	return map[string]interface{}{
		"reset-password": bodyParams,
	}
}

func resourceInstancePasswordResetRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceInstancePasswordResetUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceInstancePasswordResetDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting BMS instance password reset resource is not supported. The resource is only removed from state."
	return diag.Diagnostics{
		{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
