package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type MysqlProxyFlavorsResponseComputeFlavorGroups struct {

	// 规格组类型，如x86、arm。
	GroupType *MysqlProxyFlavorsResponseComputeFlavorGroupsGroupType `json:"group_type,omitempty"`

	// 规格信息。
	ComputeFlavors *[]MysqlProxyFlavorsResponseComputeFlavors `json:"compute_flavors,omitempty"`
}

func (o MysqlProxyFlavorsResponseComputeFlavorGroups) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "MysqlProxyFlavorsResponseComputeFlavorGroups struct{}"
	}

	return strings.Join([]string{"MysqlProxyFlavorsResponseComputeFlavorGroups", string(data)}, " ")
}

type MysqlProxyFlavorsResponseComputeFlavorGroupsGroupType struct {
	value string
}

type MysqlProxyFlavorsResponseComputeFlavorGroupsGroupTypeEnum struct {
	X86 MysqlProxyFlavorsResponseComputeFlavorGroupsGroupType
	ARM MysqlProxyFlavorsResponseComputeFlavorGroupsGroupType
}

func GetMysqlProxyFlavorsResponseComputeFlavorGroupsGroupTypeEnum() MysqlProxyFlavorsResponseComputeFlavorGroupsGroupTypeEnum {
	return MysqlProxyFlavorsResponseComputeFlavorGroupsGroupTypeEnum{
		X86: MysqlProxyFlavorsResponseComputeFlavorGroupsGroupType{
			value: "x86",
		},
		ARM: MysqlProxyFlavorsResponseComputeFlavorGroupsGroupType{
			value: "arm",
		},
	}
}

func (c MysqlProxyFlavorsResponseComputeFlavorGroupsGroupType) Value() string {
	return c.value
}

func (c MysqlProxyFlavorsResponseComputeFlavorGroupsGroupType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *MysqlProxyFlavorsResponseComputeFlavorGroupsGroupType) UnmarshalJSON(b []byte) error {
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
