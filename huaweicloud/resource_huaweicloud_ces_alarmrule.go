package huaweicloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
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
			"alarm_name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					vv := regexp.MustCompile("^[a-zA-Z0-9_]{1,128}$")
					if !vv.MatchString(value) {
						errors = append(errors, fmt.Errorf("%s must be string of 1 to 128 characters that consists of uppercase/lowercae letters, digits and underscores(_)", k))
					}
					return
				},
			},

			"alarm_description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if len(value) > 256 {
						errors = append(errors, fmt.Errorf("The length of %s must be in [0, 256]", k))
					}
					return
				},
			},

			"metric": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"namespace": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
								value := v.(string)
								vv := regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9_]{2,31}\\.[a-zA-Z][a-zA-Z0-9_]{2,31}$")
								if !vv.MatchString(value) {
									errors = append(errors, fmt.Errorf("%s is in service.item format. service and item must be a string of 3 to 32 characters that starts with a letter and consists of uppercase/lowercae letters, digits and underscores(_)", k))
								}
								return
							},
						},

						"metric_name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
								value := v.(string)
								vv := regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9_]{0,63}$")
								if !vv.MatchString(value) {
									errors = append(errors, fmt.Errorf("%s must be a string of 1 to 64 characters that starts with a letter and consists of uppercase/lowercae letters, digits and underscores(_)", k))
								}
								return
							},
						},

						"dimensions": &schema.Schema{
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 3,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": &schema.Schema{
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
											value := v.(string)
											vv := regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9_]{0,31}$")
											if !vv.MatchString(value) {
												errors = append(errors, fmt.Errorf("%s must be a string of 1 to 32 characters that starts with a letter and consists of uppercase/lowercae letters, digits and underscores(_)", k))
											}
											return
										},
									},

									"value": &schema.Schema{
										Type:     schema.TypeString,
										Required: true,
										ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
											value := v.(string)
											vv := regexp.MustCompile("^[a-zA-Z0-9][a-zA-Z0-9-]{0,63}$")
											if !vv.MatchString(value) {
												errors = append(errors, fmt.Errorf("%s must be a string of 1 to 64 characters that starts with a letter or digit and consists of uppercase/lowercae letters, digits and hyphens(-)", k))
											}
											return
										},
									},
								},
							},
						},
					},
				},
			},

			"condition": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"period": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
							ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
								value := v.(int)
								switch value {
								case 1:
								case 300:
								case 1200:
								case 3600:
								case 14400:
								case 86400:
								default:
									errors = append(errors, fmt.Errorf("%s can be 1, 300, 1200, 3600, 14400, 86400", k))
								}
								return
							},
						},

						"filter": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
								value := v.(string)
								switch value {
								case "max":
								case "min":
								case "average":
								case "sum":
								case "variance":
								default:
									errors = append(errors, fmt.Errorf("%s can be Max, Min, average, Sum, Variance", k))
								}
								return
							},
						},

						"comparison_operator": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
								value := v.(string)
								switch value {
								case ">":
								case "=":
								case "<":
								case ">=":
								case "<=":
								default:
									errors = append(errors, fmt.Errorf("%s can be >, =, <, >=, <=", k))
								}
								return
							},
						},

						"value": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
							ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
								value := v.(int)
								if value < 0 {
									errors = append(errors, fmt.Errorf("%s must be greater than or equal to 0", k))
								}
								return
							},
						},

						"unit": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},

						"count": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
							ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
								value := v.(int)
								if value < 1 || value > 5 {
									errors = append(errors, fmt.Errorf("%s must be in range [1, 5]", k))
								}
								return
							},
						},
					},
				},
			},

			"alarm_actions": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
								value := v.(string)
								switch value {
								case "notification":
								case "autoscaling":
								default:
									errors = append(errors, fmt.Errorf("%s can be notification or autoscaling", k))
								}
								return
							},
						},

						"notification_list": &schema.Schema{
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 5,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},

			"insufficientdata_actions": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
								value := v.(string)
								switch value {
								case "notification":
								case "autoscaling":
								default:
									errors = append(errors, fmt.Errorf("%s can be notification or autoscaling", k))
								}
								return
							},
						},

						"notification_list": &schema.Schema{
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 5,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},

			"ok_actions": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
								value := v.(string)
								switch value {
								case "notification":
								case "autoscaling":
								default:
									errors = append(errors, fmt.Errorf("%s can be notification or autoscaling", k))
								}
								return
							},
						},

						"notification_list": &schema.Schema{
							Type:     schema.TypeList,
							Required: true,
							MaxItems: 5,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},

			"alarm_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"alarm_action_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"update_time": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"alarm_state": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
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
	client, err := chooseCESClient(d, config)
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
	client, err := chooseCESClient(d, config)
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
	client, err := chooseCESClient(d, config)
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
	client, err := chooseCESClient(d, config)
	if err != nil {
		return fmt.Errorf("Error creating Cloud Eye Service client: %s", err)
	}

	arId := d.Id()
	log.Printf("[DEBUG] Deleting %s %s", nameCESAR, arId)

	timeout := d.Timeout(schema.TimeoutDelete)
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
