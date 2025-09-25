package kafka

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/dms/v2/kafka/instances"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var instanceBatchActionNonUpdatableParams = []string{"action", "instances", "all_failure", "force_delete"}

// @API Kafka POST /v2/{project_id}/instances/action
// @API Kafka GET /v2/{project_id}/instances/{instance_id}
func ResourceInstanceBatchAction() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceInstanceBatchActionCreate,
		ReadContext:   resourceInstanceBatchActionRead,
		UpdateContext: resourceInstanceBatchActionUpdate,
		DeleteContext: resourceInstanceBatchActionDelete,

		CustomizeDiff: config.FlexibleForceNew(instanceBatchActionNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the instances to be operated are located.`,
			},
			"action": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The type of the operation.`,
			},
			"instances": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of instance IDs to be operated.`,
			},
			"all_failure": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Whether to delete all instances that failed to be created.`,
			},
			"force_delete": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether to force delete instances.`,
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

func buildCreateInstanceBatchActionBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"action":       d.Get("action"),
		"instances":    utils.ValueIgnoreEmpty(d.Get("instances")),
		"all_failure":  utils.ValueIgnoreEmpty(d.Get("all_failure")),
		"force_delete": d.Get("force_delete"),
	}
}

func resourceInstanceBatchActionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		action = d.Get("action").(string)
	)
	client, err := cfg.NewServiceClient("dmsv2", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	createHttpUrl := "v2/{project_id}/instances/action"
	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateInstanceBatchActionBodyParams(d)),
		// `204`: Delete all created failed Kafka instances successfully.
		OkCodes: []int{
			200, 204,
		},
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("unable to batch %s instances: %s", action, err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	operatedSuccessedInstanceIds := utils.PathSearch("results[?result=='success'].instance", respBody, make([]interface{}, 0)).([]interface{})
	if len(operatedSuccessedInstanceIds) == 0 {
		return diag.Errorf("unable to batch %s instances: no successed instances", action)
	}

	for _, instanceId := range operatedSuccessedInstanceIds {
		stateConf := &resource.StateChangeConf{
			Pending:      []string{"PENDING"},
			Target:       []string{"COMPLETED"},
			Refresh:      refreshBatchActionInstancesState(client, instanceId.(string)),
			Timeout:      d.Timeout(schema.TimeoutCreate),
			Delay:        20 * time.Second,
			PollInterval: 20 * time.Second,
		}
		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			// Some instances may be operation failed, so use log to record the error.
			log.Printf("[ERROR] error waiting for the instance (%v) to be %s completed: %s", instanceId, action, err)
		}
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)
	return nil
}

func refreshBatchActionInstancesState(client *golangsdk.ServiceClient, instanceId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		// When an instance is deleted, the query instance list API returns data that no longer contains the instance,
		// but the instance is actually in the deletion state. Therefore, use the query instance details API
		respBody, err := instances.Get(client, instanceId).Extract()
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return respBody, "COMPLETED", nil
			}
			return nil, "QUERY ERROR", err
		}

		// If the recycle bin function is enabled, the status of non-forced deletion instances is `RECYCLE`.
		if respBody.Status == "RUNNING" || respBody.Status == "RECYCLE" {
			return respBody, "COMPLETED", nil
		}

		return respBody, "PENDING", nil
	}
}

func resourceInstanceBatchActionRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceInstanceBatchActionUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceInstanceBatchActionDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for restarting or deleting Kafka instances. Deleting
this resource will not clear the corresponding request record, but will only remove the resource information from the
tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
