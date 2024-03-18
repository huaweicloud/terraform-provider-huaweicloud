package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// UpgradeWorkflowTaskStatus 集群升级状态： Init: 任务初始状态 Queuing: 任务已进入执行队列 Running: 任务开始执行 Success: 任务执行成功 Failed: 任务执行失败
type UpgradeWorkflowTaskStatus struct {
	value string
}

type UpgradeWorkflowTaskStatusEnum struct {
	INIT    UpgradeWorkflowTaskStatus
	QUEUING UpgradeWorkflowTaskStatus
	RUNNING UpgradeWorkflowTaskStatus
	SUCCESS UpgradeWorkflowTaskStatus
	FAILED  UpgradeWorkflowTaskStatus
}

func GetUpgradeWorkflowTaskStatusEnum() UpgradeWorkflowTaskStatusEnum {
	return UpgradeWorkflowTaskStatusEnum{
		INIT: UpgradeWorkflowTaskStatus{
			value: "Init",
		},
		QUEUING: UpgradeWorkflowTaskStatus{
			value: "Queuing",
		},
		RUNNING: UpgradeWorkflowTaskStatus{
			value: "Running",
		},
		SUCCESS: UpgradeWorkflowTaskStatus{
			value: "Success",
		},
		FAILED: UpgradeWorkflowTaskStatus{
			value: "Failed",
		},
	}
}

func (c UpgradeWorkflowTaskStatus) Value() string {
	return c.value
}

func (c UpgradeWorkflowTaskStatus) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *UpgradeWorkflowTaskStatus) UnmarshalJSON(b []byte) error {
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
