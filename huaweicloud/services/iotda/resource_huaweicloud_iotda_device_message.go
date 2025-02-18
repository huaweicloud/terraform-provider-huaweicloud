package iotda

import (
	"context"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IoTDA POST /v5/iot/{project_id}/devices/{device_id}/messages
// @API IoTDA GET /v5/iot/{project_id}/devices/{device_id}/messages/{message_id}
func ResourceDeviceMessage() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDeviceMessageCreate,
		ReadContext:   resourceDeviceMessageRead,
		DeleteContext: resourceDeviceMessageDelete,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"device_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"message": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"message_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"properties": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"correlation_data": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"response_topic": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
							ForceNew: true,
						},
						"user_properties": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"prop_key": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
									"prop_value": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
										ForceNew: true,
									},
								},
							},
						},
					},
				},
			},
			"encoding": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"payload_format": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"topic": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"topic_full_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"ttl": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"error_info": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"error_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"error_msg": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"created_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"finished_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildCreateDeviceMessagePropertiesParams(properties []interface{}) map[string]interface{} {
	if len(properties) == 0 {
		return nil
	}

	propertiesMap := properties[0].(map[string]interface{})
	propertiesParams := map[string]interface{}{
		"correlation_data": utils.ValueIgnoreEmpty(propertiesMap["correlation_data"].(string)),
		"response_topic":   utils.ValueIgnoreEmpty(propertiesMap["response_topic"].(string)),
	}

	if userPropertiesList := propertiesMap["user_properties"].([]interface{}); len(userPropertiesList) > 0 {
		userPropertiesParams := make([]interface{}, 0, len(userPropertiesList))
		for _, v := range userPropertiesList {
			userPropMap := v.(map[string]interface{})
			userProperties := map[string]interface{}{
				"prop_key":   utils.ValueIgnoreEmpty(userPropMap["prop_key"].(string)),
				"prop_value": utils.ValueIgnoreEmpty(userPropMap["prop_value"].(string)),
			}
			userPropertiesParams = append(userPropertiesParams, userProperties)
		}
		if len(userPropertiesParams) > 0 {
			propertiesParams["user_properties"] = userPropertiesParams
		}
	}

	return propertiesParams
}

func buildDeviceMessageTTLRequestValue(value string) interface{} {
	v, err := strconv.Atoi(value)
	if err != nil {
		log.Printf("[ERROR] error converting string value to int value: %s", err)
		return nil
	}

	return v
}

func buildCreateDeviceMessageBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"message_id":      utils.ValueIgnoreEmpty(d.Get("message_id")),
		"name":            utils.ValueIgnoreEmpty(d.Get("name")),
		"message":         d.Get("message"),
		"properties":      buildCreateDeviceMessagePropertiesParams(d.Get("properties").([]interface{})),
		"encoding":        utils.ValueIgnoreEmpty(d.Get("encoding")),
		"payload_format":  utils.ValueIgnoreEmpty(d.Get("payload_format")),
		"topic":           utils.ValueIgnoreEmpty(d.Get("topic")),
		"topic_full_name": utils.ValueIgnoreEmpty(d.Get("topic_full_name")),
		"ttl":             buildDeviceMessageTTLRequestValue(d.Get("ttl").(string)),
	}

	return bodyParams
}

func resourceDeviceMessageCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		isDerived = WithDerivedAuth(cfg, region)
		product   = "iotda"
	)

	client, err := cfg.NewServiceClientWithDerivedAuth(product, region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	createPath := client.Endpoint + "v5/iot/{project_id}/devices/{device_id}/messages"
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{device_id}", d.Get("device_id").(string))
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateDeviceMessageBodyParams(d)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating IoTDA device message: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	messageId := utils.PathSearch("message_id", createRespBody, "").(string)
	if messageId == "" {
		return diag.Errorf("error creating IoTDA device message: ID is not found in API response")
	}

	d.SetId(messageId)

	return resourceDeviceMessageRead(ctx, d, meta)
}

func resourceDeviceMessageRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		isDerived = WithDerivedAuth(cfg, region)
		product   = "iotda"
	)

	client, err := cfg.NewServiceClientWithDerivedAuth(product, region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	getPath := client.Endpoint + "v5/iot/{project_id}/devices/{device_id}/messages/{message_id}"
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{device_id}", d.Get("device_id").(string))
	getPath = strings.ReplaceAll(getPath, "{message_id}", d.Id())
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	// When the parent resource (device_id) does not exist, query API will return `404` error code.
	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving IoTDA device message")
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	// When the resource does not exist, query API will return `200` status code.
	// Therefore, it is necessary to check whether the message ID is returned in the response body.
	messageIdResp := utils.PathSearch("message_id", getRespBody, "").(string)
	if messageIdResp == "" {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving IoTDA device message")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("message_id", messageIdResp),
		d.Set("name", utils.PathSearch("name", getRespBody, nil)),
		d.Set("message", utils.PathSearch("message", getRespBody, nil)),
		d.Set("properties", flattenDeviceMessageProperties(utils.PathSearch("properties", getRespBody, nil))),
		d.Set("encoding", utils.PathSearch("encoding", getRespBody, nil)),
		d.Set("payload_format", utils.PathSearch("payload_format", getRespBody, nil)),
		d.Set("status", utils.PathSearch("status", getRespBody, nil)),
		d.Set("error_info", flattenDeviceMessageErrorInfo(utils.PathSearch("error_info", getRespBody, nil))),
		d.Set("created_time", utils.PathSearch("created_time", getRespBody, nil)),
		d.Set("finished_time", utils.PathSearch("finished_time", getRespBody, nil)),
	)

	// The values of parameters `topic` and `topic_full_name` correspond to the `topic` field in the query API response.
	topicResp := utils.PathSearch("topic", getRespBody, "").(string)
	if topicResp != "" {
		if strings.HasPrefix(topicResp, "$oc/devices") {
			mErr = multierror.Append(mErr, d.Set("topic", flattenDeviceMessageTopicOrTopicFullName(topicResp)))
		} else {
			mErr = multierror.Append(mErr, d.Set("topic_full_name", flattenDeviceMessageTopicOrTopicFullName(topicResp)))
		}
	}

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDeviceMessageProperties(propertiesResp interface{}) []interface{} {
	if propertiesResp == nil {
		return nil
	}

	propertiesMap := map[string]interface{}{
		"correlation_data": utils.PathSearch("correlation_data", propertiesResp, nil),
		"response_topic":   utils.PathSearch("response_topic", propertiesResp, nil),
		"user_properties":  flattenDeviceMessageUserProperties(utils.PathSearch("user_properties", propertiesResp, nil)),
	}

	return []interface{}{propertiesMap}
}

func flattenDeviceMessageUserProperties(userPropertiesResp interface{}) []interface{} {
	if userPropertiesResp == nil {
		return nil
	}

	userPropertiesRespList := userPropertiesResp.([]interface{})
	result := make([]interface{}, len(userPropertiesRespList))
	for i, v := range userPropertiesRespList {
		result[i] = map[string]interface{}{
			"prop_key":   utils.PathSearch("prop_key", v, nil),
			"prop_value": utils.PathSearch("prop_value", v, nil),
		}
	}

	return result
}

func flattenDeviceMessageErrorInfo(errorInfoResp interface{}) []interface{} {
	if errorInfoResp == nil {
		return nil
	}

	errorInfoMap := map[string]interface{}{
		"error_code": utils.PathSearch("error_code", errorInfoResp, nil),
		"error_msg":  utils.PathSearch("error_msg", errorInfoResp, nil),
	}

	return []interface{}{errorInfoMap}
}

// The topic prefix added to the platform field is: $oc/devices/{device_id}/user/
func flattenDeviceMessageTopicOrTopicFullName(topic string) string {
	index := strings.Index(topic, "user/")
	if index != -1 {
		// For `topic` parameter, We need to extract the suffix here.
		return topic[index+len("user/"):]
	}

	return topic
}

func resourceDeviceMessageDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Delete()' method because the resource is an action resource.
	return nil
}
