package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// Request Object
type ListMediaProcessTaskRequest struct {

	// 独立加密任务ID。一次最多10个
	TaskId *[]string `json:"task_id,omitempty"`

	// 任务状态。  取值如下： - WAITING：等待启动 - PROCESSING：处理中 - SUCCEEDED：处理成功 - FAILED：处理失败 - CANCELED：已取消
	Status *ListMediaProcessTaskRequestStatus `json:"status,omitempty"`

	// 起始时间。格式为yyyymmddhhmmss。必须是与时区无关的UTC时间，指定task_id时该参数无效
	StartTime *string `json:"start_time,omitempty"`

	// 结束时间。格式为yyyymmddhhmmss。必须是与时区无关的UTC时间，指定task_id时该参数无效
	EndTime *string `json:"end_time,omitempty"`

	// 分页编号。查询指定“task_id”时，该参数无效。  默认值：0。
	Page *int32 `json:"page,omitempty"`

	// 每页记录数。查询指定“task_id”时，该参数无效。  取值范围：[1,100]。  默认值：10。
	Size *int32 `json:"size,omitempty"`
}

func (o ListMediaProcessTaskRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListMediaProcessTaskRequest struct{}"
	}

	return strings.Join([]string{"ListMediaProcessTaskRequest", string(data)}, " ")
}

type ListMediaProcessTaskRequestStatus struct {
	value string
}

type ListMediaProcessTaskRequestStatusEnum struct {
	WAITING    ListMediaProcessTaskRequestStatus
	PROCESSING ListMediaProcessTaskRequestStatus
	SUCCEEDED  ListMediaProcessTaskRequestStatus
	FAILED     ListMediaProcessTaskRequestStatus
	CANCELED   ListMediaProcessTaskRequestStatus
}

func GetListMediaProcessTaskRequestStatusEnum() ListMediaProcessTaskRequestStatusEnum {
	return ListMediaProcessTaskRequestStatusEnum{
		WAITING: ListMediaProcessTaskRequestStatus{
			value: "WAITING",
		},
		PROCESSING: ListMediaProcessTaskRequestStatus{
			value: "PROCESSING",
		},
		SUCCEEDED: ListMediaProcessTaskRequestStatus{
			value: "SUCCEEDED",
		},
		FAILED: ListMediaProcessTaskRequestStatus{
			value: "FAILED",
		},
		CANCELED: ListMediaProcessTaskRequestStatus{
			value: "CANCELED",
		},
	}
}

func (c ListMediaProcessTaskRequestStatus) Value() string {
	return c.value
}

func (c ListMediaProcessTaskRequestStatus) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ListMediaProcessTaskRequestStatus) UnmarshalJSON(b []byte) error {
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
