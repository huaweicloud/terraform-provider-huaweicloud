package kafka

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Kafka POST /v2/kafka/{project_id}/instances/{instance_id}/reassign
// @API Kafka GET /v2/{project_id}/instances/{instance_id}/tasks/{task_id}
// Corresponding OpenAPI reference of the POST method is as follows:
// @API Kafka POST /v2/{engine}/{project_id}/instances/{instance_id}/reassign
func ResourceDmsKafkaPartitionReassign() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDmsKafkaPartitionReassignCreate,
		ReadContext:   resourceDmsKafkaPartitionReassignRead,
		DeleteContext: resourceDmsKafkaPartitionReassignDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"reassignments": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"topic": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"brokers": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem:     &schema.Schema{Type: schema.TypeInt},
						},
						"replication_factor": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"assignment": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"partition": {
										Type:     schema.TypeInt,
										Optional: true,
										ForceNew: true,
									},
									"partition_brokers": {
										Type:     schema.TypeList,
										Optional: true,
										ForceNew: true,
										Elem:     &schema.Schema{Type: schema.TypeInt},
									},
								},
							},
						},
					},
				},
			},
			"throttle": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"is_schedule": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"execute_at": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"time_estimate": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"task_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"reassignment_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceDmsKafkaPartitionReassignCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dmsv2", region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)

	createHttpUrl := "v2/kafka/{project_id}/instances/{instance_id}/reassign"
	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateKafkaPartitionReassignBodyParams(d)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating kafka partition reassignment task: %s", err)
	}
	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	// set ID in UUID format for time estimation task having no `job_id` or `schedule_id` in return
	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	// just return one of the `reassignment_time`, `job_id` and `schedule_id`, depends on the value of `time_estimate` and `is_schedule`
	reassignmentTime := utils.PathSearch("reassignment_time", createRespBody, nil)
	jobID := utils.PathSearch("job_id", createRespBody, nil)
	scheduleID := utils.PathSearch("schedule_id", createRespBody, nil)

	switch {
	case jobID != nil:
		// wait for task complete
		// if it's not scheduled task, use `/v2/{project_id}/instances/{instance_id}/tasks/{task_id}` to search task.
		stateConf := &resource.StateChangeConf{
			Pending:      []string{"CREATED", "EXECUTING"},
			Target:       []string{"SUCCESS"},
			Refresh:      kafkaInstanceTaskStatusRefreshFunc(client, instanceID, jobID.(string)),
			Timeout:      d.Timeout(schema.TimeoutCreate),
			Delay:        1 * time.Second,
			PollInterval: 5 * time.Second,
		}
		if _, err := stateConf.WaitForStateContext(ctx); err != nil {
			return diag.Errorf("error waiting for the Kafka instance (%s) partition reassignment task to be finished: %s ",
				instanceID, err)
		}

		// since the resource ID is in UUID format, return `job_id`.
		d.Set("task_id", jobID)
	case scheduleID != nil:
		// if it's a scheduled task, use `/v2/{project_id}/instances/{instance_id}/scheduled-tasks` to search task,
		// but it is not an open API, so return `schedule_id` for scheduled task.
		d.Set("task_id", scheduleID)
	case reassignmentTime != nil:
		// set `reassignment_time` for time estimation task
		d.Set("reassignment_time", reassignmentTime)
	default:
		return diag.Errorf("error creating kafka partition reassignment task: `reassignment_time`, `job_id` and " +
			"`schedule_id` are all not found in response")
	}

	return resourceDmsKafkaPartitionReassignRead(ctx, d, meta)
}

func buildCreateKafkaPartitionReassignBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"reassignments": buildCreateReassignBodyParamsReassignments(d.Get("reassignments").([]interface{})),
		"throttle":      utils.ValueIgnoreEmpty(d.Get("throttle")),
		"is_schedule":   utils.ValueIgnoreEmpty(d.Get("is_schedule")),
		"execute_at":    utils.ValueIgnoreEmpty(d.Get("execute_at")),
		"time_estimate": utils.ValueIgnoreEmpty(d.Get("time_estimate")),
	}
	return bodyParams
}

func buildCreateReassignBodyParamsReassignments(rawParams []interface{}) []map[string]interface{} {
	rst := make([]map[string]interface{}, 0, len(rawParams))
	for _, val := range rawParams {
		raw := val.(map[string]interface{})
		params := map[string]interface{}{
			"topic":              raw["topic"],
			"brokers":            utils.ValueIgnoreEmpty(raw["brokers"]),
			"replication_factor": utils.ValueIgnoreEmpty(raw["replication_factor"]),
			"assignment":         buildCreateReassignBodyParamsAssignment(raw["assignment"].([]interface{})),
		}
		rst = append(rst, params)
	}

	return rst
}

func buildCreateReassignBodyParamsAssignment(rawParams []interface{}) []map[string]interface{} {
	if len(rawParams) == 0 {
		return nil
	}
	rst := make([]map[string]interface{}, 0)
	for _, val := range rawParams {
		raw, ok := val.(map[string]interface{})
		if !ok {
			continue
		}
		params := map[string]interface{}{
			"partition":         raw["partition"],
			"partition_brokers": utils.ValueIgnoreEmpty(raw["partition_brokers"]),
		}
		rst = append(rst, params)
	}
	return rst
}

func resourceDmsKafkaPartitionReassignRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDmsKafkaPartitionReassignDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting resource is not supported. The resource is only removed from the state, the task remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
