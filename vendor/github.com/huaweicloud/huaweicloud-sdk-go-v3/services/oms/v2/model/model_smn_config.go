package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// SmnConfig SMN消息通知配置。
type SmnConfig struct {

	// 当前用户所使用的管理控制台的语言。  可以选择zh-cn或者en-us。
	Language *SmnConfigLanguage `json:"language,omitempty"`

	// 迁移任务所绑定的SMN消息主题的urn号。
	TopicUrn string `json:"topic_urn"`

	// SMN消息的触发条件，取决于迁移任务状态。  迁移任务状态的取值范围为SUCCESS或者FAILURE。  - FAILURE表示任务失败后发送SMN消息。 - SUCCESS表示任务成功后发送SMN消息。
	TriggerConditions []string `json:"trigger_conditions"`

	// 如果设置此值，则表示用模板方式发送smn信息。 模板示例: {  “Task_Status”: \"\",     \"Task_Name\" : \"\",     \"Start_Time\": \"\",     \"Total_Time_Used\": \"\",     \"Transferred_Data\": \"\",     \"Average_Speed\": \"\",     \"Source_Bucket\": \"\",     \"Destination_Bucket\": \"\",     \"List_File_Bucket\": \"\",     \"List_File_Key\": \"\",     \"Success_object_list_path\": \"\",     \"Skip_object_list_path\": \"\",     \"Failed_object_list_path\": \"\" }
	MessageTemplateName *string `json:"message_template_name,omitempty"`
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
	if myConverter == nil {
		return errors.New("unsupported StringConverter type: string")
	}

	interf, err := myConverter.CovertStringToInterface(strings.Trim(string(b[:]), "\""))
	if err != nil {
		return err
	}

	if val, ok := interf.(string); ok {
		c.value = val
		return nil
	} else {
		return errors.New("convert enum data to string error")
	}
}
