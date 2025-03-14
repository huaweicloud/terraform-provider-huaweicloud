package cts

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// ResourceCTSNotification is the impl of huaweicloud_cts_notification
// @API CTS GET /v3/{project_id}/notifications/{notification_type}
// @API CTS PUT /v3/{project_id}/notifications
// @API CTS DELETE /v3/{project_id}/notifications
// @API CTS POST /v3/{project_id}/notifications
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
			"agency_name": {
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
			"filter": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"condition": {
							Type:     schema.TypeString,
							Required: true,
						},
						"rule": {
							Type:     schema.TypeList,
							Required: true,
							MinItems: 1,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"notification_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildNotificationCreateRequestBody(d *schema.ResourceData) (map[string]interface{}, error) {
	reqBody := map[string]interface{}{
		"notification_name": d.Get("name").(string),
		"operation_type":    d.Get("operation_type").(string),
		"operations":        utils.ValueIgnoreEmpty(buildKeyOperationOpts(d.Get("operations").([]interface{}))),
		"notify_user_list":  utils.ValueIgnoreEmpty(buildNotifyUserOpts(d.Get("operation_users").([]interface{}))),
		"topic_id":          d.Get("smn_topic").(string),
	}
	filter, err := buildKeyFilterOpts(d.Get("filter").([]interface{}))
	if err != nil {
		return nil, err
	}
	reqBody["filter"] = filter

	if agencyName, ok := d.GetOk("agency_name"); ok {
		reqBody["agency_name"] = agencyName
	}
	log.Printf("[DEBUG] creating CTS key events notification options: %#v", reqBody)
	return reqBody, nil
}

func buildKeyFilterOpts(params []interface{}) (map[string]interface{}, error) {
	if len(params) == 0 {
		return nil, nil
	}

	filterData := params[0].(map[string]interface{})
	filter := map[string]interface{}{
		"is_support_filter": true,
		"rule":              utils.ExpandToStringList(filterData["rule"].([]interface{})),
	}

	switch conditionStr := filterData["condition"].(string); conditionStr {
	case "AND", "OR":
		filter["condition"] = conditionStr
	default:
		return filter, fmt.Errorf("invalid condition, want 'AND' or 'OR', but got '%v'", conditionStr)
	}

	log.Printf("[DEBUG] CTS key events notification filter: %#v", filter)
	return filter, nil
}

func buildKeyOperationOpts(params []interface{}) []interface{} {
	if len(params) == 0 {
		return nil
	}

	operations := make([]interface{}, len(params))
	for i, item := range params {
		op := item.(map[string]interface{})
		operation := map[string]interface{}{
			"service_type":  op["service"].(string),
			"resource_type": op["resource"].(string),
		}

		if names, ok := op["trace_names"].([]interface{}); ok {
			traceNames := make([]string, len(names))
			for j, v := range names {
				traceNames[j] = v.(string)
			}
			operation["trace_names"] = traceNames
		}
		operations[i] = operation
	}

	log.Printf("[DEBUG] CTS key events notification operations: %#v", operations)
	return operations
}

func buildNotifyUserOpts(params []interface{}) []interface{} {
	if len(params) == 0 {
		return nil
	}

	operationUsers := make([]interface{}, len(params))
	for i, item := range params {
		user := item.(map[string]interface{})
		operationUser := map[string]interface{}{
			"user_group": user["group"].(string),
		}

		if names, ok := user["users"].([]interface{}); ok {
			userList := make([]string, len(names))
			for j, v := range names {
				userList[j] = v.(string)
			}
			operationUser["user_list"] = userList
		}

		operationUsers[i] = operationUser
	}

	log.Printf("[DEBUG] CTS key events notification users: %#v", operationUsers)
	return operationUsers
}

func resourceCTSNotificationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	ctsClient, err := cfg.NewServiceClient("cts", region)
	if err != nil {
		return diag.Errorf("error creating CTS client: %s", err)
	}

	createRequestBody, err := buildNotificationCreateRequestBody(d)
	if err != nil {
		return diag.FromErr(err)
	}
	createOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=UTF-8",
		},
		JSONBody: utils.RemoveNil(createRequestBody),
	}

	notificationHttpUrl := "v3/{project_id}/notifications"
	notificationPath := ctsClient.Endpoint + notificationHttpUrl
	notificationPath = strings.ReplaceAll(notificationPath, "{project_id}", ctsClient.ProjectID)
	resp, err := ctsClient.Request("POST", notificationPath, &createOpts)
	if err != nil {
		return diag.Errorf("error creating CTS key events notification: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	notificationName := d.Get("name").(string)
	d.SetId(notificationName)

	// update status if necessary
	status := utils.PathSearch("status", respBody, "").(string)
	expectStatus := convertNotificationStatus(d.Get("enabled").(bool))
	if status != expectStatus {
		return resourceCTSNotificationUpdate(ctx, d, meta)
	}

	return resourceCTSNotificationRead(ctx, d, meta)
}

func resourceCTSNotificationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	ctsClient, err := cfg.NewServiceClient("cts", region)
	if err != nil {
		return diag.Errorf("error creating CTS client: %s", err)
	}

	notificationName := d.Id()
	notificationType := "smn"
	listOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getNotificationUrl := "v3/{project_id}/notifications/{notification_type}?notification_name={notification_name}"
	getNotificationPath := ctsClient.Endpoint + getNotificationUrl
	getNotificationPath = strings.ReplaceAll(getNotificationPath, "{project_id}", ctsClient.ProjectID)
	getNotificationPath = strings.ReplaceAll(getNotificationPath, "{notification_type}", notificationType)
	getNotificationPath = strings.ReplaceAll(getNotificationPath, "{notification_name}", notificationName)
	resp, err := ctsClient.Request("GET", getNotificationPath, &listOpts)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CTS key events notification")
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	notification := utils.PathSearch("notifications|[0]", respBody, nil)

	if notification == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving CTS key events notification")
	}

	log.Printf("[DEBUG] retrieve CTS key events notification: %#v", notification)

	name := utils.PathSearch("notification_name", notification, "").(string)
	if name == "" {
		return diag.Errorf("error retrieving CTS key events notification: id not found in api response")
	}

	operations := utils.PathSearch("operations", notification, make([]interface{}, 0)).([]interface{})
	notifyUserList := utils.PathSearch("notify_user_list", notification, make([]interface{}, 0)).([]interface{})
	status := utils.PathSearch("status", notification, "").(string)
	createTime := utils.PathSearch("create_time", notification, float64(0)).(float64)
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", name),
		d.Set("notification_id", utils.PathSearch("notification_id", notification, nil)),
		d.Set("smn_topic", utils.PathSearch("topic_id", notification, nil)),
		d.Set("filter", flattenNotificationFilter(utils.PathSearch("filter", notification, nil))),
		d.Set("agency_name", utils.PathSearch("agency_name", notification, nil)),
		d.Set("operations", flattenNotificationOperations(operations)),
		d.Set("operation_users", flattenNotificationUsers(notifyUserList)),
		d.Set("operation_type", utils.PathSearch("operation_type", notification, nil)),
		d.Set("status", status),
		d.Set("enabled", status == "enabled"),
		d.Set("created_at", utils.FormatTimeStampRFC3339(int64(createTime)/1000, false)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceCTSNotificationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	ctsClient, err := cfg.NewServiceClient("cts", region)
	if err != nil {
		return diag.Errorf("error creating CTS client: %s", err)
	}

	status := convertNotificationStatus(d.Get("enabled").(bool))

	notificationName := d.Get("name").(string)
	updateReq := map[string]interface{}{
		"notification_name": notificationName,
		"notification_id":   d.Get("notification_id").(string),
		"operation_type":    d.Get("operation_type").(string),
		"operations":        utils.ValueIgnoreEmpty(buildKeyOperationOpts(d.Get("operations").([]interface{}))),
		"notify_user_list":  utils.ValueIgnoreEmpty(buildNotifyUserOpts(d.Get("operation_users").([]interface{}))),
		"status":            status,
		"topic_id":          d.Get("smn_topic").(string),
	}

	filter, err := buildKeyFilterOpts(d.Get("filter").([]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}
	updateReq["filter"] = filter

	if rawAgencyName, ok := d.GetOk("agency_name"); ok {
		updateReq["agency_name"] = rawAgencyName
	}

	log.Printf("[DEBUG] updating CTS key events notification options: %#v", updateReq)
	updateOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=UTF-8",
		},
		JSONBody: utils.RemoveNil(updateReq),
	}

	notificationHttpUrl := "v3/{project_id}/notifications"
	notificationPath := ctsClient.Endpoint + notificationHttpUrl
	notificationPath = strings.ReplaceAll(notificationPath, "{project_id}", ctsClient.ProjectID)
	_, err = ctsClient.Request("PUT", notificationPath, &updateOpts)
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
	region := cfg.GetRegion(d)
	ctsClient, err := cfg.NewServiceClient("cts", region)
	if err != nil {
		return diag.Errorf("error creating CTS client: %s", err)
	}

	deleteHttpUrl := "v3/{project_id}/notifications?notification_id={notification_id}"
	deletePath := ctsClient.Endpoint + deleteHttpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", ctsClient.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{notification_id}", d.Get("notification_id").(string))

	deleteOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = ctsClient.Request("DELETE", deletePath, &deleteOpts)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting CTS key events notification")
	}

	return nil
}

func convertNotificationStatus(enabled bool) string {
	if enabled {
		return "enabled"
	}
	return "disabled"
}

func flattenNotificationUsers(users []interface{}) []map[string]interface{} {
	ret := make([]map[string]interface{}, len(users))
	for i, v := range users {
		ret[i] = map[string]interface{}{
			"group": utils.PathSearch("user_group", v, nil),
			"users": utils.PathSearch("user_list", v, nil),
		}
	}

	return ret
}

func flattenNotificationOperations(ops []interface{}) []map[string]interface{} {
	ret := make([]map[string]interface{}, len(ops))
	for i, v := range ops {
		ret[i] = map[string]interface{}{
			"service":     utils.PathSearch("service_type", v, nil),
			"resource":    utils.PathSearch("resource_type", v, nil),
			"trace_names": utils.PathSearch("trace_names", v, nil),
		}
	}

	return ret
}

func flattenNotificationFilter(filter interface{}) []map[string]interface{} {
	if filter == nil {
		return nil
	}
	result := map[string]interface{}{
		"condition": utils.PathSearch("condition", filter, nil),
		"rule":      utils.PathSearch("rule", filter, nil),
	}

	return []map[string]interface{}{result}
}
