package aom

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API AOM GET /v2/{project_id}/alarm-notified-histories
func DataSourceAlarmNotifiedHistories() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAlarmNotifiedHistoriesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The region where the alarm notification histories are located.",
			},

			// Required parameters
			"event_sn": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The serial number of the alarm event.",
			},

			// Attributes
			"notified_histories": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        alarmNotifiedHistoriesSchema(),
				Description: "The list of notified histories.",
			},
		},
	}
}

func alarmNotifiedHistoriesNotificationSmnChannelInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"smn_notified_content": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The content of the notification.",
			},
			"smn_subscription_status": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The subscription status of the notification.",
			},
			"smn_subscription_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The subscription type of the notification.",
			},
		},
	}
}

func alarmNotifiedHistoriesNotificationSmnChannelSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"sent_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The timestamp when the notification was sent.",
			},
			"smn_notified_history": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        alarmNotifiedHistoriesNotificationSmnChannelInfoSchema(),
				Description: "The list of smn notification that associated the event.",
			},
			"smn_request_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The request ID of the notification detail.",
			},
			"smn_response_body": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The response body of the notification detail.",
			},
			"smn_response_code": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The response code of the notification detail.",
			},
			"smn_topic": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The SMN topic used for notification.",
			},
		},
	}
}

func alarmNotifiedHistoriesNotificationsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"action_rule": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the alarm notification rule.",
			},
			"notifier_channel": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The notification channel type.",
			},
			"smn_channel": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        alarmNotifiedHistoriesNotificationSmnChannelSchema(),
				Description: "The result detail of the notification.",
			},
		},
	}
}

func alarmNotifiedHistoriesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"event_sn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The serial number of the alarm event.",
			},
			"notifications": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        alarmNotifiedHistoriesNotificationsSchema(),
				Description: "The list of notification results that associated the event.",
			},
		},
	}
}

func buildAlarmNotifiedHistoriesQueryParams(d *schema.ResourceData) string {
	res := ""

	res = fmt.Sprintf("%s&event_sn=%v", res, d.Get("event_sn").(string))

	if len(res) > 1 {
		res = "?" + res[1:]
	}
	return res
}

func flattenNotificationSmnInfo(smnInfos []interface{}) []interface{} {
	if len(smnInfos) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(smnInfos))
	for _, smnInfo := range smnInfos {
		result = append(result, map[string]interface{}{
			"smn_notified_content":    utils.PathSearch("smn_notified_content", smnInfo, nil),
			"smn_subscription_status": utils.PathSearch("smn_subscription_status", smnInfo, nil),
			"smn_subscription_type":   utils.PathSearch("smn_subscription_type", smnInfo, nil),
		})
	}
	return result
}

func flattenNotificationSmnChannel(smnChannel interface{}) []interface{} {
	if smnChannel == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"sent_time": utils.PathSearch("sent_time", smnChannel, nil),
			"smn_notified_history": flattenNotificationSmnInfo(
				utils.PathSearch("smn_notified_history", smnChannel, make([]interface{}, 0)).([]interface{})),
			"smn_request_id":    utils.PathSearch("smn_request_id", smnChannel, nil),
			"smn_response_body": utils.PathSearch("smn_response_body", smnChannel, nil),
			"smn_response_code": utils.PathSearch("smn_response_code", smnChannel, nil),
			"smn_topic":         utils.PathSearch("smn_topic", smnChannel, nil),
		},
	}
}

func flattenNotifications(notifications []interface{}) []interface{} {
	if len(notifications) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(notifications))
	for _, notification := range notifications {
		result = append(result, map[string]interface{}{
			"action_rule":      utils.PathSearch("action_rule", notification, nil),
			"notifier_channel": utils.PathSearch("notifier_channel", notification, nil),
			"smn_channel": flattenNotificationSmnChannel(
				utils.PathSearch("smn_channel", notification, nil)),
		})
	}
	return result
}

func flattenNotifiedHistories(notifiedHistories []interface{}) []interface{} {
	if len(notifiedHistories) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(notifiedHistories))
	for _, notifiedHistory := range notifiedHistories {
		result = append(result, map[string]interface{}{
			"event_sn":      utils.PathSearch("event_sn", notifiedHistory, nil),
			"notifications": flattenNotifications(utils.PathSearch("notifications", notifiedHistory, make([]interface{}, 0)).([]interface{})),
		})
	}
	return result
}

func dataSourceAlarmNotifiedHistoriesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("aom", region)
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	listPath := client.Endpoint + "v2/{project_id}/alarm-notified-histories"
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath += buildAlarmNotifiedHistoriesQueryParams(d)
	listOpts := &golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", listPath, listOpts)
	if err != nil {
		return diag.Errorf("error querying alarm notification histories: %s", err)
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	notifiedHistories := utils.PathSearch("notified_histories", respBody, make([]interface{}, 0))
	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("notified_histories", flattenNotifiedHistories(notifiedHistories.([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
