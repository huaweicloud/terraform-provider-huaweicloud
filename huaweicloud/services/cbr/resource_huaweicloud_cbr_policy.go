package cbr

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk/openstack/cbr/v3/policies"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CBR POST /v3/{project_id}/policies
// @API CBR GET /v3/{project_id}/policies/{policy_id}
// @API CBR PUT /v3/{project_id}/policies/{policy_id}
// @API CBR DELETE /v3/{project_id}/policies/{policy_id}
func ResourcePolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicyCreate,
		ReadContext:   resourcePolicyRead,
		UpdateContext: resourcePolicyUpdate,
		DeleteContext: resourcePolicyDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region where the policy is located.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The policy name.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether to enable the CBR policy.",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The protection type of the CBR policy.",
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
							Description:  "The number of days between each backup.",
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
							Description: "The weekly backup time.",
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
							Description: "The execution time of the policy.",
						},
					},
				},
				Description: "The scheduling rule for the CBR policy backup execution.",
			},
			"destination_region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the replication destination region.",
			},
			"destination_project_id": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"destination_region"},
				Description:  "The ID of the replication destination project.",
			},
			"backup_quantity": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The maximum number of retained backups.",
			},
			"time_period": {
				Type:          schema.TypeInt,
				Optional:      true,
				ConflictsWith: []string{"backup_quantity"},
				Description:   "The duration (in days) for retained backups.",
			},
			"long_term_retention": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"daily": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The latest backup of each day is saved in the long term.",
						},
						"weekly": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The latest backup of each week is saved in the long term.",
						},
						"monthly": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The latest backup of each month is saved in the long term.",
						},
						"yearly": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The latest backup of each year is saved in the long term.",
						},
						"full_backup_interval": {
							Type:     schema.TypeInt,
							Optional: true,
							Description: "How often (after how many incremental backups) a full backup is " +
								"performed.",
						},
					},
				},
				RequiredWith: []string{"backup_quantity", "time_zone"},
				Description:  "The long-term retention rules.",
			},
			"time_zone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^UTC[+-]\d{2}:00$`),
					"The time zone must be in UTC format, such as 'UTC+08:00'."),
				Description: "The UTC time zone.",
			},
			"enable_acceleration": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Description: utils.SchemaDesc(
					"Whether to enable the acceleration function to shorten the replication time for cross-region",
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
		},
	}
}

func makeSchedule(frequency, backupType, duration, time string) (string, error) {
	timeSlice := strings.Split(time, ":")
	if len(timeSlice) != 2 {
		return "", fmt.Errorf("invalid time format (%s), should be HH:MM", time)
	}
	schedule := fmt.Sprintf("FREQ=%s;%s=%s;BYHOUR=%s;BYMINUTE=%s",
		frequency, backupType, duration, timeSlice[0], timeSlice[1])
	return schedule, nil
}

func buildPolicyBackupSchedule(backupCycleRaw []interface{}) ([]string, error) {
	var frequency, backupType, duration string

	rawInfo := backupCycleRaw[0].(map[string]interface{})
	// If 'days' is set, the value of 'interval' will be 0, relatively, the value of 'days' will be "".
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

func buildPolicyOpDefinition(d *schema.ResourceData) *policies.PolicyODCreate {
	policyODCreate := policies.PolicyODCreate{}

	if destinationProjectID, ok := d.GetOk("destination_project_id"); ok {
		policyODCreate.DestinationProjectID = destinationProjectID.(string)
		policyODCreate.DestinationRegion = d.Get("destination_region").(string)
		// Users with the permission 'OP_GATED_CSBS_REP_ACCELERATION' cannot use this parameter, and an error will be
		// reported when 'true' is configured: 'BackupService.9900', 'invalid key enable_acceleration'.
		// The default value is 'false' in the API definition (2023.06.16).
		policyODCreate.EnableAcceleration = d.Get("enable_acceleration").(bool)
	}

	// The backup_quantity and time_period are both left blank means the backups are retained permanently
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
		policyODCreate.FullBackupInterval = utils.Int(d.Get("long_term_retention.0.full_backup_interval").(int))
	}

	return &policyODCreate
}

func resourcePolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.CbrV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CBR v3 client: %s", err)
	}

	schedule, err := buildPolicyBackupSchedule(d.Get("backup_cycle").([]interface{}))
	if err != nil {
		return diag.Errorf("error parsing the format of backup cycle: %s", err)
	}
	createOpts := policies.CreateOpts{
		Name:                d.Get("name").(string),
		OperationType:       d.Get("type").(string),
		Enabled:             utils.Bool(d.Get("enabled").(bool)),
		OperationDefinition: buildPolicyOpDefinition(d),
		Trigger: &policies.Trigger{
			Properties: policies.TriggerProperties{
				Pattern: schedule,
			},
		},
	}

	cbrPolicy, err := policies.Create(client, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating CBR policy: %s", err)
	}

	d.SetId(cbrPolicy.ID)

	return resourcePolicyRead(ctx, d, meta)
}

func flattenPolicyBackupCycle(schedules []string) ([]map[string]interface{}, error) {
	schedule := make(map[string]interface{})
	// The value obtained from API is a string containing 'days', 'interval' and 'execution_times',
	// so it should be extracted when setting the corresponding value to state.
	if strings.Contains(schedules[0], "WEEKLY") {
		result := regexp.MustCompile(`BYDAY=([\w,]+);`).FindStringSubmatch(schedules[0])
		// The right length of the string match is two, first is the match result of the regex string,
		// second is the match result of the regex group.
		if len(result) != 2 {
			return nil, fmt.Errorf("invalid format of the weekly days in API response")
		}
		schedule["days"] = result[1]
	} else {
		result := regexp.MustCompile(`INTERVAL=([\d]+);`).FindStringSubmatch(schedules[0])
		if len(result) != 2 {
			return nil, fmt.Errorf("invalid format of backup interval in API response")
		}
		num, err := strconv.Atoi(result[1])
		if err != nil {
			return nil, err
		}
		schedule["interval"] = num
	}
	backupTimeList := make([]string, len(schedules))
	for i, v := range schedules {
		times := regexp.MustCompile(`BYHOUR=(\d+);BYMINUTE=(\d+)`).FindStringSubmatch(v)
		if len(times) != 3 {
			return nil, fmt.Errorf("invalid format of backup time in API response")
		}
		backupTimeList[i] = strings.Join([]string{times[1], times[2]}, ":")
	}
	schedule["execution_times"] = backupTimeList

	return []map[string]interface{}{schedule}, nil
}

func flattenLongTermRetention(resp *policies.PolicyODCreate) []map[string]interface{} {
	if resp.DailyBackups == 0 && resp.WeekBackups == 0 && resp.MonthBackups == 0 && resp.YearBackups == 0 {
		return nil
	}
	return []map[string]interface{}{
		{
			"daily":                resp.DailyBackups,
			"weekly":               resp.WeekBackups,
			"monthly":              resp.MonthBackups,
			"yearly":               resp.YearBackups,
			"full_backup_interval": resp.FullBackupInterval,
		},
	}
}

func resourcePolicyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.CbrV3Client(region)
	if err != nil {
		return diag.Errorf("error creating CBR v3 client: %s", err)
	}

	policyId := d.Id()
	resp, err := policies.Get(client, policyId).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "CBR policy")
	}

	log.Printf("[DEBUG] Retrieved policy (%s): %#v", policyId, resp)
	operationDefinition := resp.OperationDefinition
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("enabled", resp.Enabled),
		d.Set("name", resp.Name),
		d.Set("type", resp.OperationType),
		d.Set("destination_region", operationDefinition.DestinationRegion),
		d.Set("destination_project_id", operationDefinition.DestinationProjectID),
	)

	backupCycle, err := flattenPolicyBackupCycle(resp.Trigger.Properties.Pattern)
	if err != nil {
		return diag.FromErr(err)
	}
	mErr = multierror.Append(mErr, d.Set("backup_cycle", backupCycle))

	if operationDefinition.MaxBackups != -1 {
		mErr = multierror.Append(mErr,
			d.Set("backup_quantity", operationDefinition.MaxBackups),
			d.Set("long_term_retention", flattenLongTermRetention(operationDefinition)),
		)
		if operationDefinition.Timezone != "" {
			mErr = multierror.Append(mErr, d.Set("time_zone", operationDefinition.Timezone))
		}
	}
	if operationDefinition.RetentionDurationDays != -1 {
		mErr = multierror.Append(mErr, d.Set("time_period", operationDefinition.RetentionDurationDays))
	}

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving policy resource fields: %s", err)
	}
	return nil
}

func resourcePolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.CbrV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CBR v3 client: %s", err)
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
		schedule, err := buildPolicyBackupSchedule(d.Get("backup_cycle").([]interface{}))
		if err != nil {
			return diag.Errorf("error parsing the format of backup cycle: %s", err)
		}
		updateOpts.Trigger = &policies.Trigger{
			Properties: policies.TriggerProperties{
				Pattern: schedule,
			},
		}
	}
	if d.HasChangesExcept("name", "enabled", "backup_cycle") {
		updateOpts.OperationDefinition = buildPolicyOpDefinition(d)
	}

	_, err = policies.Update(client, d.Id(), updateOpts).Extract()
	if err != nil {
		return diag.Errorf("error updating CBR policy: %s", err)
	}

	return resourcePolicyRead(ctx, d, meta)
}

func resourcePolicyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.CbrV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CBR v3 client: %s", err)
	}

	if err = policies.Delete(client, d.Id()).ExtractErr(); err != nil {
		return diag.Errorf("error deleting CBR policy: %s", err)
	}

	return nil
}
