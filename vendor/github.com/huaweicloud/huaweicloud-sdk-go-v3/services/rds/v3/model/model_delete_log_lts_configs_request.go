package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// DeleteLogLtsConfigsRequest Request Object
type DeleteLogLtsConfigsRequest struct {

	// 引擎。
	Engine DeleteLogLtsConfigsRequestEngine `json:"engine"`

	// 语言。
	XLanguage *DeleteLogLtsConfigsRequestXLanguage `json:"X-Language,omitempty"`

	Body *DeleteLogConfigResponseBody `json:"body,omitempty"`
}

func (o DeleteLogLtsConfigsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteLogLtsConfigsRequest struct{}"
	}

	return strings.Join([]string{"DeleteLogLtsConfigsRequest", string(data)}, " ")
}

type DeleteLogLtsConfigsRequestEngine struct {
	value string
}

type DeleteLogLtsConfigsRequestEngineEnum struct {
	MYSQL      DeleteLogLtsConfigsRequestEngine
	POSTGRESQL DeleteLogLtsConfigsRequestEngine
	SQLSERVER  DeleteLogLtsConfigsRequestEngine
}

func GetDeleteLogLtsConfigsRequestEngineEnum() DeleteLogLtsConfigsRequestEngineEnum {
	return DeleteLogLtsConfigsRequestEngineEnum{
		MYSQL: DeleteLogLtsConfigsRequestEngine{
			value: "mysql",
		},
		POSTGRESQL: DeleteLogLtsConfigsRequestEngine{
			value: "postgresql",
		},
		SQLSERVER: DeleteLogLtsConfigsRequestEngine{
			value: "sqlserver",
		},
	}
}

func (c DeleteLogLtsConfigsRequestEngine) Value() string {
	return c.value
}

func (c DeleteLogLtsConfigsRequestEngine) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *DeleteLogLtsConfigsRequestEngine) UnmarshalJSON(b []byte) error {
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

type DeleteLogLtsConfigsRequestXLanguage struct {
	value string
}

type DeleteLogLtsConfigsRequestXLanguageEnum struct {
	ZH_CN DeleteLogLtsConfigsRequestXLanguage
	EN_US DeleteLogLtsConfigsRequestXLanguage
}

func GetDeleteLogLtsConfigsRequestXLanguageEnum() DeleteLogLtsConfigsRequestXLanguageEnum {
	return DeleteLogLtsConfigsRequestXLanguageEnum{
		ZH_CN: DeleteLogLtsConfigsRequestXLanguage{
			value: "zh-cn",
		},
		EN_US: DeleteLogLtsConfigsRequestXLanguage{
			value: "en-us",
		},
	}
}

func (c DeleteLogLtsConfigsRequestXLanguage) Value() string {
	return c.value
}

func (c DeleteLogLtsConfigsRequestXLanguage) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *DeleteLogLtsConfigsRequestXLanguage) UnmarshalJSON(b []byte) error {
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
