package aom

import (
	"context"
	"encoding/json"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/internal/entity"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/internal/httpclient_go"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
	"io"
	"regexp"
	"time"
)

func ResourceAlarmPolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceAlarmPolicyRead,
		CreateContext: resourceAlarmPolicyCreate,
		DeleteContext: resourceAlarmPolicyDelete,
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
			"project_id": {
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
			"action_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"alarm_rule_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"alarm_rule_description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"alarm_rule_enable": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"alarm_rule_status": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"alarm_rule_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"alarm_delete_list": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"metric_alarm_spec":   schemaMetricAlarmSpe(),
			"event_alarm_spec":    schemeEventAlarmSpec(),
			"alarm_notifications": schemeAlarmNotifications(),
		},
	}
}

func schemaMetricAlarmSpe() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		ForceNew: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"monitor_type": {
					Type:     schema.TypeString,
					Required: true,
					ForceNew: true,
				},
				"resource_kind": {
					Type:     schema.TypeString,
					ForceNew: true,
					Optional: true,
				},
				"metric_kind": {
					Type:     schema.TypeString,
					Optional: true,
					ForceNew: true,
				},
				"alarm_rule_template_bind_enable": {
					Type:     schema.TypeBool,
					Optional: true,
					ForceNew: true,
				},
				"alarm_rule_template_id": {
					Type:     schema.TypeString,
					Optional: true,
					ForceNew: true,
				},
				"no_data_conditions": schemeNoDataConditions(),
				"alarm_tags":              schemeAlarmTags(),
				"trigger_conditions": schemeTriggerConditions(),
				"monitor_objects":    schemeMonitorObjects(),
				"recovery_conditions": {
					Type:     schema.TypeMap,
					Optional: true,
					ForceNew: true,
					Elem:         schema.TypeInt,
				},
			},
		},
	}
}

func schemeEventAlarmSpec() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		ForceNew: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"event_source": {
					Type:     schema.TypeString,
					Optional: true,
					ForceNew: true,
				},
				"monitor_objects":    schemeMonitorObjects(),
				"no_data_conditions": schemeNoDataConditions(),
				"alarm_tags":              schemeAlarmTags(),
			},
		},
	}
}

func schemeNoDataConditions() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"no_data_timeframe": {
					Type:     schema.TypeInt,
					Optional: true,
					ForceNew: true,
				},
				"no_data_alert_state": {
					Type:     schema.TypeString,
					Optional: true,
					ForceNew: true,
				},
				"notify_no_data": {
					Type:     schema.TypeBool,
					Optional: true,
					ForceNew: true,
				},
			},
		},
	}

}

func schemeAlarmTags() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"auto_tags": {
					Type:     schema.TypeList,
					Optional: true,
					ForceNew: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
				"custom_tags": {
					Type:     schema.TypeList,
					Optional: true,
					ForceNew: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
				"custom_annotations": {
					Type:     schema.TypeList,
					Optional: true,
					ForceNew: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
			},
		},
	}
}
	
func schemeTriggerConditions() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"metric_query_mode": {
					Type:     schema.TypeString,
					Optional: true,
					ForceNew: true,
				},
				"metric_namespace": {
					Type:     schema.TypeString,
					Optional: true,
					ForceNew: true,
				},
				"metric_name": {
					Type:     schema.TypeString,
					Optional: true,
					ForceNew: true,
				},
				"metric_labels": {
					Type:     schema.TypeList,
					Optional: true,
					ForceNew: true,
					Elem:         &schema.Schema{Type: schema.TypeString},
				},
				"promql": {
					Type:     schema.TypeString,
					Optional: true,
					ForceNew: true,
				},
				"trigger_times": {
					Type:     schema.TypeInt,
					Optional: true,
					ForceNew: true,
				},
				"trigger_interval": {
					Type:     schema.TypeString,
					Optional: true,
					ForceNew: true,
				},
				"aggregation_type": {
					Type:     schema.TypeString,
					Optional: true,
					ForceNew: true,
				},
				"aggregation_window": {
					Type:     schema.TypeString,
					Optional: true,
					ForceNew: true,
				},
				"operator": {
					Type:     schema.TypeString,
					Optional: true,
					ForceNew: true,
				},
				"trigger_type": {
					Type:     schema.TypeString,
					Optional: true,
					ForceNew: true,
				},
				"thresholds": {
					Type:     schema.TypeMap,
					Optional: true,
					ForceNew: true,
				},
			},
		},
	}
}

func schemeMonitorObjects() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem:      &schema.Schema{Type:  schema.TypeMap},
	}
}

func schemeAlarmNotifications() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		ForceNew: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"notification_type": {
					Type:     schema.TypeString,
					Optional: true,
					ForceNew: true,
				},
				"route_group_enable": {
					Type:     schema.TypeBool,
					Optional: true,
					ForceNew: true,
				},
				"route_group_rule": {
					Type:     schema.TypeString,
					Optional: true,
					ForceNew: true,
				},
				"inhibit_enable": {
					Type:     schema.TypeBool,
					Optional: true,
					ForceNew: true,
				},
				"inhibit_rule": {
					Type:     schema.TypeString,
					Optional: true,
					ForceNew: true,
				},
				"notification_enable": {
					Type:     schema.TypeBool,
					Optional: true,
					ForceNew: true,
				},
				"bind_notification_rule_id": {
					Type:     schema.TypeString,
					Optional: true,
					ForceNew: true,
				},
				"notify_resolved": {
					Type:     schema.TypeBool,
					Optional: true,
					ForceNew: true,
				},
			},
		},
	}
}

func resourceAlarmPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client,dErr := httpclient_go.NewHttpClientGo(config)	
	if dErr != nil {
		return dErr	
	}
	logp.Printf("[DEBUG] create param is : %+v", d.Get("metric_alarm_spec").([]interface{}))
	createOpts := entity.AddAlarmRuleParams{
		AlarmRuleName:         d.Get("alarm_rule_name").(string),
		EnterpriseProjectId:     d.Get("enterprise_project_id").(string),
		AlarmRuleDescription: d.Get("alarm_rule_description").(string),
		AlarmRuleEnable:      d.Get("alarm_rule_enable").(bool),
		AlarmRuleStatus:      d.Get("alarm_rule_status").(string),
		AlarmRuleType:        d.Get("alarm_rule_type").(string),
		MetricAlarmSpec:      buildMetricAlarmSpec(d.Get("metric_alarm_spec").([]interface{})),
		EventAlarmSpec:       buildEventAlarmSpec(d.Get("event_alarm_spec").([]interface{})),
		AlarmNotifications:   buildAlarmNotifications(d.Get("alarm_notifications").([]interface{})),
	}
	logp.Printf("[DEBUG] send body is : %+v", createOpts)
	b, err := json.Marshal(createOpts)
	logp.Printf("[DEBUG] send body is : %+v",string(b))
	client.WithMethod(httpclient_go.MethodPost).WithUrlWithoutEndpoint(config, "aom", config.GetRegion(d),
		"v4/"+d.Get("project_id").(string)+"/alarm-rules?action_id=add-alarm-action").WithBody(createOpts)
	response, err := client.Do()

	if err != nil {
		return diag.Errorf("error creating AOM alarm rule %s: %s", createOpts.AlarmRuleName, err)
	}

	mErr := &multierror.Error{}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		mErr = multierror.Append(mErr, err)
	}

	if err = mErr.ErrorOrNil(); err != nil {
		return fmtp.DiagErrorf("error getting AOM prometheus instance fields: %w", err)
	}
	logp.Printf("[DEBUG] create result is : %+v", string(body))
	d.SetId(createOpts.AlarmRuleName)
	return resourceAlarmPolicyRead(context.TODO(), d, meta)
}

func resourceAlarmPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client,dErr := httpclient_go.NewHttpClientGo(config)	
	if dErr != nil {
		return dErr	
	}
	client.WithMethod(httpclient_go.MethodGet).WithUrlWithoutEndpoint(config, "aom", config.GetRegion(d),
		"v4/"+d.Get("project_id").(string)+"/alarm-rules")

	resp, err := client.Do()
	if err != nil {
		logp.Printf("[ERROR] %s", err)
		return common.CheckDeletedDiag(d, err, "error retrieving AOM alarm rule")
	}

	mErr := &multierror.Error{}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logp.Printf("[ERROR] %s", err)
		mErr = multierror.Append(mErr, err)
	}
	logp.Printf("[DEBUG] query result is : %+v", string(body))
	rlt := make([]entity.AddAlarmRuleParams, 0)

	err = json.Unmarshal(body, &rlt)

	if err != nil {
		mErr = multierror.Append(mErr, err)
	}
	d.Set("listAlarmRuleResponseBody", &rlt)
	if err := mErr.ErrorOrNil(); err!= nil {
		return fmtp.DiagErrorf("error getting AOM alarm policy fields: %w", err)
	}
	logp.Printf("[DEBUG] query result is : %+v", rlt)
	return nil
}

func resourceAlarmPolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	client,dErr := httpclient_go.NewHttpClientGo(config)	
	if dErr != nil {
		return dErr	
	}
	client.WithMethod(httpclient_go.MethodDelete).WithUrlWithoutEndpoint(config, "aom", config.GetRegion(d),
		"v4/"+d.Get("project_id").(string)+"/alarm-rules").WithBody(d.Get("alarm_delete_list").([]interface{}))
	resp, err := client.Do()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving AOM alarm rule")
	}

	mErr := &multierror.Error{}
	if resp.StatusCode != 200 {
		mErr = multierror.Append(mErr, fmtp.Errorf("delete alarm policy failed error code: %s", resp.StatusCode))
	} else {
		d.Set("result", "ok")
	}

	if err := mErr.ErrorOrNil(); err != nil {
		return fmtp.DiagErrorf("error delete AOM alarm policy fields: %w", err)
	}

	return nil
}

func buildMetricAlarmSpec(raw interface{}) entity.MetricAlarmSpec {
	mas := make([]entity.MetricAlarmSpec, 0)
	b, err := json.Marshal(raw)
	logp.Printf("[DEBUG] buildMetricAlarmSpec is : %+v", string(b))
	if err != nil {
		logp.Printf("[ERROR] buildMetricAlarmSpec  : %s", err)
		return entity.MetricAlarmSpec{}
	}
	json.Unmarshal(b, &mas)
	if len(mas) == 0 {
		return entity.MetricAlarmSpec{}
	}
	logp.Printf("[DEBUG] buildMetricAlarmSpec is : %+v", mas[0])
	return mas[0]
}

func buildEventAlarmSpec(raw interface{}) entity.EventAlarmSpec {
	mas := make([]entity.EventAlarmSpec, 0)
	b, err := json.Marshal(raw)
	logp.Printf("[DEBUG] buildEventAlarmSpec is : %+v", string(b))
	if err != nil {
		return entity.EventAlarmSpec{}
	}
	json.Unmarshal(b, &mas)
	if len(mas) == 0 {
		return entity.EventAlarmSpec{}
	}
	logp.Printf("[DEBUG] buildEventAlarmSpec is : %+v", mas[0])
	return mas[0]
}

func buildAlarmNotifications(raw interface{}) entity.AlarmNotifications {
	mas := make([]entity.AlarmNotifications, 0)
	b, err := json.Marshal(raw)
	logp.Printf("[DEBUG] buildAlarmNotifications is : %+v", string(b))
	if err != nil {
		return entity.AlarmNotifications{}
	}
	json.Unmarshal(b, &mas)
	if len(mas) == 0 {
		return entity.AlarmNotifications{}
	}
	logp.Printf("[DEBUG] buildAlarmNotifications is : %+v", mas[0])
	return mas[0]
}

