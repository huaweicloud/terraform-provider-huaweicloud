package dcs

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var migrationTaskRollbackIpNonUpdatableParams = []string{"task_id"}

// @API DCS POST /v2/{project_id}/migration-task/{task_id}/rollback-ip
// @API DCS GET /v2/{project_id}/migration-task/{task_id}
func ResourceDcsMigrationTaskRollbackIp() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDcsMigrationTaskRollbackIpCreate,
		ReadContext:   resourceDcsMigrationTaskRollbackIpRead,
		UpdateContext: resourceDcsMigrationTaskRollbackIpUpdate,
		DeleteContext: resourceDcsMigrationTaskRollbackIpDelete,

		CustomizeDiff: config.FlexibleForceNew(migrationTaskRollbackIpNonUpdatableParams),

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
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func resourceDcsMigrationTaskRollbackIpCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v2/{project_id}/migration-task/{task_id}/rollback-ip"
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

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes:          []int{200, 204},
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error rolling migration task(%s) IP: %s", taskId, err)
	}

	d.SetId(taskId)

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"ROLLBACK_SUCCESS"},
		Refresh:      migrationTaskSwitchIpRefreshFunc(client, taskId),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("error waiting for migration task(%s) exchange IP rollback to be completed: %s ", taskId, err)
	}

	return nil
}

func resourceDcsMigrationTaskRollbackIpRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDcsMigrationTaskRollbackIpUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDcsMigrationTaskRollbackIpDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting DCS rollbacking migration task exchange IP resource is not supported. The resource is only " +
		"removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
