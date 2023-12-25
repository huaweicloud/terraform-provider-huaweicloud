package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type BatchRestartOrDeleteInstanceReq struct {

	// 实例的ID列表。
	Instances *[]string `json:"instances,omitempty"`

	// 对实例的操作：restart、delete
	Action BatchRestartOrDeleteInstanceReqAction `json:"action"`

	// 参数值为kafka，表示删除租户所有创建失败的Kafka实例。
	AllFailure *BatchRestartOrDeleteInstanceReqAllFailure `json:"all_failure,omitempty"`
}

func (o BatchRestartOrDeleteInstanceReq) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "BatchRestartOrDeleteInstanceReq struct{}"
	}

	return strings.Join([]string{"BatchRestartOrDeleteInstanceReq", string(data)}, " ")
}

type BatchRestartOrDeleteInstanceReqAction struct {
	value string
}

type BatchRestartOrDeleteInstanceReqActionEnum struct {
	RESTART BatchRestartOrDeleteInstanceReqAction
	DELETE  BatchRestartOrDeleteInstanceReqAction
}

func GetBatchRestartOrDeleteInstanceReqActionEnum() BatchRestartOrDeleteInstanceReqActionEnum {
	return BatchRestartOrDeleteInstanceReqActionEnum{
		RESTART: BatchRestartOrDeleteInstanceReqAction{
			value: "restart",
		},
		DELETE: BatchRestartOrDeleteInstanceReqAction{
			value: "delete",
		},
	}
}

func (c BatchRestartOrDeleteInstanceReqAction) Value() string {
	return c.value
}

func (c BatchRestartOrDeleteInstanceReqAction) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *BatchRestartOrDeleteInstanceReqAction) UnmarshalJSON(b []byte) error {
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

type BatchRestartOrDeleteInstanceReqAllFailure struct {
	value string
}

type BatchRestartOrDeleteInstanceReqAllFailureEnum struct {
	KAFKA BatchRestartOrDeleteInstanceReqAllFailure
}

func GetBatchRestartOrDeleteInstanceReqAllFailureEnum() BatchRestartOrDeleteInstanceReqAllFailureEnum {
	return BatchRestartOrDeleteInstanceReqAllFailureEnum{
		KAFKA: BatchRestartOrDeleteInstanceReqAllFailure{
			value: "kafka",
		},
	}
}

func (c BatchRestartOrDeleteInstanceReqAllFailure) Value() string {
	return c.value
}

func (c BatchRestartOrDeleteInstanceReqAllFailure) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *BatchRestartOrDeleteInstanceReqAllFailure) UnmarshalJSON(b []byte) error {
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
