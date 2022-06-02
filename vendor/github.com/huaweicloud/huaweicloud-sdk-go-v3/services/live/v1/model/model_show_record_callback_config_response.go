package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// Response Object
type ShowRecordCallbackConfigResponse struct {

	// 配置id，由服务端返回。创建或修改的时候不携带
	Id *string `json:"id,omitempty"`

	// 直播推流域名
	PublishDomain *string `json:"publish_domain,omitempty"`

	// app名称。如果匹配任意需填写为*。录制规则匹配的时候，优先精确app匹配，如果匹配不到，则匹配*
	App *string `json:"app,omitempty"`

	// 录制回调通知url地址
	NotifyCallbackUrl *string `json:"notify_callback_url,omitempty"`

	// 订阅录制通知消息。消息类型。RECORD_NEW_FILE_START开始创建新的录制文件。RECORD_FILE_COMPLETE录制文件生成完成。RECORD_OVER录制结束。RECORD_FAILED表示录制失败。如果不填写,默认订阅RECORD_FILE_COMPLETE
	NotifyEventSubscription *[]ShowRecordCallbackConfigResponseNotifyEventSubscription `json:"notify_event_subscription,omitempty"`

	// 加密类型
	SignType *ShowRecordCallbackConfigResponseSignType `json:"sign_type,omitempty"`

	// 创建时间，格式：yyyy-mm-ddThh:mm:ssZ，UTC时间。 在查询的时候返回
	CreateTime *string `json:"create_time,omitempty"`

	// 修改时间，格式：yyyy-mm-ddThh:mm:ssZ，UTC时间。 在查询的时候返回
	UpdateTime     *string `json:"update_time,omitempty"`
	HttpStatusCode int     `json:"-"`
}

func (o ShowRecordCallbackConfigResponse) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ShowRecordCallbackConfigResponse struct{}"
	}

	return strings.Join([]string{"ShowRecordCallbackConfigResponse", string(data)}, " ")
}

type ShowRecordCallbackConfigResponseNotifyEventSubscription struct {
	value string
}

type ShowRecordCallbackConfigResponseNotifyEventSubscriptionEnum struct {
	RECORD_NEW_FILE_START ShowRecordCallbackConfigResponseNotifyEventSubscription
	RECORD_FILE_COMPLETE  ShowRecordCallbackConfigResponseNotifyEventSubscription
	RECORD_OVER           ShowRecordCallbackConfigResponseNotifyEventSubscription
	RECORD_FAILED         ShowRecordCallbackConfigResponseNotifyEventSubscription
}

func GetShowRecordCallbackConfigResponseNotifyEventSubscriptionEnum() ShowRecordCallbackConfigResponseNotifyEventSubscriptionEnum {
	return ShowRecordCallbackConfigResponseNotifyEventSubscriptionEnum{
		RECORD_NEW_FILE_START: ShowRecordCallbackConfigResponseNotifyEventSubscription{
			value: "RECORD_NEW_FILE_START",
		},
		RECORD_FILE_COMPLETE: ShowRecordCallbackConfigResponseNotifyEventSubscription{
			value: "RECORD_FILE_COMPLETE",
		},
		RECORD_OVER: ShowRecordCallbackConfigResponseNotifyEventSubscription{
			value: "RECORD_OVER",
		},
		RECORD_FAILED: ShowRecordCallbackConfigResponseNotifyEventSubscription{
			value: "RECORD_FAILED",
		},
	}
}

func (c ShowRecordCallbackConfigResponseNotifyEventSubscription) Value() string {
	return c.value
}

func (c ShowRecordCallbackConfigResponseNotifyEventSubscription) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ShowRecordCallbackConfigResponseNotifyEventSubscription) UnmarshalJSON(b []byte) error {
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

type ShowRecordCallbackConfigResponseSignType struct {
	value string
}

type ShowRecordCallbackConfigResponseSignTypeEnum struct {
	MD5        ShowRecordCallbackConfigResponseSignType
	HMACSHA256 ShowRecordCallbackConfigResponseSignType
}

func GetShowRecordCallbackConfigResponseSignTypeEnum() ShowRecordCallbackConfigResponseSignTypeEnum {
	return ShowRecordCallbackConfigResponseSignTypeEnum{
		MD5: ShowRecordCallbackConfigResponseSignType{
			value: "MD5",
		},
		HMACSHA256: ShowRecordCallbackConfigResponseSignType{
			value: "HMACSHA256",
		},
	}
}

func (c ShowRecordCallbackConfigResponseSignType) Value() string {
	return c.value
}

func (c ShowRecordCallbackConfigResponseSignType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ShowRecordCallbackConfigResponseSignType) UnmarshalJSON(b []byte) error {
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
