package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// SMN消息通知配置。
type SmnConfig struct {

	// 当前用户所使用的管理控制台的语言。  可以选择zh-cn或者en-us。
	Language *SmnConfigLanguage `json:"language,omitempty"`

	// 迁移任务所绑定的SMN消息主题的urn号。
	TopicUrn string `json:"topic_urn"`

	//   SMN消息的触发条件，取决于迁移任务状态。  迁移任务状态的取值范围为SUCCESS或者FAILURE。  - FAILURE表示任务失败后发送SMN消息。 - SUCCESS表示任务成功后发送SMN消息。
	TriggerConditions []string `json:"trigger_conditions"`
}

func (o SmnConfig) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SmnConfig struct{}"
	}

	return strings.Join([]string{"SmnConfig", string(data)}, " ")
}

type SmnConfigLanguage struct {
	value string
}

type SmnConfigLanguageEnum struct {
	ZH_CN SmnConfigLanguage
	EN_US SmnConfigLanguage
}

func GetSmnConfigLanguageEnum() SmnConfigLanguageEnum {
	return SmnConfigLanguageEnum{
		ZH_CN: SmnConfigLanguage{
			value: "zh-cn",
		},
		EN_US: SmnConfigLanguage{
			value: "en-us",
		},
	}
}

func (c SmnConfigLanguage) Value() string {
	return c.value
}

func (c SmnConfigLanguage) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *SmnConfigLanguage) UnmarshalJSON(b []byte) error {
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
