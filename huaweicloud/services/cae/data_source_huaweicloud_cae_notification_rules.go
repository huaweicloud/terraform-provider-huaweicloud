package cae

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CAE GET /v1/{project_id}/cae/notice-rules
func DataSourceNotificationRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNotificationRulesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the notification rules are located.`,
			},

			// Internal parameter(s).
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: utils.SchemaDesc(
					`The ID of the enterprise project to which the notification rules belong.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},

			// Attributes.
			"rules": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of the event notification rules.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the event notification rule.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the event notification rule.`,
						},
						"event_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The trigger event of the event notification.`,
						},
						"scope": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The scope in which event notification rule takes effect.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The type to which the event notification rule takes effect.`,
									},
									"environments": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: `The list of the environment IDs.`,
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"applications": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: `The list of the application IDs.`,
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"components": {
										Type:        schema.TypeList,
										Computed:    true,
										Description: `The list of the component IDs.`,
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"trigger_policy": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The trigger policy of the event notification rule.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The type of the trigger.`,
									},
									"period": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `The trigger period of the event.`,
									},
									"count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: `The number of times the event occurred.`,
									},
									"operator": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The condition of the event notification.`,
									},
								},
							},
						},
						"notification": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The configuration of the event notification.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"protocol": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The protocol of the event notification.`,
									},
									"endpoint": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The endpoint of the event notification.`,
									},
									"template": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The template language of the event notification.`,
									},
								},
							},
						},
						"enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether the event notification rule is enabled.`,
						},
					},
				},
			},
		},
	}
}

func listNotificationRules(client *golangsdk.ServiceClient, epsId string) ([]interface{}, error) {
	httpUrl := "v1/{project_id}/cae/notice-rules"
	listPath := client.Endpoint + strings.ReplaceAll(httpUrl, "{project_id}", client.ProjectID)

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildRequestMoreHeaders("", epsId),
	}

	requestResp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return nil, err
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}

	return utils.PathSearch("items", respBody, make([]interface{}, 0)).([]interface{}), nil
}

func flattenNotificationRulesScope(scope map[string]interface{}) []map[string]interface{} {
	if len(scope) < 1 {
		return nil
	}

	return []map[string]interface{}{
		{
			"type":         utils.PathSearch("type", scope, nil),
			"environments": utils.PathSearch("environments", scope, make([]interface{}, 0)),
			"applications": utils.PathSearch("applications", scope, make([]interface{}, 0)),
			"components":   utils.PathSearch("components", scope, make([]interface{}, 0)),
		},
	}
}

func flattenNotificationRulesTriggerPolicy(triggerPolicy map[string]interface{}) []map[string]interface{} {
	if len(triggerPolicy) < 1 {
		return nil
	}

	return []map[string]interface{}{
		{
			"type":     utils.PathSearch("trigger_type", triggerPolicy, nil),
			"period":   utils.PathSearch("period", triggerPolicy, nil),
			"count":    utils.PathSearch("count", triggerPolicy, nil),
			"operator": utils.PathSearch("operator", triggerPolicy, nil),
		},
	}
}

func flattenNotificationRulesNotification(notification map[string]interface{}) []map[string]interface{} {
	if len(notification) < 1 {
		return nil
	}

	return []map[string]interface{}{
		{
			"protocol": utils.PathSearch("protocol", notification, nil),
			"endpoint": utils.PathSearch("endpoint", notification, nil),
			"template": utils.PathSearch("template", notification, nil),
		},
	}
}

func flattenNotificationRules(items []interface{}) []map[string]interface{} {
	if len(items) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, len(items))
	for i, item := range items {
		itemMap := item.(map[string]interface{})
		result[i] = map[string]interface{}{
			"id":         utils.PathSearch("id", itemMap, nil),
			"name":       utils.PathSearch("name", itemMap, nil),
			"event_name": utils.PathSearch("event_name", itemMap, nil),
			"scope": flattenNotificationRulesScope(utils.PathSearch("scope", itemMap,
				make(map[string]interface{})).(map[string]interface{})),
			"trigger_policy": flattenNotificationRulesTriggerPolicy(utils.PathSearch("trigger_policy", itemMap,
				make(map[string]interface{})).(map[string]interface{})),
			"notification": flattenNotificationRulesNotification(utils.PathSearch("notification", itemMap,
				make(map[string]interface{})).(map[string]interface{})),
			"enabled": utils.PathSearch("enable", itemMap, nil),
		}
	}

	return result
}

func dataSourceNotificationRulesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cae", region)
	if err != nil {
		return diag.Errorf("error creating CAE client: %s", err)
	}

	resp, err := listNotificationRules(client, cfg.GetEnterpriseProjectID(d))
	if err != nil {
		return diag.Errorf("error querying CAE notification rules: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("rules", flattenNotificationRules(resp)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
