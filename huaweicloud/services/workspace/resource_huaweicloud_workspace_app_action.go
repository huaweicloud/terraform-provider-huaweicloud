package workspace

import (
	"context"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var appActionNonUpdatableParams = []string{"service_status"}

// @API Workspace POST /v1/{project_id}/tenant/action/active
func ResourceAppAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppActionCreate,
		ReadContext:   resourceAppActionRead,
		UpdateContext: resourceAppActionUpdate,
		DeleteContext: resourceAppActionDelete,

		CustomizeDiff: config.FlexibleForceNew(appActionNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"service_status": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The service status of the Workspace APP tenant.`,
			},
			// Internal parameter(s).
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func resourceAppActionCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		httpUrl       = "v1/{project_id}/tenant/action/active"
		serviceStatus = d.Get("service_status").(string)
	)

	client, err := cfg.NewServiceClient("appstream", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	actionPath := client.Endpoint + httpUrl
	actionPath = strings.ReplaceAll(actionPath, "{project_id}", client.ProjectID)
	actionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"service_status": serviceStatus,
		},
	}
	_, err = client.Request("POST", actionPath, &actionOpt)
	if err != nil {
		return diag.Errorf("unable to %s Workspace APP: %s", serviceStatus, err)
	}

	randomId, _ := uuid.GenerateUUID()
	d.SetId(randomId)

	return nil
}

func resourceAppActionRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceAppActionUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceAppActionDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource for activating or deactivating the Workspace APP tenant. Deleting this resource will
not clear the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
