package huaweicloud

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/huaweicloud/golangsdk/openstack/autoscaling/v1/lifecyclehooks"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

var hookTypeMap = map[string]string{
	"ADD":    "INSTANCE_LAUNCHING",
	"REMOVE": "INSTANCE_TERMINATING",
}

func ResourceASLifecycleHook() *schema.Resource {
	return &schema.Resource{
		Create: resourceASLifecycleHookCreate,
		Read:   resourceASLifecycleHookRead,
		Update: resourceASLifecycleHookUpdate,
		Delete: resourceASLifecycleHookDelete,
		Importer: &schema.ResourceImporter{
			State: resourceASLifecycleHookImportState,
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
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ADD", "REMOVE",
				}, false),
			},
			"notification_topic_urn": {
				Type:     schema.TypeString,
				Required: true,
			},
			"scaling_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"default_result": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ABANDON", "CONTINUE",
				}, false),
			},
			"timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(300, 86400),
			},
			"notification_message": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile("^[^()<>&']{1,256}$"),
					"The 'notification_message' of the lifecycle hook has special character"),
			},
			"notification_topic_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceASLifecycleHookCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.AutoscalingV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud AutoScaling client: %s", err)
	}
	groupId := d.Get("scaling_group_id").(string)
	createOpts := lifecyclehooks.CreateOpts{
		Name:                 d.Get("name").(string),
		DefaultResult:        d.Get("default_result").(string),
		Timeout:              d.Get("timeout").(int),
		NotificationTopicURN: d.Get("notification_topic_urn").(string),
		NotificationMetadata: d.Get("notification_message").(string),
	}
	hookType := d.Get("type").(string)
	v, ok := hookTypeMap[hookType]
	if !ok {
		return fmt.Errorf("Lifecycle hook type (%s) is not in the map (%#v)", hookType, hookTypeMap)
	}
	createOpts.Type = v
	hook, err := lifecyclehooks.Create(client, createOpts, groupId).Extract()
	if err != nil {
		return fmt.Errorf("Error craeting lifecycle hook: %s", err)
	}
	d.SetId(hook.Name)

	return resourceASLifecycleHookRead(d, meta)
}

func resourceASLifecycleHookRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	region := GetRegion(d, config)
	client, err := config.AutoscalingV1Client(region)
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud AutoScaling client: %s", err)
	}
	groupId := d.Get("scaling_group_id").(string)
	hook, err := lifecyclehooks.Get(client, groupId, d.Id()).Extract()
	if err != nil {
		return fmt.Errorf("Error getting the specifies lifecycle hook of the Auto Scaling service: %s", err)
	}
	d.Set("region", region)
	if err = setASLifecycleHookToState(d, hook); err != nil {
		return fmt.Errorf("Error setting lifecycle hook to state: %s", err)
	}
	return nil
}

func resourceASLifecycleHookUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.AutoscalingV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud AutoScaling client: %s", err)
	}
	//lintignore:R019
	if d.HasChanges("type", "default_result", "timeout", "notification_topic_urn", "notification_message") {
		updateOpts := lifecyclehooks.UpdateOpts{
			DefaultResult:        d.Get("default_result").(string),
			Timeout:              d.Get("timeout").(int),
			NotificationTopicURN: d.Get("notification_topic_urn").(string),
			NotificationMetadata: d.Get("notification_message").(string),
		}
		hookType := d.Get("type").(string)
		v, ok := hookTypeMap[hookType]
		if !ok {
			return fmt.Errorf("The type (%s) of hook is not in the map (%#v)", hookType, hookTypeMap)
		}
		updateOpts.Type = v
		_, err := lifecyclehooks.Update(client, updateOpts, d.Get("scaling_group_id").(string), d.Id()).Extract()
		if err != nil {
			return fmt.Errorf("Error updating the lifecycle hook of the AutoScaling service: %s", err)
		}
	}

	return resourceASLifecycleHookRead(d, meta)
}

func resourceASLifecycleHookDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.AutoscalingV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud AutoScaling client: %s", err)
	}
	err = lifecyclehooks.Delete(client, d.Get("scaling_group_id").(string), d.Id()).ExtractErr()
	if err != nil {
		return fmt.Errorf("Error deleting the lifecycle hook of the AutoScaling service: %s", err)
	}

	return nil
}

func setASLifecycleHookToState(d *schema.ResourceData, hook *lifecyclehooks.Hook) error {
	mErr := multierror.Append(
		// required && optional parameters
		d.Set("name", hook.Name),
		d.Set("default_result", hook.DefaultResult),
		d.Set("timeout", hook.Timeout),
		d.Set("notification_topic_urn", hook.NotificationTopicURN),
		d.Set("notification_message", hook.NotificationMetadata),
		setASLifecycleHookType(d, hook),
		// computed parameters
		d.Set("notification_topic_name", hook.NotificationTopicName),
		d.Set("create_time", hook.CreateTime),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return err
	}
	return nil
}

func setASLifecycleHookType(d *schema.ResourceData, hook *lifecyclehooks.Hook) error {
	for k, v := range hookTypeMap {
		if v == hook.Type {
			err := d.Set("type", k)
			return err
		}
	}
	return fmt.Errorf("The type of hook response is not in the map")
}

func resourceASLifecycleHookImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		err := fmt.Errorf("Invalid format specified for lifecycle hook, must be <scaling_group_id>/<hook_id>")
		return nil, err
	}

	config := meta.(*config.Config)
	client, err := config.AutoscalingV1Client(GetRegion(d, config))
	if err != nil {
		return nil, fmt.Errorf("Error creating HuaweiCloud AutoScaling client: %s", err)
	}

	groupId := parts[0]
	ID := parts[1]
	hook, err := lifecyclehooks.Get(client, groupId, ID).Extract()
	if err != nil {
		return nil, fmt.Errorf("Error getting the specifies lifecycle hook of the Auto Scaling service: %s", err)
	}
	d.SetId(ID)
	d.Set("scaling_group_id", groupId)
	if err = setASLifecycleHookToState(d, hook); err != nil {
		return nil, fmt.Errorf("Error setting lifecycle hook: %s", err)
	}
	return []*schema.ResourceData{d}, nil
}
