package servicestage

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

var v3ComponentActionNonUpdatableParams = []string{
	"application_id",
	"component_id",
	"action",
	"parameters",
}

// @API ServiceStage POST /v3/{project_id}/cas/applications/{application_id}/components/{component_id}/action
func ResourceV3ComponentAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceComponentActionCreate,
		ReadContext:   resourceComponentActionRead,
		UpdateContext: resourceComponentActionUpdate,
		DeleteContext: resourceComponentActionDelete,

		CustomizeDiff: config.FlexibleForceNew(v3ComponentActionNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the component is located.`,
			},
			"application_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The application ID to which the componnet belongs.`,
			},
			"component_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the component to be operated.`,
			},
			"action": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The action type of the component execution.`,
			},
			"parameters": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsJSON,
				Description:  `The job parameters.`,
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

func buildV3ComponentActionCreateBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"action":     d.Get("action").(string),
		"parameters": utils.StringToJson(d.Get("parameters").(string)),
	}
}

func resourceComponentActionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		httpUrl       = "v3/{project_id}/cas/applications/{application_id}/components/{component_id}/action"
		applicationId = d.Get("application_id").(string)
		componentId   = d.Get("component_id").(string)
	)
	client, err := cfg.NewServiceClient("servicestage", region)
	if err != nil {
		return diag.Errorf("error creating ServiceStage client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{application_id}", applicationId)
	createPath = strings.ReplaceAll(createPath, "{component_id}", componentId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildV3ComponentActionCreateBodyParams(d)),
	}

	_, err = client.Request("POST", createPath, &opt)
	if err != nil {
		return diag.Errorf("error executing component action: %s", err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	return resourceComponentActionRead(ctx, d, meta)
}

func resourceComponentActionRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceComponentActionUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceComponentActionDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for operating the component. Deleting this resource will
not clear the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
