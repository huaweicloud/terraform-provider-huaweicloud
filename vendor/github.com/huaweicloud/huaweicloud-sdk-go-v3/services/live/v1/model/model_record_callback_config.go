package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type RecordCallbackConfig struct {

	// 配置id，由服务端返回。创建或修改的时候不携带
	Id *string `json:"id,omitempty"`

	// 直播推流域名
	PublishDomain string `json:"publish_domain"`

	// app名称。如果匹配任意需填写为*。录制规则匹配的时候，优先精确app匹配，如果匹配不到，则匹配*
	App string `json:"app"`

	// 录制回调通知url地址
	NotifyCallbackUrl *string `json:"notify_callback_url,omitempty"`

	// 订阅录制通知消息。消息类型。RECORD_NEW_FILE_START开始创建新的录制文件。RECORD_FILE_COMPLETE录制文件生成完成。RECORD_OVER录制结束。RECORD_FAILED表示录制失败。如果不填写,默认订阅RECORD_FILE_COMPLETE
	NotifyEventSubscription *[]RecordCallbackConfigNotifyEventSubscription `json:"notify_event_subscription,omitempty"`

	// 加密类型
	SignType *RecordCallbackConfigSignType `json:"sign_type,omitempty"`

	// 创建时间，格式：yyyy-mm-ddThh:mm:ssZ，UTC时间。 在查询的时候返回
	CreateTime *string `json:"create_time,omitempty"`

	// 修改时间，格式：yyyy-mm-ddThh:mm:ssZ，UTC时间。 在查询的时候返回
	UpdateTime *string `json:"update_time,omitempty"`
}

func (o RecordCallbackConfig) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "RecordCallbackConfig struct{}"
	}

	return strings.Join([]string{"RecordCallbackConfig", string(data)}, " ")
}

type RecordCallbackConfigNotifyEventSubscription struct {
	value string
}

type RecordCallbackConfigNotifyEventSubscriptionEnum struct {
	RECORD_NEW_FILE_START RecordCallbackConfigNotifyEventSubscription
	RECORD_FILE_COMPLETE  RecordCallbackConfigNotifyEventSubscription
	RECORD_OVER           RecordCallbackConfigNotifyEventSubscription
	RECORD_FAILED         RecordCallbackConfigNotifyEventSubscription
}

func GetRecordCallbackConfigNotifyEventSubscriptionEnum() RecordCallbackConfigNotifyEventSubscriptionEnum {
	return RecordCallbackConfigNotifyEventSubscriptionEnum{
		RECORD_NEW_FILE_START: RecordCallbackConfigNotifyEventSubscription{
			value: "RECORD_NEW_FILE_START",
		},
		RECORD_FILE_COMPLETE: RecordCallbackConfigNotifyEventSubscription{
			value: "RECORD_FILE_COMPLETE",
		},
		RECORD_OVER: RecordCallbackConfigNotifyEventSubscription{
			value: "RECORD_OVER",
		},
		RECORD_FAILED: RecordCallbackConfigNotifyEventSubscription{
			value: "RECORD_FAILED",
		},
	}
}

func (c RecordCallbackConfigNotifyEventSubscription) Value() string {
	return c.value
}

func (c RecordCallbackConfigNotifyEventSubscription) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *RecordCallbackConfigNotifyEventSubscription) UnmarshalJSON(b []byte) error {
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

type RecordCallbackConfigSignType struct {
	value string
}

type RecordCallbackConfigSignTypeEnum struct {
	MD5        RecordCallbackConfigSignType
	HMACSHA256 RecordCallbackConfigSignType
}

func GetRecordCallbackConfigSignTypeEnum() RecordCallbackConfigSignTypeEnum {
	return RecordCallbackConfigSignTypeEnum{
		MD5: RecordCallbackConfigSignType{
			value: "MD5",
		},
		HMACSHA256: RecordCallbackConfigSignType{
			value: "HMACSHA256",
		},
	}
}

func (c RecordCallbackConfigSignType) Value() string {
	return c.value
}

func (c RecordCallbackConfigSignType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *RecordCallbackConfigSignType) UnmarshalJSON(b []byte) error {
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
