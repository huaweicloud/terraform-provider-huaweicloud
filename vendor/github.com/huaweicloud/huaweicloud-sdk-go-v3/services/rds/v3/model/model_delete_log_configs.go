package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type DeleteLogConfigs struct {

	// 实例ID。
	InstanceId string `json:"instance_id"`

	// 日志类型。
	LogType DeleteLogConfigsLogType `json:"log_type"`
}

func (o DeleteLogConfigs) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DeleteLogConfigs struct{}"
	}

	return strings.Join([]string{"DeleteLogConfigs", string(data)}, " ")
}

type DeleteLogConfigsLogType struct {
	value string
}

type DeleteLogConfigsLogTypeEnum struct {
	ERROR_LOG DeleteLogConfigsLogType
	SLOW_LOG  DeleteLogConfigsLogType
	AUDIT_LOG DeleteLogConfigsLogType
}

func GetDeleteLogConfigsLogTypeEnum() DeleteLogConfigsLogTypeEnum {
	return DeleteLogConfigsLogTypeEnum{
		ERROR_LOG: DeleteLogConfigsLogType{
			value: "error_log",
		},
		SLOW_LOG: DeleteLogConfigsLogType{
			value: "slow_log",
		},
		AUDIT_LOG: DeleteLogConfigsLogType{
			value: "audit_log",
		},
	}
}

func (c DeleteLogConfigsLogType) Value() string {
	return c.value
}

func (c DeleteLogConfigsLogType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *DeleteLogConfigsLogType) UnmarshalJSON(b []byte) error {
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
