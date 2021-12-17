package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// 任务关联的目的端信息
type TargetServerAssociatedWithTask struct {
	// 目的端在SMS数据库中的ID

	Id *string `json:"id,omitempty"`
	// 目的端虚机id

	VmId *string `json:"vm_id,omitempty"`
	// 目的端服务器名称

	Name *string `json:"name,omitempty"`
	// 目的端服务器ip

	Ip *string `json:"ip,omitempty"`
	// 目的端服务器的OS类型

	OsType *TargetServerAssociatedWithTaskOsType `json:"os_type,omitempty"`
	// 操作系统版本

	OsVersion *string `json:"os_version,omitempty"`
}

func (o TargetServerAssociatedWithTask) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TargetServerAssociatedWithTask struct{}"
	}

	return strings.Join([]string{"TargetServerAssociatedWithTask", string(data)}, " ")
}

type TargetServerAssociatedWithTaskOsType struct {
	value string
}

type TargetServerAssociatedWithTaskOsTypeEnum struct {
	WINDOWS TargetServerAssociatedWithTaskOsType
	LINUX   TargetServerAssociatedWithTaskOsType
}

func GetTargetServerAssociatedWithTaskOsTypeEnum() TargetServerAssociatedWithTaskOsTypeEnum {
	return TargetServerAssociatedWithTaskOsTypeEnum{
		WINDOWS: TargetServerAssociatedWithTaskOsType{
			value: "WINDOWS",
		},
		LINUX: TargetServerAssociatedWithTaskOsType{
			value: "LINUX",
		},
	}
}

func (c TargetServerAssociatedWithTaskOsType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *TargetServerAssociatedWithTaskOsType) UnmarshalJSON(b []byte) error {
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
