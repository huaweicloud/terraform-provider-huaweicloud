package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// Response Object
type CreateRecordCallbackConfigResponse struct {

	// 配置id，由服务端返回。创建或修改的时候不携带
	Id *string `json:"id,omitempty"`

	// 直播推流域名
	PublishDomain *string `json:"publish_domain,omitempty"`

	// app名称。如果匹配任意需填写为*。录制规则匹配的时候，优先精确app匹配，如果匹配不到，则匹配*
	App *string `json:"app,omitempty"`

	// 录制回调通知url地址
	NotifyCallbackUrl *string `json:"notify_callback_url,omitempty"`

	// 订阅录制通知消息。消息类型。RECORD_NEW_FILE_START开始创建新的录制文件。RECORD_FILE_COMPLETE录制文件生成完成。RECORD_OVER录制结束。RECORD_FAILED表示录制失败。如果不填写,默认订阅RECORD_FILE_COMPLETE
	NotifyEventSubscription *[]CreateRecordCallbackConfigResponseNotifyEventSubscription `json:"notify_event_subscription,omitempty"`

	// 加密类型
	SignType *CreateRecordCallbackConfigResponseSignType `json:"sign_type,omitempty"`

	// 创建时间，格式：yyyy-mm-ddThh:mm:ssZ，UTC时间。 在查询的时候返回
	CreateTime *string `json:"create_time,omitempty"`

	// 修改时间，格式：yyyy-mm-ddThh:mm:ssZ，UTC时间。 在查询的时候返回
	UpdateTime     *string `json:"update_time,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o CreateRecordCallbackConfigResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateRecordCallbackConfigResponse struct{}"
	}

	return strings.Join([]string{"CreateRecordCallbackConfigResponse", string(data)}, " ")
}

type CreateRecordCallbackConfigResponseNotifyEventSubscription struct {
	value string
}

type CreateRecordCallbackConfigResponseNotifyEventSubscriptionEnum struct {
	RECORD_NEW_FILE_START CreateRecordCallbackConfigResponseNotifyEventSubscription
	RECORD_FILE_COMPLETE  CreateRecordCallbackConfigResponseNotifyEventSubscription
	RECORD_OVER           CreateRecordCallbackConfigResponseNotifyEventSubscription
	RECORD_FAILED         CreateRecordCallbackConfigResponseNotifyEventSubscription
}

func GetCreateRecordCallbackConfigResponseNotifyEventSubscriptionEnum() CreateRecordCallbackConfigResponseNotifyEventSubscriptionEnum {
	return CreateRecordCallbackConfigResponseNotifyEventSubscriptionEnum{
		RECORD_NEW_FILE_START: CreateRecordCallbackConfigResponseNotifyEventSubscription{
			value: "RECORD_NEW_FILE_START",
		},
		RECORD_FILE_COMPLETE: CreateRecordCallbackConfigResponseNotifyEventSubscription{
			value: "RECORD_FILE_COMPLETE",
		},
		RECORD_OVER: CreateRecordCallbackConfigResponseNotifyEventSubscription{
			value: "RECORD_OVER",
		},
		RECORD_FAILED: CreateRecordCallbackConfigResponseNotifyEventSubscription{
			value: "RECORD_FAILED",
		},
	}
}

func (c CreateRecordCallbackConfigResponseNotifyEventSubscription) Value() string {
	return c.value
}

func (c CreateRecordCallbackConfigResponseNotifyEventSubscription) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *CreateRecordCallbackConfigResponseNotifyEventSubscription) UnmarshalJSON(b []byte) error {
	myConverter := converter.StringConverterFactory("string")
	if myConverter != nil {
		val, err := myConverter.CovertStringToInterface(strings.Trim(string(b[:]), "\""))
		if err == nil {
			c.value = val.(string)
			return nil
		}
		return err
	} else {
		return errors.New("convert enum data to string error")
	}
}

type CreateRecordCallbackConfigResponseSignType struct {
	value string
}

type CreateRecordCallbackConfigResponseSignTypeEnum struct {
	MD5        CreateRecordCallbackConfigResponseSignType
	HMACSHA256 CreateRecordCallbackConfigResponseSignType
}

func GetCreateRecordCallbackConfigResponseSignTypeEnum() CreateRecordCallbackConfigResponseSignTypeEnum {
	return CreateRecordCallbackConfigResponseSignTypeEnum{
		MD5: CreateRecordCallbackConfigResponseSignType{
			value: "MD5",
		},
		HMACSHA256: CreateRecordCallbackConfigResponseSignType{
			value: "HMACSHA256",
		},
	}
}

func (c CreateRecordCallbackConfigResponseSignType) Value() string {
	return c.value
}

func (c CreateRecordCallbackConfigResponseSignType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *CreateRecordCallbackConfigResponseSignType) UnmarshalJSON(b []byte) error {
	myConverter := converter.StringConverterFactory("string")
	if myConverter != nil {
		val, err := myConverter.CovertStringToInterface(strings.Trim(string(b[:]), "\""))
		if err == nil {
			c.value = val.(string)
			return nil
		}
		return err
	} else {
		return errors.New("convert enum data to string error")
	}
}
