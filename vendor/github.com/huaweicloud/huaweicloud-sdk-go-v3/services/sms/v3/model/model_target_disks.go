package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// 目的端磁盘信息
type TargetDisks struct {
	// 磁盘类型，普通磁盘，OS所在磁盘，BOOT所在磁盘

	DeviceUse *TargetDisksDeviceUse `json:"device_use,omitempty"`
	// 磁盘id,自动创建虚拟机不用设置

	DiskId *string `json:"disk_id,omitempty"`
	// 名称，根据磁盘顺序设置为disk X

	Name string `json:"name"`
	// 物理卷信息

	PhysicalVolumes []PhysicalVolumes `json:"physical_volumes"`
	// 大小

	Size int64 `json:"size"`
	// 使用大小

	UsedSize int64 `json:"used_size"`
}

func (o TargetDisks) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TargetDisks struct{}"
	}

	return strings.Join([]string{"TargetDisks", string(data)}, " ")
}

type TargetDisksDeviceUse struct {
	value string
}

type TargetDisksDeviceUseEnum struct {
	NORMAL TargetDisksDeviceUse
	OS     TargetDisksDeviceUse
	BOOT   TargetDisksDeviceUse
}

func GetTargetDisksDeviceUseEnum() TargetDisksDeviceUseEnum {
	return TargetDisksDeviceUseEnum{
		NORMAL: TargetDisksDeviceUse{
			value: "NORMAL",
		},
		OS: TargetDisksDeviceUse{
			value: "OS",
		},
		BOOT: TargetDisksDeviceUse{
			value: "BOOT",
		},
	}
}

func (c TargetDisksDeviceUse) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *TargetDisksDeviceUse) UnmarshalJSON(b []byte) error {
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
