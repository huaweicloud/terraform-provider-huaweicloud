package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// SmnTopics smn主题相关信息
type SmnTopics struct {

	// Topic的显示名，推送邮件消息时，作为邮件发件人显示。显示名的长度为192byte或64个中文。默认值为空。
	DisplayName *string `json:"display_name,omitempty"`

	// 创建topic的名字。Topic名称只能包含大写字母、小写字母、数字、-和_，且必须由大写字母、小写字母或数字开头，长度为1到255个字符。
	Name string `json:"name"`

	// SMN消息推送策略。取值为0或1
	PushPolicy int32 `json:"push_policy"`

	// topic中订阅者的状态。0:主题已删除或主题下订阅列表为空。1:主题下的订阅列表存在状态为“已订阅”的订阅信息。2:主题下的订阅信息状态处于“未订阅”或“已取消”。
	Status *SmnTopicsStatus `json:"status,omitempty"`

	// Topic的唯一的资源标识。
	TopicUrn string `json:"topic_urn"`
}

func (o SmnTopics) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SmnTopics struct{}"
	}

	return strings.Join([]string{"SmnTopics", string(data)}, " ")
}

type SmnTopicsStatus struct {
	value int32
}

type SmnTopicsStatusEnum struct {
	E_0 SmnTopicsStatus
	E_1 SmnTopicsStatus
	E_2 SmnTopicsStatus
}

func GetSmnTopicsStatusEnum() SmnTopicsStatusEnum {
	return SmnTopicsStatusEnum{
		E_0: SmnTopicsStatus{
			value: 0,
		}, E_1: SmnTopicsStatus{
			value: 1,
		}, E_2: SmnTopicsStatus{
			value: 2,
		},
	}
}

func (c SmnTopicsStatus) Value() int32 {
	return c.value
}

func (c SmnTopicsStatus) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *SmnTopicsStatus) UnmarshalJSON(b []byte) error {
	myConverter := converter.StringConverterFactory("int32")
	if myConverter == nil {
		return errors.New("unsupported StringConverter type: int32")
	}

	interf, err := myConverter.CovertStringToInterface(strings.Trim(string(b[:]), "\""))
	if err != nil {
		return err
	}

	if val, ok := interf.(int32); ok {
		c.value = val
		return nil
	} else {
		return errors.New("convert enum data to int32 error")
	}
}
