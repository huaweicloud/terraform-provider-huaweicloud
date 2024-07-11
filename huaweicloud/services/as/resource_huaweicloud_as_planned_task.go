package as

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/autoscaling/v1/scheduledtasks"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API AS GET /autoscaling-api/v1/{project_id}/scaling-groups/{groupID}/scheduled-tasks
// @API AS POST /autoscaling-api/v1/{project_id}/scaling-groups/{groupID}/scheduled-tasks
// @API AS DELETE /autoscaling-api/v1/{project_id}/scaling-groups/{groupID}/scheduled-tasks/{taskID}
// @API AS PUT /autoscaling-api/v1/{project_id}/scaling-groups/{groupID}/scheduled-tasks/{taskID}
func ResourcePlannedTask() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePlannedTaskCreate,
		ReadContext:   resourcePlannedTaskRead,
		UpdateContext: resourcePlannedTaskUpdate,
		DeleteContext: resourcePlannedTaskDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourcePlannedTaskImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region where the planned task resource are located.",
			},
			"scaling_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the scaling group where the planned task to create.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the planned task to create.",
			},
			"scheduled_policy": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Elem:        plannedTaskScheduledPolicySchema(),
				Description: "The policy of planned task to create.",
			},
			"instance_number": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Elem:        plannedTaskInstanceNumberSchema(),
				Description: "The numbers of scaling group instance for planned task to create.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of the planned task.",
			},
		},
	}
}

func plannedTaskScheduledPolicySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"launch_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The execution time of planned task.",
			},
			"recurrence_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The triggering type of planned task",
			},
			"recurrence_value": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The frequency at which planned task are triggered",
			},
			"start_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The effective start time of planned task.",
			},
			"end_time": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The effective end time of planned task",
			},
		},
	}
}

func plannedTaskInstanceNumberSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"max": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "The maximum number of instances for the scaling group",
				AtLeastOneOf: []string{"instance_number.0.max", "instance_number.0.min", "instance_number.0.desire"},
			},
			"min": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The minimum number of instances for the scaling group.",
			},
			"desire": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The expected number of instances for the scaling group.",
			},
		},
	}
}

func buildPlannedTaskScheduledPolicy(rawScheduledPolicy map[string]interface{}) (scheduledtasks.ScheduledPolicy, error) {
	recurrenceType := rawScheduledPolicy["recurrence_type"].(string)
	if recurrenceType == "" {
		scheduledPolicy := scheduledtasks.ScheduledPolicy{
			LaunchTime: rawScheduledPolicy["launch_time"].(string),
		}
		return scheduledPolicy, nil
	}

	startTime := rawScheduledPolicy["start_time"].(string)
	if startTime == "" {
		startTime = getCurrentUTCwithoutSec()
	}

	endTime := rawScheduledPolicy["end_time"].(string)
	if endTime == "" {
		return scheduledtasks.ScheduledPolicy{}, fmt.Errorf("the end_time must be set")
	}

	scheduledPolicy := scheduledtasks.ScheduledPolicy{
		LaunchTime:      rawScheduledPolicy["launch_time"].(string),
		RecurrenceType:  recurrenceType,
		RecurrenceValue: rawScheduledPolicy["recurrence_value"].(string),
		StartTime:       startTime,
		EndTime:         endTime,
	}

	return scheduledPolicy, nil
}

func buildPlannedTaskInstanceNumber(rawInstanceNumber map[string]interface{}) (scheduledtasks.InstanceNumber, error) {
	var instanceNumber scheduledtasks.InstanceNumber

	if maxVal, ok := rawInstanceNumber["max"].(string); ok && maxVal != "" {
		maxInt, err := strconv.Atoi(maxVal)
		if err != nil {
			return scheduledtasks.InstanceNumber{}, err
		}
		instanceNumber.Max = &maxInt
	}
	if minVal, ok := rawInstanceNumber["min"].(string); ok && minVal != "" {
		minInt, err := strconv.Atoi(minVal)
		if err != nil {
			return scheduledtasks.InstanceNumber{}, err
		}
		instanceNumber.Min = &minInt
	}
	if desireVal, ok := rawInstanceNumber["desire"].(string); ok && desireVal != "" {
		desireInt, err := strconv.Atoi(desireVal)
		if err != nil {
			return scheduledtasks.InstanceNumber{}, err
		}
		instanceNumber.Desire = &desireInt
	}

	return instanceNumber, nil
}

func resourcePlannedTaskCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.AutoscalingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating AS v1 client: %s", err)
	}

	var (
		groupID    = d.Get("scaling_group_id").(string)
		createOpts = scheduledtasks.CreateOpts{
			Name: d.Get("name").(string),
		}

		scheduledPolicyList = d.Get("scheduled_policy").([]interface{})
		scheduledPolicyMap  = scheduledPolicyList[0].(map[string]interface{})
	)
	scheduledPolicy, err := buildPlannedTaskScheduledPolicy(scheduledPolicyMap)
	if err != nil {
		return diag.FromErr(err)
	}

	createOpts.ScheduledPolicy = scheduledPolicy
	var (
		instanceNumberList = d.Get("instance_number").([]interface{})
		instanceNumberMap  = instanceNumberList[0].(map[string]interface{})
	)
	instanceNumber, err := buildPlannedTaskInstanceNumber(instanceNumberMap)
	if err != nil {
		return diag.FromErr(err)
	}
	createOpts.InstanceNumber = instanceNumber
	taskId, err := scheduledtasks.Create(client, groupID, createOpts)
	if err != nil {
		return diag.Errorf("error creating AS planned task: %s", err)
	}
	d.SetId(taskId)

	return resourcePlannedTaskRead(ctx, d, meta)
}

func flattenPlannedTaskSchedulePolicy(policy scheduledtasks.ScheduledPolicy) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"launch_time":      policy.LaunchTime,
			"recurrence_type":  policy.RecurrenceType,
			"recurrence_value": policy.RecurrenceValue,
			"start_time":       policy.StartTime,
			"end_time":         policy.EndTime,
		},
	}
}

func flattenPlannedTaskInstanceNumber(instanceNumber scheduledtasks.InstanceNumber) []map[string]interface{} {
	result := map[string]interface{}{}
	if instanceNumber.Max != nil {
		result["max"] = strconv.Itoa(*instanceNumber.Max)
	}
	if instanceNumber.Min != nil {
		result["min"] = strconv.Itoa(*instanceNumber.Min)
	}
	if instanceNumber.Desire != nil {
		result["desire"] = strconv.Itoa(*instanceNumber.Desire)
	}

	return []map[string]interface{}{result}
}

func resourcePlannedTaskRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.AutoscalingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating AS v1 client: %s", err)
	}

	var (
		groupID  = d.Get("scaling_group_id").(string)
		taskId   = d.Id()
		listOpts = scheduledtasks.ListOpts{
			GroupID: groupID,
		}
	)
	plannedTasksResp, err := scheduledtasks.List(client, listOpts)
	if err != nil {
		// When the group does not exist, the API response body is an empty list.
		// It seems that `CheckDeletedDiag` here has no practical significance.
		// In order to avoid unknown problems, this part of the code is retained.
		return common.CheckDeletedDiag(d, err, "error retrieving AS planned task")
	}

	var mErr *multierror.Error
	for _, task := range plannedTasksResp {
		if task.ID == taskId {
			mErr = multierror.Append(nil,
				d.Set("region", region),
				d.Set("name", task.Name),
				d.Set("scaling_group_id", task.GroupID),
				d.Set("scheduled_policy", flattenPlannedTaskSchedulePolicy(task.ScheduledPolicy)),
				d.Set("instance_number", flattenPlannedTaskInstanceNumber(task.InstanceNumber)),
				d.Set("created_at", task.CreateTime),
			)
			return diag.FromErr(mErr.ErrorOrNil())
		}
	}

	return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "AS planned task")
}

func resourcePlannedTaskUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.AutoscalingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating AS v1 client: %s", err)
	}

	var (
		groupID    = d.Get("scaling_group_id").(string)
		taskId     = d.Id()
		updateOpts = scheduledtasks.UpdateOpts{}
	)
	if d.HasChanges("name") {
		updateOpts.Name = d.Get("name").(string)
	}

	if d.HasChanges("scheduled_policy") {
		var (
			scheduledPolicyList = d.Get("scheduled_policy").([]interface{})
			scheduledPolicyMap  = scheduledPolicyList[0].(map[string]interface{})
		)
		scheduledPolicy, err := buildPlannedTaskScheduledPolicy(scheduledPolicyMap)
		if err != nil {
			return diag.FromErr(err)
		}

		updateOpts.ScheduledPolicy = &scheduledPolicy
	}

	if d.HasChanges("instance_number") {
		var (
			instanceNumberList = d.Get("instance_number").([]interface{})
			instanceNumberMap  = instanceNumberList[0].(map[string]interface{})
		)
		instanceNumber, err := buildPlannedTaskInstanceNumber(instanceNumberMap)
		if err != nil {
			return diag.FromErr(err)
		}

		updateOpts.InstanceNumber = &instanceNumber
	}
	log.Printf("[DEBUG] update AS planned task: %+v", updateOpts)
	if err := scheduledtasks.Update(client, groupID, taskId, updateOpts); err != nil {
		return diag.Errorf("error updating AS planned task (%s): %s", taskId, err)
	}

	return resourcePlannedTaskRead(ctx, d, meta)
}

func resourcePlannedTaskDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	asClient, err := cfg.AutoscalingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating AS v1 client: %s", err)
	}

	var (
		groupID = d.Get("scaling_group_id").(string)
		taskId  = d.Id()
	)
	if err := scheduledtasks.Delete(asClient, groupID, taskId); err != nil {
		// When the group or task does not exist, the response HTTP status code of the delete API is 404.
		return common.CheckDeletedDiag(d, err, "error deleting AS planned task")
	}

	return nil
}

func resourcePlannedTaskImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format for import ID, want '<scaling_group_id>/<id>', but got '%s'", d.Id())
	}

	d.SetId(parts[1])
	return []*schema.ResourceData{d}, d.Set("scaling_group_id", parts[0])
}
