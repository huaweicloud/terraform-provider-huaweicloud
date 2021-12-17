package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// 目的端服务器关联磁盘
type DiskIntargetServer struct {
	// 磁盘名称

	Name string `json:"name"`
	// 磁盘大小，单位：字节

	Size int64 `json:"size"`
	// 磁盘的作用

	DeviceUse *DiskIntargetServerDeviceUse `json:"device_use,omitempty"`
}

func (o DiskIntargetServer) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "DiskIntargetServer struct{}"
	}

	return strings.Join([]string{"DiskIntargetServer", string(data)}, " ")
}

type DiskIntargetServerDeviceUse struct {
	value string
}

type DiskIntargetServerDeviceUseEnum struct {
	BOOT   DiskIntargetServerDeviceUse
	OS     DiskIntargetServerDeviceUse
	NORMAL DiskIntargetServerDeviceUse
}

func GetDiskIntargetServerDeviceUseEnum() DiskIntargetServerDeviceUseEnum {
	return DiskIntargetServerDeviceUseEnum{
		BOOT: DiskIntargetServerDeviceUse{
			value: "BOOT",
		},
		OS: DiskIntargetServerDeviceUse{
			value: "OS",
		},
		NORMAL: DiskIntargetServerDeviceUse{
			value: "NORMAL",
		},
	}
}

func (c DiskIntargetServerDeviceUse) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *DiskIntargetServerDeviceUse) UnmarshalJSON(b []byte) error {
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
