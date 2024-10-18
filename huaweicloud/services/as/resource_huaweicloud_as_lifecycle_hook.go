package as

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/autoscaling/v1/lifecyclehooks"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

var hookTypeMap = map[string]string{
	"ADD":    "INSTANCE_LAUNCHING",
	"REMOVE": "INSTANCE_TERMINATING",
}

// @API AS DELETE /autoscaling-api/v1/{project_id}/scaling_lifecycle_hook/{groupID}/{hookName}
// @API AS GET /autoscaling-api/v1/{project_id}/scaling_lifecycle_hook/{groupID}/{hookName}
// @API AS PUT /autoscaling-api/v1/{project_id}/scaling_lifecycle_hook/{groupID}/{hookName}
// @API AS POST /autoscaling-api/v1/{project_id}/scaling_lifecycle_hook/{groupID}
func ResourceASLifecycleHook() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceASLifecycleHookCreate,
		ReadContext:   resourceASLifecycleHookRead,
		UpdateContext: resourceASLifecycleHookUpdate,
		DeleteContext: resourceASLifecycleHookDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceASLifecycleHookImportState,
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
			},
			"timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  3600,
			},
			"notification_message": {
				Type:     schema.TypeString,
				Optional: true,
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

func resourceASLifecycleHookCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.AutoscalingV1Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating autoscaling client: %s", err)
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
		return diag.Errorf("lifecycle hook type (%s) is not in the map (%#v)", hookType, hookTypeMap)
	}
	createOpts.Type = v

	hook, err := lifecyclehooks.Create(client, createOpts, groupId).Extract()
	if err != nil {
		return diag.Errorf("error creating lifecycle hook: %s", err)
	}

	d.SetId(hook.Name)
	return resourceASLifecycleHookRead(ctx, d, meta)
}

func resourceASLifecycleHookRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	client, err := conf.AutoscalingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating autoscaling client: %s", err)
	}

	groupId := d.Get("scaling_group_id").(string)
	hook, err := lifecyclehooks.Get(client, groupId, d.Id()).Extract()
	if err != nil {
		// When querying the lifecycle hook details, if the group does not exist, the following error will be reported:
		// {"error": {"code": "AS.2007","message": "The AS group does not exist."}}.
		// If the hook name does not exist, the response HTTP status code of the details API is 404.
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error.code", "AS.2007"),
			"error getting the specifies lifecycle hook of the autoscaling service")
	}

	if hook == nil {
		return diag.Errorf("error getting the specifies lifecycle hook of the autoscaling service: The hook response is empty")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", hook.Name),
		d.Set("default_result", hook.DefaultResult),
		d.Set("timeout", hook.Timeout),
		d.Set("notification_topic_urn", hook.NotificationTopicURN),
		d.Set("notification_message", hook.NotificationMetadata),
		setASLifecycleHookType(d, hook),
		d.Set("notification_topic_name", hook.NotificationTopicName),
		d.Set("create_time", hook.CreateTime),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceASLifecycleHookUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.AutoscalingV1Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating autoscaling client: %s", err)
	}

	updateOpts := lifecyclehooks.UpdateOpts{}
	if d.HasChange("type") {
		hookType := d.Get("type").(string)
		v, ok := hookTypeMap[hookType]
		if !ok {
			return diag.Errorf("the type (%s) of hook is not in the map (%#v)", hookType, hookTypeMap)
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
		return diag.Errorf("error updating the lifecycle hook of the autoscaling service: %s", err)
	}

	return resourceASLifecycleHookRead(ctx, d, meta)
}

func resourceASLifecycleHookDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.AutoscalingV1Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating autoscaling client: %s", err)
	}

	groupId := d.Get("scaling_group_id").(string)
	err = lifecyclehooks.Delete(client, groupId, d.Id()).ExtractErr()
	if err != nil {
		// When deleting the lifecycle hook, if the group does not exist, the following error will be reported:
		// {"error": {"code": "AS.2007","message": "The AS group does not exist."}}.
		// If the hook name does not exist, the response HTTP status code of the deletion API is 404.
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error.code", "AS.2007"),
			"error deleting the lifecycle hook of the autoscaling service")
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
	return fmt.Errorf("the type of hook response is not in the map")
}

func resourceASLifecycleHookImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for lifecycle hook, must be <scaling_group_id>/<hook_id>")
	}

	d.SetId(parts[1])
	return []*schema.ResourceData{d}, d.Set("scaling_group_id", parts[0])
}
