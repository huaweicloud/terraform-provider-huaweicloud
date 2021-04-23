package huaweicloud

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/huaweicloud/golangsdk/openstack/cbr/v3/policies"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func resourceCBRPolicyV3() *schema.Resource {
	return &schema.Resource{
		Create: resourceCBRPolicyV3Create,
		Read:   resourceCBRPolicyV3Read,
		Update: resourceCBRPolicyV3Update,
		Delete: resourceCBRPolicyV3Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 64),
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"protection_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"backup", "replication",
				}, false),
			},
			"backup_cycle": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"frequency": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"WEEKLY", "DAILY",
							}, false),
						},
						"interval": {
							Type:         schema.TypeInt,
							Optional:     true,
							ExactlyOneOf: []string{"backup_cycle.0.days"},
							ValidateFunc: validation.IntBetween(1, 30),
						},
						"days": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validateBackupWeekDay,
						},
						"execution_times": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 24,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validateTimeFormat,
							},
						},
					},
				},
			},
			"destination_region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"destination_project_id": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"destination_region"},
			},
			"backup_quantity": {
				Type:          schema.TypeInt,
				Optional:      true,
				ValidateFunc:  validation.IntBetween(2, 99999),
				ConflictsWith: []string{"time_period"},
			},
			"time_period": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(2, 99999),
			},
		},
	}
}

func validateTimeFormat(val interface{}, key string) (warns []string, errs []error) {
	time := val.(string)
	regex := regexp.MustCompile("^(\\d{2}):(\\d{2})$")
	result := regex.FindStringSubmatch(time)
	if len(result) == 0 {
		errs = append(errs, fmt.Errorf("Wrong time format %v, please check your input (HH:MM)", time))
		return
	}
	hour, err := strconv.Atoi(result[1])
	if err != nil {
		errs = append(errs, fmt.Errorf("Wrong time hour, it must be UTC time string (HH): %s", err))
		return
	}
	if hour > 23 {
		errs = append(errs, fmt.Errorf("Hour must be between 0 and 23 inclusive, got: %d", hour))
	}
	minute, err := strconv.Atoi(result[2])
	if err != nil {
		errs = append(errs, fmt.Errorf("Wrong time minute, it must be UTC time string (MM): %s", err))
		return
	}
	if minute > 59 {
		errs = append(errs, fmt.Errorf("Minute must be between 0 and 59 inclusive, got: %d", minute))
	}
	return
}

func validateBackupWeekDay(val interface{}, key string) (warns []string, errs []error) {
	weeklyDays := []string{"MO", "TU", "WE", "TH", "FR", "SA", "SU"}
	days := val.(string)
	regex := regexp.MustCompile("([^A-Z,]+)")
	result := regex.FindAllStringSubmatch(days, -1)
	if len(result) > 0 {
		errs = append(errs, fmt.Errorf("Weekly days cannot contain lowercase letters, "+
			"numbers, whitespaces and special characters except commas"))
		return
	}
	regex = regexp.MustCompile("(\\w+)")
	result = regex.FindAllStringSubmatch(days, -1)
	if len(result) == 0 {
		errs = append(errs, fmt.Errorf("Wrong weekly days format %v, please check your input", days))
		return
	}
	for _, v := range result {
		if !utils.StrSliceContains(weeklyDays, v[0]) {
			errs = append(errs, fmt.Errorf("Wrong weekly days value \"%v\", the valid values of \"days\" are: "+
				"\"MO\", \"TU\", \"WE\", \"TH\", \"FR\", \"SA\", \"SU\"", v[0]))
			return
		}
	}
	return
}

func resourceCBRPolicyV3Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.CbrV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating Huaweicloud CBR client: %s", err)
	}

	enabled := d.Get("enabled").(bool)

	createOpts := policies.CreateOpts{
		Name:                d.Get("name").(string),
		OperationDefinition: resourceCBRPolicyV3OpDefinition(d),
		Enabled:             &enabled,
		OperationType:       d.Get("protection_type").(string),
		Trigger: &policies.Trigger{
			Properties: policies.TriggerProperties{
				Pattern: resourceCBRPolicyV3TriggerPattern(d),
			},
		},
	}

	cbrPolicy, err := policies.Create(client, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating Huaweicloud CBR policy: %s", err)
	}

	d.SetId(cbrPolicy.ID)

	return resourceCBRPolicyV3Read(d, meta)
}

func resourceCBRPolicyV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.CbrV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating Huaweicloud CBR client: %s", err)
	}

	cbrPolicy, err := policies.Get(client, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "Error retrieving CBRv3 policy")
	}

	log.Printf("[DEBUG] Retrieved policy %s: %+v", d.Id(), cbrPolicy)
	operationDefinition := cbrPolicy.OperationDefinition
	mErr := multierror.Append(nil,
		d.Set("enabled", cbrPolicy.Enabled),
		d.Set("name", cbrPolicy.Name),
		d.Set("protection_type", cbrPolicy.OperationType),
		d.Set("destination_region", operationDefinition.DestinationRegion),
		d.Set("destination_project_id", operationDefinition.DestinationProjectID),
		setTriggerPattern(d, cbrPolicy.Trigger.Properties.Pattern),
		d.Set("region", GetRegion(d, config)),
	)
	if mErr.ErrorOrNil() != nil {
		return mErr
	}
	if operationDefinition.MaxBackups != -1 {
		d.Set("backup_quantity", operationDefinition.MaxBackups)
	}
	if operationDefinition.RetentionDurationDays != -1 {
		d.Set("time_period", operationDefinition.RetentionDurationDays)
	}

	return nil
}

func resourceCBRPolicyV3Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.CbrV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating Huaweicloud CBR client: %s", err)
	}

	var updateOpts policies.UpdateOpts

	if d.HasChange("name") {
		newName := d.Get("name")
		updateOpts.Name = newName.(string)
	}

	if d.HasChange("enabled") {
		enabled := d.Get("enabled").(bool)
		updateOpts.Enabled = &enabled
	}

	if d.HasChange("backup_cycle") {
		pattern := resourceCBRPolicyV3TriggerPattern(d)
		updateOpts.Trigger = &policies.Trigger{
			Properties: policies.TriggerProperties{
				Pattern: pattern,
			},
		}
	}

	if d.HasChanges("backup_quantity", "time_period", "destination_region") {
		opDefinition := resourceCBRPolicyV3OpDefinition(d)
		updateOpts.OperationDefinition = opDefinition
	}

	_, err = policies.Update(client, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error updating Huaweicloud CBR policy: %s", err)
	}

	return resourceCBRPolicyV3Read(d, meta)
}

func resourceCBRPolicyV3Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.CbrV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating Huaweicloud CBR client: %s", err)
	}

	if err = policies.Delete(client, d.Id()).ExtractErr(); err != nil {
		return fmt.Errorf("Error deleting Huaweicloud CBR policy: %s", err)
	}

	d.SetId("")
	return nil
}

func resourceCBRPolicyV3OpDefinition(d *schema.ResourceData) *policies.PolicyODCreate {
	maxBackups, ok1 := d.GetOk("backup_quantity")
	durationDays, ok2 := d.GetOk("time_period")
	policyODCreate := policies.PolicyODCreate{}
	if destinationProjectID, ok3 := d.GetOk("destination_project_id"); ok3 {
		policyODCreate.DestinationProjectID = destinationProjectID.(string)
		policyODCreate.DestinationRegion = d.Get("destination_region").(string)
	}
	if !ok1 && !ok2 {
		policyODCreate.MaxBackups = -1
		policyODCreate.RetentionDurationDays = -1
		return &policyODCreate
	}
	policyODCreate.MaxBackups = maxBackups.(int)
	policyODCreate.RetentionDurationDays = durationDays.(int)

	return &policyODCreate
}

func makeSchedule(frequency, backupType, duration, time string) string {
	timeSlice := strings.Split(time, ":")
	schedule := fmt.Sprintf("FREQ=%s;%s=%s;BYHOUR=%s;BYMINUTE=%s",
		frequency, backupType, duration, timeSlice[0], timeSlice[1])
	return schedule
}

func resourceCBRPolicyV3TriggerPattern(d *schema.ResourceData) []string {
	defaultWeekDays, defaultIntervalDays := "MO,TU,WE,TH,FR,SA,SU", "1"
	triggerPatternRaw := d.Get("backup_cycle").([]interface{})

	frequency := triggerPatternRaw[0].(map[string]interface{})["frequency"].(string)
	var backupType, duration string
	if frequency == "WEEKLY" {
		backupType = "BYDAY"
		if v, ok := triggerPatternRaw[0].(map[string]interface{})["days"]; ok {
			duration = v.(string)
		} else {
			duration = defaultWeekDays
		}
	} else {
		backupType = "INTERVAL"
		if v, ok := triggerPatternRaw[0].(map[string]interface{})["interval"]; ok {
			duration = strconv.Itoa(v.(int))
		} else {
			duration = defaultIntervalDays
		}
	}
	backupTimes := triggerPatternRaw[0].(map[string]interface{})["execution_times"].([]interface{})
	patterns := make([]string, len(backupTimes))
	for i, v := range backupTimes {
		patterns[i] = makeSchedule(frequency, backupType, duration, v.(string))
	}
	return patterns
}

func setTriggerPattern(d *schema.ResourceData, patterns []string) error {
	schedules := make([]map[string]interface{}, 1)
	schedule := make(map[string]interface{})
	if strings.Contains(patterns[0], "WEEKLY") {
		schedule["frequency"] = "WEEKLY"
		regexExp := regexp.MustCompile("BYDAY=([\\w,]+);")
		schedule["days"] = regexExp.FindStringSubmatch(patterns[0])[1]
	} else {
		schedule["frequency"] = "DAILY"
		regexExp := regexp.MustCompile("INTERVAL=([\\d]+);")
		num, err := strconv.Atoi(regexExp.FindStringSubmatch(patterns[0])[1])
		if err != nil {
			return err
		}
		schedule["interval"] = num
	}
	backupTimes := make([]string, len(patterns))
	for i, v := range patterns {
		regexExp := regexp.MustCompile("BYHOUR=(\\d+);BYMINUTE=(\\d+)")
		times := regexExp.FindStringSubmatch(v)
		backupTimes[i] = strings.Join([]string{times[1], times[2]}, ":")
	}
	schedule["execution_times"] = backupTimes
	schedules[0] = schedule
	if err := d.Set("backup_cycle", schedules); err != nil {
		return err
	}
	return nil
}
