package cae

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CAE POST /v1/{project_id}/cae/notice-rules
// @API CAE GET /v1/{project_id}/cae/notice-rules/{rule_id}
// @API CAE GET /v1/{project_id}/cae/notice-rules
// @API CAE PUT /v1/{project_id}/cae/notice-rules
// @API CAE DELETE /v1/{project_id}/cae/notice-rules/{rule_id}
func ResourceNotificationRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNotificationRuleCreate,
		ReadContext:   resourceNotificationRuleRead,
		UpdateContext: resourceNotificationRuleUpdate,
		DeleteContext: resourceNotificationRuleDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceNotificationRuleImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The name of the event notification rule.`,
			},
			"event_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The trigger event of the event notification.`,
			},
			"scope": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The type to which the event notification rule takes effect.`,
						},
						"environments": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The list of the environment IDs.`,
						},
						"applications": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The list of the applications IDs.`,
						},
						"components": {
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The list of the components IDs.`,
						},
					},
				},
				Description: `The scope in which event notification rule takes effect.`,
			},
			"trigger_policy": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The type of the trigger.`,
						},
						"period": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: `The trigger period of the event.`,
						},
						"count": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: `The number of times the event occurred.`,
						},
						"operator": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The condition of the event notification.`,
						},
					},
				},
				Description: `The trigger policy of the event notification rule.`,
			},
			"notification": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: `The protocol of the event notification.`,
						},
						"endpoint": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: `The endpoint of the event notification.`,
						},
						"template": {
							Type:        schema.TypeString,
							Required:    true,
							ForceNew:    true,
							Description: `The template language of the event notification.`,
						},
					},
				},
				Description: `The configuration of the event notification.`,
			},

			// Optional parameter(s).
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Whether to enable the event notification rule.`,
			},

			// Internal parameter(s).
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Description: utils.SchemaDesc(
					`The ID of the enterprise project to which the notification rule belongs.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
		},
	}
}

func buildCreateEventNotificationRuleBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"api_version": "v1",
		"kind":        "NoticeRule",
		"spec": map[string]interface{}{
			"name":           d.Get("name"),
			"event_name":     d.Get("event_name"),
			"scope":          buildNotificationRuleScope(d.Get("scope.0")),
			"trigger_policy": buildNotificationRuleTriggerPolicy(d.Get("trigger_policy.0")),
			"notification":   buildRuleNotification(d.Get("notification.0")),
			"enable":         d.Get("enabled"),
		},
	}
}

func buildNotificationRuleScope(scope interface{}) map[string]interface{} {
	return map[string]interface{}{
		"type": utils.PathSearch("type", scope, nil),
		"environments": utils.ValueIgnoreEmpty(utils.ExpandToStringList(utils.PathSearch("environments", scope,
			make([]interface{}, 0)).([]interface{}))),
		"applications": utils.ValueIgnoreEmpty(utils.ExpandToStringList(utils.PathSearch("applications", scope,
			make([]interface{}, 0)).([]interface{}))),
		"components": utils.ValueIgnoreEmpty(utils.ExpandToStringList(utils.PathSearch("components", scope,
			make([]interface{}, 0)).([]interface{}))),
	}
}

func buildNotificationRuleTriggerPolicy(triggerPolicy interface{}) map[string]interface{} {
	return map[string]interface{}{
		"trigger_type": utils.PathSearch("type", triggerPolicy, nil),
		"period":       utils.ValueIgnoreEmpty(utils.PathSearch("period", triggerPolicy, nil)),
		"count":        utils.ValueIgnoreEmpty(utils.PathSearch("count", triggerPolicy, nil)),
		"operator":     utils.ValueIgnoreEmpty(utils.PathSearch("operator", triggerPolicy, nil)),
	}
}

func buildRuleNotification(notification interface{}) map[string]interface{} {
	return map[string]interface{}{
		"protocol": utils.PathSearch("protocol", notification, nil),
		"endpoint": utils.PathSearch("endpoint", notification, nil),
		"template": utils.PathSearch("template", notification, nil),
	}
}

func resourceNotificationRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		httpUrl = "v1/{project_id}/cae/notice-rules"
	)

	client, err := cfg.NewServiceClient("cae", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CAE client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateEventNotificationRuleBodyParams(d)),
		MoreHeaders:      buildRequestMoreHeaders("", cfg.GetEnterpriseProjectID(d)),
	}
	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating event notification rule: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	notificationRuleId := utils.PathSearch("spec.id", respBody, "").(string)
	if notificationRuleId == "" {
		return diag.Errorf("unable to find the event notification rule ID from the API response")
	}

	d.SetId(notificationRuleId)

	return resourceNotificationRuleRead(ctx, d, meta)
}

func resourceNotificationRuleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("cae", region)
	if err != nil {
		return diag.Errorf("error creating CAE client: %s", err)
	}

	notificationRule, err := GetEventNotificationRuleById(client, d.Id(), cfg.GetEnterpriseProjectID(d))
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error retrieving event notification rule (%s)", d.Get("name").(string)))
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("spec.name", notificationRule, nil)),
		d.Set("event_name", utils.PathSearch("spec.event_name", notificationRule, nil)),
		d.Set("scope", flattenNotificationRuleScope(utils.PathSearch("spec.scope", notificationRule, nil))),
		d.Set("trigger_policy", flattenNotificationRuleTriggerPolicy(utils.PathSearch("spec.trigger_policy", notificationRule, nil))),
		d.Set("notification", flattenRuleNotification(utils.PathSearch("spec.notification", notificationRule, nil))),
		d.Set("enabled", utils.PathSearch("spec.enable", notificationRule, nil)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenNotificationRuleScope(scope interface{}) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"type":         utils.PathSearch("type", scope, nil),
			"environments": utils.PathSearch("environments", scope, nil),
			"applications": utils.PathSearch("applications", scope, nil),
			"components":   utils.PathSearch("components", scope, nil),
		},
	}
}

func flattenNotificationRuleTriggerPolicy(triggerPolicy interface{}) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"type":     utils.PathSearch("trigger_type", triggerPolicy, nil),
			"period":   utils.PathSearch("period", triggerPolicy, nil),
			"count":    utils.PathSearch("count", triggerPolicy, nil),
			"operator": utils.PathSearch("operator", triggerPolicy, nil),
		},
	}
}

func flattenRuleNotification(notification interface{}) []map[string]interface{} {
	return []map[string]interface{}{
		{
			"protocol": utils.PathSearch("protocol", notification, nil),
			"endpoint": utils.PathSearch("endpoint", notification, nil),
			"template": utils.PathSearch("template", notification, nil),
		},
	}
}

func GetEventNotificationRuleById(client *golangsdk.ServiceClient, notificationRuleId, epsId string) (interface{}, error) {
	httpUrl := "v1/{project_id}/cae/notice-rules/{rule_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{rule_id}", notificationRuleId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildRequestMoreHeaders("", epsId),
	}
	resp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func buildUpdateEventNotificationRuleBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"api_version": "v1",
		"kind":        "NoticeRule",
		"spec": map[string]interface{}{
			"event_name":     d.Get("event_name"),
			"scope":          buildNotificationRuleScope(d.Get("scope.0")),
			"trigger_policy": buildNotificationRuleTriggerPolicy(d.Get("trigger_policy.0")),
			"enable":         d.Get("enabled"),
		},
	}
}

func resourceNotificationRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		httpUrl = "v1/{project_id}/cae/notice-rules/{rule_id}"
	)

	client, err := cfg.NewServiceClient("cae", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CAE client: %s", err)
	}

	if d.HasChanges("event_name", "scope", "trigger_policy", "enabled") {
		updatePath := client.Endpoint + httpUrl
		updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
		updatePath = strings.ReplaceAll(updatePath, "{rule_id}", d.Id())
		updateOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         utils.RemoveNil(buildUpdateEventNotificationRuleBodyParams(d)),
			MoreHeaders:      buildRequestMoreHeaders("", cfg.GetEnterpriseProjectID(d)),
		}
		_, err = client.Request("PUT", updatePath, &updateOpt)
		if err != nil {
			return diag.Errorf("error updating event notification rule: %s", err)
		}
	}

	return resourceNotificationRuleRead(ctx, d, meta)
}

func resourceNotificationRuleDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		httpUrl = "v1/{project_id}/cae/notice-rules/{rule_id}"
	)

	client, err := cfg.NewServiceClient("cae", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CAE client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{rule_id}", d.Id())

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildRequestMoreHeaders("", cfg.GetEnterpriseProjectID(d)),
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	// When deleting a non-existent event notification rule, the response status code is 204, in order to avoid
	// the possibility of returning a 404 status code in the future, the CheckDeleted design is retained here.
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error deleting event notification rule (%s)", d.Get("name").(string)))
	}
	return nil
}

func getEventNotificationRules(client *golangsdk.ServiceClient, epsId string) (interface{}, error) {
	httpUrl := "v1/{project_id}/cae/notice-rules"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)

	getListOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildRequestMoreHeaders("", epsId),
	}
	resp, err := client.Request("GET", listPath, &getListOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func getEventNotificationRuleIdByName(client *golangsdk.ServiceClient, notificationRuleName, epsId string) (string, error) {
	notificationRules, err := getEventNotificationRules(client, epsId)
	if err != nil {
		return "", fmt.Errorf("error retrieving event notification rules: %s", err)
	}

	notificationRuleId := utils.PathSearch(fmt.Sprintf("items[?name=='%s']|[0].id", notificationRuleName), notificationRules, "").(string)
	if notificationRuleId == "" {
		return "", fmt.Errorf("unable to find event notification rule ID (%s) from API response : %s", notificationRuleName, err)
	}
	return notificationRuleId, nil
}

// Since the ID cannot be found on the console, so we need to import by the event notification rule name.
func resourceNotificationRuleImportState(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	var (
		cfg        = meta.(*config.Config)
		importedId = d.Id()
	)

	client, err := cfg.NewServiceClient("cae", cfg.GetRegion(d))
	if err != nil {
		return nil, fmt.Errorf("error creating CAE client: %s", err)
	}

	errorMsg := fmt.Errorf("invalid notification rule import (%s), expected format: '<id>', '<name>', '<id>/<enterprise_project_id>' "+
		"or '<name>/<enterprise_project_id>'", importedId)
	if importedId == "" {
		return nil, errorMsg
	}

	parts := strings.Split(importedId, "/")
	switch len(parts) {
	case 1:
		if !utils.IsUUID(importedId) {
			notificationRuleId, err := getEventNotificationRuleIdByName(client, importedId, cfg.GetEnterpriseProjectID(d))
			if err != nil {
				return nil, err
			}
			d.SetId(notificationRuleId)
		}
	case 2:
		if !utils.IsUUID(parts[0]) {
			notificationRuleId, err := getEventNotificationRuleIdByName(client, importedId, cfg.GetEnterpriseProjectID(d))
			if err != nil {
				return nil, err
			}
			d.SetId(notificationRuleId)
		}
		d.Set("enterprise_project_id", parts[1])
	default:
		return nil, errorMsg
	}

	return []*schema.ResourceData{d}, nil
}
