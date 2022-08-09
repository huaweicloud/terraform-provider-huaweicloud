package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// 关键操作通知高级筛选条件。
type Filter struct {

	// 多条件关系。 - AND 表示所有过滤条件满足后生效。 - OR 表示有任意一个条件满足时生效。
	Condition FilterCondition `json:"condition"`

	// 是否打开高级筛选开关。
	IsSupportFilter bool `json:"is_support_filter"`

	// 高级过滤条件规则，示例如下：\"key != value\"，格式为：字段 规则 值。 -字段取值范围：api_version,code,trace_rating,trace_type,resource_id,resource_name。 -规则：!= 或 =。 - 值：api_version正则约束：^(a-zA-Z0-9_-.){1,64}$；code：最小长度1，最大长度256；trace_rating枚举值：\"normal\", \"warning\", \"incident\"；trace_type枚举值：\"ConsoleAction\", \"ApiCall\", \"SystemAction\"；resource_id：最小长度1，最大长度350；resource_name：最小长度1，最大长度256
	Rule []string `json:"rule"`
}

func (o Filter) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Filter struct{}"
	}

	return strings.Join([]string{"Filter", string(data)}, " ")
}

type FilterCondition struct {
	value string
}

type FilterConditionEnum struct {
	AND FilterCondition
	OR  FilterCondition
}

func GetFilterConditionEnum() FilterConditionEnum {
	return FilterConditionEnum{
		AND: FilterCondition{
			value: "AND",
		},
		OR: FilterCondition{
			value: "OR",
		},
	}
}

func (c FilterCondition) Value() string {
	return c.value
}

func (c FilterCondition) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *FilterCondition) UnmarshalJSON(b []byte) error {
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
