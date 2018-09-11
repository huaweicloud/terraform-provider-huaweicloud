package huaweicloud

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/autoscaling/v1/policies"
)

func resourceASPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceASPolicyCreate,
		Read:   resourceASPolicyRead,
		Update: resourceASPolicyUpdate,
		Delete: resourceASPolicyDelete,

		Schema: map[string]*schema.Schema{
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"scaling_policy_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: resourceASPolicyValidateName,
				ForceNew:     false,
			},
			"scaling_group_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"scaling_policy_type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     false,
				ValidateFunc: resourceASPolicyValidatePolicyType,
			},
			"alarm_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"scheduled_policy": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: false,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"launch_time": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							ForceNew: false,
						},
						"recurrence_type": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							ForceNew:     false,
							ValidateFunc: resourceASPolicyValidateRecurrenceType,
						},
						"recurrence_value": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: false,
						},
						"start_time": &schema.Schema{
							Type:             schema.TypeString,
							Optional:         true,
							ForceNew:         false,
							Default:          getCurrentUTCwithoutSec(),
							DiffSuppressFunc: suppressDiffAll,
						},
						"end_time": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: false,
						},
					},
				},
			},
			"scaling_policy_action": &schema.Schema{
				Optional: true,
				Type:     schema.TypeList,
				MaxItems: 1,
				ForceNew: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"operation": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: resourceASPolicyValidateActionOperation,
						},
						"instance_number": &schema.Schema{
							Type:     schema.TypeInt,
							Optional: true,
							Default:  1,
						},
					},
				},
			},
			"cool_down_time": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: false,
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
	log.Printf("[DEBUG] validateParameters for as policy!")
	policyType := d.Get("scaling_policy_type").(string)
	alarmId := d.Get("alarm_id").(string)
	log.Printf("[DEBUG] validateParameters alarmId is :%s", alarmId)
	log.Printf("[DEBUG] validateParameters policyType is :%s", policyType)
	scheduledPolicy := d.Get("scheduled_policy").([]interface{})
	log.Printf("[DEBUG] validateParameters scheduledPolicy is :%#v", scheduledPolicy)
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
		log.Printf("[DEBUG] validateParameters scheduledPolicyMap is :%#v", scheduledPolicyMap)
		recurrenceType := scheduledPolicyMap["recurrence_type"].(string)
		endTime := scheduledPolicyMap["end_time"].(string)
		log.Printf("[DEBUG] validateParameters recurrenceType is :%#v", recurrenceType)
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

func getScheduledPolicy(rawScheduledPolicy map[string]interface{}) policies.SchedulePolicyOpts {
	scheduledPolicy := policies.SchedulePolicyOpts{
		LaunchTime:      rawScheduledPolicy["launch_time"].(string),
		RecurrenceType:  rawScheduledPolicy["recurrence_type"].(string),
		RecurrenceValue: rawScheduledPolicy["recurrence_value"].(string),
		StartTime:       rawScheduledPolicy["start_time"].(string),
		EndTime:         rawScheduledPolicy["end_time"].(string),
	}
	return scheduledPolicy
}

func getPolicyAction(rawPolicyAction map[string]interface{}) policies.ActionOpts {
	policyAction := policies.ActionOpts{
		Operation:   rawPolicyAction["operation"].(string),
		InstanceNum: rawPolicyAction["instance_number"].(int),
	}
	return policyAction
}

func resourceASPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	asClient, err := config.autoscalingV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud autoscaling client: %s", err)
	}
	log.Printf("[DEBUG] asClient: %#v", asClient)
	err = validateParameters(d)
	if err != nil {
		return fmt.Errorf("Error creating ASPolicy: %s", err)
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
		scheduledPolicy := getScheduledPolicy(scheduledPolicyMap)
		createOpts.SchedulePolicy = scheduledPolicy
	}
	policyActionList := d.Get("scaling_policy_action").([]interface{})
	if len(policyActionList) == 1 {
		policyActionMap := policyActionList[0].(map[string]interface{})
		policyAction := getPolicyAction(policyActionMap)
		createOpts.Action = policyAction
	}

	log.Printf("[DEBUG] Create AS policy Options: %#v", createOpts)
	asPolicyId, err := policies.Create(asClient, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating ASPolicy: %s", err)
	}
	d.SetId(asPolicyId)
	log.Printf("[DEBUG] Create AS Policy %q Success!", asPolicyId)
	return resourceASPolicyRead(d, meta)
}

func resourceASPolicyRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	asClient, err := config.autoscalingV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud autoscaling client: %s", err)
	}

	asPolicy, err := policies.Get(asClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "AS Policy")
	}

	log.Printf("[DEBUG] Retrieved ASPolicy %q: %+v", d.Id(), asPolicy)
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

	scheduledPolicyInfo := asPolicy.SchedulePolicy
	scheduledPolicy := map[string]interface{}{}
	scheduledPolicy["launch_time"] = scheduledPolicyInfo.LaunchTime
	scheduledPolicy["recurrence_type"] = scheduledPolicyInfo.RecurrenceType
	scheduledPolicy["recurrence_value"] = scheduledPolicyInfo.RecurrenceValue
	scheduledPolicy["start_time"] = scheduledPolicyInfo.StartTime
	scheduledPolicy["end_time"] = scheduledPolicyInfo.EndTime
	scheduledPolicies := []map[string]interface{}{}
	scheduledPolicies = append(scheduledPolicies, scheduledPolicy)
	d.Set("scheduled_policy", scheduledPolicies)

	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceASPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	asClient, err := config.autoscalingV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud autoscaling client: %s", err)
	}

	err = validateParameters(d)
	if err != nil {
		return fmt.Errorf("Error updating ASPolicy: %s", err)
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
		scheduledPolicy := getScheduledPolicy(scheduledPolicyMap)
		updateOpts.SchedulePolicy = scheduledPolicy
	}
	policyActionList := d.Get("scaling_policy_action").([]interface{})
	if len(policyActionList) == 1 {
		policyActionMap := policyActionList[0].(map[string]interface{})
		policyAction := getPolicyAction(policyActionMap)
		updateOpts.Action = policyAction
	}
	log.Printf("[DEBUG] Update AS policy Options: %#v", updateOpts)
	asPolicyID, err := policies.Update(asClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error updating ASPolicy %q: %s", asPolicyID, err)
	}

	return resourceASPolicyRead(d, meta)
}

func resourceASPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	asClient, err := config.autoscalingV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud autoscaling client: %s", err)
	}
	log.Printf("[DEBUG] Begin to delete AS policy %q", d.Id())
	if delErr := policies.Delete(asClient, d.Id()).ExtractErr(); delErr != nil {
		return fmt.Errorf("Error deleting AS policy: %s", delErr)
	}

	return nil
}

var RecurrenceTypes = [3]string{"Daily", "Weekly", "Monthly"}

func resourceASPolicyValidateRecurrenceType(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	for i := range RecurrenceTypes {
		if value == RecurrenceTypes[i] {
			return
		}
	}
	errors = append(errors, fmt.Errorf("%q must be one of %v", k, RecurrenceTypes))
	return
}

var PolicyActions = [3]string{"ADD", "REMOVE", "SET"}

func resourceASPolicyValidateActionOperation(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	for i := range PolicyActions {
		if value == PolicyActions[i] {
			return
		}
	}
	errors = append(errors, fmt.Errorf("%q must be one of %v", k, PolicyActions))
	return
}

var PolicyTypes = [3]string{"ALARM", "SCHEDULED", "RECURRENCE"}

func resourceASPolicyValidatePolicyType(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	for i := range PolicyTypes {
		if value == PolicyTypes[i] {
			return
		}
	}
	errors = append(errors, fmt.Errorf("%q must be one of %v", k, PolicyTypes))
	return
}

func resourceASPolicyValidateName(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	if len(value) > 64 || len(value) < 1 {
		errors = append(errors, fmt.Errorf("%q must contain more than 1 and less than 64 characters", value))
	}
	if !regexp.MustCompile(`^[0-9a-zA-Z-_]+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf("only alphanumeric characters, hyphens, and underscores allowed in %q", value))
	}
	return
}
