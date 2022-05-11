package aom

import (
	"context"
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	aom "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/aom/v2/model"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func ResourceAlarmRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAlarmRuleCreate,
		ReadContext:   resourceAlarmRuleRead,
		UpdateContext: resourceAlarmRuleUpdate,
		DeleteContext: resourceAlarmRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 100),
					validation.StringMatch(regexp.MustCompile(
						"^[\u4e00-\u9fa5A-Za-z0-9]([\u4e00-\u9fa5-_A-Za-z0-9]*[\u4e00-\u9fa5A-Za-z0-9])?$"),
						"The name can only consist of letters, digits, underscores (_),"+
							" hyphens (-) and chinese characters, and it must start and end with letters,"+
							" digits or chinese characters."),
				),
			},
			"metric_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"dimensions": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
			"period": {
				Type:     schema.TypeInt,
				Required: true,
				ValidateFunc: validation.IntInSlice([]int{
					60000, 300000, 900000, 3600000,
				}),
			},
			"unit": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 32),
			},
			"comparison_operator": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					">=", ">", "<=", "<", "=",
				}, false),
			},
			"statistic": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"maximum", "minimum", "average", "sum", "sampleCount",
				}, false),
			},
			"threshold": {
				Type:     schema.TypeString,
				Required: true,
			},
			"evaluation_periods": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 5),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 1000),
			},
			"alarm_level": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      2,
				ValidateFunc: validation.IntBetween(1, 4),
			},
			"alarm_actions": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"alarm_action_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  true,
			},
			"ok_actions": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"insufficient_data_actions": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"alarm_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"state_value": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state_reason": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildActionOpts(rawActions []interface{}) *[]string {
	if len(rawActions) == 0 {
		return nil
	}
	actions := make([]string, len(rawActions))
	for i, raw := range rawActions {
		actions[i] = raw.(string)
	}
	return &actions
}

func buildDimensionsOpts(rawDimensions []interface{}) []aom.Dimension {
	if len(rawDimensions) == 0 {
		return nil
	}
	dimensions := make([]aom.Dimension, len(rawDimensions))
	for i, rawdimension := range rawDimensions {
		dimension := rawdimension.(map[string]interface{})
		dimensions[i] = aom.Dimension{
			Name:  dimension["name"].(string),
			Value: dimension["value"].(string),
		}
	}
	return dimensions
}

func buildAlarmLevelOpts(alarmLevel int) *aom.UpdateAlarmRuleParamAlarmLevel {
	var alarmLevelToReq aom.UpdateAlarmRuleParamAlarmLevel
	switch alarmLevel {
	case 1:
		alarmLevelToReq = aom.GetUpdateAlarmRuleParamAlarmLevelEnum().E_1
	case 2:
		alarmLevelToReq = aom.GetUpdateAlarmRuleParamAlarmLevelEnum().E_2
	case 3:
		alarmLevelToReq = aom.GetUpdateAlarmRuleParamAlarmLevelEnum().E_3
	case 4:
		alarmLevelToReq = aom.GetUpdateAlarmRuleParamAlarmLevelEnum().E_4
	default:
		log.Printf("[WARN] alarm level invalid: %d", alarmLevel)
		return nil
	}

	return &alarmLevelToReq
}

func buildStatisticOpts(statistic string) *aom.UpdateAlarmRuleParamStatistic {
	var statisticToReq aom.UpdateAlarmRuleParamStatistic
	switch statistic {
	case "maximum":
		statisticToReq = aom.GetUpdateAlarmRuleParamStatisticEnum().MAXIMUM
	case "minimum":
		statisticToReq = aom.GetUpdateAlarmRuleParamStatisticEnum().MINIMUM
	case "average":
		statisticToReq = aom.GetUpdateAlarmRuleParamStatisticEnum().AVERAGE
	case "sum":
		statisticToReq = aom.GetUpdateAlarmRuleParamStatisticEnum().SUM
	case "sampleCount":
		statisticToReq = aom.GetUpdateAlarmRuleParamStatisticEnum().SAMPLE_COUNT
	default:
		log.Printf("[WARN] statistic invalid: %s", statistic)
		return nil
	}

	return &statisticToReq
}

func resourceAlarmRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.HcAomV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("error creating AOM client: %s", err)
	}

	createOpts := aom.AlarmRuleParam{
		AlarmRuleName:           d.Get("name").(string),
		AlarmDescription:        utils.String(d.Get("description").(string)),
		AlarmLevel:              int32(d.Get("alarm_level").(int)),
		IdTurnOn:                utils.Bool(true),
		AlarmActions:            buildActionOpts(d.Get("alarm_actions").([]interface{})),
		ActionEnabled:           utils.Bool(d.Get("alarm_action_enabled").(bool)),
		OkActions:               buildActionOpts(d.Get("ok_actions").([]interface{})),
		InsufficientDataActions: buildActionOpts(d.Get("insufficient_data_actions").([]interface{})),

		Namespace:  d.Get("namespace").(string),
		MetricName: d.Get("metric_name").(string),
		Dimensions: buildDimensionsOpts(d.Get("dimensions").([]interface{})),

		Unit:               d.Get("unit").(string),
		Threshold:          d.Get("threshold").(string),
		Statistic:          d.Get("statistic").(string),
		Period:             int32(d.Get("period").(int)),
		EvaluationPeriods:  int32(d.Get("evaluation_periods").(int)),
		ComparisonOperator: d.Get("comparison_operator").(string),
	}

	log.Printf("[DEBUG] Create %s Options: %#v", createOpts.AlarmRuleName, createOpts)

	createReq := aom.AddAlarmRuleRequest{
		Body: &createOpts,
	}
	response, err := client.AddAlarmRule(&createReq)
	if err != nil {
		return diag.Errorf("error creating AOM alarm rule %s: %s", createOpts.AlarmRuleName, err)
	}

	d.SetId(strconv.FormatInt(*response.AlarmRuleId, 10))

	return resourceAlarmRuleRead(ctx, d, meta)
}

func resourceAlarmRuleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.HcAomV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("error creating AOM client: %s", err)
	}

	response, err := client.ShowAlarmRule(&aom.ShowAlarmRuleRequest{AlarmRuleId: d.Id()})
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving AOM alarm rule")
	}

	allRules := *response.Thresholds
	if len(allRules) != 1 {
		return diag.Errorf("error retrieving AOM alarm rule %s", d.Id())
	}
	rule := allRules[0]
	log.Printf("[DEBUG] Retrieved AOM alarm rule %s: %#v", d.Id(), rule)

	alarm_level, _ := strconv.Atoi(*rule.AlarmLevel)

	mErr := multierror.Append(nil,
		d.Set("region", config.GetRegion(d)),
		d.Set("name", rule.AlarmRuleName),
		d.Set("description", rule.AlarmDescription),
		d.Set("alarm_level", alarm_level),
		d.Set("metric_name", rule.MetricName),
		d.Set("alarm_actions", rule.AlarmActions),
		d.Set("ok_actions", rule.OkActions),
		d.Set("alarm_enabled", rule.IdTurnOn),
		d.Set("alarm_action_enabled", rule.ActionEnabled),
		d.Set("comparison_operator", rule.ComparisonOperator),
		d.Set("evaluation_periods", rule.EvaluationPeriods),
		d.Set("insufficient_data_actions", rule.InsufficientDataActions),
		d.Set("namespace", rule.Namespace),
		d.Set("period", rule.Period),
		d.Set("state_value", rule.StateValue),
		d.Set("state_reason", rule.StateReason),
		d.Set("statistic", rule.Statistic),
		d.Set("threshold", rule.Threshold),
		d.Set("unit", rule.Unit),
	)

	var dimensions []map[string]interface{}
	for _, pairObject := range *rule.Dimensions {
		dimension := make(map[string]interface{})
		dimension["name"] = pairObject.Name
		dimension["value"] = pairObject.Value

		dimensions = append(dimensions, dimension)
	}
	mErr = multierror.Append(mErr, d.Set("dimensions", dimensions))

	if err := mErr.ErrorOrNil(); err != nil {
		return fmtp.DiagErrorf("error setting AOM alarm rule fields: %w", err)
	}

	return nil
}

func resourceAlarmRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.HcAomV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("error creating AOM client: %s", err)
	}

	updateOpts := aom.UpdateAlarmRuleParam{
		AlarmRuleName: d.Get("name").(string),
	}

	if d.HasChange("description") {
		updateOpts.AlarmDescription = utils.String(d.Get("description").(string))
	}
	if d.HasChange("alarm_level") {
		updateOpts.AlarmLevel = buildAlarmLevelOpts(d.Get("alarm_level").(int))
	}

	if d.HasChange("alarm_enabled") {
		updateOpts.IdTurnOn = utils.Bool(d.Get("alarm_enabled").(bool))
	}
	if d.HasChange("alarm_actions") {
		updateOpts.AlarmActions = buildActionOpts(d.Get("alarm_actions").([]interface{}))
	}
	if d.HasChange("ok_actions") {
		updateOpts.OkActions = buildActionOpts(d.Get("ok_actions").([]interface{}))
	}
	if d.HasChange("insufficient_data_actions") {
		updateOpts.InsufficientDataActions = buildActionOpts(d.Get("insufficient_data_actions").([]interface{}))
	}
	if d.HasChange("metric_name") {
		updateOpts.MetricName = utils.String(d.Get("metric_name").(string))
	}
	if d.HasChange("namespace") {
		updateOpts.Namespace = utils.String(d.Get("namespace").(string))
	}
	if d.HasChange("dimensions") {
		dimensions := buildDimensionsOpts(d.Get("dimensions").([]interface{}))
		updateOpts.Dimensions = &dimensions
	}
	if d.HasChange("period") {
		updateOpts.Period = utils.Int32(int32(d.Get("period").(int)))
	}
	if d.HasChange("unit") {
		updateOpts.Unit = utils.String(d.Get("unit").(string))
	}
	if d.HasChange("comparison_operator") {
		updateOpts.ComparisonOperator = utils.String(d.Get("comparison_operator").(string))
	}
	if d.HasChange("statistic") {
		updateOpts.Statistic = buildStatisticOpts(d.Get("statistic").(string))
	}
	if d.HasChange("threshold") {
		updateOpts.Threshold = utils.String(d.Get("threshold").(string))
	}
	if d.HasChange("evaluation_periods") {
		updateOpts.EvaluationPeriods = utils.Int32(int32(d.Get("evaluation_periods").(int)))
	}

	log.Printf("[DEBUG] Update %s Options: %#v", updateOpts.AlarmRuleName, updateOpts)

	reqOpts := aom.UpdateAlarmRuleRequest{
		Body: &updateOpts,
	}

	_, err = client.UpdateAlarmRule(&reqOpts)
	if err != nil {
		return diag.Errorf("error updating AOM alarm rule %s: %s", updateOpts.AlarmRuleName, err)
	}

	return resourceAlarmRuleRead(ctx, d, meta)
}

func resourceAlarmRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client, err := config.HcAomV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("error creating AOM client: %s", err)
	}

	_, err = client.DeleteAlarmRule(&aom.DeleteAlarmRuleRequest{AlarmRuleId: d.Id()})
	if err != nil {
		return diag.Errorf("error deleting AOM alarm rule %s: %s", d.Id(), err)
	}
	return nil
}
