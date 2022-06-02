package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// 返回列表的排序方式，可以为空。
type EventQueryParamSort struct {

	// 排序字段列表。会根据列表中定义顺序对返回列表进行排序。
	OrderBy *[]string `json:"order_by,omitempty"`

	// 排序方式枚举值。asc代表正序，desc代表倒叙。
	Order *EventQueryParamSortOrder `json:"order,omitempty"`
}

func (o EventQueryParamSort) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "EventQueryParamSort struct{}"
	}

	return strings.Join([]string{"EventQueryParamSort", string(data)}, " ")
}

type EventQueryParamSortOrder struct {
	value string
}

type EventQueryParamSortOrderEnum struct {
	ASC  EventQueryParamSortOrder
	DESC EventQueryParamSortOrder
}

func GetEventQueryParamSortOrderEnum() EventQueryParamSortOrderEnum {
	return EventQueryParamSortOrderEnum{
		ASC: EventQueryParamSortOrder{
			value: "asc",
		},
		DESC: EventQueryParamSortOrder{
			value: "desc",
		},
	}
}

func (c EventQueryParamSortOrder) Value() string {
	return c.value
}

func (c EventQueryParamSortOrder) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *EventQueryParamSortOrder) UnmarshalJSON(b []byte) error {
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
