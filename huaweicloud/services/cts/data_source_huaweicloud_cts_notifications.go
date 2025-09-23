package cts

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
						"agency_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The cloud service agency name.",
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
	ctsClient, err := cfg.NewServiceClient("cts", region)
	if err != nil {
		return diag.Errorf("error creating CTS client: %s", err)
	}

	notificationType := d.Get("type").(string)
	listHttpUrl := "v3/{project_id}/notifications/{notification_type}"
	listPath := ctsClient.Endpoint + listHttpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", ctsClient.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{notification_type}", notificationType)

	if rawName, nameExist := d.GetOk("name"); nameExist {
		listPath += fmt.Sprintf("?notification_name=%s", rawName.(string))
	}

	listOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	response, err := ctsClient.Request("GET", listPath, &listOpts)
	if err != nil {
		return diag.Errorf("error retrieving CTS key event notification: %s", err)
	}

	respBody, err := utils.FlattenResponse(response)
	if err != nil {
		return diag.FromErr(err)
	}

	notifications := utils.PathSearch("notifications", respBody, make([]interface{}, 0)).([]interface{})

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	result := flattenAllNotifications(d, notifications)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("notifications", result),
	)
	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error saving CTS key event notification data source fields: %s", mErr)
	}

	return nil
}

func flattenAllNotifications(d *schema.ResourceData, notifications []interface{}) []map[string]interface{} {
	if len(notifications) == 0 {
		return nil
	}

	expectStatus := "disabled"
	rawStatus, statusExist := d.GetOk("status")
	if statusExist && rawStatus.(string) == "enabled" {
		expectStatus = "enabled"
	}

	expectOperationType := "customized"
	rawOperationType, operationTypeExist := d.GetOk("operation_type")
	if operationTypeExist && rawOperationType.(string) == "complete" {
		expectOperationType = "complete"
	}

	expectTopicID, topicIDExist := d.GetOk("topic_id")
	expectNotificationID, notificationIDExist := d.GetOk("notification_id")

	result := make([]map[string]interface{}, 0)
	for _, notification := range notifications {
		actualStatus := utils.PathSearch("status", notification, "").(string)
		if statusExist && expectStatus != actualStatus {
			continue
		}
		actualTopicId := utils.PathSearch("topic_id", notification, "").(string)
		if topicIDExist && expectTopicID.(string) != actualTopicId {
			continue
		}
		actualNotificationId := utils.PathSearch("notification_id", notification, "").(string)
		if notificationIDExist && expectNotificationID.(string) != actualNotificationId {
			continue
		}
		actualOperationType := utils.PathSearch("operation_type", notification, "").(string)
		if operationTypeExist && expectOperationType != actualOperationType {
			continue
		}

		name := utils.PathSearch("notification_name", notification, "").(string)
		filter := utils.PathSearch("filter", notification, nil)
		createTime := utils.PathSearch("create_time", notification, float64(0)).(float64)

		notificationMap := map[string]interface{}{
			"id":          actualNotificationId,
			"name":        name,
			"topic_id":    actualTopicId,
			"filter":      flattenNotificationFilter(filter),
			"created_at":  utils.FormatTimeStampRFC3339(int64(createTime)/1000, false),
			"agency_name": utils.PathSearch("agency_name", notification, nil),
		}

		operations := utils.PathSearch("operations", notification, make([]interface{}, 0)).([]interface{})
		notificationMap["operations"] = flattenNotificationOperations(operations)
		notifyUserList := utils.PathSearch("notify_user_list", notification, make([]interface{}, 0)).([]interface{})
		notificationMap["operation_users"] = flattenNotificationUsers(notifyUserList)
		notificationMap["operation_type"] = actualOperationType
		notificationMap["status"] = actualStatus
		result = append(result, notificationMap)
	}

	return result
}
