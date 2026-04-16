package rabbitmq

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var recycleInstanceRestoreNonUpdatableParams = []string{
	"instance_id",
}

// @API RabbitMQ POST /v2/{project_id}/recycle
// @API RabbitMQ GET /v2/{project_id}/instances/{instance_id}/tasks/{task_id}
func ResourceRecycleInstanceRestore() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRecycleInstanceRestoreCreate,
		ReadContext:   resourceRecycleInstanceRestoreRead,
		UpdateContext: resourceRecycleInstanceRestoreUpdate,
		DeleteContext: resourceRecycleInstanceRestoreDelete,

		CustomizeDiff: config.FlexibleForceNew(recycleInstanceRestoreNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the instance is located to be restored.`,
			},

			// Required parameters.
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The instance ID to be restored from recycle bin.`,
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

func resourceRecycleInstanceRestoreCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		httpUrl    = "v2/{project_id}/recycle"
		instanceId = d.Get("instance_id").(string)
	)

	client, err := cfg.NewServiceClient("dms", region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody: map[string]interface{}{
			"instances": []interface{}{instanceId},
		},
	}

	createResp, err := client.Request("POST", createPath, &createOpts)
	if err != nil {
		return diag.Errorf("error restoring recycle bin instance (%s): %s", instanceId, err)
	}

	respBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	jobId := utils.PathSearch(fmt.Sprintf("results[?instance_id=='%s']|[0].job_id", instanceId), respBody, "").(string)
	if jobId == "" {
		return diag.Errorf("unable to find the job ID for restoring instance (%s) from API response", instanceId)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      recycleInstanceRestoreTaskStatusRefreshFunc(client, instanceId, jobId),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        10 * time.Second,
		PollInterval: 20 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for the instance (%s) to be restored from recycle bin: %s", instanceId, err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate resource ID: %s", err)
	}

	d.SetId(randUUID)

	return nil
}

func recycleInstanceRestoreTaskStatusRefreshFunc(client *golangsdk.ServiceClient, instanceId, taskId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := getInstanceTaskById(client, instanceId, taskId)
		if err != nil {
			return resp, "ERROR", err
		}

		status := utils.PathSearch("tasks|[0].status", resp, "").(string)
		if utils.StrSliceContains([]string{"FAILED", "DELETED"}, status) {
			return resp, status, fmt.Errorf("unexpect status (%s)", status)
		}

		if status == "SUCCESS" {
			return resp, "COMPLETED", nil
		}

		return resp, "PENDING", nil
	}
}

func resourceRecycleInstanceRestoreRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceRecycleInstanceRestoreUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceRecycleInstanceRestoreDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for restoring recycle bin instance. Deleting this resource 
will not clear the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
