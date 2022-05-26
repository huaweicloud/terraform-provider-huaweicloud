package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// Request Object
type ListTrackersRequest struct {

	// 标示追踪器名称。 在不传入该字段的情况下，将查询租户所有的追踪器。
	TrackerName *string `json:"tracker_name,omitempty"`

	// 标识追踪器类型。 目前支持系统追踪器有管理类追踪器（system）和数据类追踪器（data）。
	TrackerType *ListTrackersRequestTrackerType `json:"tracker_type,omitempty"`
}

func (o ListTrackersRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListTrackersRequest struct{}"
	}

	return strings.Join([]string{"ListTrackersRequest", string(data)}, " ")
}

type ListTrackersRequestTrackerType struct {
	value string
}

type ListTrackersRequestTrackerTypeEnum struct {
	SYSTEM ListTrackersRequestTrackerType
	DATA   ListTrackersRequestTrackerType
}

func GetListTrackersRequestTrackerTypeEnum() ListTrackersRequestTrackerTypeEnum {
	return ListTrackersRequestTrackerTypeEnum{
		SYSTEM: ListTrackersRequestTrackerType{
			value: "system",
		},
		DATA: ListTrackersRequestTrackerType{
			value: "data",
		},
	}
}

func (c ListTrackersRequestTrackerType) Value() string {
	return c.value
}

func (c ListTrackersRequestTrackerType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ListTrackersRequestTrackerType) UnmarshalJSON(b []byte) error {
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
