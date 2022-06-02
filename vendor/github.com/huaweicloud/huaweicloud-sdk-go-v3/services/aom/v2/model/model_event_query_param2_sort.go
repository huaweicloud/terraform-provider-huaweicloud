package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// 返回列表的排序方式，可以为空。
type EventQueryParam2Sort struct {

	// 排序字段列表。会根据列表中定义顺序对返回列表最排序。
	OrderBy *[]string `json:"order_by,omitempty"`

	// 排序方式枚举值。asc代表正序，desc代表倒叙。
	Order *EventQueryParam2SortOrder `json:"order,omitempty"`
}

func (o EventQueryParam2Sort) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "EventQueryParam2Sort struct{}"
	}

	return strings.Join([]string{"EventQueryParam2Sort", string(data)}, " ")
}

type EventQueryParam2SortOrder struct {
	value string
}

type EventQueryParam2SortOrderEnum struct {
	ASC  EventQueryParam2SortOrder
	DESC EventQueryParam2SortOrder
}

func GetEventQueryParam2SortOrderEnum() EventQueryParam2SortOrderEnum {
	return EventQueryParam2SortOrderEnum{
		ASC: EventQueryParam2SortOrder{
			value: "asc",
		},
		DESC: EventQueryParam2SortOrder{
			value: "desc",
		},
	}
}

func (c EventQueryParam2SortOrder) Value() string {
	return c.value
}

func (c EventQueryParam2SortOrder) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *EventQueryParam2SortOrder) UnmarshalJSON(b []byte) error {
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
