package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// This is a auto create Body Object
type Follow302StatusRequest struct {

	// follow302状态（\"off\"/\"on\"）
	Follow302Status Follow302StatusRequestFollow302Status `json:"follow302_status"`
}

func (o Follow302StatusRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Follow302StatusRequest struct{}"
	}

	return strings.Join([]string{"Follow302StatusRequest", string(data)}, " ")
}

type Follow302StatusRequestFollow302Status struct {
	value string
}

type Follow302StatusRequestFollow302StatusEnum struct {
	OFF Follow302StatusRequestFollow302Status
	ON  Follow302StatusRequestFollow302Status
}

func GetFollow302StatusRequestFollow302StatusEnum() Follow302StatusRequestFollow302StatusEnum {
	return Follow302StatusRequestFollow302StatusEnum{
		OFF: Follow302StatusRequestFollow302Status{
			value: "off",
		},
		ON: Follow302StatusRequestFollow302Status{
			value: "on",
		},
	}
}

func (c Follow302StatusRequestFollow302Status) Value() string {
	return c.value
}

func (c Follow302StatusRequestFollow302Status) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *Follow302StatusRequestFollow302Status) UnmarshalJSON(b []byte) error {
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
