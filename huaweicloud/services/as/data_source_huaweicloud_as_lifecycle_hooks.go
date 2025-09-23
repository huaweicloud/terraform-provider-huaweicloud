package as

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/autoscaling/v1/lifecyclehooks"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

var convertHookTypeMap = map[string]string{
	"INSTANCE_LAUNCHING":   "ADD",
	"INSTANCE_TERMINATING": "REMOVE",
}

// @API AS GET /autoscaling-api/v1/{project_id}/scaling_lifecycle_hook/{scaling_group_id}/list
func DataSourceLifeCycleHooks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLifeCycleHooksRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"scaling_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"default_result": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"lifecycle_hooks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"default_result": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"timeout": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"notification_topic_urn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"notification_topic_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"notification_message": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceLifeCycleHooksRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		groupID = d.Get("scaling_group_id").(string)
	)

	client, err := cfg.AutoscalingV1Client(region)
	if err != nil {
		return diag.Errorf("error creating AS v1 client: %s", err)
	}

	lifecycleHookList, err := lifecyclehooks.List(client, groupID).Extract()
	if err != nil {
		return diag.Errorf("error retrieving lifecycle hooks in AS group %s: %s", groupID, err)
	}

	lifecycleHooks, err := flattenLifecycleHooks(d, lifecycleHookList)
	if err != nil {
		return diag.FromErr(err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("lifecycle_hooks", lifecycleHooks),
	)
	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error saving lifecycle hooks data source fields: %s", mErr)
	}

	return nil
}

func flattenLifecycleHooks(d *schema.ResourceData, hooks *[]lifecyclehooks.Hook) ([]map[string]interface{}, error) {
	if hooks == nil {
		return nil, nil
	}

	rst := make([]map[string]interface{}, 0, len(*hooks))
	for _, hook := range *hooks {
		hookType, ok := convertHookTypeMap[hook.Type]
		if !ok {
			return nil, fmt.Errorf("lifecycle hook type (%s) is not in the map (%#v)", hook.Type, convertHookTypeMap)
		}

		if val, ok := d.GetOk("name"); ok && val.(string) != hook.Name {
			continue
		}
		if val, ok := d.GetOk("type"); ok && val.(string) != hookType {
			continue
		}
		if val, ok := d.GetOk("default_result"); ok && val.(string) != hook.DefaultResult {
			continue
		}

		lifecycleHookMap := map[string]interface{}{
			"name":                    hook.Name,
			"type":                    hookType,
			"default_result":          hook.DefaultResult,
			"timeout":                 hook.Timeout,
			"notification_topic_urn":  hook.NotificationTopicURN,
			"notification_topic_name": hook.NotificationTopicName,
			"notification_message":    hook.NotificationMetadata,
			"created_at":              hook.CreateTime,
		}
		rst = append(rst, lifecycleHookMap)
	}
	return rst, nil
}
