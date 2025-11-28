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

var appApplicationBatchActionNonUpdatableParams = []string{"app_group_id", "action", "application_ids"}

// @API Workspace POST /v1/{project_id}/app-groups/{app_group_id}/apps/actions/enable
// @API Workspace POST /v1/{project_id}/app-groups/{app_group_id}/apps/actions/disable
func ResourceAppApplicationBatchAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppApplicationBatchActionCreate,
		ReadContext:   resourceAppApplicationBatchActionRead,
		UpdateContext: resourceAppApplicationBatchActionUpdate,
		DeleteContext: resourceAppApplicationBatchActionDelete,

		CustomizeDiff: config.FlexibleForceNew(appApplicationBatchActionNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the applications to be operated are located.`,
			},
			"app_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the application group.`,
			},
			"action": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The type of the action.`,
				ValidateFunc: validation.StringInSlice([]string{
					"enable",
					"disable",
				}, false),
			},
			"application_ids": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of application IDs to be operated.`,
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

func resourceAppApplicationBatchActionCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		appGroupId = d.Get("app_group_id").(string)
		action     = d.Get("action").(string)
		httpUrl    = "v1/{project_id}/app-groups/{app_group_id}/apps/actions/{action}"
	)

	client, err := cfg.NewServiceClient("appstream", region)
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{app_group_id}", appGroupId)
	createPath = strings.ReplaceAll(createPath, "{action}", action)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: map[string]interface{}{
			"ids": d.Get("application_ids"),
		},
	}

	_, err = client.Request("POST", createPath, &opt)
	if err != nil {
		return diag.Errorf("unable to %s applications: %s", action, err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	return nil
}

func resourceAppApplicationBatchActionRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceAppApplicationBatchActionUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceAppApplicationBatchActionDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to batch enable or disable applications. Deleting this
resource will not clear the corresponding request record, but will only remove the resource information from the tfstate
file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
