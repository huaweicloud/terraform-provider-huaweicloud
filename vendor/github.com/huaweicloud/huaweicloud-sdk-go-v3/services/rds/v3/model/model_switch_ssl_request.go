package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// Request Object
type SwitchSslRequest struct {

	// 语言
	XLanguage *SwitchSslRequestXLanguage `json:"X-Language,omitempty"`

	// 实例ID。
	InstanceId string `json:"instance_id"`

	Body *SslOptionRequest `json:"body,omitempty"`
}

func (o SwitchSslRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SwitchSslRequest struct{}"
	}

	return strings.Join([]string{"SwitchSslRequest", string(data)}, " ")
}

type SwitchSslRequestXLanguage struct {
	value string
}

type SwitchSslRequestXLanguageEnum struct {
	ZH_CN SwitchSslRequestXLanguage
	EN_US SwitchSslRequestXLanguage
}

func GetSwitchSslRequestXLanguageEnum() SwitchSslRequestXLanguageEnum {
	return SwitchSslRequestXLanguageEnum{
		ZH_CN: SwitchSslRequestXLanguage{
			value: "zh-cn",
		},
		EN_US: SwitchSslRequestXLanguage{
			value: "en-us",
		},
	}
}

func (c SwitchSslRequestXLanguage) Value() string {
	return c.value
}

func (c SwitchSslRequestXLanguage) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *SwitchSslRequestXLanguage) UnmarshalJSON(b []byte) error {
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
