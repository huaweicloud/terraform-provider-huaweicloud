package cbr

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/chnsz/golangsdk/openstack/cbr/v3/policies"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func ResourceCBRPolicyV3() *schema.Resource {
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
			"type": {
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
						"interval": {
							Type:         schema.TypeInt,
							Optional:     true,
							ExactlyOneOf: []string{"backup_cycle.0.days"},
							ValidateFunc: validation.IntBetween(1, 30),
						},
						"days": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringMatch(
								regexp.MustCompile("^(?:MO|TU|WE|TH|FR|SA|SU)(?:,(?:MO|TU|WE|TH|FR|SA|SU))*$"),
								"the string cannot contain lowercase, letters, numbers, whitespaces and special "+
									"characters except commas. The valid string of weekly date are:"+
									"MO, TU, WE, TH, FR, SA, SU.",
							),
						},
						"execution_times": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 24,
							Elem: &schema.Schema{
								Type: schema.TypeString,
								ValidateFunc: validation.StringMatch(
									regexp.MustCompile("^[0-1][0-9]|2[0-3]:[0-5][0-9]$"),
									"the time format should be HH:MM",
								),
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
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(2, 99999),
			},
			"time_period": {
				Type:          schema.TypeInt,
				Optional:      true,
				ValidateFunc:  validation.IntBetween(2, 99999),
				ConflictsWith: []string{"backup_quantity"},
			},
			"long_term_retention": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"daily": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 100),
						},
						"weekly": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 100),
						},
						"monthly": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 100),
						},
						"yearly": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(1, 100),
						},
					},
				},
				RequiredWith: []string{"backup_quantity", "time_zone"},
			},
			"time_zone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^UTC[+-]\d{2}:00$`),
					"The time zone must be in UTC format, such as 'UTC+08:00'."),
			},
		},
	}
}

func resourceCBRPolicyV3Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.CbrV3Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating Huaweicloud CBR client: %s", err)
	}

	enabled := d.Get("enabled").(bool)

	schedule, err := resourceCBRPolicyV3BackupSchedule(d)
	if err != nil {
		return fmtp.Errorf("Error to parse the format of backup cycle: %s", err)
	}
	createOpts := policies.CreateOpts{
		Name:                d.Get("name").(string),
		OperationType:       d.Get("type").(string),
		Enabled:             &enabled,
		OperationDefinition: resourceCBRPolicyV3OpDefinition(d),
		Trigger: &policies.Trigger{
			Properties: policies.TriggerProperties{
				Pattern: schedule,
			},
		},
	}

	cbrPolicy, err := policies.Create(client, createOpts).Extract()
	if err != nil {
		return fmtp.Errorf("Error creating Huaweicloud CBR policy: %s", err)
	}

	d.SetId(cbrPolicy.ID)

	return resourceCBRPolicyV3Read(d, meta)
}

func resourceCBRPolicyV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.CbrV3Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating Huaweicloud CBR client: %s", err)
	}

	cbrPolicy, err := policies.Get(client, d.Id()).Extract()
	if err != nil {
		return common.CheckDeleted(d, err, "Error retrieving CBRv3 policy")
	}

	logp.Printf("[DEBUG] Retrieved policy %s: %+v", d.Id(), cbrPolicy)
	operationDefinition := cbrPolicy.OperationDefinition
	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("enabled", cbrPolicy.Enabled),
		d.Set("name", cbrPolicy.Name),
		d.Set("type", cbrPolicy.OperationType),
		d.Set("destination_region", operationDefinition.DestinationRegion),
		d.Set("destination_project_id", operationDefinition.DestinationProjectID),
		setCBRPolicyV3BackupCycle(d, cbrPolicy.Trigger.Properties.Pattern),
	)

	if operationDefinition.MaxBackups != -1 {
		mErr = multierror.Append(mErr,
			d.Set("backup_quantity", operationDefinition.MaxBackups),
			setCBRPolicyV3LongTermRetention(d, operationDefinition),
		)
		if operationDefinition.Timezone != "" {
			mErr = multierror.Append(mErr, d.Set("time_zone", operationDefinition.Timezone))
		}
	}
	if operationDefinition.RetentionDurationDays != -1 {
		mErr = multierror.Append(mErr, d.Set("time_period", operationDefinition.RetentionDurationDays))
	}

	return mErr.ErrorOrNil()
}

func resourceCBRPolicyV3Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.CbrV3Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating Huaweicloud CBR client: %s", err)
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
		schedule, err := resourceCBRPolicyV3BackupSchedule(d)
		if err != nil {
			return fmtp.Errorf("Error to parse the format of backup cycle: %s", err)
		}
		updateOpts.Trigger = &policies.Trigger{
			Properties: policies.TriggerProperties{
				Pattern: schedule,
			},
		}
	}
	if d.HasChanges("backup_quantity", "time_period", "destination_region", "long_term_retention", "time_zone") {
		opDefinition := resourceCBRPolicyV3OpDefinition(d)
		updateOpts.OperationDefinition = opDefinition
	}

	_, err = policies.Update(client, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmtp.Errorf("Error updating Huaweicloud CBR policy: %s", err)
	}

	return resourceCBRPolicyV3Read(d, meta)
}

func resourceCBRPolicyV3Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.CbrV3Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating Huaweicloud CBR client: %s", err)
	}

	if err = policies.Delete(client, d.Id()).ExtractErr(); err != nil {
		return fmtp.Errorf("Error deleting Huaweicloud CBR policy: %s", err)
	}

	d.SetId("")
	return nil
}

func resourceCBRPolicyV3OpDefinition(d *schema.ResourceData) *policies.PolicyODCreate {
	policyODCreate := policies.PolicyODCreate{}

	if destinationProjectID, ok3 := d.GetOk("destination_project_id"); ok3 {
		policyODCreate.DestinationProjectID = destinationProjectID.(string)
		policyODCreate.DestinationRegion = d.Get("destination_region").(string)
	}

	//the backup_quantity and time_period are both left blank means the backups are retained permanently
	maxBackups, ok1 := d.GetOk("backup_quantity")
	durationDays, ok2 := d.GetOk("time_period")
	if !ok1 && !ok2 {
		policyODCreate.MaxBackups = -1
		policyODCreate.RetentionDurationDays = -1
		return &policyODCreate
	}
	policyODCreate.MaxBackups = maxBackups.(int)
	policyODCreate.RetentionDurationDays = durationDays.(int)

	if _, ok := d.GetOk("long_term_retention"); ok {
		policyODCreate.DailyBackups = d.Get("long_term_retention.0.daily").(int)
		policyODCreate.WeekBackups = d.Get("long_term_retention.0.weekly").(int)
		policyODCreate.MonthBackups = d.Get("long_term_retention.0.monthly").(int)
		policyODCreate.YearBackups = d.Get("long_term_retention.0.yearly").(int)
		policyODCreate.Timezone = d.Get("time_zone").(string)
	}

	return &policyODCreate
}

func makeSchedule(frequency, backupType, duration, time string) (string, error) {
	timeSlice := strings.Split(time, ":")
	if len(timeSlice) != 2 {
		return "", fmtp.Errorf("Wrong time format (%s), should be HH:MM", time)
	}
	schedule := fmt.Sprintf("FREQ=%s;%s=%s;BYHOUR=%s;BYMINUTE=%s",
		frequency, backupType, duration, timeSlice[0], timeSlice[1])
	return schedule, nil
}

func resourceCBRPolicyV3BackupSchedule(d *schema.ResourceData) ([]string, error) {
	var frequency, backupType, duration string

	backupCycleRaw := d.Get("backup_cycle").([]interface{})
	rawInfo := backupCycleRaw[0].(map[string]interface{})

	//If 'days' is set, the value of 'interval' will be 0, relatively, the value of 'days' will be "".
	if rawInfo["days"] != "" {
		frequency = "WEEKLY"
		backupType = "BYDAY"
		duration = rawInfo["days"].(string)
	} else {
		frequency = "DAILY"
		backupType = "INTERVAL"
		duration = strconv.Itoa(rawInfo["interval"].(int))
	}
	backupTimes := rawInfo["execution_times"].([]interface{})
	schedules := make([]string, len(backupTimes))
	for i, v := range backupTimes {
		schedule, err := makeSchedule(frequency, backupType, duration, v.(string))
		if err != nil {
			return schedules, err
		}
		schedules[i] = schedule
	}
	return schedules, nil
}

func setCBRPolicyV3LongTermRetention(d *schema.ResourceData, resp *policies.PolicyODCreate) error {
	if resp.DailyBackups == 0 && resp.WeekBackups == 0 && resp.MonthBackups == 0 && resp.YearBackups == 0 {
		return nil
	}
	result := make([]map[string]interface{}, 1)
	retention := map[string]interface{}{
		"daily":   resp.DailyBackups,
		"weekly":  resp.WeekBackups,
		"monthly": resp.MonthBackups,
		"yearly":  resp.YearBackups,
	}
	result[0] = retention
	return d.Set("long_term_retention", result)
}

func setCBRPolicyV3BackupCycle(d *schema.ResourceData, schedules []string) error {
	schedule := make(map[string]interface{})
	//The value obtained from API is a string containing 'days', 'interval' and 'execution_times',
	//so it should be extracted when setting the corresponding value to state.
	if strings.Contains(schedules[0], "WEEKLY") {
		regexExp := regexp.MustCompile("BYDAY=([\\w,]+);")
		result := regexExp.FindStringSubmatch(schedules[0])
		//The right length of the string match is two, first is the match result of the regex string,
		//second is the match result of the regex group.
		if len(result) != 2 {
			return fmtp.Errorf("Wrong weekly days format in API response")
		}
		schedule["days"] = result[1]
	} else {
		regexExp := regexp.MustCompile("INTERVAL=([\\d]+);")
		result := regexExp.FindStringSubmatch(schedules[0])
		if len(result) != 2 {
			return fmtp.Errorf("Wrong backup interval format in API response")
		}
		num, err := strconv.Atoi(result[1])
		if err != nil {
			return err
		}
		schedule["interval"] = num
	}
	backupTimeList := make([]string, len(schedules))
	for i, v := range schedules {
		regexExp := regexp.MustCompile("BYHOUR=(\\d+);BYMINUTE=(\\d+)")
		times := regexExp.FindStringSubmatch(v)
		if len(times) != 3 {
			return fmtp.Errorf("Wrong backup time format in API response")
		}
		backupTimeList[i] = strings.Join([]string{times[1], times[2]}, ":")
	}
	schedule["execution_times"] = backupTimeList
	backupCycle := []map[string]interface{}{schedule}
	if err := d.Set("backup_cycle", backupCycle); err != nil {
		return err
	}
	return nil
}
