package cts

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	cts "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cts/v3/model"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CTS GET /v3/{project_id}/notifications/{notification_type}
func DataSourceNotifications() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNotificationsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The region in which to query the CTS key event notification.",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The type of CTS key event notification.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of CTS key event notification.",
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The status of CTS key event notification.",
			},
			"topic_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The URN of the topic which CTS key event notification uses.",
			},
			"notification_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the CTS key event notification.",
			},
			"operation_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The type of operation that will send notifications.",
			},
			"notifications": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the CTS key event notification.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CTS key event notification name.",
						},
						"operation_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of operation.",
						},
						"operations": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"service": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type of cloud service.",
									},
									"resource": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type of resource.",
									},
									"trace_names": {
										Type:        schema.TypeList,
										Computed:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "An array of trace names.",
									},
								},
							},
							Description: "An array of operations that will trigger notifications.",
						},
						"operation_users": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"group": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The IAM user group.",
									},
									"users": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Description: "An array of IAM user names in the group.",
									},
								},
							},
							Description: "An array of users.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The status of CTS key event notification.",
						},
						"topic_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URN of the topic which CTS key event notification uses.",
						},
						"filter": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"condition": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The relation between the rules.",
									},
									"rule": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
										Description: "The list of filter rules.",
									},
								},
							},
							Description: "Advanced filtering conditions for the CTS key event notification.",
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The creation time of the CTS key event notification.",
						},
					},
				},
				Description: "All CTS key event notifications that match the filter parameters.",
			},
		},
	}
}

func dataSourceNotificationsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	ctsClient, err := cfg.HcCtsV3Client(region)
	if err != nil {
		return diag.Errorf("error creating CTS client: %s", err)
	}

	var notificationType cts.ListNotificationsRequestNotificationType
	if d.Get("type").(string) == "smn" {
		notificationType = cts.GetListNotificationsRequestNotificationTypeEnum().SMN
	} else {
		notificationType = cts.GetListNotificationsRequestNotificationTypeEnum().FUN
	}

	listOpts := &cts.ListNotificationsRequest{
		NotificationType: notificationType,
	}
	if rawName, nameExist := d.GetOk("name"); nameExist {
		listOpts.NotificationName = utils.String(rawName.(string))
	}

	response, err := ctsClient.ListNotifications(listOpts)
	if err != nil {
		return diag.Errorf("error retrieving CTS key event notification: %s", err)
	}
	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	result := flattenAllNotifications(d, response)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("notifications", result),
	)
	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error saving CTS key event notification data source fields: %s", mErr)
	}

	return nil
}

func flattenAllNotifications(d *schema.ResourceData, response *cts.ListNotificationsResponse) []map[string]interface{} {
	if response.Notifications == nil || len(*response.Notifications) == 0 {
		return nil
	}

	allNotifications := *response.Notifications

	var status cts.NotificationsResponseBodyStatus
	rawStatus, statusExist := d.GetOk("status")
	if statusExist {
		v := rawStatus.(string)
		if v == "enabled" {
			status = cts.GetNotificationsResponseBodyStatusEnum().ENABLED
		} else {
			status = cts.GetNotificationsResponseBodyStatusEnum().DISABLED
		}
	}

	var operationType cts.NotificationsResponseBodyOperationType
	rawOperationType, operationTypeExist := d.GetOk("operation_type")
	if operationTypeExist {
		v := rawOperationType.(string)
		if v == "complete" {
			operationType = cts.GetNotificationsResponseBodyOperationTypeEnum().COMPLETE
		} else {
			operationType = cts.GetNotificationsResponseBodyOperationTypeEnum().CUSTOMIZED
		}
	}

	rawTopicID, topicIDExist := d.GetOk("topic_id")
	rawNotificationID, notificationIDExist := d.GetOk("notification_id")

	result := make([]map[string]interface{}, 0)
	for _, notification := range allNotifications {
		if statusExist && status != *notification.Status {
			continue
		}
		if topicIDExist && rawTopicID.(string) != *notification.TopicId {
			continue
		}
		if notificationIDExist && rawNotificationID.(string) != *notification.NotificationId {
			continue
		}
		if operationTypeExist && operationType != *notification.OperationType {
			continue
		}

		notificationMap := map[string]interface{}{
			"id":         notification.NotificationId,
			"name":       notification.NotificationName,
			"topic_id":   notification.TopicId,
			"filter":     flattenNotificationFilter(notification.Filter),
			"created_at": utils.FormatTimeStampRFC3339(*notification.CreateTime/1000, false),
		}

		if notification.Operations != nil {
			notificationMap["operations"] = flattenNotificationOperations(*notification.Operations)
		}

		if notification.NotifyUserList != nil {
			notificationMap["operation_users"] = flattenNotificationUsers(*notification.NotifyUserList)
		}

		if notification.OperationType != nil {
			notificationMap["operation_type"] = formatValue(notification.OperationType)
		}

		if notification.Status != nil {
			notificationMap["status"] = formatValue(notification.Status)
		}

		result = append(result, notificationMap)
	}

	return result
}
