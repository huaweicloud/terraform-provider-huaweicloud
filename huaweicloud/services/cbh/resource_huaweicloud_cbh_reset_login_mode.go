package cbh

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

// @API CBH PUT /v2/{project_id}/cbs/instance/login-method
func ResourceResetLoginMode() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceResetLoginModeCreate,
		UpdateContext: resourceResetLoginModeUpdate,
		ReadContext:   resourceResetLoginModeRead,
		DeleteContext: resourceResetLoginModeDelete,

		CustomizeDiff: config.FlexibleForceNew([]string{
			"server_id",
		}),

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
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func resourceResetLoginModeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		httpUrl  = "v2/{project_id}/cbs/instance/login-method"
		product  = "cbh"
		serverId = d.Get("server_id").(string)
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CBH client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"server_id": serverId,
		},
	}

	_, err = client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error resetting CBH login mode: %s", err)
	}

	d.SetId(serverId)

	return resourceResetLoginModeRead(ctx, d, meta)
}

func resourceResetLoginModeRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceResetLoginModeUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceResetLoginModeDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource using to reset login mode.
Deleting this resource will not recover the reset login mode, but will only remove the resource information from
the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
