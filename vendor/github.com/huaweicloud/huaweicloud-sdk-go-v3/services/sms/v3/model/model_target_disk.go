package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// 目的端磁盘
type TargetDisk struct {
	// 判断是普通分区，启动分区还是系统分区

	DeviceUse *TargetDiskDeviceUse `json:"device_use,omitempty"`
	// 磁盘id

	DiskId *string `json:"disk_id,omitempty"`
	// 磁盘名称

	Name *string `json:"name,omitempty"`
	// 逻辑卷信息

	PhysicalVolumes *[]TargetPhysicalVolumes `json:"physical_volumes,omitempty"`
	// 大小

	Size *int64 `json:"size,omitempty"`
	// 已使用大小

	UsedSize *int64 `json:"used_size,omitempty"`
}

func (o TargetDisk) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TargetDisk struct{}"
	}

	return strings.Join([]string{"TargetDisk", string(data)}, " ")
}

type TargetDiskDeviceUse struct {
	value string
}

type TargetDiskDeviceUseEnum struct {
	NORMAL TargetDiskDeviceUse
	OS     TargetDiskDeviceUse
	BOOT   TargetDiskDeviceUse
}

func GetTargetDiskDeviceUseEnum() TargetDiskDeviceUseEnum {
	return TargetDiskDeviceUseEnum{
		NORMAL: TargetDiskDeviceUse{
			value: "NORMAL",
		},
		OS: TargetDiskDeviceUse{
			value: "OS",
		},
		BOOT: TargetDiskDeviceUse{
			value: "BOOT",
		},
	}
}

func (c TargetDiskDeviceUse) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *TargetDiskDeviceUse) UnmarshalJSON(b []byte) error {
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
