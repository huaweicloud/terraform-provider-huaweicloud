package dsc

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var databaseWatermarkEmbedTaskActionNonUpdatableParams = []string{
	"task_id",
	"action",
}

// @API DSC PUT /v1/{project_id}/data-watermark-embed-tasks/{id}/status
func ResourceDatabaseWatermarkEmbedTaskAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDatabaseWatermarkEmbedTaskActionCreate,
		ReadContext:   resourceDatabaseWatermarkEmbedTaskActionRead,
		UpdateContext: resourceDatabaseWatermarkEmbedTaskActionUpdate,
		DeleteContext: resourceDatabaseWatermarkEmbedTaskActionDelete,

		CustomizeDiff: config.FlexibleForceNew(databaseWatermarkEmbedTaskActionNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the database watermark embed task is located.`,
			},

			// Required parameters.
			"task_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the database watermark embed task.`,
			},
			"action": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The operation type of the database watermark embed task.`,
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

func actionDatabaseWatermarkEmbedTask(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	httpUrl := "v1/{project_id}/data-watermark-embed-tasks/{id}/status"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{id}", d.Get("task_id").(string))
	createPath = fmt.Sprintf("%s?action=%v", createPath, d.Get("action"))

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	_, err := client.Request("PUT", createPath, &createOpt)
	return err
}

func resourceDatabaseWatermarkEmbedTaskActionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		taskId = d.Get("task_id").(string)
	)

	client, err := cfg.NewServiceClient("dsc", region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	if err = actionDatabaseWatermarkEmbedTask(client, d); err != nil {
		return diag.Errorf("error actioning the database watermark embed task (%s): %s", taskId, err)
	}

	if d.Get("action").(string) == "START" {
		if _, err := waitForDatabaseWatermarkEmbedTaskCompleted(ctx, client, d.Timeout(schema.TimeoutCreate), taskId); err != nil {
			return diag.Errorf("error waiting for executing the database watermark embed task to be completed: %s", err)
		}
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(randomUUID.String())

	return nil
}

func resourceDatabaseWatermarkEmbedTaskActionRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDatabaseWatermarkEmbedTaskActionUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDatabaseWatermarkEmbedTaskActionDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to operate a DSC database watermark embed task. Deleting this resource
will not clear the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
