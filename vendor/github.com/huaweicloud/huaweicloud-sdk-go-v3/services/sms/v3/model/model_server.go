package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// 源端服务器
type Server struct {
	// 源端在SMS数据库中的ID

	Id *string `json:"id,omitempty"`
	// 源端服务器ip，注册源端时必选，更新非必选

	Ip string `json:"ip"`
	// 用来区分不同源端服务器的名称

	Name string `json:"name"`
	// 源端主机名，注册源端必选，更新非必选

	Hostname *string `json:"hostname,omitempty"`
	// 源端服务器的OS类型，分为Windows和Linux，注册必选，更新非必选

	OsType ServerOsType `json:"os_type"`
	// 操作系统版本，注册必选，更新非必选

	OsVersion *string `json:"os_version,omitempty"`
	// 源端服务器启动类型，如BIOS或者UEFI

	Firmware *ServerFirmware `json:"firmware,omitempty"`
	// CPU个数，单位vCPU

	CpuQuantity *int32 `json:"cpu_quantity,omitempty"`
	// 内存大小，单位MB

	Memory *int64 `json:"memory,omitempty"`
	// 源端服务器的磁盘信息

	Disks *[]Disk `json:"disks,omitempty"`
	// Linux 必选，源端的Btrfs信息。如果源端不存在Btrfs，则为[]

	BtrfsList *[]BtrfsFileSystem `json:"btrfs_list,omitempty"`
	// 源端服务器的网卡信息

	Networks *[]NetWork `json:"networks,omitempty"`
	// 租户的domainId

	DomainId *string `json:"domain_id,omitempty"`
	// 是否安装rsync组件，Linux系统此参数为必选

	HasRsync *bool `json:"has_rsync,omitempty"`
	// Linux场景必选，源端是否是半虚拟化

	Paravirtualization *bool `json:"paravirtualization,omitempty"`
	// Linux必选，裸设备列表

	RawDevices *string `json:"raw_devices,omitempty"`
	// Windows 必选，是否缺少驱动文件

	DriverFiles *bool `json:"driver_files,omitempty"`
	// Windows必选，是否存在不正常服务

	SystemServices *bool `json:"system_services,omitempty"`
	// Windows必选，权限是否满足要求

	AccountRights *bool `json:"account_rights,omitempty"`
	// Linux必选，系统引导类型，BOOT_LOADER(GRUB/LILO)

	BootLoader *ServerBootLoader `json:"boot_loader,omitempty"`
	// Windows必选，系统目录

	SystemDir *string `json:"system_dir,omitempty"`
	// Linux必选，如果没有卷组，输入[]

	VolumeGroups *[]VolumeGroups `json:"volume_groups,omitempty"`
}

func (o Server) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "Server struct{}"
	}

	return strings.Join([]string{"Server", string(data)}, " ")
}

type ServerOsType struct {
	value string
}

type ServerOsTypeEnum struct {
	WINDOWS ServerOsType
	LINUX   ServerOsType
}

func GetServerOsTypeEnum() ServerOsTypeEnum {
	return ServerOsTypeEnum{
		WINDOWS: ServerOsType{
			value: "WINDOWS",
		},
		LINUX: ServerOsType{
			value: "LINUX",
		},
	}
}

func (c ServerOsType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ServerOsType) UnmarshalJSON(b []byte) error {
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

type ServerFirmware struct {
	value string
}

type ServerFirmwareEnum struct {
	BIOS ServerFirmware
	UEFI ServerFirmware
}

func GetServerFirmwareEnum() ServerFirmwareEnum {
	return ServerFirmwareEnum{
		BIOS: ServerFirmware{
			value: "BIOS",
		},
		UEFI: ServerFirmware{
			value: "UEFI",
		},
	}
}

func (c ServerFirmware) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ServerFirmware) UnmarshalJSON(b []byte) error {
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

type ServerBootLoader struct {
	value string
}

type ServerBootLoaderEnum struct {
	GRUB ServerBootLoader
	LILO ServerBootLoader
}

func GetServerBootLoaderEnum() ServerBootLoaderEnum {
	return ServerBootLoaderEnum{
		GRUB: ServerBootLoader{
			value: "GRUB",
		},
		LILO: ServerBootLoader{
			value: "LILO",
		},
	}
}

func (c ServerBootLoader) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *ServerBootLoader) UnmarshalJSON(b []byte) error {
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
