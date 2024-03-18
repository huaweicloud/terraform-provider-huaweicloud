package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// WorkFlowPhase 集群升级流程的执行状态： Init: 表示该升级流程中还未有任何任务开始运行 Running: 表示该升级流程中已有任务开始执行 Pending: 表示该升级流程中有任务执行失败 Success: 表示该升级流程中所有任务都已执行成功 Cancel: 表示该升级流程已被取消
type WorkFlowPhase struct {
	value string
}

type WorkFlowPhaseEnum struct {
	INIT    WorkFlowPhase
	RUNNING WorkFlowPhase
	PENDING WorkFlowPhase
	SUCCESS WorkFlowPhase
	CANCEL  WorkFlowPhase
}

func GetWorkFlowPhaseEnum() WorkFlowPhaseEnum {
	return WorkFlowPhaseEnum{
		INIT: WorkFlowPhase{
			value: "Init",
		},
		RUNNING: WorkFlowPhase{
			value: "Running",
		},
		PENDING: WorkFlowPhase{
			value: "Pending",
		},
		SUCCESS: WorkFlowPhase{
			value: "Success",
		},
		CANCEL: WorkFlowPhase{
			value: "Cancel",
		},
	}
}

func (c WorkFlowPhase) Value() string {
	return c.value
}

func (c WorkFlowPhase) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *WorkFlowPhase) UnmarshalJSON(b []byte) error {
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
