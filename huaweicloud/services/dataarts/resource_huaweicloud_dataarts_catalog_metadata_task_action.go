package dataarts

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var catalogMetadataTaskActionNonUpdatableParams = []string{
	"workspace_id",
	"task_id",
	"action",
}

// @API DataArtsStudio POST /v3/{project_id}/metadata/tasks/{task_id}/action
func ResourceCatalogMetadataTaskAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMetadataTaskActionCreate,
		ReadContext:   resourceMetadataTaskActionRead,
		UpdateContext: resourceMetadataTaskActionUpdate,
		DeleteContext: resourceMetadataTaskActionDelete,

		CustomizeDiff: config.FlexibleForceNew(catalogMetadataTaskActionNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the metadata task is located.`,
			},

			// Required parameters.
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the workspace to which the metadata task belongs.`,
			},
			"task_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The task ID of the metadata task.`,
			},
			"action": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The action of the metadata task status flag.`,
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description: utils.SchemaDesc(
					`Whether to allow parameters that do not support changes to have their change-triggered behavior set to 'ForceNew'.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
		},
	}
}

func actionMetadataTask(client *golangsdk.ServiceClient, workspaceId, taskId, action string) error {
	httpUrl := "v3/{project_id}/metadata/tasks/{task_id}/action"
	path := client.Endpoint + httpUrl
	path = strings.ReplaceAll(path, "{project_id}", client.ProjectID)
	path = strings.ReplaceAll(path, "{task_id}", taskId)
	path = fmt.Sprintf(path+"?action=%s", action)

	metadataTaskActionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      bulidCatalogMoreHeaders(workspaceId),
		OkCodes:          []int{200, 204},
	}

	_, err := client.Request("POST", path, &metadataTaskActionOpt)
	return err
}

func resourceMetadataTaskActionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
		taskId      = d.Get("task_id").(string)
		action      = d.Get("action").(string)
	)

	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	err = actionMetadataTask(client, workspaceId, taskId, action)
	if err != nil {
		return diag.Errorf("error action metadata task (%s): %s", taskId, err)
	}

	resourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(resourceId.String())

	return resourceDataServiceApiActionRead(ctx, d, meta)
}

func resourceMetadataTaskActionRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// This resource is only a one-time action resource for run/runimmediate/stop the API.
	// There is no API for the provider to query the history of this API action.
	return nil
}

func resourceMetadataTaskActionUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceMetadataTaskActionDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	msg := `This resource is only a one-time action resource for operating a metadata task. Deleting this resource will not
clear the corresponding request record, but will only remove the resource information from the tfstate file,
but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  msg,
		},
	}
}
