package as

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk/openstack/autoscaling/v1/policies"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

var (
	PolicyTypes       = []string{"ALARM", "SCHEDULED", "RECURRENCE"}
	RecurrencePeriods = []string{"Daily", "Weekly", "Monthly"}
	PolicyActions     = []string{"ADD", "REMOVE", "SET"}
)

func ResourceASPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceASPolicyCreate,
		ReadContext:   resourceASPolicyRead,
		UpdateContext: resourceASPolicyUpdate,
		DeleteContext: resourceASPolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"scaling_policy_name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 64),
					validation.StringMatch(regexp.MustCompile("^[\u4e00-\u9fa50-9a-zA-Z-_]+$"),
						"only letters, digits, underscores (_), and hyphens (-) are allowed"),
				),
			},
			"scaling_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"scaling_policy_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(PolicyTypes, false),
			},
			"alarm_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"scheduled_policy": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"launch_time": {
							Type:     schema.TypeString,
							Required: true,
						},
						"recurrence_type": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringInSlice(RecurrencePeriods, false),
						},
						"recurrence_value": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"start_time": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"end_time": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"scaling_policy_action": {
				Optional: true,
				Computed: true,
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"operation": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringInSlice(PolicyActions, false),
						},
						"instance_number": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  1,
						},
					},
				},
			},
			"cool_down_time": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(0, 86400),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func getCurrentUTCwithoutSec() string {
	utcTime := time.Now().UTC()
	return utcTime.Format("2006-01-02T15:04Z")
}

func validateParameters(d *schema.ResourceData) error {
	policyType := d.Get("scaling_policy_type").(string)
	alarmId := d.Get("alarm_id").(string)
	scheduledPolicy := d.Get("scheduled_policy").([]interface{})

	if policyType == "ALARM" && alarmId == "" {
		return fmt.Errorf("parameter alarm_id should be set if policy type is ALARM")
	}
	if policyType == "SCHEDULED" || policyType == "RECURRENCE" {
		if len(scheduledPolicy) == 0 {
			return fmt.Errorf("parameter scheduled_policy should be set if policy type is RECURRENCE or SCHEDULED")
		}
	}

	if len(scheduledPolicy) == 1 {
		scheduledPolicyMap := scheduledPolicy[0].(map[string]interface{})
		recurrenceType := scheduledPolicyMap["recurrence_type"].(string)
		endTime := scheduledPolicyMap["end_time"].(string)

		if policyType == "RECURRENCE" {
			if recurrenceType == "" {
				return fmt.Errorf("parameter recurrence_type should be set if policy type is RECURRENCE")
			}
			if endTime == "" {
				return fmt.Errorf("parameter end_time should be set if policy type is RECURRENCE")
			}
		}
	}

	return nil
}

func buildScheduledPolicy(rawScheduledPolicy map[string]interface{}) policies.SchedulePolicyOpts {
	recurrenceType := rawScheduledPolicy["recurrence_type"].(string)
	startTime := rawScheduledPolicy["start_time"].(string)
	if recurrenceType != "" && startTime == "" {
		startTime = getCurrentUTCwithoutSec()
	}

	scheduledPolicy := policies.SchedulePolicyOpts{
		LaunchTime:      rawScheduledPolicy["launch_time"].(string),
		RecurrenceValue: rawScheduledPolicy["recurrence_value"].(string),
		RecurrenceType:  recurrenceType,
		StartTime:       startTime,
		EndTime:         rawScheduledPolicy["end_time"].(string),
	}
	return scheduledPolicy
}

func buildPolicyAction(rawPolicyAction map[string]interface{}) policies.ActionOpts {
	policyAction := policies.ActionOpts{
		Operation:   rawPolicyAction["operation"].(string),
		InstanceNum: rawPolicyAction["instance_number"].(int),
	}
	return policyAction
}

func resourceASPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	asClient, err := config.AutoscalingV1Client(config.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating autoscaling client: %s", err)
	}

	err = validateParameters(d)
	if err != nil {
		return diag.Errorf("error creating AS policy: %s", err)
	}
	createOpts := policies.CreateOpts{
		Name:         d.Get("scaling_policy_name").(string),
		ID:           d.Get("scaling_group_id").(string),
		Type:         d.Get("scaling_policy_type").(string),
		AlarmID:      d.Get("alarm_id").(string),
		CoolDownTime: d.Get("cool_down_time").(int),
	}
	scheduledPolicyList := d.Get("scheduled_policy").([]interface{})
	if len(scheduledPolicyList) == 1 {
		scheduledPolicyMap := scheduledPolicyList[0].(map[string]interface{})
		scheduledPolicy := buildScheduledPolicy(scheduledPolicyMap)
		createOpts.SchedulePolicy = scheduledPolicy
	}
	policyActionList := d.Get("scaling_policy_action").([]interface{})
	if len(policyActionList) == 1 {
		policyActionMap := policyActionList[0].(map[string]interface{})
		policyAction := buildPolicyAction(policyActionMap)
		createOpts.Action = policyAction
	}

	log.Printf("[DEBUG] Create AS policy Options: %#v", createOpts)
	asPolicyId, err := policies.Create(asClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating AS policy: %s", err)
	}

	d.SetId(asPolicyId)
	return resourceASPolicyRead(ctx, d, meta)
}

func resourceASPolicyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	asClient, err := conf.AutoscalingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating autoscaling client: %s", err)
	}

	policyId := d.Id()
	asPolicy, err := policies.Get(asClient, policyId).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "AS policy")
	}

	log.Printf("[DEBUG] Retrieved AS policy %s: %+v", policyId, asPolicy)
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("scaling_policy_name", asPolicy.Name),
		d.Set("scaling_policy_type", asPolicy.Type),
		d.Set("scaling_group_id", asPolicy.ID),
		d.Set("alarm_id", asPolicy.AlarmID),
		d.Set("cool_down_time", asPolicy.CoolDownTime),
		d.Set("status", asPolicy.Status),
		d.Set("scaling_policy_action", flattenPolicyAction(asPolicy.Action)),
		d.Set("scheduled_policy", flattenSchedulePolicy(asPolicy.SchedulePolicy)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenPolicyAction(action policies.Action) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"operation":       action.Operation,
			"instance_number": action.InstanceNum,
		},
	}
}

func flattenSchedulePolicy(policy policies.SchedulePolicy) []map[string]interface{} {
	if policy.LaunchTime == "" {
		return nil
	}

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

func resourceASPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	asClient, err := conf.AutoscalingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating autoscaling client: %s", err)
	}

	err = validateParameters(d)
	if err != nil {
		return diag.Errorf("error updating AS policy: %s", err)
	}
	updateOpts := policies.UpdateOpts{
		Name:    d.Get("scaling_policy_name").(string),
		Type:    d.Get("scaling_policy_type").(string),
		AlarmID: d.Get("alarm_id").(string),
	}
	if d.HasChange("cool_down_time") {
		updateOpts.CoolDownTime = d.Get("cool_down_time").(int)
	}

	scheduledPolicyList := d.Get("scheduled_policy").([]interface{})
	if len(scheduledPolicyList) == 1 {
		scheduledPolicyMap := scheduledPolicyList[0].(map[string]interface{})
		scheduledPolicy := buildScheduledPolicy(scheduledPolicyMap)
		updateOpts.SchedulePolicy = scheduledPolicy
	}
	policyActionList := d.Get("scaling_policy_action").([]interface{})
	if len(policyActionList) == 1 {
		policyActionMap := policyActionList[0].(map[string]interface{})
		policyAction := buildPolicyAction(policyActionMap)
		updateOpts.Action = policyAction
	}

	log.Printf("[DEBUG] Update AS policy Options: %#v", updateOpts)
	asPolicyID, err := policies.Update(asClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return diag.Errorf("error updating AS policy %s: %s", asPolicyID, err)
	}

	return resourceASPolicyRead(ctx, d, meta)
}

func resourceASPolicyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	asClient, err := conf.AutoscalingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating autoscaling client: %s", err)
	}

	if delErr := policies.Delete(asClient, d.Id()).ExtractErr(); delErr != nil {
		return diag.Errorf("error deleting AS policy: %s", delErr)
	}

	return nil
}
