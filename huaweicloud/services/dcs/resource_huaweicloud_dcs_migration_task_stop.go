package dcs

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

var migrationTaskStopNonUpdatableParams = []string{"task_id"}

// @API DCS POST /v2/{project_id}/migration-task/{task_id}/stop
// @API DCS GET /v2/{project_id}/migration-task/{task_id}
func ResourceDcsMigrationTaskStop() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDcsMigrationTaskStopCreate,
		ReadContext:   resourceDcsMigrationTaskStopRead,
		UpdateContext: resourceDcsMigrationTaskStopUpdate,
		DeleteContext: resourceDcsMigrationTaskStopDelete,

		CustomizeDiff: config.FlexibleForceNew(migrationTaskStopNonUpdatableParams),

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

func resourceDcsMigrationTaskStopCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v2/{project_id}/migration-task/{task_id}/stop"
		product = "dcs"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DCS client: %s", err)
	}

	taskId := d.Get("task_id").(string)
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{task_id}", taskId)

	createOpt := golangsdk.RequestOpts{KeepResponseBody: true}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error stopping migration task(%s): %s", taskId, err)
	}

	d.SetId(taskId)

	err = checkMigrationTaskFinish(ctx, client, taskId, []string{"TERMINATED"}, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceDcsMigrationTaskStopRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDcsMigrationTaskStopUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDcsMigrationTaskStopDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting DCS stopping migration task resource is not supported. The resource is only removed from the" +
		"state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
