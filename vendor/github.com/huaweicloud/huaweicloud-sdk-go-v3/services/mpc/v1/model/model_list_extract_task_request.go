package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// Request Object
type ListExtractTaskRequest struct {

	// 客户端语言
	XLanguage *string `json:"x-language,omitempty"`

	// 任务ID。一次最多10个
	TaskId *[]string `json:"task_id,omitempty"`

	// 任务执行状态。  取值如下： - INIT：初始状态 - WAITING：等待启动 - PREPROCESSING：处理中 - SUCCEED：处理成功 - FAILED：处理失败 - CANCELED：已取消
	Status *ListExtractTaskRequestStatus `json:"status,omitempty"`

	// 起始时间。格式为yyyymmddhhmmss。必须是与时区无关的UTC时间，指定task_id时该参数无效。
	StartTime *string `json:"start_time,omitempty"`

	// 结束时间。格式为yyyymmddhhmmss。必须是与时区无关的UTC时间，指定task_id时该参数无效。
	EndTime *string `json:"end_time,omitempty"`

	// 分页编号。查询指定“task_id”时，该参数无效。  默认值：0。
	Page *int32 `json:"page,omitempty"`

	// 每页记录数。查询指定“task_id”时，该参数无效。  取值范围：[1,100]。  默认值：10。
	Size *int32 `json:"size,omitempty"`
}

func (o ListExtractTaskRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListExtractTaskRequest struct{}"
	}

	return strings.Join([]string{"ListExtractTaskRequest", string(data)}, " ")
}

type ListExtractTaskRequestStatus struct {
	value string
}

type ListExtractTaskRequestStatusEnum struct {
	INIT          ListExtractTaskRequestStatus
	WAITING       ListExtractTaskRequestStatus
	PREPROCESSING ListExtractTaskRequestStatus
	SUCCEED       ListExtractTaskRequestStatus
	FAILED        ListExtractTaskRequestStatus
	CANCELED      ListExtractTaskRequestStatus
}

func GetListExtractTaskRequestStatusEnum() ListExtractTaskRequestStatusEnum {
	return ListExtractTaskRequestStatusEnum{
		INIT: ListExtractTaskRequestStatus{
			value: "INIT",
		},
		WAITING: ListExtractTaskRequestStatus{
			value: "WAITING",
		},
		PREPROCESSING: ListExtractTaskRequestStatus{
			value: "PREPROCESSING",
		},
		SUCCEED: ListExtractTaskRequestStatus{
			value: "SUCCEED",
		},
		FAILED: ListExtractTaskRequestStatus{
			value: "FAILED",
		},
		CANCELED: ListExtractTaskRequestStatus{
			value: "CANCELED",
		},
	}
}

func (c ListExtractTaskRequestStatus) Value() string {
	return c.value
}

func (c ListExtractTaskRequestStatus) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ListExtractTaskRequestStatus) UnmarshalJSON(b []byte) error {
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
