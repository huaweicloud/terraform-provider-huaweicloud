package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

type AddLogConfigs struct {

	// 实例ID。
	InstanceId string `json:"instance_id"`

	// 日志类型。
	LogType AddLogConfigsLogType `json:"log_type"`

	// LTS日志组ID。
	LtsGroupId string `json:"lts_group_id"`

	// LTS日志流ID。
	LtsStreamId string `json:"lts_stream_id"`
}

func (o AddLogConfigs) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "AddLogConfigs struct{}"
	}

	return strings.Join([]string{"AddLogConfigs", string(data)}, " ")
}

type AddLogConfigsLogType struct {
	value string
}

type AddLogConfigsLogTypeEnum struct {
	ERROR_LOG AddLogConfigsLogType
	SLOW_LOG  AddLogConfigsLogType
	AUDIT_LOG AddLogConfigsLogType
}

func GetAddLogConfigsLogTypeEnum() AddLogConfigsLogTypeEnum {
	return AddLogConfigsLogTypeEnum{
		ERROR_LOG: AddLogConfigsLogType{
			value: "error_log",
		},
		SLOW_LOG: AddLogConfigsLogType{
			value: "slow_log",
		},
		AUDIT_LOG: AddLogConfigsLogType{
			value: "audit_log",
		},
	}
}

func (c AddLogConfigsLogType) Value() string {
	return c.value
}

func (c AddLogConfigsLogType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *AddLogConfigsLogType) UnmarshalJSON(b []byte) error {
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
