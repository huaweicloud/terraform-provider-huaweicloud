package kafka

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var instanceUpgradeNonUpdatableParams = []string{"instance_id", "is_schedule", "execute_at"}

// @API Kafka POST /v2/{project_id}/kafka/instances/{instance_id}/upgrade
// @API Kafka GET /v2/{project_id}/instances/{instance_id}/tasks/{task_id}
func ResourceInstanceUpgrade() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceInstanceUpgradeCreate,
		ReadContext:   resourceInstanceUpgradeRead,
		UpdateContext: resourceInstanceUpgradeUpdate,
		DeleteContext: resourceInstanceUpgradeDelete,

		CustomizeDiff: config.FlexibleForceNew(instanceUpgradeNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the Kafka instance to be upgraded is located.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the Kafka instance.`,
			},
			"is_schedule": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether to execute as a scheduled task.`,
			},
			"execute_at": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The scheduled time in Unix timestamp format, in milliseconds.`,
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

func resourceInstanceUpgradeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("dms", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	respBody, err := upgradeInstance(d, client, instanceId)
	if err != nil {
		return diag.Errorf("error upgrading Kafka instance: %s", err)
	}

	randomId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomId)

	isSchedule := d.Get("is_schedule").(bool)
	jobId := utils.PathSearch("job_id", respBody, "").(string)
	// The `job_id` field is not empty only when the task is not a scheduled task.
	if !isSchedule && jobId != "" {
		if err = waitForInstanceTaskStatusComplete(ctx, client, instanceId, jobId, d.Timeout(schema.TimeoutCreate)); err != nil {
			return diag.Errorf("error waiting for upgrading instance (%s) task to complete: %s", instanceId, err)
		}
	}

	return nil
}

func resourceInstanceUpgradeRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceInstanceUpgradeUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceInstanceUpgradeDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource for upgrading the Kafka instance. Deleting this resource will not 
recover the upgrade Kafka instance, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

func buildUpgradeInstanceBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"is_schedule": d.Get("is_schedule").(bool),
		"execute_at":  utils.ValueIgnoreEmpty(d.Get("execute_at")),
	}
}

func upgradeInstance(d *schema.ResourceData, client *golangsdk.ServiceClient, instanceId string) (interface{}, error) {
	httpUrl := "v2/{project_id}/kafka/instances/{instance_id}/upgrade"
	upgradePath := client.Endpoint + httpUrl
	upgradePath = strings.ReplaceAll(upgradePath, "{project_id}", client.ProjectID)
	upgradePath = strings.ReplaceAll(upgradePath, "{instance_id}", instanceId)

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildUpgradeInstanceBodyParams(d)),
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	resp, err := client.Request("POST", upgradePath, &requestOpt)
	if err != nil {
		return "", err
	}

	return utils.FlattenResponse(resp)
}
