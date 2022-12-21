package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// Request Object
type ListProvidersRequest struct {

	// 指定显示语言
	Locale *ListProvidersRequestLocale `json:"locale,omitempty"`

	// 查询记录数默认为200，limit最多为200，最小值为1。
	Limit *int32 `json:"limit,omitempty"`

	// 索引位置，从offset指定的下一条数据开始查询，必须为数字，不能为负数，默认为0。
	Offset *int32 `json:"offset,omitempty"`

	// 云服务名称
	Provider *string `json:"provider,omitempty"`
}

func (o ListProvidersRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListProvidersRequest struct{}"
	}

	return strings.Join([]string{"ListProvidersRequest", string(data)}, " ")
}

type ListProvidersRequestLocale struct {
	value string
}

type ListProvidersRequestLocaleEnum struct {
	ZH_CN ListProvidersRequestLocale
	EN_US ListProvidersRequestLocale
}

func GetListProvidersRequestLocaleEnum() ListProvidersRequestLocaleEnum {
	return ListProvidersRequestLocaleEnum{
		ZH_CN: ListProvidersRequestLocale{
			value: "zh-cn",
		},
		EN_US: ListProvidersRequestLocale{
			value: "en-us",
		},
	}
}

func (c ListProvidersRequestLocale) Value() string {
	return c.value
}

func (c ListProvidersRequestLocale) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ListProvidersRequestLocale) UnmarshalJSON(b []byte) error {
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
