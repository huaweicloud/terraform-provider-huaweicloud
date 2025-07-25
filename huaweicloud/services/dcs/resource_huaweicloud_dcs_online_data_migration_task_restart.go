package dcs

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var onlineDataMigrationTaskStartNonUpdatableParams = []string{"task_id"}

// @API DCS POST /v2/{project_id}/migration-tasks/batch-restart
// @API DCS GET /v2/{project_id}/migration-task/{task_id}
func ResourceDcsOnlineDataMigrationTaskRestart() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDcsOnlineDataMigrationTaskRestartCreate,
		ReadContext:   resourceDcsOnlineDataMigrationTaskRestartRead,
		UpdateContext: resourceDcsOnlineDataMigrationTaskRestartUpdate,
		DeleteContext: resourceDcsOnlineDataMigrationTaskRestartDelete,

		CustomizeDiff: config.FlexibleForceNew(onlineDataMigrationTaskStartNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"task_id": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceDcsOnlineDataMigrationTaskRestartCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v2/{project_id}/migration-tasks/batch-restart"
		product = "dcs"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DCS client: %s", err)
	}

	taskId := d.Get("task_id").(string)
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{KeepResponseBody: true}
	createOpt.JSONBody = utils.RemoveNil(buildOnlineDataMigrationTaskRestartBodyParams(d))

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error restarting online data migration task(%s): %s", taskId, err)
	}

	d.SetId(taskId)

	err = checkMigrationTaskFinish(ctx, client, taskId, []string{"SUCCESS", "INCRMIGEATING"}, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func buildOnlineDataMigrationTaskRestartBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"migration_tasks": []string{d.Get("task_id").(string)},
	}
	return bodyParams
}

func resourceDcsOnlineDataMigrationTaskRestartRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDcsOnlineDataMigrationTaskRestartUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDcsOnlineDataMigrationTaskRestartDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting DCS restarting online data migration task resource is not supported. The resource is only " +
		"removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
