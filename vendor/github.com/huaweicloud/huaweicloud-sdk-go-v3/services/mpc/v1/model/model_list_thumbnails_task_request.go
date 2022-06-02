package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// Request Object
type ListThumbnailsTaskRequest struct {

	// 客户端语言
	XLanguage *string `json:"x-language,omitempty"`

	// 截图服务接受任务后产生的任务ID。一次最多10个
	TaskId *[]string `json:"task_id,omitempty"`

	// 任务状态。  取值如下： - WAITING: 等待启动 - PROCESSING：截图中 - SUCCEEDED：截图成功 - FAILED：截图失败 - CANCELED：已删除
	Status *ListThumbnailsTaskRequestStatus `json:"status,omitempty"`

	// 起始时间。格式为yyyymmddhhmmss。必须是与时区无关的UTC时间，指定task_id时该参数无效
	StartTime *string `json:"start_time,omitempty"`

	// 结束时间。格式为yyyymmddhhmmss。必须是与时区无关的UTC时间，指定task_id时该参数无效
	EndTime *string `json:"end_time,omitempty"`

	// 分页编号。查询指定“task_id”时，该参数无效。  默认值：0。
	Page *int32 `json:"page,omitempty"`

	// 每页记录数。查询指定“task_id”时，该参数无效。  取值范围：[1,100]。  默认值：10。
	Size *int32 `json:"size,omitempty"`
}

func (o ListThumbnailsTaskRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListThumbnailsTaskRequest struct{}"
	}

	return strings.Join([]string{"ListThumbnailsTaskRequest", string(data)}, " ")
}

type ListThumbnailsTaskRequestStatus struct {
	value string
}

type ListThumbnailsTaskRequestStatusEnum struct {
	WAITING    ListThumbnailsTaskRequestStatus
	PROCESSING ListThumbnailsTaskRequestStatus
	SUCCEEDED  ListThumbnailsTaskRequestStatus
	FAILED     ListThumbnailsTaskRequestStatus
	CANCELED   ListThumbnailsTaskRequestStatus
}

func GetListThumbnailsTaskRequestStatusEnum() ListThumbnailsTaskRequestStatusEnum {
	return ListThumbnailsTaskRequestStatusEnum{
		WAITING: ListThumbnailsTaskRequestStatus{
			value: "WAITING",
		},
		PROCESSING: ListThumbnailsTaskRequestStatus{
			value: "PROCESSING",
		},
		SUCCEEDED: ListThumbnailsTaskRequestStatus{
			value: "SUCCEEDED",
		},
		FAILED: ListThumbnailsTaskRequestStatus{
			value: "FAILED",
		},
		CANCELED: ListThumbnailsTaskRequestStatus{
			value: "CANCELED",
		},
	}
}

func (c ListThumbnailsTaskRequestStatus) Value() string {
	return c.value
}

func (c ListThumbnailsTaskRequestStatus) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ListThumbnailsTaskRequestStatus) UnmarshalJSON(b []byte) error {
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
