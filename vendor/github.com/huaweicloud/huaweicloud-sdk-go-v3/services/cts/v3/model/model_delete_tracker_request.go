package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// Request Object
type DeleteTrackerRequest struct {

	// 标识追踪器名称。 在不传入该字段的情况下，将删除当前租户所有的数据类追踪器。
	TrackerName *string `json:"tracker_name,omitempty"`

	// 标识追踪器类型。 目前仅支持数据类追踪器（data）的删除，默认值为\"data\"。
	TrackerType *DeleteTrackerRequestTrackerType `json:"tracker_type,omitempty"`
}

func (o DeleteTrackerRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteTrackerRequest struct{}"
	}

	return strings.Join([]string{"DeleteTrackerRequest", string(data)}, " ")
}

type DeleteTrackerRequestTrackerType struct {
	value string
}

type DeleteTrackerRequestTrackerTypeEnum struct {
	DATA DeleteTrackerRequestTrackerType
}

func GetDeleteTrackerRequestTrackerTypeEnum() DeleteTrackerRequestTrackerTypeEnum {
	return DeleteTrackerRequestTrackerTypeEnum{
		DATA: DeleteTrackerRequestTrackerType{
			value: "data",
		},
	}
}

func (c DeleteTrackerRequestTrackerType) Value() string {
	return c.value
}

func (c DeleteTrackerRequestTrackerType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *DeleteTrackerRequestTrackerType) UnmarshalJSON(b []byte) error {
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
