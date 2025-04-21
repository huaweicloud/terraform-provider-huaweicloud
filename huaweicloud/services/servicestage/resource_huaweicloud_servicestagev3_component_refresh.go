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

var v3ComponentRefreshNonUpdatableParams = []string{
	"application_id",
	"component_id",
}

// @API ServiceStage PUT /v3/{project_id}/cas/applications/{application_id}/components/{component_id}/refresh
func ResourceV3ComponentRefresh() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV3ComponentRefreshCreate,
		ReadContext:   resourceV3ComponentRefreshRead,
		UpdateContext: resourceV3ComponentRefreshUpdate,
		DeleteContext: resourceV3ComponentRefreshDelete,

		CustomizeDiff: config.FlexibleForceNew(v3ComponentRefreshNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the component to be refreshed is located.`,
			},
			// Required parameters.
			"application_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The application ID to which the componnet belongs.`,
			},
			"component_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the component to be refreshed.`,
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

func resourceV3ComponentRefreshCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		httpUrl       = "v3/{project_id}/cas/applications/{application_id}/components/{component_id}/refresh"
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
			"Content-Type": "application/json;charset=utf8",
		},
	}

	_, err = client.Request("PUT", createPath, &opt)
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

func resourceV3ComponentRefreshRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceV3ComponentRefreshUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceV3ComponentRefreshDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for doing component refresh. Deleting this resource
will not clear the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
