package model

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/utils"

	"errors"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/converter"

	"strings"
)

// 目的端服务器
type TargetServer struct {
	// 源端在SMS数据库中的ID

	Id *string `json:"id,omitempty"`
	// 源端服务器ip，注册源端时必选，更新非必选

	Ip string `json:"ip"`
	// 目的端服务器名称

	Name string `json:"name"`
	// 源端主机名，注册源端必选，更新非必选

	Hostname *string `json:"hostname,omitempty"`
	// 源端服务器的OS类型，分为Windows和Linux，注册必选，更新非必选

	OsType TargetServerOsType `json:"os_type"`
	// 操作系统版本，注册必选，更新非必选

	OsVersion *string `json:"os_version,omitempty"`
	// 源端服务器启动类型，如BIOS或者UEFI

	Firmware *TargetServerFirmware `json:"firmware,omitempty"`
	// CPU个数，单位vCPU

	CpuQuantity *int32 `json:"cpu_quantity,omitempty"`
	// 内存大小，单位MB

	Memory *int64 `json:"memory,omitempty"`
	// 目的端磁盘信息，一般和源端保持一致

	Disks []TargetDisk `json:"disks"`
	// Linux 必选，源端的Btrfs信息。如果源端不存在Btrfs，则为[]

	BtrfsList *[]string `json:"btrfs_list,omitempty"`
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

	BootLoader *TargetServerBootLoader `json:"boot_loader,omitempty"`
	// Windows必选，系统目录

	SystemDir *string `json:"system_dir,omitempty"`
	// lvm信息，一般和源端保持一致

	VolumeGroups *[]VolumeGroups `json:"volume_groups,omitempty"`
	// 目的端服务器ID，自动创建虚拟机不需要这个参数

	VmId *string `json:"vm_id,omitempty"`
	// 目的端服务器的规格

	Flavor *string `json:"flavor,omitempty"`
	// 目的端代理镜像磁盘id

	ImageDiskId *string `json:"image_disk_id,omitempty"`
	// 目的端快照id

	SnapshotIds *string `json:"snapshot_ids,omitempty"`
	// 目的端回滚快照id

	CutoveredSnapshotIds *string `json:"cutovered_snapshot_ids,omitempty"`
}

func (o TargetServer) String() string {
	data, err := utils.Marshal(o)
	if err != nil {
		return "TargetServer struct{}"
	}

	return strings.Join([]string{"TargetServer", string(data)}, " ")
}

type TargetServerOsType struct {
	value string
}

type TargetServerOsTypeEnum struct {
	WINDOWS TargetServerOsType
	LINUX   TargetServerOsType
}

func GetTargetServerOsTypeEnum() TargetServerOsTypeEnum {
	return TargetServerOsTypeEnum{
		WINDOWS: TargetServerOsType{
			value: "WINDOWS",
		},
		LINUX: TargetServerOsType{
			value: "LINUX",
		},
	}
}

func (c TargetServerOsType) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *TargetServerOsType) UnmarshalJSON(b []byte) error {
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

type TargetServerFirmware struct {
	value string
}

type TargetServerFirmwareEnum struct {
	BIOS TargetServerFirmware
	UEFI TargetServerFirmware
}

func GetTargetServerFirmwareEnum() TargetServerFirmwareEnum {
	return TargetServerFirmwareEnum{
		BIOS: TargetServerFirmware{
			value: "BIOS",
		},
		UEFI: TargetServerFirmware{
			value: "UEFI",
		},
	}
}

func (c TargetServerFirmware) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *TargetServerFirmware) UnmarshalJSON(b []byte) error {
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

type TargetServerBootLoader struct {
	value string
}

type TargetServerBootLoaderEnum struct {
	GRUB TargetServerBootLoader
	LILO TargetServerBootLoader
}

func GetTargetServerBootLoaderEnum() TargetServerBootLoaderEnum {
	return TargetServerBootLoaderEnum{
		GRUB: TargetServerBootLoader{
			value: "GRUB",
		},
		LILO: TargetServerBootLoader{
			value: "LILO",
		},
	}
}

func (c TargetServerBootLoader) MarshalJSON() ([]byte, error) {
	return utils.Marshal(c.value)
}

func (c *TargetServerBootLoader) UnmarshalJSON(b []byte) error {
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
