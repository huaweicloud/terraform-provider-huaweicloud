package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// 磁盘信息
type ServerDisk struct {
	// 磁盘名称

	Name string `json:"name"`
	// 磁盘的分区类型，添加源端时源端磁盘必选

	PartitionStyle *ServerDiskPartitionStyle `json:"partition_style,omitempty"`
	// 磁盘类型

	DeviceUse ServerDiskDeviceUse `json:"device_use"`
	// 磁盘总大小，以字节为单位

	Size int64 `json:"size"`
	// 磁盘已使用大小，以字节为单位

	UsedSize int64 `json:"used_size"`
	// 磁盘上的物理分区信息

	PhysicalVolumes []PhysicalVolume `json:"physical_volumes"`
	// 是否为系统盘

	OsDisk *bool `json:"os_disk,omitempty"`
	// Linux系统 目的端ECS中与源端关联的磁盘名称

	RelationName *string `json:"relation_name,omitempty"`
}

func (o ServerDisk) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "ServerDisk struct{}"
	}

	return strings.Join([]string{"ServerDisk", string(data)}, " ")
}

type ServerDiskPartitionStyle struct {
	value string
}

type ServerDiskPartitionStyleEnum struct {
	MBR ServerDiskPartitionStyle
	GPT ServerDiskPartitionStyle
}

func GetServerDiskPartitionStyleEnum() ServerDiskPartitionStyleEnum {
	return ServerDiskPartitionStyleEnum{
		MBR: ServerDiskPartitionStyle{
			value: "MBR",
		},
		GPT: ServerDiskPartitionStyle{
			value: "GPT",
		},
	}
}

func (c ServerDiskPartitionStyle) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ServerDiskPartitionStyle) UnmarshalJSON(b []byte) error {
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

type ServerDiskDeviceUse struct {
	value string
}

type ServerDiskDeviceUseEnum struct {
	BOOT ServerDiskDeviceUse
	OS   ServerDiskDeviceUse
}

func GetServerDiskDeviceUseEnum() ServerDiskDeviceUseEnum {
	return ServerDiskDeviceUseEnum{
		BOOT: ServerDiskDeviceUse{
			value: "BOOT",
		},
		OS: ServerDiskDeviceUse{
			value: "OS",
		},
	}
}

func (c ServerDiskDeviceUse) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ServerDiskDeviceUse) UnmarshalJSON(b []byte) error {
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
