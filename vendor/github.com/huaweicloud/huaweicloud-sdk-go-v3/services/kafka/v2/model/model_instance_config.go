package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type InstanceConfig struct {

	// 配置名称。
	Name *string `json:"name,omitempty"`

	// 有效值。
	ValidValues *string `json:"valid_values,omitempty"`

	// 默认值。
	DefaultValue *string `json:"default_value,omitempty"`

	// 配置类型：static/dynamic。
	ConfigType *InstanceConfigConfigType `json:"config_type,omitempty"`

	// 配置当前值。
	Value *string `json:"value,omitempty"`

	// 值类型。
	ValueType *string `json:"value_type,omitempty"`
}

func (o InstanceConfig) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "InstanceConfig struct{}"
	}

	return strings.Join([]string{"InstanceConfig", string(data)}, " ")
}

type InstanceConfigConfigType struct {
	value string
}

type InstanceConfigConfigTypeEnum struct {
	STATIC  InstanceConfigConfigType
	DYNAMIC InstanceConfigConfigType
}

func GetInstanceConfigConfigTypeEnum() InstanceConfigConfigTypeEnum {
	return InstanceConfigConfigTypeEnum{
		STATIC: InstanceConfigConfigType{
			value: "static",
		},
		DYNAMIC: InstanceConfigConfigType{
			value: "dynamic",
		},
	}
}

func (c InstanceConfigConfigType) Value() string {
	return c.value
}

func (c InstanceConfigConfigType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *InstanceConfigConfigType) UnmarshalJSON(b []byte) error {
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
