package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// CreateInstanceIam5Request Request Object
type CreateInstanceIam5Request struct {

	// 语言
	XLanguage *CreateInstanceIam5RequestXLanguage `json:"X-Language,omitempty"`

	// 保证客户端请求幂等性的标识。 该标识为32位UUID格式，由客户端生成，且需确保72小时内不同请求之间该标识具有唯一性。
	XClientToken *string `json:"X-Client-Token,omitempty"`

	Body *CustomerCreateInstanceReq `json:"body,omitempty"`
}

func (o CreateInstanceIam5Request) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "CreateInstanceIam5Request struct{}"
	}

	return strings.Join([]string{"CreateInstanceIam5Request", string(data)}, " ")
}

type CreateInstanceIam5RequestXLanguage struct {
	value string
}

type CreateInstanceIam5RequestXLanguageEnum struct {
	ZH_CN CreateInstanceIam5RequestXLanguage
	EN_US CreateInstanceIam5RequestXLanguage
}

func GetCreateInstanceIam5RequestXLanguageEnum() CreateInstanceIam5RequestXLanguageEnum {
	return CreateInstanceIam5RequestXLanguageEnum{
		ZH_CN: CreateInstanceIam5RequestXLanguage{
			value: "zh-cn",
		},
		EN_US: CreateInstanceIam5RequestXLanguage{
			value: "en-us",
		},
	}
}

func (c CreateInstanceIam5RequestXLanguage) Value() string {
	return c.value
}

func (c CreateInstanceIam5RequestXLanguage) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *CreateInstanceIam5RequestXLanguage) UnmarshalJSON(b []byte) error {
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
