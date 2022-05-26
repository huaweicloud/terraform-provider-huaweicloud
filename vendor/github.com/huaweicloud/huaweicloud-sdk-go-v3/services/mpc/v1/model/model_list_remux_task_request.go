package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// Request Object
type ListRemuxTaskRequest struct {

	// 任务ID。一次最多10个
	TaskId *[]string `json:"task_id,omitempty"`

	// 任务执行状态。  取值如下： - INIT：初始状态 - WAITING：等待启动 - PROCESSING：处理中 - SUCCEED：处理成功 - FAILED：处理失败 - CANCELED：已取消
	Status *ListRemuxTaskRequestStatus `json:"status,omitempty"`

	// 起始时间。格式为yyyymmddhhmmss。必须是与时区无关的UTC时间，指定task_id时该参数无效
	StartTime *string `json:"start_time,omitempty"`

	// 结束时间。格式为yyyymmddhhmmss。必须是与时区无关的UTC时间，指定task_id时该参数无效
	EndTime *string `json:"end_time,omitempty"`

	// 源文件存储桶。
	InputBucket *string `json:"input_bucket,omitempty"`

	// 源对象名称.
	InputObject *string `json:"input_object,omitempty"`

	// 分页编号。查询指定“task_id”时，该参数无效。  默认值：0。
	Page *int32 `json:"page,omitempty"`

	// 每页记录数。查询指定“task_id”时，该参数无效。  取值范围：[1,100]。  默认值：10。
	Size *int32 `json:"size,omitempty"`
}

func (o ListRemuxTaskRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListRemuxTaskRequest struct{}"
	}

	return strings.Join([]string{"ListRemuxTaskRequest", string(data)}, " ")
}

type ListRemuxTaskRequestStatus struct {
	value string
}

type ListRemuxTaskRequestStatusEnum struct {
	INIT       ListRemuxTaskRequestStatus
	WAITING    ListRemuxTaskRequestStatus
	PROCESSING ListRemuxTaskRequestStatus
	SUCCEED    ListRemuxTaskRequestStatus
	FAILED     ListRemuxTaskRequestStatus
	CANCELED   ListRemuxTaskRequestStatus
}

func GetListRemuxTaskRequestStatusEnum() ListRemuxTaskRequestStatusEnum {
	return ListRemuxTaskRequestStatusEnum{
		INIT: ListRemuxTaskRequestStatus{
			value: "INIT",
		},
		WAITING: ListRemuxTaskRequestStatus{
			value: "WAITING",
		},
		PROCESSING: ListRemuxTaskRequestStatus{
			value: "PROCESSING",
		},
		SUCCEED: ListRemuxTaskRequestStatus{
			value: "SUCCEED",
		},
		FAILED: ListRemuxTaskRequestStatus{
			value: "FAILED",
		},
		CANCELED: ListRemuxTaskRequestStatus{
			value: "CANCELED",
		},
	}
}

func (c ListRemuxTaskRequestStatus) Value() string {
	return c.value
}

func (c ListRemuxTaskRequestStatus) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ListRemuxTaskRequestStatus) UnmarshalJSON(b []byte) error {
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
