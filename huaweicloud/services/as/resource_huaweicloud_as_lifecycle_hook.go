package as

import (
	"regexp"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk/openstack/autoscaling/v1/lifecyclehooks"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
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
				Default:  "ABANDON",
				ValidateFunc: validation.StringInSlice([]string{
					"ABANDON", "CONTINUE",
				}, false),
			},
			"timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      3600,
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
	client, err := config.AutoscalingV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud AutoScaling client: %s", err)
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
		return fmtp.Errorf("Lifecycle hook type (%s) is not in the map (%#v)", hookType, hookTypeMap)
	}
	createOpts.Type = v
	hook, err := lifecyclehooks.Create(client, createOpts, groupId).Extract()
	if err != nil {
		return fmtp.Errorf("Error creating lifecycle hook: %s", err)
	}
	d.SetId(hook.Name)

	return resourceASLifecycleHookRead(d, meta)
}

func resourceASLifecycleHookRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	client, err := config.AutoscalingV1Client(region)
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud AutoScaling client: %s", err)
	}
	groupId := d.Get("scaling_group_id").(string)
	hook, err := lifecyclehooks.Get(client, groupId, d.Id()).Extract()
	if err != nil {
		return common.CheckDeleted(d, err, "Error getting the specifies lifecycle hook of the AutoScaling service")
	}
	d.Set("region", region)
	if err = setASLifecycleHookToState(d, hook); err != nil {
		return fmtp.Errorf("Error setting the lifecycle hook to state: %s", err)
	}
	return nil
}

func resourceASLifecycleHookUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.AutoscalingV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud AutoScaling client: %s", err)
	}
	updateOpts := lifecyclehooks.UpdateOpts{}
	if d.HasChange("type") {
		hookType := d.Get("type").(string)
		v, ok := hookTypeMap[hookType]
		if !ok {
			return fmtp.Errorf("The type (%s) of hook is not in the map (%#v)", hookType, hookTypeMap)
		}
		updateOpts.Type = v
	}
	if d.HasChange("default_result") {
		updateOpts.DefaultResult = d.Get("default_result").(string)
	}
	if d.HasChange("timeout") {
		updateOpts.Timeout = d.Get("timeout").(int)
	}
	if d.HasChange("notification_topic_urn") {
		updateOpts.NotificationTopicURN = d.Get("notification_topic_urn").(string)
	}
	if d.HasChange("notification_message") {
		updateOpts.NotificationMetadata = d.Get("notification_message").(string)
	}
	groupId := d.Get("scaling_group_id").(string)
	_, err = lifecyclehooks.Update(client, updateOpts, groupId, d.Id()).Extract()
	if err != nil {
		return fmtp.Errorf("Error updating the lifecycle hook of the AutoScaling service: %s", err)
	}

	return resourceASLifecycleHookRead(d, meta)
}

func resourceASLifecycleHookDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	client, err := config.AutoscalingV1Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud AutoScaling client: %s", err)
	}
	groupId := d.Get("scaling_group_id").(string)
	err = lifecyclehooks.Delete(client, groupId, d.Id()).ExtractErr()
	if err != nil {
		return fmtp.Errorf("Error deleting the lifecycle hook of the AutoScaling service: %s", err)
	}

	return nil
}

func setASLifecycleHookToState(d *schema.ResourceData, hook *lifecyclehooks.Hook) error {
	mErr := multierror.Append(
		d.Set("name", hook.Name),
		d.Set("default_result", hook.DefaultResult),
		d.Set("timeout", hook.Timeout),
		d.Set("notification_topic_urn", hook.NotificationTopicURN),
		d.Set("notification_message", hook.NotificationMetadata),
		setASLifecycleHookType(d, hook),
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
	return fmtp.Errorf("The type of hook response is not in the map")
}

func resourceASLifecycleHookImportState(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, fmtp.Errorf("Invalid format specified for lifecycle hook, must be <scaling_group_id>/<hook_id>")
	}
	d.SetId(parts[1])
	d.Set("scaling_group_id", parts[0])
	return []*schema.ResourceData{d}, nil
}
