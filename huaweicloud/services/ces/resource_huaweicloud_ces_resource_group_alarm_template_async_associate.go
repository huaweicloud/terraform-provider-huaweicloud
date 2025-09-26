package ces

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CES PUT /v2/{project_id}/resource-groups/{group_id}/alarm-templates/async-association
var resourceGroupAlarmTemplateAsyncAssociateNonUpdatableParams = []string{"group_id"}

func ResourceResourceGroupAlarmTemplateAsyncAssociate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceResourceGroupAlarmTemplateAsyncAssociateCreate,
		UpdateContext: resourceResourceGroupAlarmTemplateAsyncAssociateUpdate,
		ReadContext:   resourceResourceGroupAlarmTemplateAsyncAssociateRead,
		DeleteContext: resourceResourceGroupAlarmTemplateAsyncAssociateDelete,

		CustomizeDiff: config.FlexibleForceNew(resourceGroupAlarmTemplateAsyncAssociateNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"template_ids": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"notification_enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"alarm_notifications": schemeResourceGroupAlarmTemplateAsyncAssociateNotifications(),
			"ok_notifications":    schemeResourceGroupAlarmTemplateAsyncAssociateNotifications(),
			"notification_begin_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"notification_end_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"effective_timezone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"notification_manner": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"notification_policy_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func schemeResourceGroupAlarmTemplateAsyncAssociateNotifications() *schema.Schema {
	return &schema.Schema{
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
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
			},
		},
	}
}

func resourceResourceGroupAlarmTemplateAsyncAssociateCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("ces", region)
	if err != nil {
		return diag.Errorf("error creating CES client: %s", err)
	}

	groupID := d.Get("group_id").(string)
	err = createOrUpdateOrDeleteResourceGroupAlarmTemplateAsyncAssociate(
		client, buildCreateResourceGroupAlarmTemplateAsyncAssociateBodyParams(d, cfg), groupID)
	if err != nil {
		return diag.Errorf("error creating CES resource group alarm template async association: %s", err)
	}

	d.SetId(groupID)

	return nil
}

func createOrUpdateOrDeleteResourceGroupAlarmTemplateAsyncAssociate(client *golangsdk.ServiceClient,
	bodyParams map[string]interface{}, groupID string) error {
	createHttpUrl := "v2/{project_id}/resource-groups/{group_id}/alarm-templates/async-association"

	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{group_id}", groupID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(bodyParams),
	}
	_, err := client.Request("PUT", createPath, &createOpt)
	return err
}

func buildCreateResourceGroupAlarmTemplateAsyncAssociateBodyParams(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"template_ids":         d.Get("template_ids").(*schema.Set).List(),
		"notification_enabled": d.Get("notification_enabled"),
		"alarm_notifications": buildCreateResourceGroupAlarmTemplateAsyncAssociateNotificationsBodyParams(
			d.Get("alarm_notifications")),
		"ok_notifications": buildCreateResourceGroupAlarmTemplateAsyncAssociateNotificationsBodyParams(
			d.Get("ok_notifications")),
		"notification_begin_time": utils.ValueIgnoreEmpty(d.Get("notification_begin_time")),
		"notification_end_time":   utils.ValueIgnoreEmpty(d.Get("notification_end_time")),
		"effective_timezone":      utils.ValueIgnoreEmpty(d.Get("effective_timezone")),
		"enterprise_project_id":   cfg.GetEnterpriseProjectID(d),
		"notification_manner":     utils.ValueIgnoreEmpty(d.Get("notification_manner")),
		"notification_policy_ids": d.Get("notification_policy_ids").(*schema.Set).List(),
	}
	return bodyParams
}

func buildCreateResourceGroupAlarmTemplateAsyncAssociateNotificationsBodyParams(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		rst := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			rst[i] = map[string]interface{}{
				"type":              raw["type"],
				"notification_list": raw["notification_list"],
			}
		}
		return rst
	}
	return nil
}

func resourceResourceGroupAlarmTemplateAsyncAssociateRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceResourceGroupAlarmTemplateAsyncAssociateUpdate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("ces", region)
	if err != nil {
		return diag.Errorf("error creating CES client: %s", err)
	}

	err = createOrUpdateOrDeleteResourceGroupAlarmTemplateAsyncAssociate(
		client, buildCreateResourceGroupAlarmTemplateAsyncAssociateBodyParams(d, cfg), d.Id())
	if err != nil {
		return diag.Errorf("error updating CES resource group alarm template async association: %s", err)
	}

	return nil
}

func resourceResourceGroupAlarmTemplateAsyncAssociateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("ces", region)
	if err != nil {
		return diag.Errorf("error creating CES client: %s", err)
	}

	err = createOrUpdateOrDeleteResourceGroupAlarmTemplateAsyncAssociate(
		client, buildDeleteResourceGroupAlarmTemplateAsyncAssociateBodyParams(d, cfg), d.Id())
	if err != nil {
		return diag.Errorf("error deleting CES resource group alarm template async association: %s", err)
	}

	return nil
}

func buildDeleteResourceGroupAlarmTemplateAsyncAssociateBodyParams(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"template_ids":          make([]interface{}, 0),
		"notification_enabled":  false,
		"enterprise_project_id": cfg.GetEnterpriseProjectID(d),
	}
	return bodyParams
}
