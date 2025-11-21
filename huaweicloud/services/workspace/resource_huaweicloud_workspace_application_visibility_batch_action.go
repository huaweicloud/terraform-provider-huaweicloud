package workspace

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var applicationVisibilityBatchActionNonUpdatableParams = []string{"action", "app_ids"}

// @API Workspace POST /v1/{project_id}/app-center/apps/actions/batch-enable
// @API Workspace POST /v1/{project_id}/app-center/apps/actions/batch-disable
func ResourceApplicationVisibilityBatchAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceApplicationVisibilityBatchActionCreate,
		ReadContext:   resourceApplicationVisibilityBatchActionRead,
		UpdateContext: resourceApplicationVisibilityBatchActionUpdate,
		DeleteContext: resourceApplicationVisibilityBatchActionDelete,

		CustomizeDiff: config.FlexibleForceNew(applicationVisibilityBatchActionNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the applications to be operated are located.`,
			},

			// Required parameters.
			"action": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The visibility action type.`,
			},
			"app_ids": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of application IDs to be operated.`,
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func buildApplicationVisibilityBatchActionBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"items": d.Get("app_ids"),
	}
}

func resourceApplicationVisibilityBatchActionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		action  = d.Get("action").(string)
		httpUrl = fmt.Sprintf("v1/{project_id}/app-center/apps/actions/batch-%s", action)
	)

	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: buildApplicationVisibilityBatchActionBodyParams(d),
	}

	_, err = client.Request("POST", createPath, &opt)
	if err != nil {
		return diag.Errorf("error executing application visibility batch action (%s): %s", action, err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	return resourceApplicationVisibilityBatchActionRead(ctx, d, meta)
}

func resourceApplicationVisibilityBatchActionRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceApplicationVisibilityBatchActionUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceApplicationVisibilityBatchActionDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to batch set application visibility. Deleting this
resource will not clear the corresponding request record, but will only remove the resource information from the tfstate
file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
