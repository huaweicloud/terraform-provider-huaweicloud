package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// ResourceSelector 资源标签选择器，匹配Kubernetes中node资源
type ResourceSelector struct {

	// 标签键值，取值如下 - node.uid：节点UID。
	Key ResourceSelectorKey `json:"key"`

	// 标签值列表
	Values *[]string `json:"values,omitempty"`

	// 标签逻辑运算符，当前支持如下取值 - In
	Operator ResourceSelectorOperator `json:"operator"`
}

func (o ResourceSelector) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ResourceSelector struct{}"
	}

	return strings.Join([]string{"ResourceSelector", string(data)}, " ")
}

type ResourceSelectorKey struct {
	value string
}

type ResourceSelectorKeyEnum struct {
	NODE_UID ResourceSelectorKey
}

func GetResourceSelectorKeyEnum() ResourceSelectorKeyEnum {
	return ResourceSelectorKeyEnum{
		NODE_UID: ResourceSelectorKey{
			value: "node.uid",
		},
	}
}

func (c ResourceSelectorKey) Value() string {
	return c.value
}

func (c ResourceSelectorKey) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ResourceSelectorKey) UnmarshalJSON(b []byte) error {
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

type ResourceSelectorOperator struct {
	value string
}

type ResourceSelectorOperatorEnum struct {
	IN ResourceSelectorOperator
}

func GetResourceSelectorOperatorEnum() ResourceSelectorOperatorEnum {
	return ResourceSelectorOperatorEnum{
		IN: ResourceSelectorOperator{
			value: "In",
		},
	}
}

func (c ResourceSelectorOperator) Value() string {
	return c.value
}

func (c ResourceSelectorOperator) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ResourceSelectorOperator) UnmarshalJSON(b []byte) error {
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
