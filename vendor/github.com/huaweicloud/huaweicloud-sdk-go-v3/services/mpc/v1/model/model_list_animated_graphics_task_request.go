package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// Request Object
type ListAnimatedGraphicsTaskRequest struct {

	// 客户端语言
	XLanguage *string `json:"x-language,omitempty"`

	// 任务ID。一次最多10个
	TaskId *[]string `json:"task_id,omitempty"`

	// 任务执行状态。  取值如下： - INIT：初始状态 - WAITING：待启动 - PREPROCESSING：处理中 - SUCCEED：处理成功 - FAILED：处理失败 - CANCELED：已取消
	Status *ListAnimatedGraphicsTaskRequestStatus `json:"status,omitempty"`

	// 起始时间。格式为yyyymmddhhmmss。必须是与时区无关的UTC时间，指定task_id时该参数无效。
	StartTime *string `json:"start_time,omitempty"`

	// 结束时间。格式为yyyymmddhhmmss。必须是与时区无关的UTC时间，指定task_id时该参数无效。
	EndTime *string `json:"end_time,omitempty"`

	// 分页编号。查询指定“task_id”时，该参数无效。  默认值：0。
	Page *int32 `json:"page,omitempty"`

	// 每页记录数。查询指定“task_id”时，该参数无效。  取值范围：[1,100]。  默认值：10。
	Size *int32 `json:"size,omitempty"`
}

func (o ListAnimatedGraphicsTaskRequest) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ListAnimatedGraphicsTaskRequest struct{}"
	}

	return strings.Join([]string{"ListAnimatedGraphicsTaskRequest", string(data)}, " ")
}

type ListAnimatedGraphicsTaskRequestStatus struct {
	value string
}

type ListAnimatedGraphicsTaskRequestStatusEnum struct {
	INIT          ListAnimatedGraphicsTaskRequestStatus
	WAITING       ListAnimatedGraphicsTaskRequestStatus
	PREPROCESSING ListAnimatedGraphicsTaskRequestStatus
	SUCCEED       ListAnimatedGraphicsTaskRequestStatus
	FAILED        ListAnimatedGraphicsTaskRequestStatus
	CANCELED      ListAnimatedGraphicsTaskRequestStatus
}

func GetListAnimatedGraphicsTaskRequestStatusEnum() ListAnimatedGraphicsTaskRequestStatusEnum {
	return ListAnimatedGraphicsTaskRequestStatusEnum{
		INIT: ListAnimatedGraphicsTaskRequestStatus{
			value: "INIT",
		},
		WAITING: ListAnimatedGraphicsTaskRequestStatus{
			value: "WAITING",
		},
		PREPROCESSING: ListAnimatedGraphicsTaskRequestStatus{
			value: "PREPROCESSING",
		},
		SUCCEED: ListAnimatedGraphicsTaskRequestStatus{
			value: "SUCCEED",
		},
		FAILED: ListAnimatedGraphicsTaskRequestStatus{
			value: "FAILED",
		},
		CANCELED: ListAnimatedGraphicsTaskRequestStatus{
			value: "CANCELED",
		},
	}
}

func (c ListAnimatedGraphicsTaskRequestStatus) Value() string {
	return c.value
}

func (c ListAnimatedGraphicsTaskRequestStatus) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ListAnimatedGraphicsTaskRequestStatus) UnmarshalJSON(b []byte) error {
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
