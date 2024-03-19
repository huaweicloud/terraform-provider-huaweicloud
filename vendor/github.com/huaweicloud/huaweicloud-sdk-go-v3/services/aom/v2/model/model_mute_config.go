package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// MuteConfig 静默规则的生效时间配置
type MuteConfig struct {

	// 静默规则结束时间
	EndsAt *int64 `json:"ends_at,omitempty"`

	// 当type为每周或者每月时，scope不能为空
	Scope *[]int32 `json:"scope,omitempty"`

	// 静默规则开始时间
	StartsAt int64 `json:"starts_at"`

	// 静默规则生效时间种类。FIXED：固定方式统计，DAILY：按日合计，WEEKLY：按周统计，MONTHLY：按月统计
	Type MuteConfigType `json:"type"`
}

func (o MuteConfig) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "MuteConfig struct{}"
	}

	return strings.Join([]string{"MuteConfig", string(data)}, " ")
}

type MuteConfigType struct {
	value string
}

type MuteConfigTypeEnum struct {
	FIXED   MuteConfigType
	DAILY   MuteConfigType
	WEEKLY  MuteConfigType
	MONTHLY MuteConfigType
}

func GetMuteConfigTypeEnum() MuteConfigTypeEnum {
	return MuteConfigTypeEnum{
		FIXED: MuteConfigType{
			value: "FIXED",
		},
		DAILY: MuteConfigType{
			value: "DAILY",
		},
		WEEKLY: MuteConfigType{
			value: "WEEKLY",
		},
		MONTHLY: MuteConfigType{
			value: "MONTHLY",
		},
	}
}

func (c MuteConfigType) Value() string {
	return c.value
}

func (c MuteConfigType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *MuteConfigType) UnmarshalJSON(b []byte) error {
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
