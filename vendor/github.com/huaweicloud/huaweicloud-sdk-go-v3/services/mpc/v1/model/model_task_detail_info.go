package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type TaskDetailInfo struct {

	// 任务ID。
	TaskId *string `json:"task_id,omitempty"`

	// 任务执行状态，取值如下。
	Status *TaskDetailInfoStatus `json:"status,omitempty"`

	// 转码任务启动时间
	CreateTime *string `json:"create_time,omitempty"`

	// 下发xcode任务成功时间
	StartTime *string `json:"start_time,omitempty"`

	// 转码任务结束时间
	EndTime *string `json:"end_time,omitempty"`

	Input *ObsObjInfo `json:"input,omitempty"`

	Output *ObsObjInfo `json:"output,omitempty"`

	// 用户数据。
	UserData *string `json:"user_data,omitempty"`

	// 转码任务错误码。
	ErrorCode *string `json:"error_code,omitempty"`

	// 转码任务描述，当转码出现异常时，此字段为异常的原因。
	Description *string `json:"description,omitempty"`

	MediaDetail *MediaDetail `json:"media_detail,omitempty"`

	XcodeError *ErrorResponse `json:"xcode_error,omitempty"`
}

func (o TaskDetailInfo) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TaskDetailInfo struct{}"
	}

	return strings.Join([]string{"TaskDetailInfo", string(data)}, " ")
}

type TaskDetailInfoStatus struct {
	value string
}

type TaskDetailInfoStatusEnum struct {
	NO_TASK          TaskDetailInfoStatus
	WAITING          TaskDetailInfoStatus
	TRANSCODING      TaskDetailInfoStatus
	SUCCEEDED        TaskDetailInfoStatus
	FAILED           TaskDetailInfoStatus
	CANCELED         TaskDetailInfoStatus
	NEED_TO_BE_AUDIT TaskDetailInfoStatus
}

func GetTaskDetailInfoStatusEnum() TaskDetailInfoStatusEnum {
	return TaskDetailInfoStatusEnum{
		NO_TASK: TaskDetailInfoStatus{
			value: "NO_TASK",
		},
		WAITING: TaskDetailInfoStatus{
			value: "WAITING",
		},
		TRANSCODING: TaskDetailInfoStatus{
			value: "TRANSCODING",
		},
		SUCCEEDED: TaskDetailInfoStatus{
			value: "SUCCEEDED",
		},
		FAILED: TaskDetailInfoStatus{
			value: "FAILED",
		},
		CANCELED: TaskDetailInfoStatus{
			value: "CANCELED",
		},
		NEED_TO_BE_AUDIT: TaskDetailInfoStatus{
			value: "NEED_TO_BE_AUDIT",
		},
	}
}

func (c TaskDetailInfoStatus) Value() string {
	return c.value
}

func (c TaskDetailInfoStatus) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *TaskDetailInfoStatus) UnmarshalJSON(b []byte) error {
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
