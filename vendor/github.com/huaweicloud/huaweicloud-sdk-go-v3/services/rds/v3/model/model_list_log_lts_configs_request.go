package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// ListLogLtsConfigsRequest Request Object
type ListLogLtsConfigsRequest struct {

	// 引擎。
	Engine ListLogLtsConfigsRequestEngine `json:"engine"`

	// 企业项目ID。默认为空。
	EnterpriseProjectId *string `json:"enterprise_project_id,omitempty"`

	// 实例ID。默认为空。
	InstanceId *string `json:"instance_id,omitempty"`

	// 实例名称。默认为空。
	InstanceName *string `json:"instance_name,omitempty"`

	// 查询记录数。默认10。
	Limit *int32 `json:"limit,omitempty"`

	// 索引位置，偏移量。默认0。
	Offset *int32 `json:"offset,omitempty"`

	// 排序
	Sort *string `json:"sort,omitempty"`

	// 实例状态
	InstanceStatus *string `json:"instance_status,omitempty"`

	// 语言。
	XLanguage *ListLogLtsConfigsRequestXLanguage `json:"X-Language,omitempty"`
}

func (o ListLogLtsConfigsRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListLogLtsConfigsRequest struct{}"
	}

	return strings.Join([]string{"ListLogLtsConfigsRequest", string(data)}, " ")
}

type ListLogLtsConfigsRequestEngine struct {
	value string
}

type ListLogLtsConfigsRequestEngineEnum struct {
	MYSQL      ListLogLtsConfigsRequestEngine
	POSTGRESQL ListLogLtsConfigsRequestEngine
	SQLSERVER  ListLogLtsConfigsRequestEngine
}

func GetListLogLtsConfigsRequestEngineEnum() ListLogLtsConfigsRequestEngineEnum {
	return ListLogLtsConfigsRequestEngineEnum{
		MYSQL: ListLogLtsConfigsRequestEngine{
			value: "mysql",
		},
		POSTGRESQL: ListLogLtsConfigsRequestEngine{
			value: "postgresql",
		},
		SQLSERVER: ListLogLtsConfigsRequestEngine{
			value: "sqlserver",
		},
	}
}

func (c ListLogLtsConfigsRequestEngine) Value() string {
	return c.value
}

func (c ListLogLtsConfigsRequestEngine) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ListLogLtsConfigsRequestEngine) UnmarshalJSON(b []byte) error {
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

type ListLogLtsConfigsRequestXLanguage struct {
	value string
}

type ListLogLtsConfigsRequestXLanguageEnum struct {
	ZH_CN ListLogLtsConfigsRequestXLanguage
	EN_US ListLogLtsConfigsRequestXLanguage
}

func GetListLogLtsConfigsRequestXLanguageEnum() ListLogLtsConfigsRequestXLanguageEnum {
	return ListLogLtsConfigsRequestXLanguageEnum{
		ZH_CN: ListLogLtsConfigsRequestXLanguage{
			value: "zh-cn",
		},
		EN_US: ListLogLtsConfigsRequestXLanguage{
			value: "en-us",
		},
	}
}

func (c ListLogLtsConfigsRequestXLanguage) Value() string {
	return c.value
}

func (c ListLogLtsConfigsRequestXLanguage) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ListLogLtsConfigsRequestXLanguage) UnmarshalJSON(b []byte) error {
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
