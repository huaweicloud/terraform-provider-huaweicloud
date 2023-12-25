package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// SetLogLtsConfigsRequest Request Object
type SetLogLtsConfigsRequest struct {

	// 引擎。
	Engine SetLogLtsConfigsRequestEngine `json:"engine"`

	// 语言。
	XLanguage *SetLogLtsConfigsRequestXLanguage `json:"X-Language,omitempty"`

	Body *AddLogConfigResponseBody `json:"body,omitempty"`
}

func (o SetLogLtsConfigsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "SetLogLtsConfigsRequest struct{}"
	}

	return strings.Join([]string{"SetLogLtsConfigsRequest", string(data)}, " ")
}

type SetLogLtsConfigsRequestEngine struct {
	value string
}

type SetLogLtsConfigsRequestEngineEnum struct {
	MYSQL      SetLogLtsConfigsRequestEngine
	POSTGRESQL SetLogLtsConfigsRequestEngine
	SQLSERVER  SetLogLtsConfigsRequestEngine
}

func GetSetLogLtsConfigsRequestEngineEnum() SetLogLtsConfigsRequestEngineEnum {
	return SetLogLtsConfigsRequestEngineEnum{
		MYSQL: SetLogLtsConfigsRequestEngine{
			value: "mysql",
		},
		POSTGRESQL: SetLogLtsConfigsRequestEngine{
			value: "postgresql",
		},
		SQLSERVER: SetLogLtsConfigsRequestEngine{
			value: "sqlserver",
		},
	}
}

func (c SetLogLtsConfigsRequestEngine) Value() string {
	return c.value
}

func (c SetLogLtsConfigsRequestEngine) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *SetLogLtsConfigsRequestEngine) UnmarshalJSON(b []byte) error {
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

type SetLogLtsConfigsRequestXLanguage struct {
	value string
}

type SetLogLtsConfigsRequestXLanguageEnum struct {
	ZH_CN SetLogLtsConfigsRequestXLanguage
	EN_US SetLogLtsConfigsRequestXLanguage
}

func GetSetLogLtsConfigsRequestXLanguageEnum() SetLogLtsConfigsRequestXLanguageEnum {
	return SetLogLtsConfigsRequestXLanguageEnum{
		ZH_CN: SetLogLtsConfigsRequestXLanguage{
			value: "zh-cn",
		},
		EN_US: SetLogLtsConfigsRequestXLanguage{
			value: "en-us",
		},
	}
}

func (c SetLogLtsConfigsRequestXLanguage) Value() string {
	return c.value
}

func (c SetLogLtsConfigsRequestXLanguage) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *SetLogLtsConfigsRequestXLanguage) UnmarshalJSON(b []byte) error {
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
