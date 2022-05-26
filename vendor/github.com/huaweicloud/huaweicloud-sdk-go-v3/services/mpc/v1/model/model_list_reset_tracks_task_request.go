package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// Request Object
type ListResetTracksTaskRequest struct {

	// 任务ID。一次最多10个
	TaskId *[]string `json:"task_id,omitempty"`

	// 任务执行状态
	Status *ListResetTracksTaskRequestStatus `json:"status,omitempty"`

	// 起始时间。格式为yyyymmddhhmmss。必须是与时区无关的UTC时间，指定task_id时该参数无效
	StartTime *string `json:"start_time,omitempty"`

	// 结束时间。格式为yyyymmddhhmmss。必须是与时区无关的UTC时间，指定task_id时该参数无效
	EndTime *string `json:"end_time,omitempty"`

	// 分页编号。查询指定“task_id”时，该参数无效。  默认值：0。
	Page *int32 `json:"page,omitempty"`

	// 每页记录数。查询指定“task_id”时，该参数无效。  取值范围：[1,100]。  默认值：10。
	Size *int32 `json:"size,omitempty"`
}

func (o ListResetTracksTaskRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListResetTracksTaskRequest struct{}"
	}

	return strings.Join([]string{"ListResetTracksTaskRequest", string(data)}, " ")
}

type ListResetTracksTaskRequestStatus struct {
	value string
}

type ListResetTracksTaskRequestStatusEnum struct {
	WAITING    ListResetTracksTaskRequestStatus
	PROCESSING ListResetTracksTaskRequestStatus
	SUCCEEDED  ListResetTracksTaskRequestStatus
	FAILED     ListResetTracksTaskRequestStatus
	CANCELED   ListResetTracksTaskRequestStatus
}

func GetListResetTracksTaskRequestStatusEnum() ListResetTracksTaskRequestStatusEnum {
	return ListResetTracksTaskRequestStatusEnum{
		WAITING: ListResetTracksTaskRequestStatus{
			value: "WAITING",
		},
		PROCESSING: ListResetTracksTaskRequestStatus{
			value: "PROCESSING",
		},
		SUCCEEDED: ListResetTracksTaskRequestStatus{
			value: "SUCCEEDED",
		},
		FAILED: ListResetTracksTaskRequestStatus{
			value: "FAILED",
		},
		CANCELED: ListResetTracksTaskRequestStatus{
			value: "CANCELED",
		},
	}
}

func (c ListResetTracksTaskRequestStatus) Value() string {
	return c.value
}

func (c ListResetTracksTaskRequestStatus) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ListResetTracksTaskRequestStatus) UnmarshalJSON(b []byte) error {
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
