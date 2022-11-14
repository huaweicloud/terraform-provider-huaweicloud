package as

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk/openstack/autoscaling/v1/policies"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
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
							Type:             schema.TypeString,
							Optional:         true,
							Default:          getCurrentUTCwithoutSec(),
							DiffSuppressFunc: utils.SuppressDiffAll,
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
				Type:     schema.TypeInt,
				Optional: true,
				Default:  900,
			},
		},
	}
}

func getCurrentUTCwithoutSec() string {
	utcTime := time.Now().UTC().Format(time.RFC3339)
	splits := strings.SplitN(utcTime, ":", 3)
	resultTime := strings.Join(splits[0:2], ":") + "Z"
	return resultTime
}

func validateParameters(d *schema.ResourceData) error {
	policyType := d.Get("scaling_policy_type").(string)
	alarmId := d.Get("alarm_id").(string)
	scheduledPolicy := d.Get("scheduled_policy").([]interface{})

	if policyType == "ALARM" {
		if alarmId == "" {
			return fmt.Errorf("Parameter alarm_id should be set if policy type is ALARM.")
		}
	}
	if policyType == "SCHEDULED" || policyType == "RECURRENCE" {
		if len(scheduledPolicy) == 0 {
			return fmt.Errorf("Parameter scheduled_policy should be set if policy type is RECURRENCE or SCHEDULED.")
		}
	}

	if len(scheduledPolicy) == 1 {
		scheduledPolicyMap := scheduledPolicy[0].(map[string]interface{})
		recurrenceType := scheduledPolicyMap["recurrence_type"].(string)
		endTime := scheduledPolicyMap["end_time"].(string)

		if policyType == "RECURRENCE" {
			if recurrenceType == "" {
				return fmt.Errorf("Parameter recurrence_type should be set if policy type is RECURRENCE.")
			}
			if endTime == "" {
				return fmt.Errorf("Parameter end_time should be set if policy type is RECURRENCE.")
			}
		}
	}

	return nil
}

func buildScheduledPolicy(rawScheduledPolicy map[string]interface{}) policies.SchedulePolicyOpts {
	scheduledPolicy := policies.SchedulePolicyOpts{
		LaunchTime:      rawScheduledPolicy["launch_time"].(string),
		RecurrenceType:  rawScheduledPolicy["recurrence_type"].(string),
		RecurrenceValue: rawScheduledPolicy["recurrence_value"].(string),
		StartTime:       rawScheduledPolicy["start_time"].(string),
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
		return diag.Errorf("Error creating autoscaling client: %s", err)
	}

	err = validateParameters(d)
	if err != nil {
		return diag.Errorf("Error creating ASPolicy: %s", err)
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
		return diag.Errorf("Error creating ASPolicy: %s", err)
	}

	d.SetId(asPolicyId)
	return resourceASPolicyRead(ctx, d, meta)
}

func resourceASPolicyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	asClient, err := conf.AutoscalingV1Client(region)
	if err != nil {
		return diag.Errorf("Error creating autoscaling client: %s", err)
	}

	policyId := d.Id()
	asPolicy, err := policies.Get(asClient, policyId).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "AS Policy")
	}

	log.Printf("[DEBUG] Retrieved ASPolicy %q: %+v", policyId, asPolicy)
	d.Set("scaling_policy_name", asPolicy.Name)
	d.Set("scaling_policy_type", asPolicy.Type)
	d.Set("alarm_id", asPolicy.AlarmID)
	d.Set("cool_down_time", asPolicy.CoolDownTime)

	policyActionInfo := asPolicy.Action
	policyAction := map[string]interface{}{}
	policyAction["operation"] = policyActionInfo.Operation
	policyAction["instance_number"] = policyActionInfo.InstanceNum
	policyActionList := []map[string]interface{}{}
	policyActionList = append(policyActionList, policyAction)
	d.Set("scaling_policy_action", policyActionList)

	scheduledInfo := asPolicy.SchedulePolicy
	if scheduledInfo.LaunchTime != "" {
		scheduledMap := map[string]interface{}{
			"launch_time":      scheduledInfo.LaunchTime,
			"recurrence_type":  scheduledInfo.RecurrenceType,
			"recurrence_value": scheduledInfo.RecurrenceValue,
			"start_time":       scheduledInfo.StartTime,
			"end_time":         scheduledInfo.EndTime,
		}
		scheduledPolicies := []map[string]interface{}{}
		scheduledPolicies = append(scheduledPolicies, scheduledMap)
		d.Set("scheduled_policy", scheduledPolicies)
	} else {
		d.Set("scheduled_policy", nil)
	}

	d.Set("region", region)

	return nil
}

func resourceASPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	asClient, err := conf.AutoscalingV1Client(region)
	if err != nil {
		return diag.Errorf("Error creating autoscaling client: %s", err)
	}

	err = validateParameters(d)
	if err != nil {
		return diag.Errorf("Error updating ASPolicy: %s", err)
	}
	updateOpts := policies.UpdateOpts{
		Name:         d.Get("scaling_policy_name").(string),
		Type:         d.Get("scaling_policy_type").(string),
		AlarmID:      d.Get("alarm_id").(string),
		CoolDownTime: d.Get("cool_down_time").(int),
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
		return diag.Errorf("Error updating ASPolicy %q: %s", asPolicyID, err)
	}

	return resourceASPolicyRead(ctx, d, meta)
}

func resourceASPolicyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	asClient, err := conf.AutoscalingV1Client(region)
	if err != nil {
		return diag.Errorf("Error creating autoscaling client: %s", err)
	}

	log.Printf("[DEBUG] Begin to delete AS policy %q", d.Id())
	if delErr := policies.Delete(asClient, d.Id()).ExtractErr(); delErr != nil {
		return diag.Errorf("Error deleting AS policy: %s", delErr)
	}

	return nil
}
