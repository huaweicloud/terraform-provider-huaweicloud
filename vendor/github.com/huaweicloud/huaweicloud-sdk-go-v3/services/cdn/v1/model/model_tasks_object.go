package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type TasksObject struct {

	// 任务id。
	Id *string `json:"id,omitempty"`

	// 任务的类型， 其值可以为refresh或preheating。
	TaskType *TasksObjectTaskType `json:"task_type,omitempty"`

	// 刷新结果。task_done表示刷新成功  ，task_inprocess表示刷新中。
	Status *string `json:"status,omitempty"`

	// 处理中的url个数。
	Processing *int32 `json:"processing,omitempty"`

	// 成功处理的url个数。
	Succeed *int32 `json:"succeed,omitempty"`

	// 处理失败的url个数。
	Failed *int32 `json:"failed,omitempty"`

	// url总数。
	Total *int32 `json:"total,omitempty"`

	// 任务的创建时间，相对于UTC 1970-01-01到当前时间相隔的毫秒数。
	CreateTime *int64 `json:"create_time,omitempty"`

	// 默认是文件file。file：文件,directory：目录。
	FileType *TasksObjectFileType `json:"file_type,omitempty"`
}

func (o TasksObject) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TasksObject struct{}"
	}

	return strings.Join([]string{"TasksObject", string(data)}, " ")
}

type TasksObjectTaskType struct {
	value string
}

type TasksObjectTaskTypeEnum struct {
	REFRESH    TasksObjectTaskType
	PREHEATING TasksObjectTaskType
}

func GetTasksObjectTaskTypeEnum() TasksObjectTaskTypeEnum {
	return TasksObjectTaskTypeEnum{
		REFRESH: TasksObjectTaskType{
			value: "refresh",
		},
		PREHEATING: TasksObjectTaskType{
			value: "preheating",
		},
	}
}

func (c TasksObjectTaskType) Value() string {
	return c.value
}

func (c TasksObjectTaskType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *TasksObjectTaskType) UnmarshalJSON(b []byte) error {
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

type TasksObjectFileType struct {
	value string
}

type TasksObjectFileTypeEnum struct {
	FILE      TasksObjectFileType
	DIRECTORY TasksObjectFileType
}

func GetTasksObjectFileTypeEnum() TasksObjectFileTypeEnum {
	return TasksObjectFileTypeEnum{
		FILE: TasksObjectFileType{
			value: "file",
		},
		DIRECTORY: TasksObjectFileType{
			value: "directory",
		},
	}
}

func (c TasksObjectFileType) Value() string {
	return c.value
}

func (c TasksObjectFileType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *TasksObjectFileType) UnmarshalJSON(b []byte) error {
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
