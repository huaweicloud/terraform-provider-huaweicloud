package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// TaskType 集群升级任务类型： Cluster: 集群升级任务 PreCheck: 集群升级预检查任务 Rollback: 集群升级回归任务 Snapshot: 集群升级快照任务 PostCheck: 集群升级后检查任务
type TaskType struct {
	value string
}

type TaskTypeEnum struct {
	CLUSTER    TaskType
	PRE_CHECK  TaskType
	ROLLBACK   TaskType
	SNAPSHOT   TaskType
	POST_CHECK TaskType
}

func GetTaskTypeEnum() TaskTypeEnum {
	return TaskTypeEnum{
		CLUSTER: TaskType{
			value: "Cluster",
		},
		PRE_CHECK: TaskType{
			value: "PreCheck",
		},
		ROLLBACK: TaskType{
			value: "Rollback",
		},
		SNAPSHOT: TaskType{
			value: "Snapshot",
		},
		POST_CHECK: TaskType{
			value: "PostCheck",
		},
	}
}

func (c TaskType) Value() string {
	return c.value
}

func (c TaskType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *TaskType) UnmarshalJSON(b []byte) error {
	myConverter := converter.StringConverterFactory("string")
	if myConverter == nil {
		return errors.New("unsupported StringConverter type: string")
	}

	interf, err := myConverter.CovertStringToInterface(strings.Trim(string(b[:]), "\""))
	if err != nil {
		return err
	}

	if val, ok := interf.(string); ok {
		c.value = val
		return nil
	} else {
		return errors.New("convert enum data to string error")
	}
}
