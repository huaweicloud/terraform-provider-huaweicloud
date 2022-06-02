package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type RecordCallbackConfigRequest struct {

	// 直播推流域名
	PublishDomain string `json:"publish_domain"`

	// app名称。如果需要匹配任意应用则需填写*。录制规则匹配的时候，优先精确app匹配，如果匹配不到，则匹配*
	App string `json:"app"`

	// 录制回调通知url地址
	NotifyCallbackUrl *string `json:"notify_callback_url,omitempty"`

	// 订阅录制通知消息。消息类型。RECORD_NEW_FILE_START开始创建新的录制文件。RECORD_FILE_COMPLETE录制文件生成完成。RECORD_OVER录制结束。RECORD_FAILED表示录制失败。如果不填写,默认订阅RECORD_FILE_COMPLETE
	NotifyEventSubscription *[]RecordCallbackConfigRequestNotifyEventSubscription `json:"notify_event_subscription,omitempty"`

	// 加密类型
	SignType *RecordCallbackConfigRequestSignType `json:"sign_type,omitempty"`
}

func (o RecordCallbackConfigRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RecordCallbackConfigRequest struct{}"
	}

	return strings.Join([]string{"RecordCallbackConfigRequest", string(data)}, " ")
}

type RecordCallbackConfigRequestNotifyEventSubscription struct {
	value string
}

type RecordCallbackConfigRequestNotifyEventSubscriptionEnum struct {
	RECORD_NEW_FILE_START RecordCallbackConfigRequestNotifyEventSubscription
	RECORD_FILE_COMPLETE  RecordCallbackConfigRequestNotifyEventSubscription
	RECORD_OVER           RecordCallbackConfigRequestNotifyEventSubscription
	RECORD_FAILED         RecordCallbackConfigRequestNotifyEventSubscription
}

func GetRecordCallbackConfigRequestNotifyEventSubscriptionEnum() RecordCallbackConfigRequestNotifyEventSubscriptionEnum {
	return RecordCallbackConfigRequestNotifyEventSubscriptionEnum{
		RECORD_NEW_FILE_START: RecordCallbackConfigRequestNotifyEventSubscription{
			value: "RECORD_NEW_FILE_START",
		},
		RECORD_FILE_COMPLETE: RecordCallbackConfigRequestNotifyEventSubscription{
			value: "RECORD_FILE_COMPLETE",
		},
		RECORD_OVER: RecordCallbackConfigRequestNotifyEventSubscription{
			value: "RECORD_OVER",
		},
		RECORD_FAILED: RecordCallbackConfigRequestNotifyEventSubscription{
			value: "RECORD_FAILED",
		},
	}
}

func (c RecordCallbackConfigRequestNotifyEventSubscription) Value() string {
	return c.value
}

func (c RecordCallbackConfigRequestNotifyEventSubscription) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *RecordCallbackConfigRequestNotifyEventSubscription) UnmarshalJSON(b []byte) error {
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

type RecordCallbackConfigRequestSignType struct {
	value string
}

type RecordCallbackConfigRequestSignTypeEnum struct {
	MD5        RecordCallbackConfigRequestSignType
	HMACSHA256 RecordCallbackConfigRequestSignType
}

func GetRecordCallbackConfigRequestSignTypeEnum() RecordCallbackConfigRequestSignTypeEnum {
	return RecordCallbackConfigRequestSignTypeEnum{
		MD5: RecordCallbackConfigRequestSignType{
			value: "MD5",
		},
		HMACSHA256: RecordCallbackConfigRequestSignType{
			value: "HMACSHA256",
		},
	}
}

func (c RecordCallbackConfigRequestSignType) Value() string {
	return c.value
}

func (c RecordCallbackConfigRequestSignType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *RecordCallbackConfigRequestSignType) UnmarshalJSON(b []byte) error {
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
