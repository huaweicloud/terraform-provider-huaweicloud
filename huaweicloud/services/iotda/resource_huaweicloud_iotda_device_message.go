package iotda

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/iotda/v5/model"

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

func buildDeviceMessagePropertiesParams(properties []interface{}) *model.PropertiesDto {
	if len(properties) == 0 {
		return nil
	}

	propertiesMap := properties[0].(map[string]interface{})
	propertiesParams := &model.PropertiesDto{
		CorrelationData: utils.StringIgnoreEmpty(propertiesMap["correlation_data"].(string)),
		ResponseTopic:   utils.StringIgnoreEmpty(propertiesMap["response_topic"].(string)),
	}

	if userPropertiesList := propertiesMap["user_properties"].([]interface{}); len(userPropertiesList) > 0 {
		userPropertiesParams := make([]model.UserPropDto, 0, len(userPropertiesList))
		for _, v := range userPropertiesList {
			userPropMap := v.(map[string]interface{})
			userProperties := model.UserPropDto{
				PropKey:   utils.StringIgnoreEmpty(userPropMap["prop_key"].(string)),
				PropValue: utils.StringIgnoreEmpty(userPropMap["prop_value"].(string)),
			}
			userPropertiesParams = append(userPropertiesParams, userProperties)
		}
		if len(userPropertiesParams) > 0 {
			propertiesParams.UserProperties = &userPropertiesParams
		}
	}

	return propertiesParams
}

func buildDeviceMessageBodyParams(d *schema.ResourceData) *model.CreateMessageRequest {
	message := d.Get("message")
	bodyParams := model.CreateMessageRequest{
		DeviceId: d.Get("device_id").(string),
		Body: &model.DeviceMessageRequest{
			MessageId:     utils.StringIgnoreEmpty(d.Get("message_id").(string)),
			Name:          utils.StringIgnoreEmpty(d.Get("name").(string)),
			Message:       &message,
			Properties:    buildDeviceMessagePropertiesParams(d.Get("properties").([]interface{})),
			Encoding:      utils.StringIgnoreEmpty(d.Get("encoding").(string)),
			PayloadFormat: utils.StringIgnoreEmpty(d.Get("payload_format").(string)),
			Topic:         utils.StringIgnoreEmpty(d.Get("topic").(string)),
			TopicFullName: utils.StringIgnoreEmpty(d.Get("topic_full_name").(string)),
			Ttl:           convertStringValueToInt32(d.Get("ttl").(string)),
		},
	}

	return &bodyParams
}

func resourceDeviceMessageCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		isDerived = WithDerivedAuth(cfg, region)
	)
	client, err := cfg.HcIoTdaV5Client(region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA v5 client: %s", err)
	}

	createOpts := buildDeviceMessageBodyParams(d)
	resp, err := client.CreateMessage(createOpts)
	if err != nil {
		return diag.Errorf("error creating device message: %s", err)
	}

	if resp == nil || resp.MessageId == nil {
		return diag.Errorf("error creating device message: ID is not found in API response")
	}

	d.SetId(*resp.MessageId)

	return resourceDeviceMessageRead(ctx, d, meta)
}

func buildDeviceMessageQueryParams(d *schema.ResourceData) *model.ShowDeviceMessageRequest {
	return &model.ShowDeviceMessageRequest{
		DeviceId:  d.Get("device_id").(string),
		MessageId: d.Id(),
	}
}

func resourceDeviceMessageRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		isDerived = WithDerivedAuth(cfg, region)
	)
	client, err := cfg.HcIoTdaV5Client(region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA v5 client: %s", err)
	}

	// When the parent resource (device_id) does not exist, query API will return `404` error code.
	response, err := client.ShowDeviceMessage(buildDeviceMessageQueryParams(d))
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving IoTDA device message")
	}

	// When the resource does not exist, query API will return `200` status code.
	// Therefore, it is necessary to check whether the message ID is returned in the response body.
	if response == nil || response.MessageId == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving IoTDA device message")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("message_id", response.MessageId),
		d.Set("name", response.Name),
		d.Set("message", response.Message),
		d.Set("properties", flattenDeviceMessageProperties(response.Properties)),
		d.Set("encoding", response.Encoding),
		d.Set("payload_format", response.PayloadFormat),
		d.Set("status", response.Status),
		d.Set("error_info", flattenDeviceMessageErrorInfo(response.ErrorInfo)),
		d.Set("created_time", response.CreatedTime),
		d.Set("finished_time", response.FinishedTime),
	)

	// The values of parameters `topic` and `topic_full_name` correspond to the `topic` field in the query API response.
	if response.Topic != nil {
		topicStr := *response.Topic
		if strings.HasPrefix(topicStr, "$oc/devices") {
			mErr = multierror.Append(mErr, d.Set("topic", flattenDeviceMessageTopicOrTopicFullName(topicStr)))
		} else {
			mErr = multierror.Append(mErr, d.Set("topic_full_name", flattenDeviceMessageTopicOrTopicFullName(topicStr)))
		}
	}

	return diag.FromErr(mErr.ErrorOrNil())
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

func flattenDeviceMessageProperties(properties *model.PropertiesDto) []interface{} {
	if properties == nil {
		return nil
	}

	propertiesMap := map[string]interface{}{
		"correlation_data": properties.CorrelationData,
		"response_topic":   properties.ResponseTopic,
		"user_properties":  flattenDeviceMessageUserProperties(properties.UserProperties),
	}

	return []interface{}{propertiesMap}
}

func flattenDeviceMessageUserProperties(userProperties *[]model.UserPropDto) []interface{} {
	if userProperties == nil {
		return nil
	}

	result := make([]interface{}, len(*userProperties))
	for i, v := range *userProperties {
		result[i] = map[string]interface{}{
			"prop_key":   v.PropKey,
			"prop_value": v.PropValue,
		}
	}

	return result
}

func flattenDeviceMessageErrorInfo(errorInfo *model.ErrorInfoDto) []interface{} {
	if errorInfo == nil {
		return nil
	}

	errorInfoMap := map[string]interface{}{
		"error_code": errorInfo.ErrorCode,
		"error_msg":  errorInfo.ErrorMsg,
	}

	return []interface{}{errorInfoMap}
}

func resourceDeviceMessageDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Delete()' method because the resource is an action resource.
	return nil
}
