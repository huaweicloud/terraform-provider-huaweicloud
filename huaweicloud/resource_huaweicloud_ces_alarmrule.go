package huaweicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/cloudeyeservice/alarmrule"
)

const nameCESAR = "CES-AlarmRule"

func resourceAlarmRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlarmRuleCreate,
		Read:   resourceAlarmRuleRead,
		Update: resourceAlarmRuleUpdate,
		Delete: resourceAlarmRuleDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"alarm_name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"alarm_description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"metric": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"namespace": {
							Type:     schema.TypeString,
							Required: true,
						},

						"metric_name": {
							Type:     schema.TypeString,
							Required: true,
						},

						"dimensions": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 3,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},

									"value": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},

			"condition": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"period": {
							Type:     schema.TypeInt,
							Required: true,
						},

						"filter": {
							Type:     schema.TypeString,
							Required: true,
						},

						"comparison_operator": {
							Type:     schema.TypeString,
							Required: true,
						},

						"value": {
							Type:     schema.TypeInt,
							Required: true,
						},

						"unit": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"count": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},

			"alarm_actions": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},

						"notification_list": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 5,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},

			"insufficientdata_actions": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},

						"notification_list": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 5,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},

			"ok_actions": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},

						"notification_list": {
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 5,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},

			"alarm_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"alarm_action_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"update_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"alarm_state": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func getMetricOpts(d *schema.ResourceData) (alarmrule.MetricOpts, error) {
	mos, ok := d.Get("metric").([]interface{})
	if !ok {
		return alarmrule.MetricOpts{}, fmt.Errorf("Error converting opt of metric:%v", d.Get("metric"))
	}
	mo := mos[0].(map[string]interface{})

	mod := mo["dimensions"].([]interface{})
	dopts := make([]alarmrule.DimensionOpts, len(mod))
	for i, v := range mod {
		v1 := v.(map[string]interface{})
		dopts[i] = alarmrule.DimensionOpts{
			Name:  v1["name"].(string),
			Value: v1["value"].(string),
		}
	}
	return alarmrule.MetricOpts{
		Namespace:  mo["namespace"].(string),
		MetricName: mo["metric_name"].(string),
		Dimensions: dopts,
	}, nil
}

func getAlarmAction(d *schema.ResourceData, name string) []alarmrule.ActionOpts {
	aos := d.Get(name).([]interface{})
	if len(aos) == 0 {
		return nil
	}
	opts := make([]alarmrule.ActionOpts, len(aos))
	for i, v := range aos {
		v1 := v.(map[string]interface{})

		v2 := v1["notification_list"].([]interface{})
		nl := make([]string, len(v2))
		for j, v3 := range v2 {
			nl[j] = v3.(string)
		}

		opts[i] = alarmrule.ActionOpts{
			Type:             v1["type"].(string),
			NotificationList: nl,
		}
	}
	return opts
}

func resourceAlarmRuleCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.newCESClient(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating Cloud Eye Service client: %s", err)
	}

	metric, err := getMetricOpts(d)
	if err != nil {
		return err
	}
	cos := d.Get("condition").([]interface{})
	co := cos[0].(map[string]interface{})
	createOpts := alarmrule.CreateOpts{
		AlarmName:        d.Get("alarm_name").(string),
		AlarmDescription: d.Get("alarm_description").(string),
		Metric:           metric,
		Condition: alarmrule.ConditionOpts{
			Period:             co["period"].(int),
			Filter:             co["filter"].(string),
			ComparisonOperator: co["comparison_operator"].(string),
			Value:              co["value"].(int),
			Unit:               co["unit"].(string),
			Count:              co["count"].(int),
		},
		AlarmActions:            getAlarmAction(d, "alarm_actions"),
		InsufficientdataActions: getAlarmAction(d, "insufficientdata_actions"),
		OkActions:               getAlarmAction(d, "ok_actions"),
		AlarmEnabled:            d.Get("alarm_enabled").(bool),
		AlarmActionEnabled:      d.Get("alarm_action_enabled").(bool),
	}
	log.Printf("[DEBUG] Create %s Options: %#v", nameCESAR, createOpts)

	r, err := alarmrule.Create(client, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating %s: %s", nameCESAR, err)
	}
	log.Printf("[DEBUG] Create %s: %#v", nameCESAR, *r)

	d.SetId(r.AlarmID)

	return resourceAlarmRuleRead(d, meta)
}

func resourceAlarmRuleRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.newCESClient(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating Cloud Eye Service client: %s", err)
	}

	r, err := alarmrule.Get(client, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "alarmrule")
	}
	log.Printf("[DEBUG] Retrieved %s %s: %#v", nameCESAR, d.Id(), r)

	m, err := convertStructToMap(r, map[string]string{"notificationList": "notification_list"})
	if err != nil {
		return err
	}
	d.Set("alarm_name", m["alarm_name"])
	d.Set("alarm_description", m["alarm_description"])
	d.Set("metric", m["metric"])
	d.Set("condition", m["condition"])
	d.Set("alarm_actions", m["alarm_actions"])
	d.Set("insufficientdata_actions", m["insufficientdata_actions"])
	d.Set("ok_actions", m["ok_actions"])
	d.Set("alarm_enabled", m["alarm_enabled"])
	d.Set("alarm_action_enabled", m["alarm_action_enabled"])
	d.Set("update_time", m["update_time"])
	d.Set("alarm_state", m["alarm_state"])
	return nil
}

func resourceAlarmRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.newCESClient(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating Cloud Eye Service client: %s", err)
	}

	arId := d.Id()

	if !d.HasChange("alarm_enabled") {
		log.Printf("[WARN] %s Nothing will be updated", nameCESAR)
		return nil
	}
	updateOpts := alarmrule.UpdateOpts{AlarmEnabled: d.Get("alarm_enabled").(bool)}
	log.Printf("[DEBUG] Updating %s %s with options: %#v", nameCESAR, arId, updateOpts)

	timeout := d.Timeout(schema.TimeoutUpdate)
	//lintignore:R006
	err = resource.Retry(timeout, func() *resource.RetryError {
		err := alarmrule.Update(client, arId, updateOpts).ExtractErr()
		if err != nil {
			return checkForRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("Error updating %s %s: %s", nameCESAR, arId, err)
	}

	return resourceAlarmRuleRead(d, meta)
}

func resourceAlarmRuleDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.newCESClient(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating Cloud Eye Service client: %s", err)
	}

	arId := d.Id()
	log.Printf("[DEBUG] Deleting %s %s", nameCESAR, arId)

	timeout := d.Timeout(schema.TimeoutDelete)
	//lintignore:R006
	err = resource.Retry(timeout, func() *resource.RetryError {
		err := alarmrule.Delete(client, arId).ExtractErr()
		if err != nil {
			return checkForRetryableError(err)
		}
		return nil
	})
	if err != nil {
		if isResourceNotFound(err) {
			log.Printf("[INFO] deleting an unavailable %s: %s", nameCESAR, arId)
			return nil
		}
		return fmt.Errorf("Error deleting %s %s: %s", nameCESAR, arId, err)
	}

	return nil
}
