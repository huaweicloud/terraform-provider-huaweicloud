package as

import (
	"context"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/autoscaling/v1/lifecyclehooks"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API AS PUT /autoscaling-api/v1/{project_id}/scaling_instance_hook/{scaling_group_id}/callback
func ResourceLifecycleHookCallBack() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLifecycleHookCallBackCreate,
		ReadContext:   resourceLifecycleHookCallBackRead,
		DeleteContext: resourceLifecycleHookCallBackDelete,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"scaling_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"lifecycle_action_result": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"lifecycle_action_key": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"lifecycle_action_key"},
			},
			"lifecycle_hook_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"lifecycle_action_key"},
			},
		},
	}
}

func resourceLifecycleHookCallBackCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		conf   = meta.(*config.Config)
		region = conf.GetRegion(d)
	)
	client, err := conf.AutoscalingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating autoscaling client: %s", err)
	}

	callBackOpts := lifecyclehooks.CallBackOpts{
		LifecycleActionResult: d.Get("lifecycle_action_result").(string),
		LifecycleActionKey:    d.Get("lifecycle_action_key").(string),
		InstanceId:            d.Get("instance_id").(string),
		LifecycleHookName:     d.Get("lifecycle_hook_name").(string),
	}

	err = lifecyclehooks.CallBack(client, callBackOpts, d.Get("scaling_group_id").(string)).ExtractErr()
	if err != nil {
		return diag.Errorf("callback the AS lifecycle hook failed: %s", err)
	}

	uuId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(uuId)

	return resourceLifecycleHookCallBackRead(ctx, d, meta)
}

func resourceLifecycleHookCallBackRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceLifecycleHookCallBackDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}
