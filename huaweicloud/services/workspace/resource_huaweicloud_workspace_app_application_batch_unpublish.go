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

var appApplicationBatchUnpublishNonUpdatableParams = []string{
	"app_group_id",
	"application_ids",
}

// @API Workspace POST /v1/{project_id}/app-groups/{app_group_id}/apps/batch-unpublish
func ResourceAppApplicationBatchUnpublish() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAppApplicationBatchUnpublishCreate,
		ReadContext:   resourceAppApplicationBatchUnpublishRead,
		UpdateContext: resourceAppApplicationBatchUnpublishUpdate,
		DeleteContext: resourceAppApplicationBatchUnpublishDelete,

		CustomizeDiff: config.FlexibleForceNew(appApplicationBatchUnpublishNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the applications to be unpublished are located.`,
			},
			"app_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the application group.`,
			},
			"application_ids": {
				Type:        schema.TypeList,
				Required:    true,
				Description: `The list of application IDs to be unpublished.`,
				Elem:        &schema.Schema{Type: schema.TypeString},
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

func resourceAppApplicationBatchUnpublishCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                = meta.(*config.Config)
		region             = cfg.GetRegion(d)
		applicationGroupId = d.Get("app_group_id").(string)
		httpUrl            = "v1/{project_id}/app-groups/{app_group_id}/apps/batch-unpublish"
	)

	client, err := cfg.NewServiceClient("appstream", region)
	if err != nil {
		return diag.Errorf("error creating Workspace APP client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{app_group_id}", applicationGroupId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
		JSONBody: map[string]interface{}{
			"ids": d.Get("application_ids"),
		},
	}

	_, err = client.Request("POST", createPath, &opt)
	if err != nil {
		return diag.Errorf("error batch unpublishing applications under application group (%s): %s", applicationGroupId, err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	return nil
}

func resourceAppApplicationBatchUnpublishRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceAppApplicationBatchUnpublishUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceAppApplicationBatchUnpublishDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to batch unpublish applications. Deleting this resource will
not clear the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
