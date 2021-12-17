package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// This is a auto create Body Object
type UpdateTaskStatusReq struct {
	// 操作任务的具体动作

	Operation UpdateTaskStatusReqOperation `json:"operation"`
	// 操作参数

	Param map[string]string `json:"param,omitempty"`
}

func (o UpdateTaskStatusReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "UpdateTaskStatusReq struct{}"
	}

	return strings.Join([]string{"UpdateTaskStatusReq", string(data)}, " ")
}

type UpdateTaskStatusReqOperation struct {
	value string
}

type UpdateTaskStatusReqOperationEnum struct {
	START                UpdateTaskStatusReqOperation
	STOP                 UpdateTaskStatusReqOperation
	COLLECT_LOG          UpdateTaskStatusReqOperation
	TEST                 UpdateTaskStatusReqOperation
	CLONE_TEST           UpdateTaskStatusReqOperation
	RESTART              UpdateTaskStatusReqOperation
	SYNC_FAILED_ROLLBACK UpdateTaskStatusReqOperation
}

func GetUpdateTaskStatusReqOperationEnum() UpdateTaskStatusReqOperationEnum {
	return UpdateTaskStatusReqOperationEnum{
		START: UpdateTaskStatusReqOperation{
			value: "start",
		},
		STOP: UpdateTaskStatusReqOperation{
			value: "stop",
		},
		COLLECT_LOG: UpdateTaskStatusReqOperation{
			value: "collect_log",
		},
		TEST: UpdateTaskStatusReqOperation{
			value: "test",
		},
		CLONE_TEST: UpdateTaskStatusReqOperation{
			value: "clone_test",
		},
		RESTART: UpdateTaskStatusReqOperation{
			value: "restart",
		},
		SYNC_FAILED_ROLLBACK: UpdateTaskStatusReqOperation{
			value: "sync_failed_rollback",
		},
	}
}

func (c UpdateTaskStatusReqOperation) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *UpdateTaskStatusReqOperation) UnmarshalJSON(b []byte) error {
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
