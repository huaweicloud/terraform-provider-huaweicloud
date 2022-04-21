package cts

import (
	"context"
	"log"
	"regexp"
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	cts "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/cts/v3/model"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// ResourceCTSNotification is the impl of huaweicloud_cts_notification
func ResourceCTSNotification() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCTSNotificationCreate,
		ReadContext:   resourceCTSNotificationRead,
		UpdateContext: resourceCTSNotificationUpdate,
		DeleteContext: resourceCTSNotificationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.All(
					validation.StringLenBetween(1, 64),
					validation.StringMatch(regexp.MustCompile("^[\u4e00-\u9fa50-9a-zA-Z_]+$"),
						"only letters, digits, underscores(_), and Chinese characters are allowed"),
				),
			},
			"operation_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"complete", "customized"}, false),
			},
			"smn_topic": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"operations": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service": {
							Type:     schema.TypeString,
							Required: true,
						},
						"resource": {
							Type:     schema.TypeString,
							Required: true,
						},
						"trace_names": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"operation_users": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group": {
							Type:     schema.TypeString,
							Required: true,
						},
						"users": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"notification_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

}

func buildNotificationCreateRequestBody(d *schema.ResourceData) *cts.CreateNotificationRequestBody {
	reqBody := cts.CreateNotificationRequestBody{
		NotificationName: d.Get("name").(string),
		OperationType:    formatCreateNotificationType(d.Get("operation_type").(string)),
		Operations:       buildKeyOperationOpts(d),
		NotifyUserList:   buildNotifyUserOpts(d),
		TopicId:          utils.String(d.Get("smn_topic").(string)),
	}

	log.Printf("[DEBUG] creating CTS key events notification options: %#v", reqBody)
	return &reqBody
}

func buildKeyOperationOpts(d *schema.ResourceData) *[]cts.Operations {
	rawOperations := d.Get("operations").([]interface{})
	if len(rawOperations) == 0 {
		return nil
	}

	operations := make([]cts.Operations, len(rawOperations))
	for i, item := range rawOperations {
		op := item.(map[string]interface{})
		operations[i] = cts.Operations{
			ServiceType:  op["service"].(string),
			ResourceType: op["resource"].(string),
		}

		if names, ok := op["trace_names"].([]interface{}); ok {
			traceNames := make([]string, len(names))
			for j, v := range names {
				traceNames[j] = v.(string)
			}
			operations[i].TraceNames = traceNames
		}
	}

	log.Printf("[DEBUG] CTS key events notification operations: %#v", operations)
	return &operations
}

func buildNotifyUserOpts(d *schema.ResourceData) *[]cts.NotificationUsers {
	rawOperationUsers := d.Get("operation_users").([]interface{})
	if len(rawOperationUsers) == 0 {
		return nil
	}

	operationUsers := make([]cts.NotificationUsers, len(rawOperationUsers))
	for i, item := range rawOperationUsers {
		user := item.(map[string]interface{})
		operationUsers[i] = cts.NotificationUsers{
			UserGroup: user["group"].(string),
		}

		if names, ok := user["users"].([]interface{}); ok {
			userList := make([]string, len(names))
			for j, v := range names {
				userList[j] = v.(string)
			}
			operationUsers[i].UserList = userList
		}
	}

	log.Printf("[DEBUG] CTS key events notification users: %#v", operationUsers)
	return &operationUsers
}

func resourceCTSNotificationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	ctsClient, err := cfg.HcCtsV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CTS client: %s", err)
	}

	createOpts := cts.CreateNotificationRequest{
		Body: buildNotificationCreateRequestBody(d),
	}

	notification, err := ctsClient.CreateNotification(&createOpts)
	if err != nil {
		return diag.Errorf("error creating CTS key events notification: %s", err)
	}

	notificationName := d.Get("name").(string)
	d.SetId(notificationName)

	// update status if necessary
	status := formatValue(notification.Status)
	expectStatus := convertNotificationStatus(d.Get("enabled").(bool))
	if status != expectStatus {
		return resourceCTSNotificationUpdate(ctx, d, meta)
	}

	return resourceCTSNotificationRead(ctx, d, meta)
}

func resourceCTSNotificationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	ctsClient, err := cfg.HcCtsV3Client(region)
	if err != nil {
		return diag.Errorf("error creating CTS client: %s", err)
	}

	notificationName := d.Id()
	notificationType := cts.GetListNotificationsRequestNotificationTypeEnum().SMN
	listOpts := &cts.ListNotificationsRequest{
		NotificationType: notificationType,
		NotificationName: &notificationName,
	}

	response, err := ctsClient.ListNotifications(listOpts)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CTS key events notification")
	}

	if response.Notifications == nil || len(*response.Notifications) == 0 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving CTS key events notification")
	}

	allNotifications := *response.Notifications
	ctsNotification := allNotifications[0]
	log.Printf("[DEBUG] retrieve CTS key events notification: %#v", ctsNotification)

	if ctsNotification.NotificationName == nil {
		return diag.Errorf("error retrieving CTS key events notification: id not found in api response")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", ctsNotification.NotificationName),
		d.Set("notification_id", ctsNotification.NotificationId),
		d.Set("smn_topic", ctsNotification.TopicId),
	)

	if ctsNotification.Operations != nil {
		keyOperations := flattenNotificationOperations(*ctsNotification.Operations)
		mErr = multierror.Append(mErr, d.Set("operations", keyOperations))
	}

	if ctsNotification.NotifyUserList != nil {
		operationUsers := flattenNotificationUsers(*ctsNotification.NotifyUserList)
		mErr = multierror.Append(mErr, d.Set("operation_users", operationUsers))
	}

	if ctsNotification.OperationType != nil {
		mErr = multierror.Append(mErr, d.Set("operation_type", formatValue(ctsNotification.OperationType)))
	}
	if ctsNotification.Status != nil {
		status := formatValue(ctsNotification.Status)
		mErr = multierror.Append(mErr,
			d.Set("status", status),
			d.Set("enabled", status == "enabled"),
		)
	}

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error setting resource: %s", mErr)
	}
	return nil
}

func resourceCTSNotificationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	ctsClient, err := cfg.HcCtsV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CTS client: %s", err)
	}

	status := convertNotificationStatus(d.Get("enabled").(bool))
	enabledStatus := cts.UpdateNotificationRequestBodyStatus{}
	if err := enabledStatus.UnmarshalJSON([]byte(status)); err != nil {
		return diag.Errorf("failed to parse status %s: %s", status, err)
	}

	notificationName := d.Get("name").(string)
	updateReq := cts.UpdateNotificationRequestBody{
		NotificationName: notificationName,
		NotificationId:   d.Get("notification_id").(string),
		OperationType:    formatUpdateNotificationType(d.Get("operation_type").(string)),
		Operations:       buildKeyOperationOpts(d),
		NotifyUserList:   buildNotifyUserOpts(d),
		Status:           enabledStatus,
		TopicId:          utils.String(d.Get("smn_topic").(string)),
	}

	log.Printf("[DEBUG] updating CTS key events notification options: %#v", updateReq)
	updateOpts := cts.UpdateNotificationRequest{
		Body: &updateReq,
	}

	_, err = ctsClient.UpdateNotification(&updateOpts)
	if err != nil {
		return diag.Errorf("error updating CTS key events notification: %s", err)
	}

	if d.HasChange("name") {
		d.SetId(notificationName)
	}

	return resourceCTSNotificationRead(ctx, d, meta)
}

func resourceCTSNotificationDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	ctsClient, err := cfg.HcCtsV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CTS client: %s", err)
	}

	deleteOpts := cts.DeleteNotificationRequest{
		NotificationId: d.Get("notification_id").(string),
	}

	_, err = ctsClient.DeleteNotification(&deleteOpts)
	if err != nil {
		return diag.Errorf("error deleting CTS key events notification: %s", err)
	}

	return nil
}

func formatCreateNotificationType(operationType string) cts.CreateNotificationRequestBodyOperationType {
	allTypes := cts.GetCreateNotificationRequestBodyOperationTypeEnum()
	if operationType == "complete" {
		return allTypes.COMPLETE
	}
	return allTypes.CUSTOMIZED

}

func formatUpdateNotificationType(operationType string) cts.UpdateNotificationRequestBodyOperationType {
	allTypes := cts.GetUpdateNotificationRequestBodyOperationTypeEnum()
	if operationType == "complete" {
		return allTypes.COMPLETE
	}
	return allTypes.CUSTOMIZED
}

func convertNotificationStatus(enabled bool) string {
	var status = "enabled"
	if !enabled {
		status = "disabled"
	}

	return status
}

func flattenNotificationUsers(users []cts.NotificationUsers) []map[string]interface{} {
	ret := make([]map[string]interface{}, len(users))
	for i, v := range users {
		ret[i] = map[string]interface{}{
			"group": v.UserGroup,
			"users": v.UserList,
		}
	}

	return ret
}

func flattenNotificationOperations(ops []cts.Operations) []map[string]interface{} {
	ret := make([]map[string]interface{}, len(ops))
	for i, v := range ops {
		ret[i] = map[string]interface{}{
			"service":     v.ServiceType,
			"resource":    v.ResourceType,
			"trace_names": v.TraceNames,
		}
	}

	return ret
}
