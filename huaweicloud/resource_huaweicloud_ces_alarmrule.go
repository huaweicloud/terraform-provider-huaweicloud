package huaweicloud

import (
	"fmt"
	"log"
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
			},

			"alarm_description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
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
						},

						"metric_name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
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
									},

									"value": &schema.Schema{
										Type:     schema.TypeString,
										Required: true,
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
						},

						"filter": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},

						"comparison_operator": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},

						"value": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
						},

						"unit": &schema.Schema{
							Type:     schema.TypeString,
							Optional: true,
						},

						"count": &schema.Schema{
							Type:     schema.TypeInt,
							Required: true,
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

func resourceAlarmRuleCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := chooseCESClient(d, config)
	if err != nil {
		return fmt.Errorf("Error creating Cloud Eye Service client: %s", err)
	}

	var createOpts alarmrule.CreateOpts
	_, err = buildCreateParam(&createOpts, d, &map[string]string{"notificationList": "notification_list"})
	if err != nil {
		return fmt.Errorf("Error creating %s: building parameter failed:%s", nameCESAR, err)
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

	return refreshResourceData(r, d, &map[string]string{"notificationList": "notification_list"})
}

func resourceAlarmRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := chooseCESClient(d, config)
	if err != nil {
		return fmt.Errorf("Error creating Cloud Eye Service client: %s", err)
	}

	arId := d.Id()

	var updateOpts alarmrule.UpdateOpts
	_, err = buildUpdateParam(&updateOpts, d, &map[string]string{"notificationList": "notification_list"})
	if err != nil {
		return fmt.Errorf("Error updating %s %s: building parameter failed:%s", nameCESAR, arId, err)
	}
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
