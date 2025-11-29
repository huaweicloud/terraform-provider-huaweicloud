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

var userActionNonUpdatableParams = []string{
	"user_id",
	"op_type",
}

// @API Workspace POST /v2/{project_id}/users/{user_id}/actions
func ResourceUserAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserActionCreate,
		ReadContext:   resourceUserActionRead,
		UpdateContext: resourceUserActionUpdate,
		DeleteContext: resourceUserActionDelete,

		CustomizeDiff: config.FlexibleForceNew(userActionNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the user is located.`,
			},

			// Required parameters.
			"user_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the user to be operated.`,
			},
			"op_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The operation type.`,
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

func buildUserActionBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"op_type": d.Get("op_type"),
	}
}

func resourceUserActionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	var (
		httpUrl = "v2/{project_id}/users/{user_id}/actions"
		userId  = d.Get("user_id").(string)
	)

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{user_id}", userId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: buildUserActionBodyParams(d),
		OkCodes:  []int{204},
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error operating Workspace user: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	return resourceUserActionRead(ctx, d, meta)
}

func resourceUserActionRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceUserActionUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceUserActionDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource for operating user. Deleting this resource will not clear
the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
