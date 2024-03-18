package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// ProxyInfoFlavorInfo 数据库代理规格信息。
type ProxyInfoFlavorInfo struct {

	// 规格类型。
	GroupType *ProxyInfoFlavorInfoGroupType `json:"group_type,omitempty"`

	// 规格码。
	Code *string `json:"code,omitempty"`
}

func (o ProxyInfoFlavorInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ProxyInfoFlavorInfo struct{}"
	}

	return strings.Join([]string{"ProxyInfoFlavorInfo", string(data)}, " ")
}

type ProxyInfoFlavorInfoGroupType struct {
	value string
}

type ProxyInfoFlavorInfoGroupTypeEnum struct {
	X86 ProxyInfoFlavorInfoGroupType
	RAM ProxyInfoFlavorInfoGroupType
}

func GetProxyInfoFlavorInfoGroupTypeEnum() ProxyInfoFlavorInfoGroupTypeEnum {
	return ProxyInfoFlavorInfoGroupTypeEnum{
		X86: ProxyInfoFlavorInfoGroupType{
			value: "X86",
		},
		RAM: ProxyInfoFlavorInfoGroupType{
			value: "RAM",
		},
	}
}

func (c ProxyInfoFlavorInfoGroupType) Value() string {
	return c.value
}

func (c ProxyInfoFlavorInfoGroupType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ProxyInfoFlavorInfoGroupType) UnmarshalJSON(b []byte) error {
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
